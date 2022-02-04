package core

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"path"
	"time"

	"github.com/oslokommune/okctl/pkg/binaries"
	"github.com/oslokommune/okctl/pkg/clients/kubectl/binary"
	"github.com/oslokommune/okctl/pkg/credentials"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/scaffold"
	"github.com/spf13/afero"

	"github.com/asdine/storm/v3"

	"github.com/oslokommune/okctl/pkg/helm/charts/argocd"

	"github.com/oslokommune/okctl/pkg/config/constant"

	"github.com/google/uuid"
	"github.com/oslokommune/okctl/pkg/api"

	"github.com/oslokommune/okctl/pkg/client"
)

type argoCDService struct {
	state               client.ArgoCDState
	identity            client.IdentityManagerService
	cert                client.CertificateService
	manifest            client.ManifestService
	param               client.ParameterService
	helm                client.HelmService
	fs                  *afero.Afero
	absoluteRepoDir     string
	binaryProvider      binaries.Provider
	credentialsProvider credentials.Provider
}

// nolint: gosec
const (
	argoClientSecretName = "argocd/client_secret"
	argoSecretKeyName    = "argocd/secret_key"
	argoPurpose          = "argocd"
	argoPrivateKeyName   = "argocd-privatekey"
	argoSecretName       = "argocd-secret"
	// argoChartTimeout does not work right now. See https://trello.com/c/zrS1xDXz for details
	argoChartTimeout                     = 15 * time.Minute
	argoRepositoryTypeGit                = "git"
	defaultArgoCDApplicationManifestName = "applications"
)

