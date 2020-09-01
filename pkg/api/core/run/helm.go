package run

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/helm/charts/argocd"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/helm"
	"github.com/oslokommune/okctl/pkg/helm/charts/awsalbingresscontroller"
	"github.com/oslokommune/okctl/pkg/helm/charts/externalsecrets"
)

type helmRun struct {
	helm            helm.Helmer
	kubeConfigStore api.KubeConfigStore
}

func (r *helmRun) CreateArgoCD(opts api.CreateArgoCDOpts) (*api.Helm, error) {
	chart := argocd.New(argocd.NewDefaultValues(argocd.ValuesOpts{
		URL:                  fmt.Sprintf("https://%s", opts.ArgoDomain),
		HostName:             opts.ArgoDomain,
		CertificateARN:       opts.ArgoCertificateARN,
		ClientID:             opts.GithubOauthClientID,
		Organisation:         opts.GithubOrganisation,
		Team:                 opts.GithubTeam,
		RepoURL:              opts.GithubRepoURL,
		RepoName:             opts.GithubRepoName,
		PrivateKeySecretName: opts.ExternalSecretName,
		PrivateKeySecretKey:  opts.ExternalSecretKey,
	}))

	return r.createHelmChart(opts.ClusterName, opts.Repository, opts.Environment, chart)
}

func (r *helmRun) CreateAlbIngressControllerHelmChart(opts api.CreateAlbIngressControllerHelmChartOpts) (*api.Helm, error) {
	chart := awsalbingresscontroller.New(awsalbingresscontroller.NewDefaultValues(opts.ClusterName, opts.VpcID, opts.Region))

	return r.createHelmChart(opts.ClusterName, opts.Repository, opts.Environment, chart)
}

func (r *helmRun) CreateExternalSecretsHelmChart(opts api.CreateExternalSecretsHelmChartOpts) (*api.Helm, error) {
	chart := externalsecrets.ExternalSecrets(externalsecrets.DefaultExternalSecretsValues())

	return r.createHelmChart(opts.ClusterName, opts.Repository, opts.Environment, chart)
}

func (r *helmRun) createHelmChart(clusterName, repository, env string, chart *helm.Chart) (*api.Helm, error) {
	err := r.helm.RepoAdd(chart.RepositoryName, chart.RepositoryURL)
	if err != nil {
		return nil, fmt.Errorf("failed to add repository: %w", err)
	}

	err = r.helm.RepoUpdate()
	if err != nil {
		return nil, fmt.Errorf("failed to update repository: %w", err)
	}

	cfg, err := chart.InstallConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to create install config: %w", err)
	}

	kubeConf, err := r.kubeConfigStore.GetKubeConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get kubeconfig: %w", err)
	}

	release, err := r.helm.Install(kubeConf.Path, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to install chart: %w", err)
	}

	return &api.Helm{
		ClusterName: clusterName,
		Repository:  repository,
		Environment: env,
		Release:     release,
		Chart:       chart,
	}, nil
}

// NewHelmRun returns an initialised helm runner
func NewHelmRun(helm helm.Helmer, kubeConfigStore api.KubeConfigStore) api.HelmRun {
	return &helmRun{
		helm:            helm,
		kubeConfigStore: kubeConfigStore,
	}
}
