package reconciliation

import (
	"context"
	"fmt"

	clientCore "github.com/oslokommune/okctl/pkg/client/core"
	"github.com/oslokommune/okctl/pkg/controller/common/reconciliation"

	"github.com/oslokommune/okctl/pkg/client"
)

const promtailReconcilerIdentifier = "Promtail"

type promtailReconciler struct {
	client client.MonitoringService
}

// Reconcile knows how to do what is necessary to ensure the desired state is achieved
func (z *promtailReconciler) Reconcile(ctx context.Context, meta reconciliation.Metadata, state *clientCore.StateHandlers) (reconciliation.Result, error) {
	action, err := z.determineAction(meta, state)
	if err != nil {
		return reconciliation.Result{}, fmt.Errorf("determining course of action: %w", err)
	}

	switch action {
	case reconciliation.ActionCreate:
		_, err = z.client.CreatePromtail(
			ctx,
			reconciliation.ClusterMetaAsID(meta.ClusterDeclaration.Metadata),
		)
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf("creating promtail: %w", err)
		}

		return reconciliation.Result{Requeue: false}, nil
	case reconciliation.ActionDelete:
		err = z.client.DeletePromtail(
			ctx,
			reconciliation.ClusterMetaAsID(meta.ClusterDeclaration.Metadata),
		)
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf("deleting promtail: %w", err)
		}

		return reconciliation.Result{Requeue: false}, nil
	case reconciliation.ActionWait:
		return reconciliation.Result{Requeue: true}, nil
	case reconciliation.ActionNoop:
		return reconciliation.Result{Requeue: false}, nil
	}

	return reconciliation.Result{}, fmt.Errorf("action %s is not implemented", string(action))
}

func (z *promtailReconciler) determineAction(meta reconciliation.Metadata, state *clientCore.StateHandlers) (reconciliation.Action, error) {
	userIndiciation := reconciliation.DetermineUserIndication(meta, meta.ClusterDeclaration.Integrations.Promtail)

	clusterExists, err := state.Cluster.HasCluster(meta.ClusterDeclaration.Metadata.Name)
	if err != nil {
		return reconciliation.ActionNoop, fmt.Errorf("checking cluster existence: %w", err)
	}

	switch userIndiciation {
	case reconciliation.ActionCreate:
		if !clusterExists {
			return reconciliation.ActionWait, nil
		}

		componentExists, err := state.Promtail.HasPromtail()
		if err != nil {
			return reconciliation.ActionNoop, fmt.Errorf("checking component existence: %w", err)
		}

		if componentExists {
			return reconciliation.ActionNoop, nil
		}

		return reconciliation.ActionCreate, nil
	case reconciliation.ActionDelete:
		if !clusterExists {
			return reconciliation.ActionNoop, nil
		}

		componentExists, err := state.Promtail.HasPromtail()
		if err != nil {
			return reconciliation.ActionNoop, fmt.Errorf("checking component existence: %w", err)
		}

		if !componentExists {
			return reconciliation.ActionNoop, nil
		}

		return reconciliation.ActionDelete, nil
	}

	return reconciliation.ActionNoop, reconciliation.ErrIndecisive
}

// String returns the identifier for this reconciler
func (z *promtailReconciler) String() string {
	return promtailReconcilerIdentifier
}

// NewPromtailReconciler creates a new reconciler for the Promtail resource
func NewPromtailReconciler(client client.MonitoringService) reconciliation.Reconciler {
	return &promtailReconciler{
		client: client,
	}
}
