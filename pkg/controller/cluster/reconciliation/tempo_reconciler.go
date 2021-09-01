package reconciliation

import (
	"context"
	"fmt"
	"github.com/oslokommune/okctl/pkg/config/constant"

	clientCore "github.com/oslokommune/okctl/pkg/client/core"
	"github.com/oslokommune/okctl/pkg/controller/common/reconciliation"

	"github.com/oslokommune/okctl/pkg/client"
)

const tempoReconcilerIdentifier = "Tempo"

type tempoReconciler struct {
	client client.MonitoringService
}

// Reconcile knows how to do what is necessary to ensure the desired state is achieved
func (z *tempoReconciler) Reconcile(ctx context.Context, meta reconciliation.Metadata, state *clientCore.StateHandlers) (reconciliation.Result, error) {
	action, err := z.determineAction(meta, state)
	if err != nil {
		return reconciliation.Result{Requeue: false}, fmt.Errorf(constant.ReconcilerDetermineActionError, err)
	}

	switch action {
	case reconciliation.ActionCreate:
		_, err = z.client.CreateTempo(ctx, reconciliation.ClusterMetaAsID(meta.ClusterDeclaration.Metadata))
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf(constant.CreateTempoError, err)
		}

		return reconciliation.Result{Requeue: false}, nil
	case reconciliation.ActionDelete:
		err = z.client.DeleteTempo(ctx, reconciliation.ClusterMetaAsID(meta.ClusterDeclaration.Metadata))
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf(constant.DeleteTempoError, err)
		}

		return reconciliation.Result{Requeue: false}, nil
	case reconciliation.ActionWait:
		return reconciliation.Result{Requeue: true}, nil
	case reconciliation.ActionNoop:
		return reconciliation.Result{Requeue: false}, nil
	}

	return reconciliation.Result{}, fmt.Errorf(constant.ActionNotImplementedError, string(action))
}

func (z *tempoReconciler) determineAction(meta reconciliation.Metadata, state *clientCore.StateHandlers) (reconciliation.Action, error) {
	userIndication := reconciliation.DetermineUserIndication(meta, meta.ClusterDeclaration.Integrations.Tempo)

	clusterExists, err := state.Cluster.HasCluster(meta.ClusterDeclaration.Metadata.Name)
	if err != nil {
		return reconciliation.ActionNoop, fmt.Errorf(constant.CheckClusterExistanceError, err)
	}

	switch userIndication {
	case reconciliation.ActionCreate:
		if !clusterExists {
			return reconciliation.ActionWait, nil
		}

		componentExists, err := state.Tempo.HasTempo()
		if err != nil {
			return reconciliation.ActionNoop, fmt.Errorf(constant.CheckComponentExistenceError, err)
		}

		if componentExists {
			return reconciliation.ActionNoop, nil
		}

		return reconciliation.ActionCreate, nil
	case reconciliation.ActionDelete:
		if !clusterExists {
			return reconciliation.ActionNoop, nil
		}

		componentExists, err := state.Tempo.HasTempo()
		if err != nil {
			return reconciliation.ActionNoop, fmt.Errorf(constant.CheckComponentExistenceError, err)
		}

		if !componentExists {
			return reconciliation.ActionNoop, nil
		}

		return reconciliation.ActionDelete, nil
	}

	return reconciliation.ActionNoop, reconciliation.ErrIndecisive
}

// String returns the identifier for this reconciler
func (z *tempoReconciler) String() string {
	return tempoReconcilerIdentifier
}

// NewTempoReconciler creates a new reconciler for the tempo resource
func NewTempoReconciler(client client.MonitoringService) reconciliation.Reconciler {
	return &tempoReconciler{
		client: client,
	}
}
