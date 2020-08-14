// Package helm implements a helm client, this is based on code from:
// - https://github.com/PrasadG193/helm-clientgo-example
// - https://github.com/helm/helm
package helm

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/strvals"
	"k8s.io/cli-runtime/pkg/genericclioptions"

	"github.com/gofrs/flock"
	"github.com/pkg/errors"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/cli/values"
	"helm.sh/helm/v3/pkg/downloader"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/repo"
)

const (
	// DefaultHelmDriver is set to secrets, which is the default
	// for Helm 3: https://helm.sh/docs/topics/advanced/#storage-backends
	DefaultHelmDriver = "secrets"
	// DefaultHelmLockExt is the extension used to create a file lock
	DefaultHelmLockExt = ".lock"
)

type RepoAdder interface {
	RepoAdd(repoName, url string) error
}

type RepoUpdater interface {
	RepoUpdate() error
}

type ChartInstaller interface {
	Install(releaseName, repoName, chartName string, args map[string]string) (*release.Release, error)
}

type Helmer interface {
	RepoAdder
	RepoUpdater
	ChartInstaller
}

type Helm struct {
	restClient genericclioptions.RESTClientGetter
	config     *Config
	fs         *afero.Afero
}

// Config lists all configuration variables that must be set
type Config struct {
	Namespace  string
	KubeConfig string

	HelmPluginsDirectory string
	HelmRegistryConfig   string
	HelmRepositoryConfig string
	HelmRepositoryCache  string
	HelmBaseDir          string

	Debug bool
}

// Envs returns the config as a helm compatible
// set of env vars
func (c *Config) Envs() map[string]string {
	return map[string]string{
		"KUBECONFIG": c.KubeConfig,

		"HELM_CACHE_HOME":        c.HelmBaseDir,
		"HELM_CONFIG_HOME":       c.HelmBaseDir,
		"HELM_CONFIG_DATA_HOME":  c.HelmBaseDir,
		"HELM_PLUGINS":           c.HelmPluginsDirectory,
		"HELM_REGISTRY_CONFIG":   c.HelmRegistryConfig,
		"HELM_REPOSITORY_CONFIG": c.HelmRepositoryConfig,
		"HELM_REPOSITORY_CACHE":  c.HelmRepositoryCache,
		"HELM_NAMESPACE":         c.Namespace,

		"HELM_DEBUG": fmt.Sprintf("%t", c.Debug),
	}
}

// New initialises a new Helm operator
func New(config *Config, fs *afero.Afero) *Helm {
	return &Helm{
		config: config,
		restClient: &genericclioptions.ConfigFlags{
			KubeConfig: &config.KubeConfig,
			Namespace:  &config.Namespace,
		},
		fs: fs,
	}
}

// RepoAdd adds repo with given name and url
func (h *Helm) RepoAdd(name, url string) error {
	restoreFn, err := EstablishEnv(h.config.Envs())
	if err != nil {
		return err
	}

	defer func() {
		err = restoreFn()
	}()

	settings := cli.New()
	repoFile := settings.RepositoryConfig

	err = h.fs.MkdirAll(h.config.HelmBaseDir, 0744)
	if err != nil {
		return err
	}

	unlockFn, err := Lock(repoFile)
	if err != nil {
		return err
	}

	defer func() {
		err = unlockFn()
	}()

	b, err := h.fs.ReadFile(repoFile)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	var f repo.File

	err = yaml.Unmarshal(b, &f)
	if err != nil {
		return err
	}

	// We already have this repository, so we are done
	if f.Has(name) {
		return nil
	}

	c := repo.Entry{
		Name: name,
		URL:  url,
	}

	r, err := repo.NewChartRepository(&c, getter.All(settings))
	if err != nil {
		log.Fatal(err)
	}

	_, err = r.DownloadIndexFile()
	if err != nil {
		return err
	}

	f.Update(&c)

	err = f.WriteFile(repoFile, 0644)
	if err != nil {
		return err
	}

	return nil
}

// RepoUpdate updates charts for all helm repos
func (h *Helm) RepoUpdate() error {
	restoreFn, err := EstablishEnv(h.config.Envs())
	if err != nil {
		return err
	}

	defer func() {
		err = restoreFn()
	}()

	settings := cli.New()
	repoFile := settings.RepositoryConfig

	exists, err := h.fs.Exists(repoFile)
	if err != nil {
		return err
	}

	if !exists {
		return nil
	}

	f, err := repo.LoadFile(repoFile)
	if err != nil {
		return err
	}

	providers := getter.Providers{
		{
			Schemes: []string{"https"},
			New:     getter.NewHTTPGetter,
		},
	}

	for _, cfg := range f.Repositories {
		r, err := repo.NewChartRepository(cfg, providers)
		if err != nil {
			return err
		}

		_, err = r.DownloadIndexFile()
		if err != nil {
			return err
		}
	}

	return nil
}

