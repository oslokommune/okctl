package reconciliation

import (
	"context"
	"fmt"

	"github.com/oslokommune/okctl/pkg/lib/paths"

	clientCore "github.com/oslokommune/okctl/pkg/client/core"
	"github.com/oslokommune/okctl/pkg/controller/common/reconciliation"
)

// Reconcile knows how to do what is necessary to ensure the desired state is achieved
func (c applicationReconciler) Reconcile(_ context.Context, meta reconciliation.Metadata, state *clientCore.StateHandlers) (reconciliation.Result, error) {
	action, err := c.determineAction(meta, state)
	if err != nil {
		return reconciliation.Result{}, fmt.Errorf("determining course of action: %w", err)
	}

	switch action {
	case reconciliation.ActionCreate:
		absoluteIACRepositoryRootDirectoryPath, err := paths.GetAbsoluteIACRepositoryRootDirectory()
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf("acquiring absolute IAC repository directory path: %w", err)
		}

		err = state.Application.Initialize(*meta.ClusterDeclaration, absoluteIACRepositoryRootDirectoryPath)
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf("initializing applications state: %w", err)
		}

		return reconciliation.Result{Requeue: false}, nil
	case reconciliation.ActionDelete:
		return reconciliation.Result{Requeue: false}, nil
	}

	return reconciliation.NoopWaitIndecisiveHandler(action)
}

// determineAction knows how to determine if a resource should be created, deleted or updated
func (c applicationReconciler) determineAction(meta reconciliation.Metadata, state *clientCore.StateHandlers) (reconciliation.Action, error) {
	userIndication := reconciliation.DetermineUserIndication(meta, true)

	hasCluster, err := state.Cluster.HasCluster(meta.ClusterDeclaration.Metadata.Name)
	if err != nil {
		return "", fmt.Errorf("checking cluster existence: %w", err)
	}

	switch userIndication {
	case reconciliation.ActionCreate:
		if !hasCluster {
			return reconciliation.ActionWait, nil
		}

		return reconciliation.ActionCreate, nil
	case reconciliation.ActionDelete:
		if !hasCluster {
			return reconciliation.ActionNoop, nil
		}

		return reconciliation.ActionDelete, nil
	}

	return "", reconciliation.ErrIndecisive
}

// String returns a descriptive identifier of the domain that this reconciler represents
func (c applicationReconciler) String() string {
	return "Applications"
}

// NewApplicationReconciler returns an initialized application reconciler
func NewApplicationReconciler() reconciliation.Reconciler {
	return &applicationReconciler{}
}

type applicationReconciler struct{}
