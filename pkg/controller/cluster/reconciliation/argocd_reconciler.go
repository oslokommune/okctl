package reconciliation

import (
	"context"
	"fmt"
	"strings"

	"github.com/oslokommune/okctl/pkg/controller/common/reconciliation"

	"github.com/oslokommune/okctl/pkg/cfn"

	clientCore "github.com/oslokommune/okctl/pkg/client/core"

	"github.com/oslokommune/okctl/pkg/config/constant"

	"github.com/miekg/dns"
	"github.com/oslokommune/okctl/pkg/client"
)

const argocdReconcilerIdentifier = "ArgoCD"

// argocdReconciler contains service and metadata for the relevant resource
type argocdReconciler struct {
	argocdClient client.ArgoCDService
	githubClient client.GithubService
}

// Reconcile knows how to do what is necessary to ensure the desired state is achieved
func (z *argocdReconciler) Reconcile(ctx context.Context, meta reconciliation.Metadata, state *clientCore.StateHandlers) (result reconciliation.Result, err error) {
	action, err := z.determineAction(meta, state)
	if err != nil {
		return reconciliation.Result{}, fmt.Errorf("determining course of action: %w", err)
	}

	switch action {
	case reconciliation.ActionCreate:
		return z.createArgoCD(ctx, meta, state)
	case reconciliation.ActionDelete:
		err := z.argocdClient.DeleteArgoCD(ctx, client.DeleteArgoCDOpts{
			ID: reconciliation.ClusterMetaAsID(meta.ClusterDeclaration.Metadata),
		})
		if err != nil {
			return result, fmt.Errorf("deleting argocd: %w", err)
		}

		err = z.githubClient.DeleteGithubRepository(ctx, client.DeleteGithubRepositoryOpts{
			ID:           reconciliation.ClusterMetaAsID(meta.ClusterDeclaration.Metadata),
			Organisation: meta.ClusterDeclaration.Github.Organisation,
			Name:         meta.ClusterDeclaration.Github.Repository,
		})
		if err != nil {
			return result, fmt.Errorf("deleting github repository: %w", err)
		}

		return reconciliation.Result{Requeue: false}, nil
	case reconciliation.ActionWait:
		return reconciliation.Result{Requeue: true}, nil
	case reconciliation.ActionNoop:
		return reconciliation.Result{Requeue: false}, nil
	}

	return reconciliation.Result{}, fmt.Errorf("action %s is not implemented", string(action))
}

func (z *argocdReconciler) createArgoCD(ctx context.Context, meta reconciliation.Metadata, state *clientCore.StateHandlers) (reconciliation.Result, error) {
	repo, err := z.githubClient.CreateGithubRepository(ctx, client.CreateGithubRepositoryOpts{
		ID:           reconciliation.ClusterMetaAsID(meta.ClusterDeclaration.Metadata),
		Host:         constant.DefaultGithubHost,
		Organization: meta.ClusterDeclaration.Github.Organisation,
		Name:         meta.ClusterDeclaration.Github.Repository,
	})
	if err != nil {
		return reconciliation.Result{}, fmt.Errorf("fetching deploy key: %w", err)
	}

	hz, err := state.Domain.GetPrimaryHostedZone()
	if err != nil {
		return reconciliation.Result{}, fmt.Errorf("getting primary hosted zone: %w", err)
	}

	im, err := state.IdentityManager.GetIdentityPool(
		cfn.NewStackNamer().IdentityPool(meta.ClusterDeclaration.Metadata.Name),
	)
	if err != nil {
		return reconciliation.Result{}, fmt.Errorf("getting identity pool: %w", err)
	}

	_, err = z.argocdClient.CreateArgoCD(ctx, client.CreateArgoCDOpts{
		ID:                 reconciliation.ClusterMetaAsID(meta.ClusterDeclaration.Metadata),
		Domain:             meta.ClusterDeclaration.ClusterRootDomain,
		FQDN:               dns.Fqdn(meta.ClusterDeclaration.ClusterRootDomain),
		HostedZoneID:       hz.HostedZoneID,
		GithubOrganisation: meta.ClusterDeclaration.Github.Organisation,
		UserPoolID:         im.UserPoolID,
		AuthDomain:         im.AuthDomain,
		Repository:         repo,
	})

	if err != nil {
		// nolint: godox
		// TODO: Need to identify the correct error
		if strings.Contains(strings.ToLower(err.Error()), "timeout") {
			fmt.Println(meta.Out, fmt.Errorf("got ArgoCD timeout: %w", err).Error())

			return reconciliation.Result{Requeue: true}, nil
		}

		return reconciliation.Result{}, fmt.Errorf("creating argocd: %w", err)
	}

	return reconciliation.Result{Requeue: false}, nil
}

func (z *argocdReconciler) determineAction(meta reconciliation.Metadata, state *clientCore.StateHandlers) (reconciliation.Action, error) {
	userIndication := reconciliation.DetermineUserIndication(meta, meta.ClusterDeclaration.Integrations.ArgoCD)

	componentExists, err := state.ArgoCD.HasArgoCD()
	if err != nil {
		return reconciliation.ActionNoop, fmt.Errorf("acquiring ArgoCD existence: %w", err)
	}

	switch userIndication {
	case reconciliation.ActionCreate:
		if componentExists {
			return reconciliation.ActionNoop, nil
		}

		dependenciesReady, err := reconciliation.AssertDependencyExistence(true,
			reconciliation.GenerateClusterExistenceTest(state, meta.ClusterDeclaration.Metadata.Name),
			state.Domain.HasPrimaryHostedZone,
			state.IdentityManager.HasIdentityPool,
		)
		if err != nil {
			return reconciliation.ActionNoop, fmt.Errorf("checking dependencies: %w", err)
		}

		if !dependenciesReady {
			return reconciliation.ActionWait, nil
		}

		return reconciliation.ActionCreate, nil
	case reconciliation.ActionDelete:
		if !componentExists {
			return reconciliation.ActionNoop, nil
		}

		return reconciliation.ActionDelete, nil
	}

	return reconciliation.ActionNoop, reconciliation.ErrIndecisive
}

// String returns the identifier for this reconciler
func (z *argocdReconciler) String() string {
	return argocdReconcilerIdentifier
}

// NewArgocdReconciler creates a new reconciler for the ArgoCD resource
func NewArgocdReconciler(
	argocdClient client.ArgoCDService,
	githubClient client.GithubService,
) reconciliation.Reconciler {
	return &argocdReconciler{
		argocdClient: argocdClient,
		githubClient: githubClient,
	}
}
