package reconciliation

import (
	"context"
	"fmt"

	"github.com/oslokommune/okctl/pkg/config/constant"
	"github.com/oslokommune/okctl/pkg/controller/common/reconciliation"

	clientCore "github.com/oslokommune/okctl/pkg/client/core"

	"github.com/oslokommune/okctl/pkg/client"
)

const nameserverDelegationReconcilerIdentifier = "nameserver delegation"

// nameserverDelegationReconciler handles creation (later edit and deletion) of nameserver delegation resources.
// A nameserver delegation consists of creating a request to add a NS record to the top level domain and verifying
// that the delegation has happened.
type nameserverDelegationReconciler struct {
	client client.NSRecordDelegationService
}

// Reconcile knows how to do what is necessary to ensure the desired state is achieved
//nolint:funlen
func (z *nameserverDelegationReconciler) Reconcile(_ context.Context, meta reconciliation.Metadata, state *clientCore.StateHandlers) (reconciliation.Result, error) {
	var labels []string

	if meta.ClusterDeclaration.Experimental != nil {
		if meta.ClusterDeclaration.Experimental.AutomatizeZoneDelegation {
			labels = append(labels, constant.DefaultAutomaticPullRequestMergeLabel)
		}
	}

	action, err := z.determineAction(meta, state)
	if err != nil {
		return reconciliation.Result{}, fmt.Errorf("determining course of action: %w", err)
	}

	switch action {
	case reconciliation.ActionCreate:
		hz, err := state.Domain.GetPrimaryHostedZone()
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf("getting primary hosted zone: %w", err)
		}

		err = z.client.InitiateDomainDelegation(client.InitiateDomainDelegationOpts{
			ClusterID:             reconciliation.ClusterMetaAsID(meta.ClusterDeclaration.Metadata),
			PrimaryHostedZoneFQDN: hz.FQDN,
			Nameservers:           hz.NameServers,
			Labels:                labels,
		})
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf("initiating dns zone delegation: %w", err)
		}

		return reconciliation.Result{Requeue: false}, nil
	case reconciliation.ActionDelete:
		hz, err := state.Domain.GetPrimaryHostedZone()
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf("getting primary hosted zone: %w", err)
		}

		err = z.client.RevokeDomainDelegation(client.RevokeDomainDelegationOpts{
			ClusterID:             reconciliation.ClusterMetaAsID(meta.ClusterDeclaration.Metadata),
			PrimaryHostedZoneFQDN: hz.FQDN,
			Labels:                labels,
		})
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf("revoking dns zone delegation: %w", err)
		}

		return reconciliation.Result{Requeue: false}, nil
	case reconciliation.ActionWait:
		return reconciliation.Result{Requeue: true}, nil
	case reconciliation.ActionNoop:
		return reconciliation.Result{}, nil
	}

	return reconciliation.Result{}, fmt.Errorf("action %s is not implemented", string(action))
}

func (z *nameserverDelegationReconciler) determineAction(meta reconciliation.Metadata, state *clientCore.StateHandlers) (reconciliation.Action, error) {
	userIndication := reconciliation.DetermineUserIndication(meta, true)

	primaryHZExists, err := state.Domain.HasPrimaryHostedZone()
	if err != nil {
		return reconciliation.ActionNoop, fmt.Errorf("acquiring primary hosted zone state: %w", err)
	}

	switch userIndication {
	case reconciliation.ActionCreate:
		if !primaryHZExists {
			return reconciliation.ActionWait, nil
		}

		hz, err := state.Domain.GetPrimaryHostedZone()
		if err != nil {
			return reconciliation.ActionNoop, fmt.Errorf("acquiring state: %w", err)
		}

		if hz.IsDelegated {
			return reconciliation.ActionNoop, nil
		}

		return reconciliation.ActionCreate, nil
	case reconciliation.ActionDelete:
		if !primaryHZExists {
			return reconciliation.ActionNoop, nil
		}

		hz, err := state.Domain.GetPrimaryHostedZone()
		if err != nil {
			return reconciliation.ActionNoop, fmt.Errorf("acquiring state: %w", err)
		}

		if !hz.IsDelegated {
			return reconciliation.ActionNoop, nil
		}

		return reconciliation.ActionDelete, nil
	}

	return reconciliation.ActionNoop, reconciliation.ErrIndecisive
}

// String returns the identifier for this reconciler
func (z *nameserverDelegationReconciler) String() string {
	return nameserverDelegationReconcilerIdentifier
}

// NewNameserverDelegationReconciler creates a new reconciler for the nameserver record delegation resource
func NewNameserverDelegationReconciler(client client.NSRecordDelegationService) reconciliation.Reconciler {
	return &nameserverDelegationReconciler{
		client: client,
	}
}
