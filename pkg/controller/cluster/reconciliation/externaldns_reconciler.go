package reconciliation

import (
	"context"
	"fmt"

	clientCore "github.com/oslokommune/okctl/pkg/client/core"
	"github.com/oslokommune/okctl/pkg/controller/common/reconciliation"

	"github.com/mishudark/errors"

	"github.com/oslokommune/okctl/pkg/client"
)

// externalDNSReconciler contains service and metadata for the relevant resource
type externalDNSReconciler struct {
	client client.ExternalDNSService
}

const externalDNSReconcilerIdentifier = "DNS controller"

// Reconcile knows how to do what is necessary to ensure the desired state is achieved
func (z *externalDNSReconciler) Reconcile(ctx context.Context, meta reconciliation.Metadata, state *clientCore.StateHandlers) (result reconciliation.Result, err error) {
	action, err := z.determineAction(meta, state)
	if err != nil {
		return reconciliation.Result{}, fmt.Errorf("determining course of action: %w", err)
	}

	switch action {
	case reconciliation.ActionCreate:
		hz, err := state.Domain.GetPrimaryHostedZone()
		if err != nil {
			return result, fmt.Errorf("getting primary hosted zone: %w", err)
		}

		_, err = z.client.CreateExternalDNS(ctx, client.CreateExternalDNSOpts{
			ID:           reconciliation.ClusterMetaAsID(meta.ClusterDeclaration.Metadata),
			HostedZoneID: hz.HostedZoneID,
			Domain:       meta.ClusterDeclaration.ClusterRootDomain,
		})
		if err != nil {
			result.Requeue = errors.IsKind(err, errors.Timeout)

			return result, fmt.Errorf("creating external DNS: %w", err)
		}

		return reconciliation.Result{Requeue: false}, nil
	case reconciliation.ActionDelete:
		err = z.client.DeleteExternalDNS(ctx, reconciliation.ClusterMetaAsID(meta.ClusterDeclaration.Metadata))
		if err != nil {
			return result, fmt.Errorf("deleting external DNS: %w", err)
		}

		return reconciliation.Result{Requeue: false}, nil
	case reconciliation.ActionWait:
		return reconciliation.Result{Requeue: true}, nil
	case reconciliation.ActionNoop:
		return reconciliation.Result{Requeue: false}, nil
	}

	return reconciliation.Result{}, fmt.Errorf("action %s is not implemented", string(action))
}

//nolint:funlen
func (z *externalDNSReconciler) determineAction(meta reconciliation.Metadata, state *clientCore.StateHandlers) (reconciliation.Action, error) {
	userIndication := reconciliation.DetermineUserIndication(meta, meta.ClusterDeclaration.Integrations.ExternalDNS)

	clusterExistenceTest := reconciliation.GenerateClusterExistenceTest(state, meta.ClusterDeclaration.Metadata.Name)

	switch userIndication {
	case reconciliation.ActionCreate:
		dependenciesReady, err := reconciliation.AssertDependencyExistence(true,
			reconciliation.GeneratePrimaryDomainDelegationTest(state),
			clusterExistenceTest,
		)
		if err != nil {
			return reconciliation.ActionNoop, fmt.Errorf("checking dependencies: %w", err)
		}

		if !dependenciesReady {
			return reconciliation.ActionWait, nil
		}

		componentExists, err := state.ExternalDNS.HasExternalDNS()
		if err != nil {
			return reconciliation.ActionNoop, fmt.Errorf("acquiring DNS controller existence: %w", err)
		}

		if componentExists {
			return reconciliation.ActionNoop, nil
		}

		return reconciliation.ActionCreate, nil
	case reconciliation.ActionDelete:
		clusterExists, err := clusterExistenceTest()
		if err != nil {
			return reconciliation.ActionNoop, fmt.Errorf("checking cluster existence: %w", err)
		}

		if !clusterExists {
			return reconciliation.ActionNoop, nil
		}

		componentExists, err := state.ExternalDNS.HasExternalDNS()
		if err != nil {
			return reconciliation.ActionNoop, fmt.Errorf("acquiring DNS controller existence: %w", err)
		}

		if !componentExists {
			return reconciliation.ActionNoop, nil
		}

		dependenciesReady, err := reconciliation.AssertDependencyExistence(false,
			state.ArgoCD.HasArgoCD,
			state.Monitoring.HasKubePromStack,
		)
		if err != nil {
			return reconciliation.ActionNoop, fmt.Errorf("checking dependencies: %w", err)
		}

		if !dependenciesReady {
			return reconciliation.ActionWait, nil
		}

		return reconciliation.ActionDelete, nil
	}

	return reconciliation.ActionNoop, reconciliation.ErrIndecisive
}

// String returns the identifier for this reconciler
func (z *externalDNSReconciler) String() string {
	return externalDNSReconcilerIdentifier
}

// NewExternalDNSReconciler creates a new reconciler for the ExternalDNS resource
func NewExternalDNSReconciler(client client.ExternalDNSService) reconciliation.Reconciler {
	return &externalDNSReconciler{
		client: client,
	}
}
