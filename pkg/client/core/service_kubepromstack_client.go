package core

import (
	"context"
	"fmt"

	"github.com/oslokommune/okctl/pkg/config"

	"github.com/oslokommune/okctl/pkg/client/store"

	"github.com/google/uuid"

	"github.com/oslokommune/okctl/pkg/spinner"

	"github.com/oslokommune/okctl/pkg/api"

	"github.com/oslokommune/okctl/pkg/client"
)

type kubePromStackService struct {
	spinner spinner.Spinner
	api     client.KubePromStackAPI
	store   client.KubePromStackStore
	state   client.KubePromStackState
	report  client.KubePromStackReport

	cert     client.CertificateService
	ident    client.IdentityManagerService
	param    client.ParameterService
	manifest client.ManifestService
}

const (
	grafanaSubDomain = "grafana"
	grafanaPurpose   = "grafana"
	clientSecretName = "client-secret"
	secretKeyName    = "secret-key"
	adminUserName    = "admin-user"
	adminPassName    = "admin-pass"
	secretsCfgName   = "grafana-secrets-cm"
)

func grafanaDomain(baseDomain string) string {
	return fmt.Sprintf("%s.%s", grafanaSubDomain, baseDomain)
}

func (s *kubePromStackService) DeleteKubePromStack(ctx context.Context, opts client.DeleteKubePromStackOpts) error {
	err := s.spinner.Start("kubepromstack")
	if err != nil {
		return err
	}

	defer func() {
		err = s.spinner.Stop()
	}()

	// This is too heavy handed, we should only remove the chart
	// and the secrets manifest. We can, however, wait until Loki
	// with doing this.
	err = s.manifest.DeleteNamespace(ctx, api.DeleteNamespaceOpts{
		ID:        opts.ID,
		Namespace: config.DefaultMonitoringNamespace,
	})
	if err != nil {
		return err
	}

	err = s.ident.DeleteIdentityPoolClient(ctx, api.DeleteIdentityPoolClientOpts{
		ID:      opts.ID,
		Purpose: grafanaPurpose,
	})

	for _, secretName := range []string{clientSecretName, secretKeyName, adminUserName, adminPassName} {
		if err = s.param.DeleteSecret(ctx, api.DeleteSecretOpts{Name: secretName}); err != nil {
			return err
		}
	}

	err = s.cert.DeleteCertificate(ctx, api.DeleteCertificateOpts{
		ID:     opts.ID,
		Domain: grafanaDomain(opts.Domain),
	})
	if err != nil {
		return err
	}

	r1, err := s.store.RemoveKubePromStack(opts.ID)
	if err != nil {
		return err
	}

	r2, err := s.state.RemoveKubePromStack(opts.ID)
	if err != nil {
		return err
	}

	return s.report.ReportRemoveKubePromStack([]*store.Report{r1, r2})
}

// nolint: funlen
func (s *kubePromStackService) CreateKubePromStack(ctx context.Context, opts client.CreateKubePromStackOpts) (*client.KubePromStack, error) {
	err := s.spinner.Start("kubepromstack")
	if err != nil {
		return nil, err
	}

	defer func() {
		err = s.spinner.Stop()
	}()

	cert, err := s.cert.CreateCertificate(ctx, api.CreateCertificateOpts{
		ID:           opts.ID,
		FQDN:         fmt.Sprintf("%s.", grafanaDomain(opts.Domain)),
		Domain:       grafanaDomain(opts.Domain),
		HostedZoneID: opts.HostedZoneID,
	})
	if err != nil {
		return nil, err
	}

	poolClient, err := s.ident.CreateIdentityPoolClient(ctx, api.CreateIdentityPoolClientOpts{
		ID:          opts.ID,
		UserPoolID:  opts.UserPoolID,
		Purpose:     grafanaPurpose,
		CallbackURL: fmt.Sprintf("https://%s/login/generic_oauth", cert.Domain),
	})
	if err != nil {
		return nil, err
	}

	clientSecret, err := s.param.CreateSecret(ctx, api.CreateSecretOpts{
		ID:     opts.ID,
		Name:   clientSecretName,
		Secret: poolClient.ClientSecret,
	})
	if err != nil {
		return nil, err
	}

	cookieSecret, err := s.param.CreateSecret(ctx, api.CreateSecretOpts{
		ID:     opts.ID,
		Name:   secretKeyName,
		Secret: uuid.NewString(),
	})
	if err != nil {
		return nil, err
	}

	adminUser, err := s.param.CreateSecret(ctx, api.CreateSecretOpts{
		ID:     opts.ID,
		Name:   adminUserName,
		Secret: uuid.NewString(),
	})
	if err != nil {
		return nil, err
	}

	adminPass, err := s.param.CreateSecret(ctx, api.CreateSecretOpts{
		ID:     opts.ID,
		Name:   adminPassName,
		Secret: uuid.NewString(),
	})
	if err != nil {
		return nil, err
	}

	manifest, err := s.manifest.CreateExternalSecret(ctx, client.CreateExternalSecretOpts{
		ID: opts.ID,
		Manifests: []api.Manifest{
			{
				Name:      secretsCfgName,
				Namespace: config.DefaultMonitoringNamespace,
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

	chart, err := s.api.CreateKubePromStack(api.CreateKubePrometheusStackOpts{
		ID:                     opts.ID,
		CertificateARN:         cert.CertificateARN,
		Hostname:               cert.Domain,
		AuthHostname:           opts.AuthDomain,
		ClientID:               poolClient.ClientID,
		SecretsConfigName:      secretsCfgName,
		SecretsCookieSecretKey: cookieSecret.Name,
		SecretsClientSecretKey: clientSecret.Name,
		SecretsAdminUserKey:    adminUser.Name,
		SecretsAdminPassKey:    adminPass.Name,
	})
	if err != nil {
		return nil, err
	}

	stack := &client.KubePromStack{
		ID:                     opts.ID,
		CertificateARN:         cert.CertificateARN,
		Hostname:               cert.Domain,
		AuthHostname:           opts.AuthDomain,
		ClientID:               poolClient.ClientID,
		SecretsConfigName:      secretsCfgName,
		SecretsCookieSecretKey: cookieSecret.Name,
		SecretsClientSecretKey: clientSecret.Name,
		SecretsAdminUserKey:    adminUser.Name,
		SecretsAdminPassKey:    adminPass.Name,
		Chart:                  chart,
		Certificate:            cert,
		IdentityPoolClient:     poolClient,
		ExternalSecret:         manifest,
	}

	r1, err := s.store.SaveKubePromStack(stack)
	if err != nil {
		return nil, err
	}

	r2, err := s.state.SaveKubePromStack(stack)
	if err != nil {
		return nil, err
	}

	err = s.report.ReportSaveKubePromStack(stack, []*store.Report{r1, r2})
	if err != nil {
		return nil, err
	}

	return stack, nil
}

// NewKubePrometheusStackService returns an initialised service
func NewKubePrometheusStackService(
	spinner spinner.Spinner,
	api client.KubePromStackAPI,
	store client.KubePromStackStore,
	state client.KubePromStackState,
	report client.KubePromStackReport,
	cert client.CertificateService,
	ident client.IdentityManagerService,
	manifest client.ManifestService,
	param client.ParameterService,
) client.KubePromStackService {
	return &kubePromStackService{
		spinner:  spinner,
		api:      api,
		store:    store,
		state:    state,
		report:   report,
		cert:     cert,
		ident:    ident,
		manifest: manifest,
		param:    param,
	}
}
