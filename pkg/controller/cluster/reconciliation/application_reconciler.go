package reconciliation

import (
	"context"
	"fmt"

	"github.com/oslokommune/okctl/cmd/okctl/handlers"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/spinner"

	"github.com/oslokommune/okctl/pkg/lib/paths"

	clientCore "github.com/oslokommune/okctl/pkg/client/core"
	"github.com/oslokommune/okctl/pkg/controller/common/reconciliation"
)

// Reconcile knows how to do what is necessary to ensure the desired state is achieved
func (c applicationReconciler) Reconcile(ctx context.Context, meta reconciliation.Metadata, state *clientCore.StateHandlers) (reconciliation.Result, error) {
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
		applications, err := state.Application.List()
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf("listing applications: %w", err)
		}

		for _, app := range applications {
			err = deleteApplication(deleteApplicationOpts{
				Ctx:                 ctx,
				Meta:                meta,
				Services:            c.services,
				State:               state,
				ClusterManifest:     *meta.ClusterDeclaration,
				ApplicationManifest: app,
			})
			if err != nil {
				return reconciliation.Result{}, fmt.Errorf("deleting application: %w", err)
			}
		}

		return reconciliation.Result{Requeue: false}, nil
	}

	return reconciliation.NoopWaitIndecisiveHandler(action)
}

type deleteApplicationOpts struct {
	Ctx                 context.Context
	Meta                reconciliation.Metadata
	Services            *clientCore.Services
	State               *clientCore.StateHandlers
	ClusterManifest     v1alpha1.Cluster
	ApplicationManifest v1alpha1.Application
}

func deleteApplication(opts deleteApplicationOpts) error {
	spin, err := spinner.New("deleting", opts.Meta.Out)
	if err != nil {
		return fmt.Errorf("creating spinner: %w", err)
	}

	scheduler := handlers.CreateScheduler(handlers.CreateSchedulerOpts{
		Out:                 opts.Meta.Out,
		Services:            opts.Services,
		State:               opts.State,
		Spinner:             spin,
		ClusterManifest:     opts.ClusterManifest,
		ApplicationManifest: opts.ApplicationManifest,
		Purge:               true,
		DelayFunction:       reconciliation.DefaultDelayFunction,
	})

	_, err = scheduler.Run(opts.Ctx, opts.State)
	if err != nil {
		return fmt.Errorf("deleting application: %w", err)
	}

	return nil
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
func NewApplicationReconciler(services *clientCore.Services) reconciliation.Reconciler {
	return &applicationReconciler{services: services}
}

type applicationReconciler struct {
	services *clientCore.Services
}
