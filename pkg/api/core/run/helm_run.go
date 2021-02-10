package run

import (
	"fmt"
	"github.com/oslokommune/okctl/pkg/helm/charts/kube_prometheus_stack"

	"github.com/oslokommune/okctl/pkg/helm/charts/blockstorage"

	"github.com/oslokommune/okctl/pkg/helm/charts/autoscaler"

	"github.com/oslokommune/okctl/pkg/helm/charts/awslbc"

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

func (r *helmRun) CreateAlbIngressControllerHelmChart(opts api.CreateAlbIngressControllerHelmChartOpts) (*api.Helm, error) {
	chart := awsalbingresscontroller.New(awsalbingresscontroller.NewDefaultValues(opts.ID.ClusterName, opts.VpcID, opts.ID.Region))

	return r.createHelmChart(opts.ID, chart)
}

func (r *helmRun) CreateAWSLoadBalancerControllerHelmChart(opts api.CreateAWSLoadBalancerControllerHelmChartOpts) (*api.Helm, error) {
	chart := awslbc.New(awslbc.NewDefaultValues(opts.ID.ClusterName, opts.VpcID, opts.ID.Region))

	return r.createHelmChart(opts.ID, chart)
}

func (r *helmRun) CreateExternalSecretsHelmChart(opts api.CreateExternalSecretsHelmChartOpts) (*api.Helm, error) {
	chart := externalsecrets.ExternalSecrets(externalsecrets.DefaultExternalSecretsValues())

	return r.createHelmChart(opts.ID, chart)
}

func (r *helmRun) createHelmChart(id api.ID, chart *helm.Chart) (*api.Helm, error) {
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
		ID:      id,
		Release: release,
		Chart:   chart,
	}, nil
}

func (r *helmRun) CreateKubePrometheusStack(opts api.CreateKubePrometheusStackOpts) (*api.Helm, error) {
	chart := kube_prometheus_stack.KubePrometheusStack(kube_prometheus_stack.DefaultKubePrometheusStackValues())

	return r.createHelmChart(opts.ID, chart)
}

// NewHelmRun returns an initialised helm runner
func NewHelmRun(helm helm.Helmer, kubeConfigStore api.KubeConfigStore) api.HelmRun {
	return &helmRun{
		helm:            helm,
		kubeConfigStore: kubeConfigStore,
	}
}
