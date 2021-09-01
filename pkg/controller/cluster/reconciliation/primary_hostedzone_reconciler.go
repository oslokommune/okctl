package reconciliation

import (
	"context"
	"fmt"
	"github.com/oslokommune/okctl/pkg/config/constant"

	"github.com/oslokommune/okctl/pkg/controller/common/reconciliation"

	clientCore "github.com/oslokommune/okctl/pkg/client/core"

	"github.com/miekg/dns"
	"github.com/oslokommune/okctl/pkg/client"
)

// zoneReconciler contains service and metadata for the relevant resource
type zoneReconciler struct {
	client client.DomainService
}

const zoneReconcilerIdentifier = "hosted zones"

// Reconcile knows how to do what is necessary to ensure the desired state is achieved
func (z *zoneReconciler) Reconcile(ctx context.Context, meta reconciliation.Metadata, state *clientCore.StateHandlers) (reconciliation.Result, error) {
	action, err := z.determineAction(meta, state)
	if err != nil {
		return reconciliation.Result{}, fmt.Errorf(constant.ReconcilerDetermineActionError, err)
	}

	switch action {
	case reconciliation.ActionCreate:
		_, err = z.client.CreatePrimaryHostedZone(ctx, client.CreatePrimaryHostedZoneOpts{
			ID:     reconciliation.ClusterMetaAsID(meta.ClusterDeclaration.Metadata),
			Domain: meta.ClusterDeclaration.ClusterRootDomain,
			FQDN:   dns.Fqdn(meta.ClusterDeclaration.ClusterRootDomain),
		})
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf(constant.CreateHostedZoneError, err)
		}

		return reconciliation.Result{Requeue: false}, nil
	case reconciliation.ActionDelete:
		hz, err := state.Domain.GetPrimaryHostedZone()
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf(constant.GetPrimaryHostedZoneError, err)
		}

		err = z.client.DeletePrimaryHostedZone(ctx, client.DeletePrimaryHostedZoneOpts{
			ID:           reconciliation.ClusterMetaAsID(meta.ClusterDeclaration.Metadata),
			HostedZoneID: hz.HostedZoneID,
		})
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf(constant.DeletePrimaryHostedZoneError, err)
		}

		return reconciliation.Result{Requeue: false}, nil
	case reconciliation.ActionWait:
		return reconciliation.Result{Requeue: true}, nil
	case reconciliation.ActionNoop:
		return reconciliation.Result{Requeue: false}, nil
	}

	return reconciliation.Result{}, fmt.Errorf(constant.ActionNotImplementedError, string(action))
}

func (z *zoneReconciler) determineAction(meta reconciliation.Metadata, state *clientCore.StateHandlers) (reconciliation.Action, error) {
	userIndication := reconciliation.DetermineUserIndication(meta, true)

	primaryHostedZoneExists, err := state.Domain.HasPrimaryHostedZone()
	if err != nil {
		return reconciliation.ActionNoop, fmt.Errorf(constant.QueryStateError, err)
	}

	switch userIndication {
	case reconciliation.ActionCreate:
		if primaryHostedZoneExists {
			return reconciliation.ActionNoop, nil
		}

		return reconciliation.ActionCreate, nil
	case reconciliation.ActionDelete:
		if !primaryHostedZoneExists {
			return reconciliation.ActionNoop, nil
		}

		dependenciesReady, err := reconciliation.AssertDependencyExistence(false,
			reconciliation.GenerateClusterExistenceTest(state, meta.ClusterDeclaration.Metadata.Name),
		)
		if err != nil {
			return reconciliation.ActionNoop, err
		}

		if !dependenciesReady {
			return reconciliation.ActionWait, nil
		}

		return reconciliation.ActionDelete, nil
	}

	return reconciliation.ActionNoop, reconciliation.ErrIndecisive
}

// String returns the identifier for this reconciler
func (z *zoneReconciler) String() string {
	return zoneReconcilerIdentifier
}

// NewZoneReconciler creates a new reconciler for the Hosted Zone resource
func NewZoneReconciler(client client.DomainService) reconciliation.Reconciler {
	return &zoneReconciler{
		client: client,
	}
}
