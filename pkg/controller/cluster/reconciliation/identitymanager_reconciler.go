package reconciliation

import (
	"context"
	"fmt"

	clientCore "github.com/oslokommune/okctl/pkg/client/core"
	"github.com/oslokommune/okctl/pkg/controller/common/reconciliation"

	"github.com/miekg/dns"
	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/client"
)

const identityManagerReconcilerIdentifier = "identity pool"

// identityManagerReconciler contains service and metadata for the relevant resource
type identityManagerReconciler struct {
	client client.IdentityManagerService
}

/*
Reconcile knows how to do what is necessary to ensure the desired state is achieved
Requires:
- Hosted Zone
- Nameservers setup
*/
func (z *identityManagerReconciler) Reconcile(ctx context.Context, meta reconciliation.Metadata, state *clientCore.StateHandlers) (reconciliation.Result, error) {
	action, err := z.determineAction(meta, state)
	if err != nil {
		return reconciliation.Result{}, fmt.Errorf("determining course of action: %w", err)
	}

	switch action {
	case reconciliation.ActionCreate:
		authDomain := fmt.Sprintf("auth.%s", meta.ClusterDeclaration.ClusterRootDomain)
		authFQDN := dns.Fqdn(authDomain)

		hz, err := state.Domain.GetPrimaryHostedZone()
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf("getting primary hosted zone: %w", err)
		}

		_, err = z.client.CreateIdentityPool(ctx, client.CreateIdentityPoolOpts{
			ID:           reconciliation.ClusterMetaAsID(meta.ClusterDeclaration.Metadata),
			AuthDomain:   authDomain,
			AuthFQDN:     authFQDN,
			HostedZoneID: hz.HostedZoneID,
		})
		if err != nil {
			if errors.IsKind(err, errors.Timeout) {
				return reconciliation.Result{Requeue: true}, nil
			}

			return reconciliation.Result{}, fmt.Errorf("creating identity manager resource: %w", err)
		}

		return reconciliation.Result{Requeue: false}, nil
	case reconciliation.ActionDelete:
		err = z.client.DeleteIdentityPool(
			ctx, reconciliation.ClusterMetaAsID(meta.ClusterDeclaration.Metadata),
		)
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf("deleting identity manager: %w", err)
		}

		return reconciliation.Result{Requeue: false}, nil
	case reconciliation.ActionWait:
		return reconciliation.Result{Requeue: true}, nil
	case reconciliation.ActionNoop:
		return reconciliation.Result{Requeue: false}, nil
	}

	return reconciliation.Result{}, fmt.Errorf("action %s is not implemented", string(action))
}

func (z *identityManagerReconciler) hasCreateDependenciesMet(meta reconciliation.Metadata, state *clientCore.StateHandlers) (bool, error) {
	dependenciesReady, err := reconciliation.AssertDependencyExistence(true,
		reconciliation.GenerateClusterExistenceTest(state, meta.ClusterDeclaration.Metadata.Name),
		reconciliation.GeneratePrimaryDomainDelegationTest(state),
		state.Domain.HasPrimaryHostedZone,
	)
	if err != nil {
		return false, fmt.Errorf("checking dependency existence: %w", err)
	}

	if !dependenciesReady {
		return false, nil
	}

	return true, nil
}

func (z *identityManagerReconciler) determineAction(meta reconciliation.Metadata, state *clientCore.StateHandlers) (reconciliation.Action, error) {
	userIndication := reconciliation.DetermineUserIndication(meta, meta.ClusterDeclaration.Integrations.Cognito)

	componentExists, err := state.IdentityManager.HasIdentityPool()
	if err != nil {
		return reconciliation.ActionNoop, fmt.Errorf("acquiring Identity Pool existence: %w", err)
	}

	switch userIndication {
	case reconciliation.ActionCreate:
		dependenciesReady, err := z.hasCreateDependenciesMet(meta, state)
		if err != nil {
			return reconciliation.ActionNoop, fmt.Errorf("checking dependencies: %w", err)
		}

		if !dependenciesReady {
			return reconciliation.ActionWait, nil
		}

		if componentExists {
			return reconciliation.ActionNoop, nil
		}

		return reconciliation.ActionCreate, nil
	case reconciliation.ActionDelete:
		clusterExists, err := state.Cluster.HasCluster(meta.ClusterDeclaration.Metadata.Name)
		if err != nil {
			return reconciliation.ActionNoop, fmt.Errorf("acquiring cluster existence: %w", err)
		}

		if !clusterExists || !componentExists {
			return reconciliation.ActionNoop, nil
		}

		return reconciliation.ActionDelete, nil
	}

	return reconciliation.ActionNoop, reconciliation.ErrIndecisive
}

// String returns the identifier for this reconciler
func (z *identityManagerReconciler) String() string {
	return identityManagerReconcilerIdentifier
}

// NewIdentityManagerReconciler creates a new reconciler for the Identity Manager resource
func NewIdentityManagerReconciler(client client.IdentityManagerService) reconciliation.Reconciler {
	return &identityManagerReconciler{
		client: client,
	}
}
