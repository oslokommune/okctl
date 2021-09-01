package reconciliation

import (
	"context"
	"fmt"

	"github.com/oslokommune/okctl/pkg/config/constant"
	"github.com/oslokommune/okctl/pkg/servicequota"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	clientCore "github.com/oslokommune/okctl/pkg/client/core"
	"github.com/oslokommune/okctl/pkg/controller/common/reconciliation"

	"github.com/oslokommune/okctl/pkg/client"
)

const vpcReconcilerIdentifier = "virtual private network"

// vpcReconciler contains service and metadata for the relevant resource
type vpcReconciler struct {
	client        client.VPCService
	cloudProvider v1alpha1.CloudProvider
}

// Reconcile knows how to do what is necessary to ensure the desired state is achieved
func (z *vpcReconciler) Reconcile(ctx context.Context, meta reconciliation.Metadata, state *clientCore.StateHandlers) (reconciliation.Result, error) {
	action, err := z.determineAction(meta, state)
	if err != nil {
		return reconciliation.Result{}, fmt.Errorf(constant.ReconcilerDetermineActionError, err)
	}

	switch action {
	case reconciliation.ActionCreate:
		_, err = z.client.CreateVpc(ctx, client.CreateVpcOpts{
			ID:      reconciliation.ClusterMetaAsID(meta.ClusterDeclaration.Metadata),
			Cidr:    meta.ClusterDeclaration.VPC.CIDR,
			Minimal: !meta.ClusterDeclaration.VPC.HighAvailability,
		})
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf(constant.CreateVPCError, err)
		}

		return reconciliation.Result{Requeue: false}, nil
	case reconciliation.ActionDelete:
		err = z.client.DeleteVpc(ctx, client.DeleteVpcOpts{
			ID: reconciliation.ClusterMetaAsID(meta.ClusterDeclaration.Metadata),
		})
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf(constant.DeleteVPCError, err)
		}

		return reconciliation.Result{Requeue: false}, nil
	case reconciliation.ActionWait:
		return reconciliation.Result{Requeue: true}, nil
	case reconciliation.ActionNoop:
		return reconciliation.Result{Requeue: false}, nil
	}

	return reconciliation.Result{}, fmt.Errorf(constant.ActionNotImplementedError, string(action))
}

func (z *vpcReconciler) determineAction(meta reconciliation.Metadata, state *clientCore.StateHandlers) (reconciliation.Action, error) {
	userIndication := reconciliation.DetermineUserIndication(meta, true)

	componentExists, err := state.Vpc.HasVPC(meta.ClusterDeclaration.Metadata.Name)
	if err != nil {
		return reconciliation.ActionNoop, fmt.Errorf(constant.QueryStateError, err)
	}

	switch userIndication {
	case reconciliation.ActionCreate:
		if componentExists {
			return reconciliation.ActionNoop, nil
		}

		err = servicequota.CheckQuotas(
			servicequota.NewVpcCheck(constant.DefaultRequiredVpcs, z.cloudProvider),
			servicequota.NewEipCheck(constant.DefaultRequiredEpis, z.cloudProvider),
			servicequota.NewIgwCheck(constant.DefaultRequiredIgws, z.cloudProvider),
		)
		if err != nil {
			return reconciliation.ActionNoop, fmt.Errorf(constant.CheckServiceQuotasError, err)
		}

		return reconciliation.ActionCreate, nil
	case reconciliation.ActionDelete:
		dependenciesReady, err := reconciliation.AssertDependencyExistence(false,
			reconciliation.GenerateClusterExistenceTest(state, meta.ClusterDeclaration.Metadata.Name),
			generatePostgresDBExistenceTest(state),
		)
		if err != nil {
			return reconciliation.ActionNoop, fmt.Errorf(constant.CheckDependenciesError, err)
		}

		if !dependenciesReady {
			return reconciliation.ActionWait, nil
		}

		if !componentExists {
			return reconciliation.ActionNoop, nil
		}

		return reconciliation.ActionDelete, nil
	}

	return reconciliation.ActionNoop, reconciliation.ErrIndecisive
}

func generatePostgresDBExistenceTest(state *clientCore.StateHandlers) func() (bool, error) {
	return func() (bool, error) {
		dbs, err := state.Component.GetPostgresDatabases()
		if err != nil {
			return false, fmt.Errorf(constant.CheckPostgresDatabaseStateError, err)
		}

		return len(dbs) != 0, nil
	}
}

// String returns the identifier for this reconciler
func (z *vpcReconciler) String() string {
	return vpcReconcilerIdentifier
}

// NewVPCReconciler creates a new reconciler for the VPC resource
func NewVPCReconciler(client client.VPCService, cloudProvider v1alpha1.CloudProvider) reconciliation.Reconciler {
	return &vpcReconciler{
		client:        client,
		cloudProvider: cloudProvider,
	}
}
