package core

import (
	"context"
	"fmt"

	"github.com/oslokommune/okctl/pkg/config/constant"

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

	identity client.IdentityManagerService
	cert     client.CertificateService
	manifest client.ManifestService
	param    client.ParameterService
}

// nolint: gosec
const (
	argoClientSecretName = "argocd/client_secret"
	argoSecretKeyName    = "argocd/secret_key"
)

func (s *argoCDService) DeleteArgoCD(ctx context.Context, opts client.DeleteArgoCDOpts) error {
	info := s.state.GetArgoCD(opts.ID)

	err := s.manifest.DeleteNamespace(ctx, api.DeleteNamespaceOpts{
		ID:        opts.ID,
		Namespace: constant.DefaultArgoCDNamespace,
	})
	if err != nil {
		return err
	}

	err = s.cert.DeleteCertificate(ctx, client.DeleteCertificateOpts{
		ID:     opts.ID,
		Domain: info.ArgoDomain,
	})
	if err != nil {
		return err
	}

	err = s.identity.DeleteIdentityPoolClient(ctx, api.DeleteIdentityPoolClientOpts{
		ID:      opts.ID,
		Purpose: "argocd",
	})
	if err != nil {
		return err
	}

	for _, secret := range []string{argoSecretKeyName, argoClientSecretName} {
		err = s.param.DeleteSecret(ctx, client.DeleteSecretOpts{
			ID:   opts.ID,
			Name: secret,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// nolint: funlen
func (s *argoCDService) CreateArgoCD(ctx context.Context, opts client.CreateArgoCDOpts) (*client.ArgoCD, error) {
	cert, err := s.cert.CreateCertificate(ctx, client.CreateCertificateOpts{
		ID:           opts.ID,
		FQDN:         fmt.Sprintf("argocd.%s", opts.FQDN),
		Domain:       fmt.Sprintf("argocd.%s", opts.Domain),
		HostedZoneID: opts.HostedZoneID,
	})
	if err != nil {
		return nil, err
	}

	identityClient, err := s.identity.CreateIdentityPoolClient(ctx, api.CreateIdentityPoolClientOpts{
		ID:          opts.ID,
		UserPoolID:  opts.UserPoolID,
		Purpose:     "argocd",
		CallbackURL: fmt.Sprintf("https://%s/api/dex/callback", cert.Domain),
	})
	if err != nil {
		return nil, err
	}

	_, err = s.manifest.CreateNamespace(ctx, api.CreateNamespaceOpts{
		ID:        opts.ID,
		Namespace: constant.DefaultArgoCDNamespace,
	})
	if err != nil {
		return nil, err
	}

	clientSecret, err := s.param.CreateSecret(ctx, client.CreateSecretOpts{
		ID:     opts.ID,
		Name:   argoClientSecretName,
		Secret: identityClient.ClientSecret,
	})
	if err != nil {
		return nil, err
	}

	secretKey, err := s.param.CreateSecret(ctx, client.CreateSecretOpts{
		ID:     opts.ID,
		Name:   argoSecretKeyName,
		Secret: uuid.New().String(),
	})
	if err != nil {
		return nil, err
	}

	privateKeyName := "argocd-privatekey"
	privateKeyDataName := "ssh-private-key"

	priv, err := s.manifest.CreateExternalSecret(ctx, client.CreateExternalSecretOpts{
		ID:        opts.ID,
		Namespace: constant.DefaultArgoCDNamespace,
		Manifest: api.Manifest{
			Name:      privateKeyName,
			Namespace: constant.DefaultArgoCDNamespace,
			Backend:   api.BackendTypeParameterStore,
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
	})
	if err != nil {
		return nil, err
	}

	sec, err := s.manifest.CreateExternalSecret(ctx, client.CreateExternalSecretOpts{
		ID:        opts.ID,
		Namespace: constant.DefaultArgoCDNamespace,
		Manifest: api.Manifest{
			Name:      "argocd-secret",
			Namespace: "argocd",
			Backend:   api.BackendTypeParameterStore,
			Annotations: map[string]string{
				"meta.helm.sh/release-name":      "argocd",
				"meta.helm.sh/release-namespace": "argocd",
			},
			Labels: map[string]string{
				"app.kubernetes.io/managed-by": "Helm",
			},
			Data: []api.Data{
				{
					Name: "dex.cognito.clientSecret",
					Key:  clientSecret.Path,
				},
				{
					Name: "server.secretkey",
					Key:  secretKey.Path,
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	chartOpts := api.CreateArgoCDOpts{
		ID:                 opts.ID,
		ArgoDomain:         cert.Domain,
		ArgoCertificateARN: cert.ARN,
		GithubOrganisation: opts.GithubOrganisation,
		GithubRepoURL:      opts.Repository.GitURL,
		GithubRepoName:     opts.Repository.Repository,
		ClientID:           identityClient.ClientID,
		AuthDomain:         opts.AuthDomain,
		UserPoolID:         opts.UserPoolID,
		PrivateKeyName:     privateKeyName,
		PrivateKeyKey:      privateKeyDataName,
	}

	argo, err := s.api.CreateArgoCD(chartOpts)
	if err != nil {
		return nil, err
	}

	argo.ID = opts.ID
	argo.ArgoDomain = cert.Domain
	argo.ArgoURL = fmt.Sprintf("https://%s", cert.Domain)
	argo.Certificate = cert
	argo.IdentityClient = identityClient
	argo.ClientSecret = clientSecret
	argo.PrivateKey = priv
	argo.Secret = sec
	argo.SecretKey = secretKey
	argo.AuthDomain = opts.AuthDomain

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
	identity client.IdentityManagerService,
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
		identity: identity,
		cert:     cert,
		manifest: manifest,
		param:    param,
	}
}
