package reconciliation

import (
	"context"
	"fmt"

	"github.com/oslokommune/okctl/pkg/logging"

	"github.com/oslokommune/okctl/pkg/client"

	"github.com/oslokommune/okctl/pkg/controller/common/reconciliation"

	clientCore "github.com/oslokommune/okctl/pkg/client/core"
)

const (
	postgresReconcilerName = "postgres integration"
	postgresLogTag         = "appReconciliation/postgres"
)

// containerRepositoryReconciler contains service and metadata for the relevant resource
type postgresReconciler struct {
	applicationPostgresService client.ApplicationPostgresService
}

// Reconcile knows how to do what is necessary to ensure the desired state is achieved
func (c *postgresReconciler) Reconcile(ctx context.Context, meta reconciliation.Metadata, state *clientCore.StateHandlers) (reconciliation.Result, error) {
	log := logging.GetLogger(postgresLogTag, "reconcile")

	action, err := c.determineAction(ctx, meta, state)
	if err != nil {
		return reconciliation.Result{}, fmt.Errorf("determining course of action: %w", err)
	}

	log.Debug(fmt.Sprintf("determined action: %s", string(action)))

	switch action {
	case reconciliation.ActionCreate:
		err = c.applicationPostgresService.AddPostgresToApplication(ctx, client.AddPostgresToApplicationOpts{
			Cluster:      *meta.ClusterDeclaration,
			Application:  meta.ApplicationDeclaration,
			DatabaseName: meta.ApplicationDeclaration.Postgres,
		})
		if err != nil {
			return reconciliation.Result{Requeue: false}, fmt.Errorf("adding postgres to application: %w", err)
		}

		return reconciliation.Result{Requeue: false}, nil
	case reconciliation.ActionDelete:
		err = c.applicationPostgresService.RemovePostgresFromApplication(ctx, client.RemovePostgresFromApplicationOpts{
			Cluster:      *meta.ClusterDeclaration,
			Application:  meta.ApplicationDeclaration,
			DatabaseName: meta.ApplicationDeclaration.Postgres,
		})
		if err != nil {
			return reconciliation.Result{Requeue: false}, fmt.Errorf("removing Postgres from application: %w", err)
		}

		return reconciliation.Result{Requeue: false}, nil
	case reconciliation.ActionWait:
		return reconciliation.Result{Requeue: true}, nil
	case reconciliation.ActionNoop:
		return reconciliation.Result{Requeue: false}, nil
	}

	return reconciliation.Result{}, fmt.Errorf("action %s is not implemented", string(action))
}

func (c *postgresReconciler) determineAction(
	ctx context.Context,
	meta reconciliation.Metadata,
	_ *clientCore.StateHandlers,
) (reconciliation.Action, error) {
	userIndication := reconciliation.DetermineUserIndication(meta, meta.ApplicationDeclaration.Postgres != "")

	exists, err := c.applicationPostgresService.HasPostgresIntegration(ctx, client.HasPostgresIntegrationOpts{
		Cluster:      *meta.ClusterDeclaration,
		Application:  meta.ApplicationDeclaration,
		DatabaseName: meta.ApplicationDeclaration.Postgres,
	})
	if err != nil {
		return reconciliation.ActionNoop, fmt.Errorf("checking integration existence: %w", err)
	}

	switch userIndication {
	case reconciliation.ActionCreate:
		if exists {
			return reconciliation.ActionNoop, nil
		}

		return reconciliation.ActionCreate, nil
	case reconciliation.ActionDelete:
		if !exists {
			return reconciliation.ActionNoop, nil
		}

		return reconciliation.ActionDelete, nil
	}

	return reconciliation.ActionNoop, reconciliation.ErrIndecisive
}

// String returns an identifier for this reconciler
func (c *postgresReconciler) String() string {
	return postgresReconcilerName
}

// NewPostgresReconciler creates a new reconciler for the VPC resource
func NewPostgresReconciler(applicationPostgresService client.ApplicationPostgresService) reconciliation.Reconciler {
	return &postgresReconciler{
		applicationPostgresService: applicationPostgresService,
	}
}
