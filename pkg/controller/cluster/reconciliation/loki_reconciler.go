package reconciliation

import (
	"context"
	"fmt"
	"github.com/oslokommune/okctl/pkg/config/constant"

	clientCore "github.com/oslokommune/okctl/pkg/client/core"
	"github.com/oslokommune/okctl/pkg/controller/common/reconciliation"

	"github.com/oslokommune/okctl/pkg/client"
)

type lokiReconciler struct {
	client client.MonitoringService
}

const lokiReconcilerIdentifier = "Loki"

// Reconcile knows how to do what is necessary to ensure the desired state is achieved
func (z *lokiReconciler) Reconcile(ctx context.Context, meta reconciliation.Metadata, state *clientCore.StateHandlers) (reconciliation.Result, error) {
	action, err := z.determineAction(meta, state)
	if err != nil {
		return reconciliation.Result{Requeue: false}, fmt.Errorf(constant.ReconcilerDetermineActionError, err)
	}

	switch action {
	case reconciliation.ActionCreate:
		_, err = z.client.CreateLoki(ctx, reconciliation.ClusterMetaAsID(meta.ClusterDeclaration.Metadata))
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf(constant.CreateLokiError, err)
		}

		return reconciliation.Result{Requeue: false}, nil
	case reconciliation.ActionDelete:
		err = z.client.DeleteLoki(ctx, reconciliation.ClusterMetaAsID(meta.ClusterDeclaration.Metadata))
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf(constant.DeleteLokiError, err)
		}

		return reconciliation.Result{Requeue: false}, nil
	case reconciliation.ActionWait:
		return reconciliation.Result{Requeue: true}, nil
	case reconciliation.ActionNoop:
		return reconciliation.Result{Requeue: false}, nil
	}

	return reconciliation.Result{}, fmt.Errorf(constant.ActionNotImplementedError, string(action))
}

func (z *lokiReconciler) determineAction(meta reconciliation.Metadata, state *clientCore.StateHandlers) (reconciliation.Action, error) {
	userIndication := reconciliation.DetermineUserIndication(meta, meta.ClusterDeclaration.Integrations.Loki)

	clusterExistenceTest := reconciliation.GenerateClusterExistenceTest(state, meta.ClusterDeclaration.Metadata.Name)

	switch userIndication {
	case reconciliation.ActionCreate:
		dependenciesReady, err := reconciliation.AssertDependencyExistence(true, clusterExistenceTest)
		if err != nil {
			return reconciliation.ActionNoop, fmt.Errorf(constant.CheckDepedencyError, err)
		}

		if !dependenciesReady {
			return reconciliation.ActionWait, nil
		}

		componentExists, err := state.Loki.HasLoki()
		if err != nil {
			return reconciliation.ActionNoop, fmt.Errorf(constant.CheckLokiExistenceError, err)
		}

		if componentExists {
			return reconciliation.ActionNoop, nil
		}

		return reconciliation.ActionCreate, nil
	case reconciliation.ActionDelete:
		clusterExists, err := clusterExistenceTest()
		if err != nil {
			return reconciliation.ActionNoop, fmt.Errorf(constant.CheckClusterExistanceError, err)
		}

		if !clusterExists {
			return reconciliation.ActionNoop, nil
		}

		componentExists, err := state.Loki.HasLoki()
		if err != nil {
			return reconciliation.ActionNoop, fmt.Errorf(constant.CheckLokiExistenceError, err)
		}

		if !componentExists {
			return reconciliation.ActionNoop, nil
		}

		return reconciliation.ActionDelete, nil
	}

	return reconciliation.ActionNoop, reconciliation.ErrIndecisive
}

// String returns the identifier for this reconciler
func (z *lokiReconciler) String() string {
	return lokiReconcilerIdentifier
}

// NewLokiReconciler creates a new reconciler for the loki resource
func NewLokiReconciler(client client.MonitoringService) reconciliation.Reconciler {
	return &lokiReconciler{
		client: client,
	}
}
