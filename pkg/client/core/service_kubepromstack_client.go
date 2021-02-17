package core

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/oslokommune/okctl/pkg/spinner"

	"github.com/oslokommune/okctl/pkg/api"

	"github.com/oslokommune/okctl/pkg/client"
)

type kubePrometheusStackService struct {
	spinner spinner.Spinner
	api     client.KubePrometheusStackAPI
	store   client.KubePrometheusStackStore
	report  client.KubePrometheusStackReport

	cert     client.CertificateService
	ident    client.IdentityManagerService
	param    client.ParameterService
	manifest client.ManifestService
}

// nolint: funlen
func (s *kubePrometheusStackService) CreateKubePrometheusStack(ctx context.Context, opts client.CreateKubePrometheusStackOpts) (*client.KubePrometheusStack, error) {
	err := s.spinner.Start("kubepromstack")
	if err != nil {
		return nil, err
	}

	defer func() {
		err = s.spinner.Stop()
	}()

	cert, err := s.cert.CreateCertificate(ctx, api.CreateCertificateOpts{
		ID:           opts.ID,
		FQDN:         fmt.Sprintf("grafana.%s.", opts.Domain),
		Domain:       fmt.Sprintf("grafana.%s", opts.Domain),
		HostedZoneID: opts.HostedZoneID,
	})
	if err != nil {
		return nil, err
	}

	poolClient, err := s.ident.CreateIdentityPoolClient(ctx, api.CreateIdentityPoolClientOpts{
		ID:          opts.ID,
		UserPoolID:  opts.UserPoolID,
		Purpose:     "grafana",
		CallbackURL: fmt.Sprintf("https://%s/login/generic_oauth", cert.Domain),
	})
	if err != nil {
		return nil, err
	}

	clientSecret, err := s.param.CreateSecret(ctx, api.CreateSecretOpts{
		ID:     opts.ID,
		Name:   "client-secret",
		Secret: poolClient.ClientSecret,
	})
	if err != nil {
		return nil, err
	}

	cookieSecret, err := s.param.CreateSecret(ctx, api.CreateSecretOpts{
		ID:     opts.ID,
		Name:   "secret-key",
		Secret: uuid.NewString(),
	})
	if err != nil {
		return nil, err
	}

	adminUser, err := s.param.CreateSecret(ctx, api.CreateSecretOpts{
		ID:     opts.ID,
		Name:   "admin-user",
		Secret: uuid.NewString(),
	})
	if err != nil {
		return nil, err
	}

	adminPass, err := s.param.CreateSecret(ctx, api.CreateSecretOpts{
		ID:     opts.ID,
		Name:   "admin-pass",
		Secret: uuid.NewString(),
	})
	if err != nil {
		return nil, err
	}

	manifest, err := s.manifest.CreateExternalSecret(ctx, client.CreateExternalSecretOpts{
		ID: opts.ID,
		Manifests: []api.Manifest{
			{
				Name:      "grafana-secrets-cm",
				Namespace: "monitoring",
				Data: []api.Data{
					{
						Key:  clientSecret.Path,
						Name: clientSecret.Name,
					},
					{
						Key:  cookieSecret.Path,
						Name: cookieSecret.Name,
					},
					{
						Key:  adminUser.Path,
						Name: adminUser.Name,
					},
					{
						Key:  adminPass.Path,
						Name: adminPass.Name,
					},
				},
			},
		},
	})

	chart, err := s.api.CreateKubePrometheusStackHelmChart(api.CreateKubePrometheusStackOpts{
		ID:                     opts.ID,
		CertificateARN:         cert.CertificateARN,
		Hostname:               cert.Domain,
		AuthHostname:           opts.AuthDomain,
		ClientID:               poolClient.ClientID,
		SecretsConfigName:      "grafana-secrets-cm",
		SecretsCookieSecretKey: cookieSecret.Name,
		SecretsClientSecretKey: clientSecret.Name,
		SecretsAdminUserKey:    adminUser.Name,
		SecretsAdminPassKey:    "admin-pass-key",
	})
	if err != nil {
		return nil, err
	}

	stack := &client.KubePrometheusStack{
		CertificateARN:         cert.CertificateARN,
		Hostname:               cert.Domain,
		AuthHostname:           opts.AuthDomain,
		ClientID:               poolClient.ClientID,
		SecretsConfigName:      "grafana-secrets-cm",
		SecretsCookieSecretKey: cookieSecret.Name,
		SecretsClientSecretKey: clientSecret.Name,
		SecretsAdminUserKey:    adminUser.Name,
		SecretsAdminPassKey:    adminPass.Name,
		Chart:                  chart,
		Certificate:            cert,
		IdentityPoolClient:     poolClient,
		ExternalSecret:         manifest,
	}

	report, err := s.store.SaveKubePrometheusStack(stack)
	if err != nil {
		return nil, err
	}

	err = s.report.ReportCreateKubePrometheusStack(stack, report)
	if err != nil {
		return nil, err
	}

	return stack, nil
}

// NewKubePrometheusStackService returns an initialised service
func NewKubePrometheusStackService(
	spinner spinner.Spinner,
	api client.KubePrometheusStackAPI,
	store client.KubePrometheusStackStore,
	report client.KubePrometheusStackReport,
	cert client.CertificateService,
	ident client.IdentityManagerService,
	manifest client.ManifestService,
	param client.ParameterService,
) client.KubePrometheusStackService {
	return &kubePrometheusStackService{
		spinner:  spinner,
		api:      api,
		store:    store,
		report:   report,
		cert:     cert,
		ident:    ident,
		manifest: manifest,
		param:    param,
	}
}
