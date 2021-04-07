package run

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/config/constant"

	"github.com/oslokommune/okctl/pkg/helm/charts/promtail"

	"github.com/oslokommune/okctl/pkg/helm/charts/loki"

	"github.com/oslokommune/okctl/pkg/helm/charts/kubepromstack"

	"github.com/oslokommune/okctl/pkg/helm/charts/blockstorage"

	"github.com/oslokommune/okctl/pkg/helm/charts/autoscaler"

	"github.com/oslokommune/okctl/pkg/helm/charts/awslbc"

	"github.com/oslokommune/okctl/pkg/helm/charts/argocd"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/helm"
)

type helmRun struct {
	helm            helm.Helmer
	kubeConfigStore api.KubeConfigStore
}

func (r *helmRun) CreateHelmRelease(opts api.CreateHelmReleaseOpts) (*api.Helm, error) {
	return r.createHelmChart(opts.ID, &helm.Chart{
		RepositoryName: opts.RepositoryName,
		RepositoryURL:  opts.RepositoryURL,
		ReleaseName:    opts.ReleaseName,
		Version:        opts.Version,
		Chart:          opts.Chart,
		Namespace:      opts.Namespace,
		Timeout:        constant.DefaultChartApplyTimeout,
		Values:         opts.Values,
	})
}

func (r *helmRun) DeleteHelmRelease(opts api.DeleteHelmReleaseOpts) error {
	kubeConf, err := r.kubeConfigStore.GetKubeConfig()
	if err != nil {
		return fmt.Errorf("getting kubeconfig: %w", err)
	}

	err = r.helm.Delete(kubeConf.Path, &helm.DeleteConfig{
		ReleaseName: opts.ReleaseName,
		Namespace:   opts.Namespace,
		Timeout:     constant.DefaultChartRemoveTimeout,
	})
	if err != nil {
		return fmt.Errorf("removing chart: %w", err)
	}

	return nil
}

func (r *helmRun) CreatePromtailHelmChart(opts api.CreatePromtailHelmChartOpts) (*api.Helm, error) {
	chart := promtail.New(promtail.NewDefaultValues())

	return r.createHelmChart(opts.ID, chart)
}

func (r *helmRun) CreateLokiHelmChart(opts api.CreateLokiHelmChartOpts) (*api.Helm, error) {
	chart := loki.New(loki.NewDefaultValues())

	return r.createHelmChart(opts.ID, chart)
}

func (r *helmRun) CreateAutoscalerHelmChart(opts api.CreateAutoscalerHelmChartOpts) (*api.Helm, error) {
	chart := autoscaler.New(autoscaler.NewDefaultValues(opts.ID.Region, opts.ID.ClusterName, "autoscaler"))

	return r.createHelmChart(opts.ID, chart)
}

func (r *helmRun) CreateBlockstorageHelmChart(opts api.CreateBlockstorageHelmChartOpts) (*api.Helm, error) {
	chart := blockstorage.New(blockstorage.NewDefaultValues(opts.ID.Region, opts.ID.ClusterName, "blockstorage"))

	return r.createHelmChart(opts.ID, chart)
}

func (r *helmRun) CreateArgoCD(opts api.CreateArgoCDOpts) (*api.Helm, error) {
	chart := argocd.New(argocd.NewDefaultValues(argocd.ValuesOpts{
		URL:                  fmt.Sprintf("https://%s", opts.ArgoDomain),
		HostName:             opts.ArgoDomain,
		Region:               opts.ID.Region,
		CertificateARN:       opts.ArgoCertificateARN,
		ClientID:             opts.ClientID,
		Organisation:         opts.GithubOrganisation,
		AuthDomain:           opts.AuthDomain,
		UserPoolID:           opts.UserPoolID,
		RepoURL:              opts.GithubRepoURL,
		RepoName:             opts.GithubRepoName,
		PrivateKeySecretName: opts.PrivateKeyName,
		PrivateKeySecretKey:  opts.PrivateKeyKey,
	}))

	return r.createHelmChart(opts.ID, chart)
}

func (r *helmRun) CreateAWSLoadBalancerControllerHelmChart(opts api.CreateAWSLoadBalancerControllerHelmChartOpts) (*api.Helm, error) {
	chart := awslbc.New(awslbc.NewDefaultValues(opts.ID.ClusterName, opts.VpcID, opts.ID.Region))

	return r.createHelmChart(opts.ID, chart)
}

func (r *helmRun) createHelmChart(id api.ID, chart *helm.Chart) (*api.Helm, error) {
	err := r.helm.RepoAdd(chart.RepositoryName, chart.RepositoryURL)
	if err != nil {
		return nil, fmt.Errorf("adding repository: %w", err)
	}

	err = r.helm.RepoUpdate()
	if err != nil {
		return nil, fmt.Errorf("updating repository: %w", err)
	}

	cfg, err := chart.InstallConfig()
	if err != nil {
		return nil, fmt.Errorf("creating install config: %w", err)
	}

	kubeConf, err := r.kubeConfigStore.GetKubeConfig()
	if err != nil {
		return nil, fmt.Errorf("getting kubeconfig: %w", err)
	}

	release, err := r.helm.Install(kubeConf.Path, cfg)
	if err != nil {
		return nil, fmt.Errorf("installing chart: %w", err)
	}

	return &api.Helm{
		ID:      id,
		Release: release,
		Chart:   chart,
	}, nil
}

func (r *helmRun) CreateKubePromStack(opts api.CreateKubePrometheusStackOpts) (*api.Helm, error) {
	chart := kubepromstack.New(constant.DefaultChartApplyTimeout, &kubepromstack.Values{
		GrafanaServiceAccountName:          opts.GrafanaCloudWatchServiceAccountName,
		GrafanaCertificateARN:              opts.CertificateARN,
		GrafanaHostname:                    opts.Hostname,
		AuthHostname:                       opts.AuthHostname,
		ClientID:                           opts.ClientID,
		SecretsConfigName:                  opts.SecretsConfigName,
		SecretsGrafanaCookieSecretKey:      opts.SecretsCookieSecretKey,
		SecretsGrafanaOauthClientSecretKey: opts.SecretsClientSecretKey,
		SecretsGrafanaAdminUserKey:         opts.SecretsAdminUserKey,
		SecretsGrafanaAdminPassKey:         opts.SecretsAdminPassKey,
	})

	return r.createHelmChart(opts.ID, chart)
}

// NewHelmRun returns an initialised helm runner
func NewHelmRun(helm helm.Helmer, kubeConfigStore api.KubeConfigStore) api.HelmRun {
	return &helmRun{
		helm:            helm,
		kubeConfigStore: kubeConfigStore,
	}
}
