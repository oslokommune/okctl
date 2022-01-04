package reconciliation

import (
	"context"
	"fmt"

	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/client"
	clientCore "github.com/oslokommune/okctl/pkg/client/core"
	"github.com/oslokommune/okctl/pkg/controller/common/reconciliation"
)

const argoCDApplicationReconcilerIdentifier = "ArgoCD"

// argoCDApplicationReconciler handles reconciliation for this feature
type argoCDApplicationReconciler struct {
	client client.ApplicationService
}

// Reconcile knows how to do what is necessary to ensure the desired state is achieved
func (a *argoCDApplicationReconciler) Reconcile(_ context.Context, meta reconciliation.Metadata, _ *clientCore.StateHandlers) (reconciliation.Result, error) {
	action, err := a.determineAction(meta)
	if err != nil {
		return reconciliation.Result{}, errors.E(err, "determining course of action")
	}

	switch action {
	case reconciliation.ActionCreate:
		err = a.client.CreateArgoCDApplicationManifest(client.CreateArgoCDApplicationManifestOpts{
			Cluster:     *meta.ClusterDeclaration,
			Application: meta.ApplicationDeclaration,
		})
		if err != nil {
			return reconciliation.Result{}, errors.E(err, "creating ArgoCD application manifest")
		}

		return reconciliation.Result{}, nil
	case reconciliation.ActionDelete:
		return reconciliation.Result{}, errors.New("deletion of an ArgoCD Application is not implemented")
	case reconciliation.ActionNoop:
		return reconciliation.Result{Requeue: false}, nil
	case reconciliation.ActionWait:
		return reconciliation.Result{Requeue: true}, nil
	}

	return reconciliation.Result{}, fmt.Errorf("action %s is not implemented", action)
}

func (a *argoCDApplicationReconciler) String() string {
	return argoCDApplicationReconcilerIdentifier
}

func (a *argoCDApplicationReconciler) determineAction(meta reconciliation.Metadata) (reconciliation.Action, error) {
	userIndication := reconciliation.DetermineUserIndication(meta, meta.ClusterDeclaration.Integrations.ArgoCD)

	switch userIndication {
	case reconciliation.ActionCreate:
		return reconciliation.ActionCreate, nil
	case reconciliation.ActionDelete:
		return reconciliation.ActionDelete, nil
	}

	return reconciliation.ActionNoop, reconciliation.ErrIndecisive
}

// NewArgoCDApplicationReconciler initializes a new ArgoCDApplicationReconciler
func NewArgoCDApplicationReconciler(client client.ApplicationService) reconciliation.Reconciler {
	return &argoCDApplicationReconciler{
		client: client,
	}
}
