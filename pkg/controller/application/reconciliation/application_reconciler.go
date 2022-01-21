package reconciliation

import (
	"context"
	"errors"
	"fmt"

	stormpkg "github.com/asdine/storm/v3"
	"github.com/oslokommune/okctl/pkg/client"
	clientCore "github.com/oslokommune/okctl/pkg/client/core"
	"github.com/oslokommune/okctl/pkg/controller/common/reconciliation"
)

const applicationReconcilerIdentifier = "kubernetes manifests"

// applicationReconciler contains service and metadata for the relevant resource
type applicationReconciler struct {
	client client.ApplicationService
}

// Reconcile knows how to do what is necessary to ensure the desired state is achieved
func (a *applicationReconciler) Reconcile(ctx context.Context, meta reconciliation.Metadata, state *clientCore.StateHandlers) (reconciliation.Result, error) {
	opts := generalOpts{
		Ctx:   ctx,
		Meta:  meta,
		State: state,
	}

	action, err := a.determineAction(opts)
	if err != nil {
		return reconciliation.Result{}, fmt.Errorf("determining course of action: %w", err)
	}

	switch action {
	case reconciliation.ActionCreate:
		return a.createApplication(opts)
	case reconciliation.ActionDelete:
		err = a.client.DeleteApplicationManifests(ctx, client.DeleteApplicationManifestsOpts{
			Cluster:     *meta.ClusterDeclaration,
			Application: meta.ApplicationDeclaration,
		})
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf("deleting application manifests: %w", err)
		}

		return reconciliation.Result{Requeue: false}, nil
	case reconciliation.ActionWait:
		return reconciliation.Result{Requeue: true}, nil
	case reconciliation.ActionNoop:
		return reconciliation.Result{Requeue: false}, nil
	}

	return reconciliation.Result{}, fmt.Errorf("action %s is not implemented", string(action))
}

func (a *applicationReconciler) createApplication(opts generalOpts) (reconciliation.Result, error) {
	hz, err := opts.State.Domain.GetPrimaryHostedZone()
	if err != nil {
		return reconciliation.Result{}, fmt.Errorf("getting primary hosted zone: %w", err)
	}

	if opts.Meta.ApplicationDeclaration.Image.HasName() {
		repo, err := opts.State.ContainerRepository.GetContainerRepository(opts.Meta.ApplicationDeclaration.Image.Name)
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf("getting container repository: %w", err)
		}

		opts.Meta.ApplicationDeclaration.Image.Name = ""
		opts.Meta.ApplicationDeclaration.Image.URI = repo.URI()
	}

	err = a.client.ScaffoldApplication(opts.Ctx, &client.ScaffoldApplicationOpts{
		Cluster:      *opts.Meta.ClusterDeclaration,
		Application:  opts.Meta.ApplicationDeclaration,
		HostedZoneID: hz.HostedZoneID,
	})
	if err != nil {
		return reconciliation.Result{}, err
	}

	return reconciliation.Result{Requeue: false}, nil
}

func (a *applicationReconciler) determineAction(opts generalOpts) (reconciliation.Action, error) {
	userIndication := reconciliation.DetermineUserIndication(opts.Meta, true)

	switch userIndication {
	case reconciliation.ActionCreate:
		dependenciesReady, err := a.hasCreateDependenciesMet(opts.Meta, opts.State)
		if err != nil {
			return reconciliation.ActionNoop, fmt.Errorf("acquiring dependency state: %w", err)
		}

		if !dependenciesReady {
			return reconciliation.ActionWait, nil
		}

		return reconciliation.ActionCreate, nil
	case reconciliation.ActionDelete:
		return reconciliation.ActionDelete, nil
	}

	return reconciliation.ActionNoop, reconciliation.ErrIndecisive
}

// hasCreateDependenciesMet checks dependencies
func (a *applicationReconciler) hasCreateDependenciesMet(meta reconciliation.Metadata, state *clientCore.StateHandlers) (bool, error) {
	if exists, err := state.Domain.HasPrimaryHostedZone(); err == nil {
		if !exists {
			return false, nil
		}
	} else {
		return false, fmt.Errorf("determining existence of primary hosted zone for %s: %w", a.String(), err)
	}

	if _, err := state.Github.GetGithubRepository(meta.ClusterDeclaration.Github.Path()); err != nil {
		if errors.Is(err, stormpkg.ErrNotFound) {
			return false, nil
		}

		return false, fmt.Errorf("determining existence of a Github repository for %s: %w", a.String(), err)
	}

	if meta.ApplicationDeclaration.Image.HasName() {
		exists, err := state.ContainerRepository.ApplicationHasImage(meta.ApplicationDeclaration.Metadata.Name)
		if err != nil {
			return false, fmt.Errorf("determining existence of a ECR repository: %w", err)
		}

		if !exists {
			return false, nil
		}
	}

	return true, nil
}

// NodeType returns the relevant NodeType for this reconciler
func (a *applicationReconciler) String() string {
	return applicationReconcilerIdentifier
}

// NewApplicationReconciler creates a new reconciler for the VPC resource
func NewApplicationReconciler(client client.ApplicationService) reconciliation.Reconciler {
	return &applicationReconciler{
		client: client,
	}
}

type generalOpts struct {
	Ctx   context.Context
	Meta  reconciliation.Metadata
	State *clientCore.StateHandlers
}
