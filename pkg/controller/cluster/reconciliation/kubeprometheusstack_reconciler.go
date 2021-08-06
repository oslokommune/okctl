package reconciliation

import (
	"context"
	"fmt"

	"github.com/oslokommune/okctl/pkg/cfn"
	"github.com/oslokommune/okctl/pkg/controller/common/reconciliation"

	clientCore "github.com/oslokommune/okctl/pkg/client/core"

	"github.com/oslokommune/okctl/pkg/client"
)

const kubePrometheusStackReconcilerIdentifier = "Prometheus & Grafana"

type kubePrometheusStackReconciler struct {
	client client.MonitoringService
}

// Reconcile knows how to do what is necessary to ensure the desired state is achieved
func (z *kubePrometheusStackReconciler) Reconcile(ctx context.Context, meta reconciliation.Metadata, state *clientCore.StateHandlers) (reconciliation.Result, error) {
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

		im, err := state.IdentityManager.GetIdentityPool(
			cfn.NewStackNamer().IdentityPool(meta.ClusterDeclaration.Metadata.Name),
		)
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf("getting identity pool: %w", err)
		}

		_, err = z.client.CreateKubePromStack(ctx, client.CreateKubePromStackOpts{
			ID:           reconciliation.ClusterMetaAsID(meta.ClusterDeclaration.Metadata),
			Domain:       meta.ClusterDeclaration.ClusterRootDomain,
			HostedZoneID: hz.HostedZoneID,
			AuthDomain:   im.AuthDomain,
			UserPoolID:   im.UserPoolID,
		})
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf("creating kubepromstack: %w", err)
		}

		return reconciliation.Result{Requeue: false}, nil
	case reconciliation.ActionDelete:
		err = z.client.DeleteKubePromStack(ctx, client.DeleteKubePromStackOpts{
			ID:     reconciliation.ClusterMetaAsID(meta.ClusterDeclaration.Metadata),
			Domain: meta.ClusterDeclaration.ClusterRootDomain,
		})
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf("deleting kubepromstack: %w", err)
		}

		return reconciliation.Result{Requeue: false}, nil
	case reconciliation.ActionWait:
		return reconciliation.Result{Requeue: true}, nil
	case reconciliation.ActionNoop:
		return reconciliation.Result{Requeue: false}, nil
	}

	return reconciliation.Result{}, fmt.Errorf("action %s is not implemented", string(action))
}

func (z *kubePrometheusStackReconciler) determineAction(meta reconciliation.Metadata, state *clientCore.StateHandlers) (reconciliation.Action, error) {
	userIndication := reconciliation.DetermineUserIndication(meta, meta.ClusterDeclaration.Integrations.KubePromStack)

	componentExists, err := state.Monitoring.HasKubePromStack()
	if err != nil {
		return reconciliation.ActionNoop, fmt.Errorf("checking component existence: %w", err)
	}

	clusterExistenceTest := reconciliation.GenerateClusterExistenceTest(
		state,
		meta.ClusterDeclaration.Metadata.Name,
	)

	switch userIndication {
	case reconciliation.ActionCreate:
		if componentExists {
			return reconciliation.ActionNoop, nil
		}

		dependenciesReady, err := reconciliation.AssertDependencyExistence(true,
			clusterExistenceTest,
			reconciliation.GeneratePrimaryDomainDelegationTest(state),
		)
		if err != nil {
			return reconciliation.ActionNoop, fmt.Errorf("checking dependencies: %w", err)
		}

		if !dependenciesReady {
			return reconciliation.ActionWait, nil
		}

		return reconciliation.ActionCreate, nil
	case reconciliation.ActionDelete:
		clusterExists, err := clusterExistenceTest()
		if err != nil {
			return reconciliation.ActionNoop, fmt.Errorf("checking cluster existence: %w", err)
		}

		if !componentExists || !clusterExists {
			return reconciliation.ActionNoop, nil
		}

		return reconciliation.ActionDelete, nil
	}

	return reconciliation.ActionNoop, reconciliation.ErrIndecisive
}

// String returns the identifier type
func (z *kubePrometheusStackReconciler) String() string {
	return kubePrometheusStackReconcilerIdentifier
}

// NewKubePrometheusStackReconciler creates a new reconciler for the KubePromStack resource
func NewKubePrometheusStackReconciler(client client.MonitoringService) reconciliation.Reconciler {
	return &kubePrometheusStackReconciler{
		client: client,
	}
}
