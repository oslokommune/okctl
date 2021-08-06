package reconciliation

import (
	"context"
	"fmt"

	clientCore "github.com/oslokommune/okctl/pkg/client/core"
	"github.com/oslokommune/okctl/pkg/controller/common/reconciliation"

	"github.com/oslokommune/okctl/pkg/client"
)

const externalSecretsReconcilerIdentifier = "secrets controller"

// externalSecretsReconciler contains service and metadata for the relevant resource
type externalSecretsReconciler struct {
	client client.ExternalSecretsService
}

// Reconcile knows how to do what is necessary to ensure the desired state is achieved
func (z *externalSecretsReconciler) Reconcile(ctx context.Context, meta reconciliation.Metadata, state *clientCore.StateHandlers) (reconciliation.Result, error) {
	action, err := z.determineAction(ctx, meta, state)
	if err != nil {
		return reconciliation.Result{}, fmt.Errorf("determining course of action: %w", err)
	}

	switch action {
	case reconciliation.ActionCreate:
		_, err = z.client.CreateExternalSecrets(ctx, client.CreateExternalSecretsOpts{
			ID: reconciliation.ClusterMetaAsID(meta.ClusterDeclaration.Metadata),
		})
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf("creating external secrets: %w", err)
		}

		return reconciliation.Result{Requeue: false}, nil
	case reconciliation.ActionDelete:
		err = z.client.DeleteExternalSecrets(
			ctx,
			reconciliation.ClusterMetaAsID(meta.ClusterDeclaration.Metadata),
		)
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf("deleting external secrets: %w", err)
		}

		return reconciliation.Result{Requeue: false}, nil
	case reconciliation.ActionWait:
		return reconciliation.Result{Requeue: true}, nil
	case reconciliation.ActionNoop:
		return reconciliation.Result{Requeue: false}, nil
	}

	return reconciliation.Result{}, fmt.Errorf("action %s is not implemented", string(action))
}

func (z *externalSecretsReconciler) determineAction(_ context.Context, meta reconciliation.Metadata, state *clientCore.StateHandlers) (reconciliation.Action, error) {
	userIndication := reconciliation.DetermineUserIndication(meta, meta.ClusterDeclaration.Integrations.ExternalSecrets)

	clusterExists, err := state.Cluster.HasCluster(meta.ClusterDeclaration.Metadata.Name)
	if err != nil {
		return reconciliation.ActionNoop, fmt.Errorf("acquiring cluster existence: %w", err)
	}

	componentExists := false
	if clusterExists {
		componentExists, err = state.ExternalSecrets.HasExternalSecrets()
		if err != nil {
			return reconciliation.ActionNoop, fmt.Errorf("acquiring secrets controller existence: %w", err)
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

// String returns the identifier for this reconciler
func (z *externalSecretsReconciler) String() string {
	return externalSecretsReconcilerIdentifier
}

// NewExternalSecretsReconciler creates a new reconciler for the ExternalSecrets resource
func NewExternalSecretsReconciler(client client.ExternalSecretsService) reconciliation.Reconciler {
	return &externalSecretsReconciler{
		client: client,
	}
}