// nolint: funlen
func (s *argoCDService) DeleteArgoCD(ctx context.Context, opts client.DeleteArgoCDOpts) error {
	cd, err := s.state.GetArgoCD()
	if err != nil {
		if errors.Is(err, storm.ErrNotFound) {
			return nil
		}

		return err
	}

	err = s.helm.DeleteHelmRelease(ctx, client.DeleteHelmReleaseOpts{
		ID:          opts.ID,
		ReleaseName: argocd.ReleaseName,
		Namespace:   argocd.Namespace,
	})
	if err != nil {
		return err
	}

	err = s.manifest.DeleteNamespace(ctx, api.DeleteNamespaceOpts{
		ID:        opts.ID,
		Namespace: constant.DefaultArgoCDNamespace,
	})
	if err != nil {
		return err
	}

	for _, name := range []string{argoSecretName, argoPrivateKeyName} {
		err = s.manifest.DeleteExternalSecret(ctx, client.DeleteExternalSecretOpts{
			ID:   opts.ID,
			Name: name,
			Secrets: map[string]string{
				name: constant.DefaultArgoCDNamespace,
			},
		})
		if err != nil {
			return err
		}
	}

	err = s.cert.DeleteCertificate(ctx, client.DeleteCertificateOpts{
		ID:     opts.ID,
		Domain: cd.ArgoDomain,
	})
	if err != nil {
		return err
	}

	err = s.identity.DeleteIdentityPoolClient(ctx, client.DeleteIdentityPoolClientOpts{
		ID:      opts.ID,
		Purpose: argoPurpose,
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

	err = s.state.RemoveArgoCD()
	if err != nil {
		return err
	}

	return nil
}

// nolint: funlen
func (s *argoCDService) CreateArgoCD(ctx context.Context, opts client.CreateArgoCDOpts) (*client.ArgoCD, error) {
	clusterID := clusterMetaAsID(opts.ClusterManifest.Metadata)

	cert, err := s.cert.CreateCertificate(ctx, client.CreateCertificateOpts{
		ID:           clusterID,
		FQDN:         fmt.Sprintf("argocd.%s", opts.FQDN),
		Domain:       fmt.Sprintf("argocd.%s", opts.Domain),
		HostedZoneID: opts.HostedZoneID,
	})
	if err != nil {
		return nil, fmt.Errorf("creating certificate: %w", err)
	}

	identityClient, err := s.identity.CreateIdentityPoolClient(ctx, client.CreateIdentityPoolClientOpts{
		ID:          clusterID,
		UserPoolID:  opts.UserPoolID,
		Purpose:     argoPurpose,
		CallbackURL: fmt.Sprintf("https://%s/api/dex/callback", cert.Domain),
	})
	if err != nil {
		return nil, fmt.Errorf("creating IdentityPool client: %w", err)
	}

	_, err = s.manifest.CreateNamespace(ctx, api.CreateNamespaceOpts{
		ID:        clusterID,
		Namespace: constant.DefaultArgoCDNamespace,
	})
	if err != nil {
		return nil, fmt.Errorf("creating k8s namespace: %w", err)
	}

	clientSecret, err := s.param.CreateSecret(ctx, client.CreateSecretOpts{
		ID:     clusterID,
		Name:   argoClientSecretName,
		Secret: identityClient.ClientSecret,
	})
	if err != nil {
		return nil, fmt.Errorf("creating IdentityPool client secret: %w", err)
	}

	secretKey, err := s.param.CreateSecret(ctx, client.CreateSecretOpts{
		ID:     clusterID,
		Name:   argoSecretKeyName,
		Secret: uuid.New().String(),
	})
	if err != nil {
		return nil, fmt.Errorf("creating Argo secret key: %w", err)
	}

	privateKeyDataName := "ssh-private-key"

	priv, err := s.manifest.CreateExternalSecret(ctx, client.CreateExternalSecretOpts{
		ID:        clusterID,
		Name:      argoPrivateKeyName,
		Namespace: constant.DefaultArgoCDNamespace,
		Manifest: api.Manifest{
			Name:      argoPrivateKeyName,
			Namespace: constant.DefaultArgoCDNamespace,
			Backend:   api.BackendTypeParameterStore,
			Annotations: map[string]string{
				"meta.helm.sh/release-name":      "argocd",
				"meta.helm.sh/release-namespace": "argocd",
			},
			Labels: map[string]string{
				"app.kubernetes.io/managed-by":   "Helm",
				"argocd.argoproj.io/secret-type": "repository",
			},
			Data: []api.Data{
				{
					Name: privateKeyDataName,
					Key:  opts.Repository.DeployKey.PrivateKeySecret.Path,
				},
			},
			Template: api.ExternalSecretSpecTemplate{
				StringData: map[string]interface{}{
					"name":          opts.Repository.FullName,
					"url":           opts.Repository.GitURL,
					"type":          argoRepositoryTypeGit,
					"sshPrivateKey": fmt.Sprintf("%s data[\"%s\"] %s", "<%=", privateKeyDataName, "%>"),
				},
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("creating external secret for deploy key: %w", err)
	}

	sec, err := s.manifest.CreateExternalSecret(ctx, client.CreateExternalSecretOpts{
		ID:        clusterID,
		Name:      argoSecretName,
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
		return nil, fmt.Errorf("creating external secret for IdentityPool client: %w", err)
	}

	chart := argocd.New(argocd.NewDefaultValues(argocd.ValuesOpts{
		URL:                  fmt.Sprintf("https://%s", cert.Domain),
		HostName:             cert.Domain,
		Region:               clusterID.Region,
		CertificateARN:       cert.ARN,
		ClientID:             identityClient.ClientID,
		Organisation:         opts.ClusterManifest.Github.Organisation,
		AuthDomain:           opts.AuthDomain,
		UserPoolID:           opts.UserPoolID,
		RepoURL:              opts.Repository.GitURL,
		RepoName:             opts.Repository.Repository,
		PrivateKeySecretName: argoPrivateKeyName,
		PrivateKeySecretKey:  privateKeyDataName,
	}), argoChartTimeout)

	values, err := chart.ValuesYAML()
	if err != nil {
		return nil, fmt.Errorf("preparing chart values: %w", err)
	}

	a, err := s.helm.CreateHelmRelease(ctx, client.CreateHelmReleaseOpts{
		ID:             clusterID,
		RepositoryName: chart.RepositoryName,
		RepositoryURL:  chart.RepositoryURL,
		ReleaseName:    chart.ReleaseName,
		Version:        chart.Version,
		Chart:          chart.Chart,
		Namespace:      chart.Namespace,
		Values:         values,
	})
	if err != nil {
		return nil, fmt.Errorf("creating Helm release: %w", err)
	}

	argo := &client.ArgoCD{
		ID:             clusterID,
		ArgoDomain:     cert.Domain,
		ArgoURL:        fmt.Sprintf("https://%s", cert.Domain),
		AuthDomain:     opts.AuthDomain,
		Certificate:    cert,
		IdentityClient: identityClient,
		PrivateKey:     priv,
		Secret:         sec,
		ClientSecret:   clientSecret,
		SecretKey:      secretKey,
		Chart: &client.Helm{
			ID:      a.ID,
			Release: a.Release,
			Chart:   a.Chart,
		},
	}

	err = s.SetupApplicationsSync(ctx, opts.ClusterManifest)
	if err != nil {
		return nil, fmt.Errorf("setting up application sync: %w", err)
	}

	err = s.state.SaveArgoCD(argo)
	if err != nil {
		return nil, fmt.Errorf("storing state: %w", err)
	}

	return argo, nil
}

// SetupApplicationsSync applies an ArgoCD Application manifest which ensures syncing of an application folder
func (s *argoCDService) SetupApplicationsSync(_ context.Context, cluster v1alpha1.Cluster) error {
	relativeArgoCDManifestPath := path.Join(
		cluster.Github.OutputPath,
		cluster.Metadata.Name,
		constant.DefaultArgoCDClusterConfigDir,
		fmt.Sprintf("%s.yaml", defaultArgoCDApplicationManifestName),
	)

	absoluteArgoCDManifestPath := path.Join(s.absoluteRepoDir, relativeArgoCDManifestPath)

	originalManifest, err := scaffold.GenerateArgoCDApplicationManifest(scaffold.GenerateArgoCDApplicationManifestOpts{
		Name:          defaultArgoCDApplicationManifestName,
		Namespace:     argocd.Namespace,
		IACRepoURL:    cluster.Github.URL(),
		SourceSyncDir: path.Dir(relativeArgoCDManifestPath),
		Prune:         true,
	})
	if err != nil {
		return fmt.Errorf("generating ArgoCD application manifest: %w", err)
	}

	manifestCopy := bytes.Buffer{}
	manifest := io.TeeReader(originalManifest, &manifestCopy)

	err = s.fs.WriteReader(absoluteArgoCDManifestPath, manifest)
	if err != nil {
		return fmt.Errorf("writing manifest: %w", err)
	}

	kubectlClient := binary.New(s.fs, s.binaryProvider, s.credentialsProvider, cluster)

	err = kubectlClient.Apply(&manifestCopy)
	if err != nil {
		return fmt.Errorf("applying ArgoCD application manifest: %w", err)
	}

	return nil
}

// NewArgoCDServiceOpts defines required data for an ArgoCD service
type NewArgoCDServiceOpts struct {
	Fs                  *afero.Afero
	BinaryProvider      binaries.Provider
	CredentialsProvider credentials.Provider
	AbsoluteRepoDir     string
	Identity            client.IdentityManagerService
	Cert                client.CertificateService
	Manifest            client.ManifestService
	Param               client.ParameterService
	Helm                client.HelmService
	State               client.ArgoCDState
}

// NewArgoCDService returns an initialised service
func NewArgoCDService(opts NewArgoCDServiceOpts) client.ArgoCDService {
	return &argoCDService{
		fs:                  opts.Fs,
		absoluteRepoDir:     opts.AbsoluteRepoDir,
		binaryProvider:      opts.BinaryProvider,
		credentialsProvider: opts.CredentialsProvider,
		state:               opts.State,
		identity:            opts.Identity,
		cert:                opts.Cert,
		manifest:            opts.Manifest,
		param:               opts.Param,
		helm:                opts.Helm,
	}
}
