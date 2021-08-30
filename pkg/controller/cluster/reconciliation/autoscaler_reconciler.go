package reconciliation

import (
	"context"
	"fmt"

	"github.com/oslokommune/okctl/pkg/config/constant"

	clientCore "github.com/oslokommune/okctl/pkg/client/core"
	"github.com/oslokommune/okctl/pkg/controller/common/reconciliation"

	"github.com/oslokommune/okctl/pkg/client"
)

// autoscalerReconciler contains service and metadata for the relevant resource
type autoscalerReconciler struct {
	client client.AutoscalerService
}

const autoscalerReconcilerIdentifier = "autoscaler"

// Reconcile knows how to do what is necessary to ensure the desired state is achieved
func (z *autoscalerReconciler) Reconcile(ctx context.Context, meta reconciliation.Metadata, state *clientCore.StateHandlers) (reconciliation.Result, error) {
	action, err := z.determineAction(ctx, meta, state)
	if err != nil {
		return reconciliation.Result{}, fmt.Errorf(constant.ReconcilerDetermineActionError, err)
	}

	switch action {
	case reconciliation.ActionCreate:
		_, err = z.client.CreateAutoscaler(ctx, client.CreateAutoscalerOpts{
			ID: reconciliation.ClusterMetaAsID(meta.ClusterDeclaration.Metadata),
		})
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf(constant.CreateAutoScalerError, err)
		}

		return reconciliation.Result{Requeue: false}, nil
	case reconciliation.ActionDelete:
		err = z.client.DeleteAutoscaler(
			ctx,
			reconciliation.ClusterMetaAsID(meta.ClusterDeclaration.Metadata),
		)
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf(constant.DeleteAutoScalerError, err)
		}

		return reconciliation.Result{Requeue: false}, nil
	case reconciliation.ActionWait:
		return reconciliation.Result{Requeue: true}, nil
	case reconciliation.ActionNoop:
		return reconciliation.Result{Requeue: false}, nil
	}

	return reconciliation.Result{}, fmt.Errorf(constant.ActionNotImplementedError, string(action))
}

func (z *autoscalerReconciler) determineAction(_ context.Context, meta reconciliation.Metadata, state *clientCore.StateHandlers) (reconciliation.Action, error) {
	userIndication := reconciliation.DetermineUserIndication(meta, meta.ClusterDeclaration.Integrations.Autoscaler)

	clusterExists, err := state.Cluster.HasCluster(meta.ClusterDeclaration.Metadata.Name)
	if err != nil {
		return reconciliation.ActionNoop, fmt.Errorf(constant.CheckIfClusterExistsError, err)
	}

	autoscalerExists := false
	if clusterExists {
		autoscalerExists, err = state.Autoscaler.HasAutoscaler()
		if err != nil {
			return reconciliation.ActionNoop, fmt.Errorf(constant.CheckIfAutoScalerExistsError, err)
		}
	}

	switch userIndication {
	case reconciliation.ActionCreate:
		if !clusterExists {
			return reconciliation.ActionWait, nil
		}

		if autoscalerExists {
			return reconciliation.ActionNoop, nil
		}

		return reconciliation.ActionCreate, nil
	case reconciliation.ActionDelete:
		if !clusterExists || !autoscalerExists {
			return reconciliation.ActionNoop, nil
		}

		return reconciliation.ActionDelete, nil
	}

	return reconciliation.ActionNoop, reconciliation.ErrIndecisive
}

// String returns the identifier for this reconciler
func (z *autoscalerReconciler) String() string {
	return autoscalerReconcilerIdentifier
}

// NewAutoscalerReconciler creates a new reconciler for the autoscaler resource
func NewAutoscalerReconciler(client client.AutoscalerService) reconciliation.Reconciler {
	return &autoscalerReconciler{
		client: client,
	}
}
