package reconciliation

import (
	"context"
	"fmt"

	"github.com/oslokommune/okctl/pkg/config/constant"

	clientCore "github.com/oslokommune/okctl/pkg/client/core"
	"github.com/oslokommune/okctl/pkg/controller/common/reconciliation"

	"github.com/oslokommune/okctl/pkg/client"
)

// blockstorageReconciler contains service and metadata for the relevant resource
type blockstorageReconciler struct {
	client client.BlockstorageService
}

const blockstorageReconcilerIdentifier = "persistent storage"

// NodeType returns the relevant NodeType for this reconciler
func (z *blockstorageReconciler) String() string {
	return blockstorageReconcilerIdentifier
}

// Reconcile knows how to do what is necessary to ensure the desired state is achieved
func (z *blockstorageReconciler) Reconcile(ctx context.Context, meta reconciliation.Metadata, state *clientCore.StateHandlers) (result reconciliation.Result, err error) {
	action, err := z.determineAction(ctx, meta, state)
	if err != nil {
		return reconciliation.Result{}, fmt.Errorf(constant.ReconcilerDetermineActionError, err)
	}

	switch action {
	case reconciliation.ActionCreate:
		_, err = z.client.CreateBlockstorage(ctx, client.CreateBlockstorageOpts{
			ID: reconciliation.ClusterMetaAsID(meta.ClusterDeclaration.Metadata),
		})
		if err != nil {
			return result, fmt.Errorf(constant.CreateBlockStorageError, err)
		}

		return reconciliation.Result{Requeue: false}, nil
	case reconciliation.ActionDelete:
		err = z.client.DeleteBlockstorage(
			ctx,
			reconciliation.ClusterMetaAsID(meta.ClusterDeclaration.Metadata),
		)
		if err != nil {
			return result, fmt.Errorf(constant.DeleteBlockStorageError, err)
		}

		return reconciliation.Result{Requeue: false}, nil
	case reconciliation.ActionWait:
		return reconciliation.Result{Requeue: true}, nil
	case reconciliation.ActionNoop:
		return reconciliation.Result{Requeue: false}, nil
	}

	return reconciliation.Result{}, fmt.Errorf(constant.ActionNotImplementedError, string(action))
}

func (z *blockstorageReconciler) determineAction(_ context.Context, meta reconciliation.Metadata, state *clientCore.StateHandlers) (reconciliation.Action, error) {
	userIndication := reconciliation.DetermineUserIndication(meta, meta.ClusterDeclaration.Integrations.Blockstorage)

	clusterExists, err := state.Cluster.HasCluster(meta.ClusterDeclaration.Metadata.Name)
	if err != nil {
		return reconciliation.ActionNoop, fmt.Errorf(constant.CheckClusterExistanceError, err)
	}

	componentExists := false
	if clusterExists {
		componentExists, err = state.Blockstorage.HasBlockstorage()
		if err != nil {
			return reconciliation.ActionNoop, fmt.Errorf(constant.CheckIfAWSLoadBalancerExistsError, err)
		}
	}

	switch userIndication {
	case reconciliation.ActionCreate:
		if !clusterExists {
			return reconciliation.ActionWait, nil
		}

		if componentExists {
			return reconciliation.ActionNoop, nil
		}

		return reconciliation.ActionCreate, nil
	case reconciliation.ActionDelete:
		if !clusterExists || !componentExists {
			return reconciliation.ActionNoop, nil
		}

		return reconciliation.ActionDelete, nil
	}

	return reconciliation.ActionNoop, reconciliation.ErrIndecisive
}

// NewBlockstorageReconciler creates a new reconciler for the Blockstorage resource
func NewBlockstorageReconciler(client client.BlockstorageService) reconciliation.Reconciler {
	return &blockstorageReconciler{
		client: client,
	}
}
