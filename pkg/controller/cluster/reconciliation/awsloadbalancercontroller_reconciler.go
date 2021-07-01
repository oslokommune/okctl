package reconciliation

import (
	"context"
	"fmt"

	"github.com/oslokommune/okctl/pkg/cfn"
	clientCore "github.com/oslokommune/okctl/pkg/client/core"
	"github.com/oslokommune/okctl/pkg/controller/common/reconciliation"

	"github.com/oslokommune/okctl/pkg/client"
)

const awsLoadBalancerControllerIdentifier = "AWS Load Balancer controller"

// albIngressReconciler contains service and metadata for the relevant resource
type awsLoadBalancerControllerReconciler struct {
	client client.AWSLoadBalancerControllerService
}

// Reconcile knows how to do what is necessary to ensure the desired state is achieved
func (z *awsLoadBalancerControllerReconciler) Reconcile(
	ctx context.Context,
	meta reconciliation.Metadata,
	state *clientCore.StateHandlers,
) (reconciliation.Result, error) {
	action, err := z.determineAction(ctx, meta, state)
	if err != nil {
		return reconciliation.Result{}, fmt.Errorf("determining course of action: %w", err)
	}

	switch action {
	case reconciliation.ActionCreate:
		vpc, err := state.Vpc.GetVpc(cfn.NewStackNamer().Vpc(meta.ClusterDeclaration.Metadata.Name))
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf("getting vpc: %w", err)
		}

		_, err = z.client.CreateAWSLoadBalancerController(ctx, client.CreateAWSLoadBalancerControllerOpts{
			ID:    reconciliation.ClusterMetaAsID(meta.ClusterDeclaration.Metadata),
			VPCID: vpc.VpcID,
		})
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf("creating aws load balancer controller: %w", err)
		}

		return reconciliation.Result{Requeue: false}, nil
	case reconciliation.ActionDelete:
		err = z.client.DeleteAWSLoadBalancerController(
			ctx,
			reconciliation.ClusterMetaAsID(meta.ClusterDeclaration.Metadata),
		)
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf("deleting aws load balancer controller: %w", err)
		}

		return reconciliation.Result{Requeue: false}, nil
	case reconciliation.ActionWait:
		return reconciliation.Result{Requeue: true}, nil
	case reconciliation.ActionNoop:
		return reconciliation.Result{Requeue: false}, nil
	}

	return reconciliation.Result{}, fmt.Errorf("action %s is not implemented", string(action))
}

func (z *awsLoadBalancerControllerReconciler) determineAction(
	_ context.Context,
	meta reconciliation.Metadata,
	state *clientCore.StateHandlers,
) (reconciliation.Action, error) {
	userIndication := reconciliation.DetermineUserIndication(
		meta,
		meta.ClusterDeclaration.Integrations.AWSLoadBalancerController,
	)

	clusterExists, err := state.Cluster.HasCluster(meta.ClusterDeclaration.Metadata.Name)
	if err != nil {
		return reconciliation.ActionNoop, fmt.Errorf("acquiring cluster existence: %w", err)
	}

	awsLoadBalancerControllerExists := false
	if clusterExists {
		awsLoadBalancerControllerExists, err = state.AWSLoadBalancerController.HasAWSLoadBalancerController()
		if err != nil {
			return reconciliation.ActionNoop, fmt.Errorf("acquiring AWS Load Balancer Controller existence: %w", err)
		}
	}

	switch userIndication {
	case reconciliation.ActionCreate:
		if !clusterExists {
			return reconciliation.ActionWait, nil
		}

		if awsLoadBalancerControllerExists {
			return reconciliation.ActionNoop, nil
		}

		return reconciliation.ActionCreate, nil
	case reconciliation.ActionDelete:
		if !clusterExists || !awsLoadBalancerControllerExists {
			return reconciliation.ActionNoop, nil
		}

		dependenciesReady, err := reconciliation.AssertDependencyExistence(false,
			state.ArgoCD.HasArgoCD,
			state.Monitoring.HasKubePromStack,
		)
		if err != nil {
			return reconciliation.ActionNoop, fmt.Errorf("checking deletion dependencies: %w", err)
		}

		if !dependenciesReady {
			return reconciliation.ActionWait, nil
		}

		return reconciliation.ActionDelete, nil
	}

	return reconciliation.ActionNoop, reconciliation.ErrIndecisive
}

// String returns the identifier for this reconciler
func (z *awsLoadBalancerControllerReconciler) String() string {
	return awsLoadBalancerControllerIdentifier
}

// NewAWSLoadBalancerControllerReconciler creates a new reconciler for the aws load balancer controller resource
func NewAWSLoadBalancerControllerReconciler(client client.AWSLoadBalancerControllerService) reconciliation.Reconciler {
	return &awsLoadBalancerControllerReconciler{
		client: client,
	}
}
