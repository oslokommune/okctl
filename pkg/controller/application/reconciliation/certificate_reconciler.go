package reconciliation

import (
	"context"
	"errors"
	"fmt"

	"github.com/asdine/storm/v3"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
	clientCore "github.com/oslokommune/okctl/pkg/client/core"
	"github.com/oslokommune/okctl/pkg/controller/common/reconciliation"
)

// Reconcile knows how to do what is necessary to ensure the desired state is achieved
func (c certificateReconciler) Reconcile(ctx context.Context, meta reconciliation.Metadata, state *clientCore.StateHandlers) (reconciliation.Result, error) {
	action, err := c.determineAction(meta, state)
	if err != nil {
		return reconciliation.Result{}, fmt.Errorf("determining course of action: %w", err)
	}

	appURL, err := meta.ApplicationDeclaration.URL()
	if err != nil {
		return reconciliation.Result{}, fmt.Errorf("acquiring application URL: %w", err)
	}

	clusterID := api.ID{
		Region:       meta.ClusterDeclaration.Metadata.Region,
		AWSAccountID: meta.ClusterDeclaration.Metadata.AccountID,
		ClusterName:  meta.ClusterDeclaration.Metadata.Name,
	}

	switch action {
	case reconciliation.ActionCreate:
		hz, err := c.domainService.GetPrimaryHostedZone(ctx)
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf("acquiring primary hosted zone: %w", err)
		}

		_, err = c.certificateService.CreateCertificate(ctx, client.CreateCertificateOpts{
			ID:           clusterID,
			FQDN:         appURL.String(),
			Domain:       appURL.String(),
			HostedZoneID: hz.HostedZoneID,
		})
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf("creating certificate: %w", err)
		}

		return reconciliation.Result{Requeue: false}, nil
	case reconciliation.ActionDelete:
		err = c.certificateService.DeleteCertificate(ctx, client.DeleteCertificateOpts{
			ID:     clusterID,
			Domain: appURL.String(),
		})
		if err != nil {
			return reconciliation.Result{}, fmt.Errorf("deleting certificate: %w", err)
		}

		return reconciliation.Result{Requeue: false}, nil
	case reconciliation.ActionNoop:
		return reconciliation.Result{Requeue: false}, nil
	case reconciliation.ActionWait:
		return reconciliation.Result{Requeue: true}, nil
	}

	return reconciliation.Result{}, fmt.Errorf("action %s not implemented", action)
}

func (c certificateReconciler) determineAction(meta reconciliation.Metadata, state *clientCore.StateHandlers) (reconciliation.Action, error) {
	userIndication := reconciliation.DetermineUserIndication(meta, meta.ApplicationDeclaration.HasIngress())

	appURL, err := meta.ApplicationDeclaration.URL()
	if err != nil {
		return reconciliation.ActionNoop, fmt.Errorf("acquiring application URL: %w", err)
	}

	switch userIndication {
	case reconciliation.ActionCreate:
		_, err = state.Certificate.GetCertificate(appURL.String())
		if err != nil && !errors.Is(err, storm.ErrNotFound) {
			return "", fmt.Errorf("acquiring certificate: %w", err)
		}

		if err == nil {
			return reconciliation.ActionNoop, nil
		}

		return reconciliation.ActionCreate, nil
	case reconciliation.ActionDelete:
		_, err = state.Certificate.GetCertificate(appURL.String())
		if err != nil {
			if errors.Is(err, storm.ErrNotFound) {
				return reconciliation.ActionNoop, nil
			}

			return "", fmt.Errorf("acquiring certificate: %w", err)
		}

		ingressExists, err := state.Kubernetes.HasResource(
			"ingress",
			meta.ApplicationDeclaration.Metadata.Namespace,
			meta.ApplicationDeclaration.Metadata.Name,
		)
		if err != nil {
			return "", fmt.Errorf("checking ingress existence: %w", err)
		}

		if ingressExists {
			return reconciliation.ActionWait, nil
		}

		return reconciliation.ActionDelete, nil
	}

	return "", reconciliation.ErrIndecisive
}

// String returns an identifier for this reconciler
func (c certificateReconciler) String() string {
	return certificateReconcilerName
}

// NewCertificateReconciler initializes a new certificate reconciler
func NewCertificateReconciler(cert client.CertificateService, domain client.DomainService) reconciliation.Reconciler {
	return &certificateReconciler{
		certificateService: cert,
		domainService:      domain,
	}
}

const certificateReconcilerName = "certificates"

type certificateReconciler struct {
	certificateService client.CertificateService
	domainService      client.DomainService
}
