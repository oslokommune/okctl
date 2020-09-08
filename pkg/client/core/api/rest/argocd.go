package rest

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
)

// TargetHelmArgoCD matches the REST API route
const TargetHelmArgoCD = "helm/argocd/"

type argoCDAPI struct {
	client   *HTTPClient
	gh       client.GithubService // We use the Service, because this updates state
	cert     client.CertificateAPI
	manifest client.ManifestAPI
	param    client.ParameterAPI
}

// nolint: funlen
func (a *argoCDAPI) CreateArgoCD(opts client.CreateArgoCDOpts) (*client.ArgoCD, error) {
	ctx := context.Background()

	cert, err := a.cert.CreateCertificate(api.CreateCertificateOpts{
		ID:           opts.ID,
		FQDN:         fmt.Sprintf("argocd.%s", opts.FQDN),
		Domain:       fmt.Sprintf("argocd.%s", opts.Domain),
		HostedZoneID: opts.HostedZoneID,
	})
	if err != nil {
		return nil, err
	}

	app, err := a.gh.CreateGithubOauthApp(ctx, client.CreateGithubOauthAppOpts{
		ID:           opts.ID,
		Organisation: opts.GithubOrganisation,
		Name:         fmt.Sprintf("okctl-argocd-%s", opts.ID.ClusterName),
		SiteURL:      fmt.Sprintf("https://%s", cert.Domain),
		CallbackURL:  fmt.Sprintf("https://%s/api/dex/calback", cert.Domain),
	})
	if err != nil {
		return nil, err
	}

	secretKey, err := a.param.CreateSecret(api.CreateSecretOpts{
		ID:     opts.ID,
		Name:   "argocd/secret_key",
		Secret: uuid.New().String(),
	})
	if err != nil {
		return nil, err
	}

	privateKeyName := "argocd-privatekey"
	privateKeyDataName := "ssh-private-key"

	manifest, err := a.manifest.CreateExternalSecret(client.CreateExternalSecretOpts{
		ID: opts.ID,
		Manifests: []api.Manifest{
			{
				Name:      privateKeyName,
				Namespace: "argocd",
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

	chartOpts := &api.CreateArgoCDOpts{
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

	chart := &api.Helm{}

	err = a.client.DoPost(TargetHelmArgoCD, chartOpts, chart)
	if err != nil {
		return nil, err
	}

	return &client.ArgoCD{
		ID:             opts.ID,
		ArgoDomain:     cert.Domain,
		ArgoURL:        fmt.Sprintf("https://%s", cert.Domain),
		Certificate:    cert,
		GithubOauthApp: app,
		ExternalSecret: manifest,
		Chart:          chart,
		SecretKey:      secretKey,
	}, nil
}

// NewArgoCDAPI returns an initialised service
func NewArgoCDAPI(gh client.GithubService, param client.ParameterAPI, manifest client.ManifestAPI, cert client.CertificateAPI, client *HTTPClient) client.ArgoCDAPI { // nolint: lll
	return &argoCDAPI{
		client:   client,
		gh:       gh,
		cert:     cert,
		manifest: manifest,
		param:    param,
	}
}
