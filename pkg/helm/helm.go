// Package helm implements a helm client, this is based on code from:
// - https://github.com/PrasadG193/helm-clientgo-example
// - https://github.com/helm/helm
package helm

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/oslokommune/okctl/pkg/credentials/aws"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
	"helm.sh/helm/v3/pkg/release"
	"k8s.io/cli-runtime/pkg/genericclioptions"

	"github.com/gofrs/flock"
	"helm.sh/helm/v3/pkg/action"
	chartPkg "helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	valuesPkg "helm.sh/helm/v3/pkg/cli/values"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/repo"
)

const (
	// DefaultHelmDriver is set to secrets, which is the default
	// for Helm 3: https://helm.sh/docs/topics/advanced/#storage-backends
	DefaultHelmDriver = "secrets"
	// DefaultHelmLockExt is the extension used to create a file lock
	DefaultHelmLockExt = ".lock"
	// DefaultHelmLockTimeout is the default timeout in seconds
	DefaultHelmLockTimeout = 120 * time.Second
)

// RepoAdder defines the interface for adding a helm repository
type RepoAdder interface {
	RepoAdd(repoName, url string) error
}

// RepoUpdater defines the interface for updating the helm repositories
type RepoUpdater interface {
	RepoUpdate() error
}

// Installer defines the interface for installing a helm chart
type Installer interface {
	Install(kubeConfigPath string, cfg *InstallConfig) (*release.Release, error)
}

// Helmer defines all available helm operations
type Helmer interface {
	RepoAdder
	RepoUpdater
	Installer
}

// Helm stores all state required for running helm tasks
type Helm struct {
	providers getter.Providers
	config    *Config
	auth      aws.Authenticator
	fs        *afero.Afero
}

// Ensure that Helm implements the Helmer interface
var _ Helmer = &Helm{}

// Config lists all configuration variables that must be set
// to configure the environment for helm correctly
type Config struct {
	// HomeDir allows us to modify where $HOME/.kube
	// ends up
	HomeDir string
	Path    string

	HelmPluginsDirectory string
	HelmRegistryConfig   string
	HelmRepositoryConfig string
	HelmRepositoryCache  string
	HelmBaseDir          string

	Debug       bool
	DebugOutput io.Writer
}

// Envs returns the config as a helm compatible
// set of env vars
func (c *Config) Envs() map[string]string {
	return map[string]string{

		"HOME": c.HomeDir,
		"PATH": c.Path,

		"HELM_CACHE_HOME":        c.HelmBaseDir,
		"HELM_CONFIG_HOME":       c.HelmBaseDir,
		"HELM_CONFIG_DATA_HOME":  c.HelmBaseDir,
		"HELM_PLUGINS":           c.HelmPluginsDirectory,
		"HELM_REGISTRY_CONFIG":   c.HelmRegistryConfig,
		"HELM_REPOSITORY_CONFIG": c.HelmRepositoryConfig,
		"HELM_REPOSITORY_CACHE":  c.HelmRepositoryCache,

		"HELM_DEBUG": fmt.Sprintf("%t", c.Debug),
	}
}

// New initialises a new Helm operator
func New(config *Config, auth aws.Authenticator, fs *afero.Afero) *Helm {
	return &Helm{
		config: config,
		providers: getter.Providers{
			{
				Schemes: []string{"https"},
				New:     getter.NewHTTPGetter,
			},
		},
		auth: auth,
		fs:   fs,
	}
}

// RepoAdd adds repo with given name and url
// nolint: funlen
func (h *Helm) RepoAdd(name, url string) error {
	restoreFn, err := EstablishEnv(h.config.Envs())

	defer func() {
		err = restoreFn()
	}()

	if err != nil {
		return err
	}

	settings := cli.New()

	err = h.fs.MkdirAll(h.config.HelmBaseDir, 0744)
	if err != nil {
		return err
	}

	unlockFn, err := Lock(settings.RepositoryConfig)
	if err != nil {
		return err
	}

	defer func() {
		err = unlockFn()
	}()

	b, err := h.fs.ReadFile(settings.RepositoryConfig)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	var f repo.File

	err = yaml.Unmarshal(b, &f)
	if err != nil {
		return err
	}

	if f.Has(name) {
		return nil
	}

	entry := repo.Entry{
		Name: name,
		URL:  url,
	}

	r, err := repo.NewChartRepository(&entry, h.providers)
	if err != nil {
		return err
	}

	_, err = r.DownloadIndexFile()
	if err != nil {
		return err
	}

	f.Update(&entry)

	err = f.WriteFile(settings.RepositoryConfig, 0644)
	if err != nil {
		return err
	}

	return nil
}

