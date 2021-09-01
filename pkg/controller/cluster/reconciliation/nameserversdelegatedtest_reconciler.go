package reconciliation

import (
	"context"
	"fmt"
	"github.com/oslokommune/okctl/pkg/config/constant"
	"time"

	"github.com/oslokommune/okctl/pkg/client"

	"github.com/oslokommune/okctl/pkg/controller/common/reconciliation"

	clientCore "github.com/oslokommune/okctl/pkg/client/core"

	"github.com/logrusorgru/aurora"
	"github.com/oslokommune/okctl/pkg/domain"
)

const (
	defaultTestingIntervalMinutes                = 5 * time.Minute
	nameserversDelegatedTestReconcilerIdentifier = "nameserver delegation verification"
)

type nameserversDelegatedTestReconciler struct {
	domainService client.DomainService
}

// Reconcile knows how to do what is necessary to ensure the desired state is achieved
func (n *nameserversDelegatedTestReconciler) Reconcile(
	ctx context.Context,
	meta reconciliation.Metadata,
	state *clientCore.StateHandlers,
) (reconciliation.Result, error) {
	action, err := n.determineAction(meta, state)
	if err != nil {
		return reconciliation.Result{}, fmt.Errorf(constant.ReconcilerDetermineActionError, err)
	}

	switch action {
	case reconciliation.ActionCreate:
		_, _ = fmt.Fprintf(
			meta.Out,
			delegationRequestMessage,
			aurora.Green("nameserver delegation request"),
			aurora.Bold("#kjøremiljø-support"),
		)

		hz, err := state.Domain.GetPrimaryHostedZone()
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf(constant.GetPrimaryHostedZoneError, err)
		}

		err = domain.ShouldHaveNameServers(hz.FQDN, hz.NameServers)
		if err != nil {
			_, _ = fmt.Fprintf(meta.Out, "failed to validate nameservers: %s", err)

			return reconciliation.Result{
				Requeue:      true,
				RequeueAfter: defaultTestingIntervalMinutes,
			}, fmt.Errorf(constant.ValidateNameserversError, err)
		}

		err = n.domainService.SetHostedZoneDelegation(ctx, meta.ClusterDeclaration.ClusterRootDomain, true)
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf(constant.SetHostedZoneDelegationStatus, err)
		}

		return reconciliation.Result{Requeue: false}, nil
	case reconciliation.ActionDelete:
		return reconciliation.Result{Requeue: false}, nil
	case reconciliation.ActionWait:
		return reconciliation.Result{Requeue: true}, nil
	case reconciliation.ActionNoop:
		return reconciliation.Result{Requeue: false}, nil
	}

	return reconciliation.Result{}, fmt.Errorf(constant.ActionNotImplementedError, string(action))
}

func (n *nameserversDelegatedTestReconciler) determineAction(meta reconciliation.Metadata, state *clientCore.StateHandlers) (reconciliation.Action, error) {
	userIndication := reconciliation.DetermineUserIndication(meta, true)

	delegationTest := reconciliation.GeneratePrimaryDomainDelegationTest(state)

	switch userIndication {
	case reconciliation.ActionCreate:
		dependenciesReady, err := reconciliation.AssertDependencyExistence(true,
			reconciliation.GenerateClusterExistenceTest(state, meta.ClusterDeclaration.Metadata.Name),
			state.Domain.HasPrimaryHostedZone,
		)
		if err != nil {
			return reconciliation.ActionNoop, fmt.Errorf(constant.TestDependenciesError, err)
		}

		if !dependenciesReady {
			return reconciliation.ActionWait, nil
		}

		isDelegated, err := delegationTest()
		if err != nil {
			return reconciliation.ActionNoop, fmt.Errorf(constant.CheckHostedZoneDelegationStateError, err)
		}

		if isDelegated {
			return reconciliation.ActionNoop, nil
		}

		return reconciliation.ActionCreate, nil
	case reconciliation.ActionDelete:
		primaryHZExists, err := state.Domain.HasPrimaryHostedZone()
		if err != nil {
			return reconciliation.ActionNoop, fmt.Errorf(constant.CheckDomainExistenceError, err)
		}

		if !primaryHZExists {
			return reconciliation.ActionNoop, nil
		}

		isDelegated, err := delegationTest()
		if err != nil {
			return reconciliation.ActionNoop, fmt.Errorf(constant.CheckHostedZoneDelegationStateError, err)
		}

		if !isDelegated {
			return reconciliation.ActionNoop, nil
		}

		return reconciliation.ActionDelete, nil
	}

	return reconciliation.ActionNoop, reconciliation.ErrIndecisive
}

// String returns the identifier for this reconciler
func (n *nameserversDelegatedTestReconciler) String() string {
	return nameserversDelegatedTestReconcilerIdentifier
}

// NewNameserverDelegatedTestReconciler creates a new reconciler for the nameserver record delegation test resource
func NewNameserverDelegatedTestReconciler(domainService client.DomainService) reconciliation.Reconciler {
	return &nameserversDelegatedTestReconciler{domainService: domainService}
}

const delegationRequestMessage = `


A %s has been submitted. We'll process this request as soon as possible.
Let us know in %s if it takes too long


`
