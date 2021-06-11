package reconciliation

import (
	"context"
	"fmt"

	"github.com/oslokommune/okctl/pkg/controller/common/reconciliation"

	"github.com/oslokommune/okctl/pkg/client"
	clientCore "github.com/oslokommune/okctl/pkg/client/core"
)

const containerRepositoryReconcilerName = "container repositories"

// containerRepositoryReconciler contains service and metadata for the relevant resource
type containerRepositoryReconciler struct {
	client client.ContainerRepositoryService
}

// Reconcile knows how to do what is necessary to ensure the desired state is achieved
func (c *containerRepositoryReconciler) Reconcile(ctx context.Context, meta reconciliation.Metadata, state *clientCore.StateHandlers) (reconciliation.Result, error) {
	action, err := c.determineAction(meta, state)
	if err != nil {
		return reconciliation.Result{}, fmt.Errorf("determining course of action: %w", err)
	}

	switch action {
	case reconciliation.ActionCreate:
		_, err = c.client.CreateContainerRepository(ctx, client.CreateContainerRepositoryOpts{
			ClusterID:       reconciliation.ClusterMetaAsID(meta.ClusterDeclaration.Metadata),
			ApplicationName: meta.ApplicationDeclaration.Metadata.Name,
			ImageName:       meta.ApplicationDeclaration.Image.Name,
		})
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf("creating container repository: %w", err)
		}

		return reconciliation.Result{Requeue: false}, nil
	case reconciliation.ActionDelete:
		err := c.client.DeleteContainerRepository(ctx, client.DeleteContainerRepositoryOpts{
			ClusterID: reconciliation.ClusterMetaAsID(meta.ClusterDeclaration.Metadata),
			ImageName: meta.ApplicationDeclaration.Image.Name,
		})
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf("deleting container repository: %w", err)
		}

		return reconciliation.Result{Requeue: false}, nil
	case reconciliation.ActionWait:
		return reconciliation.Result{Requeue: true}, nil
	case reconciliation.ActionNoop:
		return reconciliation.Result{Requeue: false}, nil
	}

	return reconciliation.Result{}, fmt.Errorf("action %s is not implemented", string(action))
}

func (c *containerRepositoryReconciler) determineAction(meta reconciliation.Metadata, state *clientCore.StateHandlers) (reconciliation.Action, error) {
	userIndication := reconciliation.DetermineUserIndication(meta, meta.ApplicationDeclaration.Image.HasName())

	hasExistingImage, err := state.ContainerRepository.ApplicationHasImage(meta.ApplicationDeclaration.Metadata.Name)
	if err != nil {
		return reconciliation.ActionNoop, fmt.Errorf("acquiring existence from state %w", err)
	}

	switch userIndication {
	case reconciliation.ActionCreate:
		if hasExistingImage {
			return reconciliation.ActionNoop, nil
		}

		return reconciliation.ActionCreate, nil
	case reconciliation.ActionDelete:
		if !hasExistingImage {
			return reconciliation.ActionNoop, nil
		}

		return reconciliation.ActionDelete, nil
	case reconciliation.ActionNoop:
		return reconciliation.ActionNoop, nil
	}

	return reconciliation.ActionNoop, reconciliation.ErrIndecisive
}

// String returns an identifier for this reconciler
func (c *containerRepositoryReconciler) String() string {
	return containerRepositoryReconcilerName
}

// NewContainerRepositoryReconciler creates a new reconciler for the VPC resource
func NewContainerRepositoryReconciler(client client.ContainerRepositoryService) reconciliation.Reconciler {
	return &containerRepositoryReconciler{
		client: client,
	}
}
