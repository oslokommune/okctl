package core

import (
	"context"
	"fmt"
	"path"

	"github.com/oslokommune/okctl/pkg/config/constant"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/scaffold"
)

type applicationService struct {
	cert  client.CertificateService
	store client.ApplicationStore
}

func (s *applicationService) createCertificate(ctx context.Context, id *api.ID, hostedZoneID, fqdn string) (string, error) {
	cert, certFnErr := s.cert.CreateCertificate(ctx, client.CreateCertificateOpts{
		ID:           *id,
		FQDN:         fqdn,
		Domain:       fqdn,
		HostedZoneID: hostedZoneID,
	})
	if certFnErr != nil {
		return "", certFnErr
	}

	return cert.ARN, nil
}

// ScaffoldApplication turns a file path into Kubernetes resources
// nolint: funlen
func (s *applicationService) ScaffoldApplication(ctx context.Context, opts *client.ScaffoldApplicationOpts) error {
	err := opts.Validate()
	if err != nil {
		return err
	}

	relativeApplicationDir := path.Join(opts.OutputDir, constant.DefaultApplicationsOutputDir, opts.Application.Metadata.Name)
	relativeArgoCDSourcePath := path.Join(relativeApplicationDir, constant.DefaultApplicationOverlayDir, opts.ID.ClusterName)

	base, err := scaffold.GenerateApplicationBase(opts.Application, opts.IACRepoURL, relativeArgoCDSourcePath)
	if err != nil {
		return fmt.Errorf("creating a new application deployment: %w", err)
	}

	certArn, err := s.createCertificate(
		ctx,
		opts.ID,
		opts.HostedZoneID,
		fmt.Sprintf("%s.%s", opts.Application.SubDomain, opts.HostedZoneDomain),
	)
	if err != nil {
		return fmt.Errorf("create certificate: %w", err)
	}

	overlay, err := scaffold.GenerateApplicationOverlay(opts.Application, opts.HostedZoneDomain, certArn)
	if err != nil {
		return fmt.Errorf("generating application overlay: %w", err)
	}

	applicationScaffold := &client.ScaffoldedApplication{
		ApplicationName:      opts.Application.Metadata.Name,
		ClusterName:          opts.ID.ClusterName,
		BaseKustomization:    base.Kustomization,
		Deployment:           base.Deployment,
		Service:              base.Service,
		Ingress:              base.Ingress,
		Volume:               base.Volumes,
		OverlayKustomization: overlay.Kustomization,
		ArgoCDResource:       base.ArgoApplication,
		IngressPatch:         overlay.IngressPatch,
		ServicePatch:         overlay.ServicePatch,
		DeploymentPatch:      overlay.DeploymentPatch,
	}

	_, err = s.store.SaveApplication(applicationScaffold)
	if err != nil {
		return err
	}

	return nil
}

// NewApplicationService initializes a new Scaffold application service
func NewApplicationService(
	cert client.CertificateService,
	store client.ApplicationStore,
) client.ApplicationService {
	return &applicationService{
		cert:  cert,
		store: store,
	}
}
