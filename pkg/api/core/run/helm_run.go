package run

import (
	"fmt"

	merrors "github.com/mishudark/errors"

	"github.com/oslokommune/okctl/pkg/config/constant"

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
	kubeConf, err := r.kubeConfigStore.GetKubeConfig(opts.ID.ClusterName)
	if err != nil {
		return fmt.Errorf(constant.GetKubeConfigError, err)
	}

	err = r.helm.Delete(kubeConf.Path, &helm.DeleteConfig{
		ReleaseName: opts.ReleaseName,
		Namespace:   opts.Namespace,
		Timeout:     constant.DefaultChartRemoveTimeout,
	})
	if err != nil {
		return fmt.Errorf(constant.RemoveChartError, err)
	}

	return nil
}

func (r *helmRun) GetHelmRelease(opts api.GetHelmReleaseOpts) (*api.Helm, error) {
	kubeConf, err := r.kubeConfigStore.GetKubeConfig(opts.ClusterID.ClusterName)
	if err != nil {
		return nil, fmt.Errorf(constant.GetKubeConfigError, err)
	}

	release, err := r.helm.Find(kubeConf.Path, &helm.FindConfig{
		ReleaseName: opts.ReleaseName,
		Namespace:   opts.Namespace,
		Timeout:     constant.DefaultChartFindTimeout,
	})
	if err != nil {
		return nil, merrors.E(err, "finding release")
	}

	return &api.Helm{
		ID:      opts.ClusterID,
		Release: release,
		Chart: &helm.Chart{
			RepositoryName: release.Chart.Name(),
			RepositoryURL:  "n/a",
			ReleaseName:    release.Name,
			Version:        release.Chart.Metadata.Version,
			Chart:          "n/a",
			Namespace:      release.Namespace,
			Values:         release.Chart.Values,
		},
	}, nil
}

func (r *helmRun) createHelmChart(id api.ID, chart *helm.Chart) (*api.Helm, error) {
	err := r.helm.RepoAdd(chart.RepositoryName, chart.RepositoryURL)
	if err != nil {
		return nil, fmt.Errorf(constant.AddRepositoryError, err)
	}

	err = r.helm.RepoUpdate()
	if err != nil {
		return nil, fmt.Errorf(constant.UpdateRepositoryError, err)
	}

	cfg, err := chart.InstallConfig()
	if err != nil {
		return nil, fmt.Errorf(constant.CreateInstallConfigError, err)
	}

	kubeConf, err := r.kubeConfigStore.GetKubeConfig(id.ClusterName)
	if err != nil {
		return nil, fmt.Errorf(constant.GetKubeConfigError, err)
	}

	release, err := r.helm.Install(kubeConf.Path, cfg)
	if err != nil {
		return nil, fmt.Errorf(constant.InstallChartError, err)
	}

	return &api.Helm{
		ID:      id,
		Release: release,
		Chart:   chart,
	}, nil
}

// NewHelmRun returns an initialised helm runner
func NewHelmRun(helm helm.Helmer, kubeConfigStore api.KubeConfigStore) api.HelmRun {
	return &helmRun{
		helm:            helm,
		kubeConfigStore: kubeConfigStore,
	}
}
