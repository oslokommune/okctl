package core

import (
	"context"
	"fmt"

	"github.com/oslokommune/okctl/pkg/client/store"

	"github.com/google/uuid"
	"github.com/oslokommune/okctl/pkg/api"

	"github.com/oslokommune/okctl/pkg/client"
)

type argoCDService struct {
	api    client.ArgoCDAPI
	store  client.ArgoCDStore
	report client.ArgoCDReport
	state  client.ArgoCDState

	gh       client.GithubService
	cert     client.CertificateService
	manifest client.ManifestService
	param    client.ParameterService
}

// nolint: funlen
func (s *argoCDService) CreateArgoCD(ctx context.Context, opts client.CreateArgoCDOpts) (*client.ArgoCD, error) {
	cert, err := s.cert.CreateCertificate(ctx, api.CreateCertificateOpts{
		ID:           opts.ID,
		FQDN:         fmt.Sprintf("argocd.%s", opts.FQDN),
		Domain:       fmt.Sprintf("argocd.%s", opts.Domain),
		HostedZoneID: opts.HostedZoneID,
	})
	if err != nil {
		return nil, err
	}

	app, err := s.gh.CreateGithubOauthApp(ctx, client.CreateGithubOauthAppOpts{
		ID:           opts.ID,
		Organisation: opts.GithubOrganisation,
		Name:         fmt.Sprintf("okctl-argocd-%s", opts.ID.ClusterName),
		SiteURL:      fmt.Sprintf("https://%s", cert.Domain),
		CallbackURL:  fmt.Sprintf("https://%s/api/dex/callback", cert.Domain),
	})
	if err != nil {
		return nil, err
	}

	secretKey, err := s.param.CreateSecret(ctx, api.CreateSecretOpts{
		ID:     opts.ID,
		Name:   "argocd/secret_key",
		Secret: uuid.New().String(),
	})
	if err != nil {
		return nil, err
	}

	privateKeyName := "argocd-privatekey"
	privateKeyDataName := "ssh-private-key"

	manifest, err := s.manifest.CreateExternalSecret(ctx, client.CreateExternalSecretOpts{
		ID: opts.ID,
		Manifests: []api.Manifest{
			{
				Name:      privateKeyName,
				Namespace: "argocd",
				Annotations: map[string]string{
					"meta.helm.sh/release-name":      "argocd",
					"meta.helm.sh/release-namespace": "argocd",
				},
				Labels: map[string]string{
					"app.kubernetes.io/managed-by": "Helm",
				},
				Data: []api.Data{
					{
						Name: privateKeyDataName,
						Key:  opts.Repository.DeployKey.PrivateKeySecret.Path,
					},
				},
			},
			{
				Name:      "argocd-secret",
				Namespace: "argocd",
				Annotations: map[string]string{
					"meta.helm.sh/release-name":      "argocd",
					"meta.helm.sh/release-namespace": "argocd",
				},
				Labels: map[string]string{
					"app.kubernetes.io/managed-by": "Helm",
				},
				Data: []api.Data{
					{
						Name: "dex.github.clientSecret",
						Key:  app.ClientSecret.Path,
					},
					{
						Name: "server.secretkey",
						Key:  secretKey.Path,
					},
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	chartOpts := api.CreateArgoCDOpts{
		ID:                  opts.ID,
		ArgoDomain:          cert.Domain,
		ArgoCertificateARN:  cert.CertificateARN,
		GithubOrganisation:  opts.GithubOrganisation,
		GithubTeam:          app.Team.Name,
		GithubRepoURL:       opts.Repository.GitURL,
		GithubRepoName:      opts.Repository.Repository,
		GithubOauthClientID: app.ClientID,
		PrivateKeyName:      privateKeyName,
		PrivateKeyKey:       privateKeyDataName,
	}

	argo, err := s.api.CreateArgoCD(chartOpts)
	if err != nil {
		return nil, err
	}

	argo.ID = opts.ID
	argo.ArgoDomain = cert.Domain
	argo.ArgoURL = fmt.Sprintf("https://%s", cert.Domain)
	argo.Certificate = cert
	argo.GithubOauthApp = app
	argo.ExternalSecret = &api.ExternalSecretsKube{
		ID:        manifest.ID,
		Manifests: manifest.Manifests,
	}
	argo.SecretKey = secretKey

	r1, err := s.store.SaveArgoCD(argo)
	if err != nil {
		return nil, err
	}

	r2, err := s.state.SaveArgoCD(argo)
	if err != nil {
		return nil, err
	}

	err = s.report.CreateArgoCD(argo, []*store.Report{r1, r2})
	if err != nil {
		return nil, err
	}

	return argo, nil
}

// NewArgoCDService returns an initialised service
func NewArgoCDService(
	gh client.GithubService,
	cert client.CertificateService,
	manifest client.ManifestService,
	param client.ParameterService,
	api client.ArgoCDAPI,
	store client.ArgoCDStore,
	report client.ArgoCDReport,
	state client.ArgoCDState,
) client.ArgoCDService {
	return &argoCDService{
		api:      api,
		store:    store,
		report:   report,
		state:    state,
		gh:       gh,
		cert:     cert,
		manifest: manifest,
		param:    param,
	}
}