// RepoUpdate updates charts for all helm repos
func (h *Helm) RepoUpdate() error {
	restoreFn, err := EstablishEnv(h.config.Envs())

	defer func() {
		err = restoreFn()
	}()

	if err != nil {
		return err
	}

	settings := cli.New()

	exists, err := h.fs.Exists(settings.RepositoryConfig)
	if err != nil {
		return err
	}

	if !exists {
		return nil
	}

	f, err := repo.LoadFile(settings.RepositoryConfig)
	if err != nil {
		return err
	}

	for _, cfg := range f.Repositories {
		r, err := repo.NewChartRepository(cfg, h.providers)
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

// InstallConfig defines the variables that must be set to install a chart
type InstallConfig struct {
	// ReleaseName is the name given to the release in Kubernetes
	ReleaseName string
	// Version is the version of the chart to install
	Version string
	// Chart is the name of the chart to install
	Chart string
	// Repo is the name of the repository from which to install
	Repo string
	// Namespace is the namespace the chart will be installed in
	Namespace string
	// Values is a yaml encoded byte array of the values.yaml file
	ValuesBody []byte
}

// RepoChart returns the [repo]/[chart] string
func (c *InstallConfig) RepoChart() string {
	return fmt.Sprintf("%s/%s", c.Repo, c.Chart)
}

func (h *Helm) findRelease(releaseName string, config *action.Configuration) (*release.Release, error) {
	client := action.NewList(config)
	client.All = true

	client.SetStateMask()

	releases, err := client.Run()
	if err != nil {
		return nil, err
	}

	for _, r := range releases {
		if r.Name == releaseName {
			return r, nil
		}
	}

	return nil, nil
}

// Install a chart, comparable to: https://helm.sh/docs/helm/helm_install/
// though we have not implemented all the functionality found there
// Some details to consider about CRDs:
// - https://helm.sh/docs/chart_best_practices/custom_resource_definitions/#some-caveats-and-explanations
// nolint: funlen gocyclo
func (h *Helm) Install(kubeConfigPath string, cfg *InstallConfig) (*release.Release, error) {
	envs := h.config.Envs()
	envs["HELM_NAMESPACE"] = cfg.Namespace
	envs["KUBECONFIG"] = kubeConfigPath

	awsEnvs, err := h.auth.AsEnv()
	if err != nil {
		return nil, err
	}

	for k, v := range SplitEnv(awsEnvs) {
		envs[k] = v
	}

	restoreFn, err := EstablishEnv(envs)

	defer func() {
		err = restoreFn()
	}()

	if err != nil {
		return nil, err
	}

	settings := cli.New()

	actionConfig := new(action.Configuration)

	debug := func(format string, v ...interface{}) {
		if h.config.Debug {
			_, _ = fmt.Fprintf(h.config.DebugOutput, format, v...)
		}
	}

	restClient := &genericclioptions.ConfigFlags{
		KubeConfig: &kubeConfigPath,
		Namespace:  &cfg.Namespace,
	}

	err = actionConfig.Init(restClient, settings.Namespace(), DefaultHelmDriver, debug)
	if err != nil {
		return nil, err
	}

	rel, err := h.findRelease(cfg.ReleaseName, actionConfig)
	if err != nil {
		return nil, err
	}

	if rel != nil {
		switch rel.Info.Status {
		case release.StatusDeployed, release.StatusSuperseded, release.StatusPendingUpgrade, release.StatusPendingRollback, release.StatusPendingInstall:
			return rel, nil
		case release.StatusUninstalled:
			break
		case release.StatusFailed, release.StatusUninstalling, release.StatusUnknown:
			return nil, fmt.Errorf("release is in state: %s, cannot continue", rel.Info.Status)
		}
	}

	client := action.NewInstall(actionConfig)

	client.Version = cfg.Version
	client.ReleaseName = cfg.ReleaseName
	client.Namespace = settings.Namespace()
	client.CreateNamespace = true
	client.Wait = true
	client.Timeout = 150 * time.Second // nolint: gomnd Need to make this configurable
	client.Atomic = true

	cp, err := client.ChartPathOptions.LocateChart(cfg.RepoChart(), settings)
	if err != nil {
		return nil, err
	}

	valuesFile, err := StageValuesYaml(cfg.ValuesBody, h.fs)

	defer func() {
		_ = h.fs.Remove(valuesFile)
	}()

	if err != nil {
		return nil, err
	}

	valueOpts := &valuesPkg.Options{
		ValueFiles: []string{
			valuesFile,
		},
	}

	values, err := valueOpts.MergeValues(h.providers)
	if err != nil {
		return nil, err
	}

	chart, err := loader.Load(cp)
	if err != nil {
		return nil, err
	}

	err = IsChartInstallable(chart)
	if err != nil {
		return nil, err
	}

	r, err := client.Run(chart, values)
	if err != nil {
		return nil, err
	}

	return r, nil
}

// IsChartInstallable determines if a chart can be installed or not
func IsChartInstallable(ch *chartPkg.Chart) error {
	switch ch.Metadata.Type {
	case "", "application":
		return nil
	default:
		return fmt.Errorf("chart: %s is not installable", ch.Metadata.Type)
	}
}

// UnlockFn can be deferred in the calling function to
// unlock the file
type UnlockFn func() error

// Lock a file to ensure no concurrent access
func Lock(file string) (UnlockFn, error) {
	lockFile := strings.Replace(file, filepath.Ext(file), DefaultHelmLockExt, 1)
	lock := flock.New(lockFile)

	lockCtx, cancel := context.WithTimeout(context.Background(), DefaultHelmLockTimeout)
	defer cancel()

	locked, err := lock.TryLockContext(lockCtx, 1*time.Second)
	if err != nil {
		return nil, err
	}

	if locked {
		return lock.Unlock, nil
	}

	return nil, fmt.Errorf("failed to create lock: %s", lockFile)
}

// StageValuesYaml returns the path to the values.yaml
// temporary file we have created
func StageValuesYaml(body []byte, fs *afero.Afero) (string, error) {
	f, err := fs.TempFile("", "values")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary file for values.yaml: %w", err)
	}

	_, err = f.Write(body)
	if err != nil {
		return "", fmt.Errorf("failed to write content to values.yaml: %w", err)
	}

	err = f.Close()
	if err != nil {
		return "", fmt.Errorf("failed to close values.yaml: %w", err)
	}

	return f.Name(), nil
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

	fn := func() error {
		originalEnvVars := originalEnvVars

		os.Clearenv()

		originalSplit := SplitEnv(originalEnvVars)

		for k, v := range originalSplit {
			err := os.Setenv(k, v)
			if err != nil {
				return err
			}
		}

		return nil
	}

	for key, val := range envs {
		err := os.Setenv(key, val)
		if err != nil {
			return fn, err
		}
	}

	return fn, nil
}

// SplitEnv returns the split envvars
func SplitEnv(envs []string) map[string]string {
	out := map[string]string{}

	for _, envVar := range envs {
		e := strings.SplitN(envVar, "=", 2)

		var key, val string

		switch len(e) {
		case 0:
			continue
		case 1:
			key = e[0]
			val = ""
		case 2: // nolint: gomnd
			key = e[0]
			val = e[1]
		}

		out[key] = val
	}

	return out
}

// Chart contains the state for installing a chart
type Chart struct {
	RepositoryName string
	RepositoryURL  string

	ReleaseName string
	Version     string
	Chart       string
	Namespace   string

	Values interface{}
}

// InstallConfig returns a valid install config
func (c *Chart) InstallConfig() (*InstallConfig, error) {
	values, err := yaml.Marshal(c.Values)
	if err != nil {
		return nil, fmt.Errorf("failed to serialise values: %w", err)
	}

	return &InstallConfig{
		ReleaseName: c.ReleaseName,
		Version:     c.Version,
		Chart:       c.Chart,
		Repo:        c.RepositoryName,
		Namespace:   c.Namespace,
		ValuesBody:  values,
	}, nil
}

// MysqlValues demonstrates how the values can be set
type MysqlValues struct {
	MysqlRootPassword string           `yaml:"mysqlRootPassword"`
	Persistence       MysqlPersistence `yaml:"persistence"`
	ImagePullPolicy   string           `yaml:"imagePullPolicy"`
}

// MysqlPersistence demonstrates how the values can be set
type MysqlPersistence struct {
	Enabled bool `yaml:"enabled"`
}

// Mysql demonstrates how a chart can be specified
func Mysql(values interface{}) *Chart {
	return &Chart{
		RepositoryName: "stable",
		RepositoryURL:  "https://kubernetes-charts.storage.googleapis.com",
		ReleaseName:    "mysql",
		Version:        "1.6.6",
		Chart:          "mysql",
		Namespace:      "test-helm",
		Values:         values,
	}
}