// InstallChart
func (h *Helm) Install(name, repo, chart string, args map[string]string) (*release.Release, error) {
	restoreFn, err := EstablishEnv(h.config.Envs())
	if err != nil {
		return nil, err
	}

	defer func() {
		err = restoreFn()
	}()

	settings := cli.New()

	actionConfig := new(action.Configuration)
	err = actionConfig.Init(h.restClient, settings.Namespace(), DefaultHelmDriver, debug)
	if err != nil {
		return nil, err
	}
	client := action.NewInstall(actionConfig)

	if client.Version == "" && client.Devel {
		client.Version = ">0.0.0-0"
	}
	//name, chart, err := client.NameAndChart(args)
	client.ReleaseName = name
	cp, err := client.ChartPathOptions.LocateChart(fmt.Sprintf("%s/%s", repo, chart), settings)
	if err != nil {
		return nil, err
	}

	p := getter.All(settings)
	valueOpts := &values.Options{}

	v, err := valueOpts.MergeValues(p)
	if err != nil {
		return nil, err
	}

	err = strvals.ParseInto(args["set"], v)
	if err != nil {
		return nil, err
	}

	// Check chart dependencies to make sure all are present in /charts
	chartRequested, err := loader.Load(cp)
	if err != nil {
		return nil, err
	}

	validInstallableChart, err := isChartInstallable(chartRequested)
	if err != nil {
		return nil, err
	}

	if !validInstallableChart {
		return nil, fmt.Errorf("invalid chart")
	}

	if req := chartRequested.Metadata.Dependencies; req != nil {
		// If CheckDependencies returns an error, we have unfulfilled dependencies.
		// As of Helm 2.4.0, this is treated as a stopping condition:
		// https://github.com/helm/helm/issues/2209
		if err := action.CheckDependencies(chartRequested, req); err != nil {
			if client.DependencyUpdate {
				man := &downloader.Manager{
					Out:              os.Stderr,
					ChartPath:        cp,
					Keyring:          client.ChartPathOptions.Keyring,
					SkipUpdate:       false,
					Getters:          p,
					RepositoryConfig: settings.RepositoryConfig,
					RepositoryCache:  settings.RepositoryCache,
				}
				if err := man.Update(); err != nil {
					return nil, err
				}
			} else {
				return nil, err
			}
		}
	}

	client.Namespace = settings.Namespace()
	r, err := client.Run(chartRequested, v)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func isChartInstallable(ch *chart.Chart) (bool, error) {
	switch ch.Metadata.Type {
	case "", "application":
		return true, nil
	}
	return false, errors.Errorf("%s charts are not installable", ch.Metadata.Type)
}

func debug(format string, v ...interface{}) {
	format = fmt.Sprintf("[debug] %s\n", format)
	_ = log.Output(2, fmt.Sprintf(format, v...))
}

// UnlockFn can be deferred in the calling function to
// unlock the file
type UnlockFn func() error

// Lock a file to ensure no concurrent access
func Lock(file string) (UnlockFn, error) {
	lockFile := strings.Replace(file, filepath.Ext(file), DefaultHelmLockExt, 1)
	lock := flock.New(lockFile)

	lockCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	locked, err := lock.TryLockContext(lockCtx, 1*time.Second)
	if err != nil {
		return nil, err
	}

	if locked {
		return func() error {
			return lock.Unlock()
		}, nil
	}

	return nil, fmt.Errorf("failed to create lock: %s", lockFile)
}

// RestoreEnvFn can be deferred in the calling function
// and will return the environment to its original state
type RestoreEnvFn func() error

// EstablishEnv provides functionality for setting a safe environment,
// this is required, because helm for some reason, loves fetching
// everything from environment variables
func EstablishEnv(envs map[string]string) (RestoreEnvFn, error) {
	originalEnvVars := os.Environ()
	os.Clearenv()

	for key, val := range envs {
		err := os.Setenv(key, val)
		if err != nil {
			return nil, err
		}
	}

	return func() error {
		originalEnvVars := originalEnvVars

		os.Clearenv()

		for _, envVar := range originalEnvVars {
			e := strings.SplitN(envVar, "=", 2)

			var key, val string

			switch len(e) {
			case 0:
				continue
			case 1:
				key = e[0]
				val = ""
			case 2:
				key = e[0]
				val = e[1]
			}

			err := os.Setenv(key, val)
			if err != nil {
				return err
			}
		}

		return nil
	}, nil
}
