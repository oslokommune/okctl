package run

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/helm"
	"github.com/oslokommune/okctl/pkg/helm/charts/externalsecrets"
)

type helmRun struct {
	helm            helm.Helmer
	kubeConfigStore api.KubeConfigStore
}

func (r *helmRun) CreateExternalSecretsHelmChart(opts api.CreateExternalSecretsHelmChartOpts) (*api.Helm, error) {
	chart := externalsecrets.ExternalSecrets(externalsecrets.DefaultExternalSecretsValues())

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
		Repository:  opts.Repository,
		Environment: opts.Environment,
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
