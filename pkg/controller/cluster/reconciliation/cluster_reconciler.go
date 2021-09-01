package reconciliation

import (
	"context"
	"fmt"

	"github.com/oslokommune/okctl/pkg/cleaner"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/servicequota"

	"github.com/oslokommune/okctl/pkg/cfn"
	"github.com/oslokommune/okctl/pkg/controller/common/reconciliation"

	clientCore "github.com/oslokommune/okctl/pkg/client/core"

	"github.com/oslokommune/okctl/pkg/config/constant"

	"github.com/oslokommune/okctl/pkg/client"
)

// clusterReconciler contains service and metadata for the relevant resource
type clusterReconciler struct {
	client        client.ClusterService
	cloudProvider v1alpha1.CloudProvider
}

const clusterReconcilerIdentifier = "kubernetes cluster"

// Reconcile knows how to do what is necessary to ensure the desired state is achieved
func (z *clusterReconciler) Reconcile(ctx context.Context, meta reconciliation.Metadata, state *clientCore.StateHandlers) (reconciliation.Result, error) {
	action, err := z.determineAction(meta, state)
	if err != nil {
		return reconciliation.Result{}, fmt.Errorf(constant.ReconcilerDetermineActionError, err)
	}

	switch action {
	case reconciliation.ActionCreate:
		vpc, err := state.Vpc.GetVpc(cfn.NewStackNamer().Vpc(meta.ClusterDeclaration.Metadata.Name))
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf(constant.GetVpcError, err)
		}

		_, err = z.client.CreateCluster(ctx, client.ClusterCreateOpts{
			ID:                reconciliation.ClusterMetaAsID(meta.ClusterDeclaration.Metadata),
			Cidr:              meta.ClusterDeclaration.VPC.CIDR,
			Version:           constant.DefaultEKSKubernetesVersion,
			VpcID:             vpc.VpcID,
			VpcPrivateSubnets: vpc.PrivateSubnets,
			VpcPublicSubnets:  vpc.PublicSubnets,
		})
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf(constant.CreateClusterError, err)
		}

		return reconciliation.Result{Requeue: false}, nil
	case reconciliation.ActionDelete:
		err = z.client.DeleteCluster(ctx, client.ClusterDeleteOpts{
			ID: reconciliation.ClusterMetaAsID(meta.ClusterDeclaration.Metadata),
		})
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf(constant.DeleteClusterError, err)
		}

		err = z.cleanUpDanglers(meta, state)
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf(constant.DeleteDanglingALBError, err)
		}

		return reconciliation.Result{Requeue: false}, nil
	case reconciliation.ActionWait:
		return reconciliation.Result{Requeue: true}, nil
	case reconciliation.ActionNoop:
		return reconciliation.Result{}, nil
	}

	return reconciliation.Result{}, fmt.Errorf(constant.ActionNotImplementedError, string(action))
}

func (z *clusterReconciler) determineAction(meta reconciliation.Metadata, state *clientCore.StateHandlers) (reconciliation.Action, error) {
	userIndication := reconciliation.DetermineUserIndication(meta, true)

	componentExists, err := state.Cluster.HasCluster(meta.ClusterDeclaration.Metadata.Name)
	if err != nil {
		return reconciliation.ActionNoop, fmt.Errorf(constant.CheckClusterExistanceError, err)
	}

	switch userIndication {
	case reconciliation.ActionCreate:
		if componentExists {
			return reconciliation.ActionNoop, nil
		}

		dependenciesReady, err := z.hasCreateDependenciesMet(meta, state)
		if err != nil {
			return reconciliation.ActionNoop, fmt.Errorf(constant.CheckDepedencyReadyError, err)
		}

		if !dependenciesReady {
			return reconciliation.ActionWait, nil
		}

		return reconciliation.ActionCreate, nil
	case reconciliation.ActionDelete:
		if !componentExists {
			return reconciliation.ActionNoop, nil
		}

		dependenciesReady, err := z.hasDeleteDependenciesMet(state)
		if err != nil {
			return reconciliation.ActionNoop, fmt.Errorf(constant.CheckDepedencyReadyError, err)
		}

		if !dependenciesReady {
			return reconciliation.ActionWait, nil
		}

		return reconciliation.ActionDelete, nil
	}

	return reconciliation.ActionNoop, reconciliation.ErrIndecisive
}

func (z *clusterReconciler) cleanUpDanglers(meta reconciliation.Metadata, state *clientCore.StateHandlers) error {
	vpc, err := state.Vpc.GetVpc(cfn.NewStackNamer().Vpc(meta.ClusterDeclaration.Metadata.Name))
	if err != nil {
		return fmt.Errorf(constant.GetVpcError, err)
	}

	clean := cleaner.New(z.cloudProvider)

	err = clean.DeleteDanglingALBs(vpc.VpcID)
	if err != nil {
		return fmt.Errorf(constant.CleanUpALBError, err)
	}

	err = clean.DeleteDanglingTargetGroups(meta.ClusterDeclaration.Metadata.Name)
	if err != nil {
		return fmt.Errorf(constant.CleanUpTragetGroupError, err)
	}

	return nil
}

func (z *clusterReconciler) hasCreateDependenciesMet(meta reconciliation.Metadata, state *clientCore.StateHandlers) (bool, error) {
	hasVPC, err := state.Vpc.HasVPC(meta.ClusterDeclaration.Metadata.Name)
	if err != nil {
		return false, fmt.Errorf(constant.GetVPCStateError, err)
	}

	err = servicequota.CheckQuotas(
		servicequota.NewFargateCheck(constant.DefaultRequiredFargateOnDemandPods, z.cloudProvider),
	)
	if err != nil {
		return false, fmt.Errorf(constant.CheckServiceQuotasError, err)
	}

	return hasVPC, nil
}

func (z *clusterReconciler) hasDeleteDependenciesMet(state *clientCore.StateHandlers) (bool, error) {
	ok, err := reconciliation.AssertDependencyExistence(false,
		state.ExternalDNS.HasExternalDNS,
		state.Monitoring.HasKubePromStack,
		state.Autoscaler.HasAutoscaler,
		state.AWSLoadBalancerController.HasAWSLoadBalancerController,
		state.Blockstorage.HasBlockstorage,
		state.ExternalSecrets.HasExternalSecrets,
		state.ArgoCD.HasArgoCD,
		state.Loki.HasLoki,
		state.Promtail.HasPromtail,
		state.Tempo.HasTempo,
	)
	if err != nil {
		return false, fmt.Errorf(constant.AssertExistenceError, err)
	}

	return ok, nil
}

// String returns the identifier for this reconciler
func (z *clusterReconciler) String() string {
	return clusterReconcilerIdentifier
}

// NewClusterReconciler creates a new reconciler for the cluster resource
func NewClusterReconciler(client client.ClusterService, cloudProvider v1alpha1.CloudProvider) reconciliation.Reconciler {
	return &clusterReconciler{
		client:        client,
		cloudProvider: cloudProvider,
	}
}
