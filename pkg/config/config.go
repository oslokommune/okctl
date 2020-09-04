// Package config interacts with all configuration state
package config

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/mitchellh/go-homedir"
	"github.com/oslokommune/okctl/pkg/api/core"
	"github.com/oslokommune/okctl/pkg/config/repository"
	"github.com/oslokommune/okctl/pkg/config/user"
	"github.com/oslokommune/okctl/pkg/context"
	"github.com/oslokommune/okctl/pkg/rotatefilehook"
	"github.com/oslokommune/okctl/pkg/storage"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	// DefaultDir is the default location directory for the okctl application config
	DefaultDir = ".okctl"
	// DefaultConfig is the default filename of the okctl application config
	DefaultConfig = "conf.yml"
	// DefaultConfigName is the default name of the okctl application config
	DefaultConfigName = "conf"
	// DefaultConfigType is the default type of the okctl application config
	DefaultConfigType = "yml"
	// DefaultLogDir is the default directory name for logs
	DefaultLogDir = "logs"
	// DefaultLogName is the default name of the file to log to
	DefaultLogName = "console.log"
	// DefaultLogDays determines how many days we keep the logs
	DefaultLogDays = 28
	// DefaultLogBackups determines how many backups we will keep
	DefaultLogBackups = 3
	// DefaultLogSizeInMb determines how much storage we will consume
	DefaultLogSizeInMb = 50
	// DefaultCredentialsDirName sets the name of the directory for creds
	DefaultCredentialsDirName = "credentials"

	// DefaultRepositoryConfig is the default filename of the okctl repository config
	DefaultRepositoryConfig = ".okctl.yml"
	// DefaultRepositoryConfigName is the default name of the okctl repository config
	DefaultRepositoryConfigName = ".okctl"
	// DefaultRepositoryConfigType is the default type of the okctl repository config
	DefaultRepositoryConfigType = "yml"

	// DefaultClusterConfig is the default filename of the eksctl cluster config
	DefaultClusterConfig = "cluster.yml"
	// DefaultClusterKubeConfig is the default filename of the kubectl kubeconfig
	DefaultClusterKubeConfig = "kubeconfig"
	// DefaultClusterAwsConfig is the default filename of the aws config
	DefaultClusterAwsConfig = "aws-config"
	// DefaultClusterAwsCredentials is the default filename of the aws credentials
	DefaultClusterAwsCredentials = "aws-credentials"
	// DefaultClusterBaseDir is the default directory name of the eksctl cluster config
	DefaultClusterBaseDir = "cluster"

	// DefaultHelmBaseDir is the default directory for storing helm related stuff
	DefaultHelmBaseDir = "helm"
	// DefaultHelmRegistryConfig is the name of the registry config file
	DefaultHelmRegistryConfig = "registry.json"
	// DefaultHelmRepositoryConfig is the name of the repositories file
	DefaultHelmRepositoryConfig = "repositories.yaml"
	// DefaultHelmRepositoryCache is the name of the repository cache directory
	DefaultHelmRepositoryCache = "repository"
	// DefaultHelmPluginsDirectory is the name of the plugins directory
	DefaultHelmPluginsDirectory = "plugins"

	// DefaultVpcOutputs is the default filename of the vpc outputs information
	DefaultVpcOutputs = "vpc-outputs.json"
	// DefaultVpcCloudFormationTemplate is the default filename of the vpc cloud formation template
	DefaultVpcCloudFormationTemplate = "vpc-cf.yml"
	// DefaultVpcBaseDir is the default directory of the vpc resources
	DefaultVpcBaseDir = "vpc"

	// DefaultPolicyCloudFormationTemplateFile is the default filename of the cloud formation template
	DefaultPolicyCloudFormationTemplateFile = "policy-cf.yml"
	// DefaultPolicyOutputFile is the default filename of the outputs
	DefaultPolicyOutputFile = "policy-outputs.json"
	// DefaultServiceAccountOutputsFile is the default file name of the service account outputs
	DefaultServiceAccountOutputsFile = "service-account-outputs.json"
	// DefaultServiceAccountConfigFile is the default file name of the service account config
	DefaultServiceAccountConfigFile = "service-account-config.yml"
	// DefaultHelmOutputsFile is the default file name of the helm output
	DefaultHelmOutputsFile = "helm-outputs.json"
	// DefaultHelmReleaseFile is the default file name of the helm release
	DefaultHelmReleaseFile = "helm-release.json"
	// DefaultHelmChartFile is the default file name of the chart
	DefaultHelmChartFile = "helm-chart.json"
	// DefaultKubeOutputsFile is the default file name of the kube output
	DefaultKubeOutputsFile = "kube-outputs.json"
	// DefaultDomainOutputsFile is the default file name of the domain output
	DefaultDomainOutputsFile = "domains-outputs.json"
	// DefaultDomainCloudFormationTemplate is the default file name of the cloud formation template
	DefaultDomainCloudFormationTemplate = "domains-cf.yml"
	// DefaultCertificateOutputsFile is the default file name of the domain output
	DefaultCertificateOutputsFile = "certificate-outputs.json"
	// DefaultCertificateCloudFormationTemplate is the default file name of the cloud formation template
	DefaultCertificateCloudFormationTemplate = "certificate-cf.yml"
	// DefaultParameterOutputsFile is the default file name of the outputs
	DefaultParameterOutputsFile = "parameter-outputs.json"

	// DefaultExternalSecretsBaseDir is the default directory of the external secrets resources
	DefaultExternalSecretsBaseDir = "external-secrets"
	// DefaultAlbIngressControllerBaseDir is the default directory of the external secrets resources
	DefaultAlbIngressControllerBaseDir = "alb-ingress-controller"
	// DefaultExternalDNSBaseDir is the default directory of the external dns resources
	DefaultExternalDNSBaseDir = "external-dns"
	// DefaultDomainBaseDir is the default directory for domains
	DefaultDomainBaseDir = "domains"
	// DefaultCertificateBaseDir is the default directory for certificates
	DefaultCertificateBaseDir = "certificates"
	// DefaultParameterBaseDir is the default directory for parameters
	DefaultParameterBaseDir = "parameters"
	// DefaultArgoCDBaseDir is the default directory for argo cd
	DefaultArgoCDBaseDir = "argocd"

	// EnvPrefix of environment variables that will be processed by okctl
	EnvPrefix = "OKCTL"
	// EnvHome is the default env var parsed for determining the application home
	EnvHome = "OKCTL_HOME"

	// DefaultKeyringServiceName is the name of the keyring or encrypted file used to store client secrets
	DefaultKeyringServiceName = "okctlService"
)

// DataLoaderFn is the type for loading configuration data
type DataLoaderFn func(*Config) error

// NoopDataLoader does nothing
func NoopDataLoader(_ *Config) error {
	return nil
}

// Config stores state for representing and interacting
// with okctl state
type Config struct {
	*context.Context

	UserDataLoader DataLoaderFn
	UserData       *user.Data

	RepoDataLoader DataLoaderFn
	RepoData       *repository.Data

	Destination string
	ServerURL   string

	format  core.EncodeResponseType
	homeDir string
	repoDir string
}

// New Config initialises a default okctl configuration
func New() *Config {
	dest := "127.0.0.1:8085"

	return &Config{
		Context:        context.New(),
		UserDataLoader: NoopDataLoader,
		RepoDataLoader: NoopDataLoader,
		Destination:    dest,
		ServerURL:      fmt.Sprintf("http://%s/v1/", dest),
	}
}

// EnableFileLog turns on logging to files in addition to console
func (c *Config) EnableFileLog() error {
	logFile, err := c.GetLogName()
	if err != nil {
		return err
	}

	rotateFileHook, err := rotatefilehook.NewRotateFileHook(rotatefilehook.RotateFileConfig{
		Filename:   logFile,
		MaxSize:    DefaultLogSizeInMb,
		MaxBackups: DefaultLogBackups,
		MaxAge:     DefaultLogDays,
		Level:      logrus.DebugLevel,
		Formatter: &logrus.JSONFormatter{
			TimestampFormat: time.RFC822,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to initialize file rotate hook: %v", err)
	}

	c.Logger.AddHook(rotateFileHook)

	return nil
}

// SetFormat sets the encode response type
func (c *Config) SetFormat(responseType core.EncodeResponseType) {
	c.format = responseType
}

// Format returns the encode response type
func (c *Config) Format() core.EncodeResponseType {
	return c.format
}

// LoadRepoData will attempt to load repository data
func (c *Config) LoadRepoData() error {
	c.RepoData = nil

	if c.RepoDataLoader == nil {
		c.RepoDataLoader = NoopDataLoader
	}

	return c.RepoDataLoader(c)
}

// LoadUserData will attempt to load okctl application data
func (c *Config) LoadUserData() error {
	c.UserData = nil

	if c.UserDataLoader == nil {
		c.UserDataLoader = NoopDataLoader
	}

	return c.UserDataLoader(c)
}

// WriteUserData will store the current app data state to disk
func (c *Config) WriteUserData(b []byte) error {
	home, err := c.GetHomeDir()
	if err != nil {
		return err
	}

	store := storage.NewFileSystemStorage(home)

	writer, err := store.Recreate(DefaultDir, DefaultConfig, 0o644)
	if err != nil {
		return err
	}

	defer func() {
		err = writer.Close()
	}()

	_, err = io.Copy(writer, bytes.NewReader(b))
	if err != nil {
		return err
	}

	return nil
}

// WriteCurrentUserData writes the current app data state
// to disk
func (c *Config) WriteCurrentUserData() error {
	b, err := c.UserData.YAML()
	if err != nil {
		return err
	}

	return c.WriteUserData(b)
}

// GetRepoDir will return the currently active repository directory
func (c *Config) GetRepoDir() (string, error) {
	if len(c.repoDir) != 0 {
		return c.repoDir, nil
	}

	repoDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	absoluteRepo, err := filepath.Abs(repoDir)
	if err != nil {
		return "", err
	}

	_, err = git.PlainOpen(absoluteRepo)
	if err != nil {
		return "", errors.Wrap(err, "working directory must be a git repository")
	}

	c.repoDir = absoluteRepo

	return c.repoDir, nil
}

// GetRepoDataDir will return the directory where repo data should be read/written
func (c *Config) GetRepoDataDir() (string, error) {
	return c.GetRepoDir()
}

// GetRepoDataPath will return the filename where repo data should be read/written
func (c *Config) GetRepoDataPath() (string, error) {
	base, err := c.GetRepoDataDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(base, DefaultRepositoryConfig), nil
}

// WriteCurrentRepoData will write current repo state to disk
func (c *Config) WriteCurrentRepoData() error {
	data, err := c.RepoData.YAML()
	if err != nil {
		return err
	}

	c.Logger.Debugf("write current repo data: %s", string(data))

	return c.WriteRepoData(data)
}

// WriteRepoData will write the provided repo state to disk
func (c *Config) WriteRepoData(b []byte) error {
	repo, err := c.GetRepoDir()
	if err != nil {
		return err
	}

	store := storage.NewFileSystemStorage(repo)

	writer, err := store.Recreate("", DefaultRepositoryConfig, 0o644)
	if err != nil {
		return err
	}

	defer func() {
		err = writer.Close()
	}()

	_, err = io.Copy(writer, bytes.NewReader(b))
	if err != nil {
		return err
	}

	return nil
}

// GetHomeDir will get the okctl application home dir
func (c *Config) GetHomeDir() (string, error) {
	if len(c.homeDir) != 0 {
		return c.homeDir, nil
	}

	homeDir := os.Getenv(EnvHome)

	if len(homeDir) == 0 {
		dir, err := homedir.Dir()
		if err != nil {
			return "", err
		}

		homeDir = dir
	}

	absoluteHome, err := filepath.Abs(homeDir)
	if err != nil {
		return "", err
	}

	c.homeDir = absoluteHome

	return c.homeDir, nil
}

// GetUserDataDir will get the directory to where okctl
// application data should be written
func (c *Config) GetUserDataDir() (string, error) {
	home, err := c.GetHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, DefaultDir), nil
}

// GetUserDataPath returns the path to the okctl application
// config path
func (c *Config) GetUserDataPath() (string, error) {
	base, err := c.GetUserDataDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(base, DefaultConfig), nil
}

// GetLogName returns the path to a logfile
func (c *Config) GetLogName() (string, error) {
	base, err := c.GetUserDataDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(base, DefaultLogDir, DefaultLogName), nil
}

// GetRepoOutputDir return the repository output directory,
// where cloud formation stacks, etc., should be written
func (c *Config) GetRepoOutputDir(env string) (string, error) {
	base, err := c.GetRepoDataDir()
	if err != nil {
		return "", err
	}

	return path.Join(base, c.RepoData.OutputDir, env), nil
}

// ClusterName returns a consistent cluster name
func (c *Config) ClusterName(env string) string {
	return fmt.Sprintf("%s-%s", c.RepoData.Name, env)
}

// AWSAccountID returns the aws account ID for the given env
func (c *Config) AWSAccountID(env string) (string, error) {
	for _, cluster := range c.RepoData.Clusters {
		if cluster.Environment == env {
			return cluster.AWS.AccountID, nil
		}
	}

	return "", fmt.Errorf("could not find configuration for cluster: %s", env)
}

// Domain returns the domain for the given env
func (c *Config) Domain(env string) string {
	for _, cluster := range c.RepoData.Clusters {
		if cluster.Environment == env {
			return cluster.HostedZone.Domain
		}
	}

	return ""
}

// FQDN returns the fully qualified domain name for a given env
func (c *Config) FQDN(env string) string {
	for _, cluster := range c.RepoData.Clusters {
		if cluster.Environment == env {
			return cluster.HostedZone.FQDN
		}
	}

	return ""
}

// HostedZoneIsDelegated returns whether the hosted zone has been delegated or not
func (c *Config) HostedZoneIsDelegated(env string) bool {
	for _, cluster := range c.RepoData.Clusters {
		if cluster.Environment == env {
			return cluster.HostedZone.IsDelegated
		}
	}

	return false
}

// SetHostedZoneIsDelegated to the provided value
func (c *Config) SetHostedZoneIsDelegated(val bool, env string) {
	for i, cluster := range c.RepoData.Clusters {
		if cluster.Environment == env {
			cluster.HostedZone.IsDelegated = val
		}

		c.RepoData.Clusters[i] = cluster
	}
}

// HostedZoneIsCreated returns whether the hosted zone has been created or not
func (c *Config) HostedZoneIsCreated(env string) bool {
	for _, cluster := range c.RepoData.Clusters {
		if cluster.Environment == env {
			return cluster.HostedZone.IsCreated
		}
	}

	return false
}

// SetHostedZoneIsCreated to the provided value
func (c *Config) SetHostedZoneIsCreated(val bool, env string) {
	for i, cluster := range c.RepoData.Clusters {
		if cluster.Environment == env {
			cluster.HostedZone.IsCreated = val
		}

		c.RepoData.Clusters[i] = cluster
	}
}

// GithubRepository returns the selected github repository name and url
func (c *Config) GithubRepository(env string) repository.Repository {
	cluster := c.RepoData.ClusterForEnv(env)
	if cluster == nil {
		return repository.Repository{}
	}

	return cluster.Github.Repository
}

// SetGithubRepository sets the github repository name and url
func (c *Config) SetGithubRepository(repository repository.Repository, env string) {
	cluster := c.RepoData.ClusterForEnv(env)
	if cluster == nil {
		return
	}

	cluster.Github.Repository = repository

	c.RepoData.SetClusterForEnv(cluster, env)
}

// GithubTeamName returns the selected github team
func (c *Config) GithubTeamName(env string) string {
	cluster := c.RepoData.ClusterForEnv(env)
	if cluster == nil {
		return ""
	}

	return cluster.Github.Team
}

// SetGithubTeamName sets the github team
func (c *Config) SetGithubTeamName(name, env string) {
	cluster := c.RepoData.ClusterForEnv(env)
	if cluster == nil {
		return
	}

	cluster.Github.Team = name

	c.RepoData.SetClusterForEnv(cluster, env)
}

// GithubOrganisationName returns the selected github organisation
func (c *Config) GithubOrganisationName(env string) string {
	cluster := c.RepoData.ClusterForEnv(env)
	if cluster == nil {
		return ""
	}

	return cluster.Github.Organisation
}

// SetGithubOrganisationName sets the github organisation
func (c *Config) SetGithubOrganisationName(name, env string) {
	cluster := c.RepoData.ClusterForEnv(env)
	if cluster == nil {
		return
	}

	cluster.Github.Organisation = name

	c.RepoData.SetClusterForEnv(cluster, env)
}

// GithubOauthApp returns the name and client id
func (c *Config) GithubOauthApp(env string) repository.OauthApp {
	cluster := c.RepoData.ClusterForEnv(env)
	if cluster == nil {
		return repository.OauthApp{}
	}

	return cluster.Github.OauthApp
}

// SetGithubOauthApp sets the name and client id
func (c *Config) SetGithubOauthApp(app repository.OauthApp, env string) {
	cluster := c.RepoData.ClusterForEnv(env)
	if cluster == nil {
		return
	}

	cluster.Github.OauthApp = app

	c.RepoData.SetClusterForEnv(cluster, env)
}

// GithubDeployKey returns the github deploy key title and id
func (c *Config) GithubDeployKey(env string) repository.DeployKey {
	cluster := c.RepoData.ClusterForEnv(env)
	if cluster == nil {
		return repository.DeployKey{}
	}

	return cluster.Github.DeployKey
}

// SetGithubDeployKey sets the github deploy key title and id
func (c *Config) SetGithubDeployKey(key repository.DeployKey, env string) {
	cluster := c.RepoData.ClusterForEnv(env)
	if cluster == nil {
		return
	}

	cluster.Github.DeployKey = key

	c.RepoData.SetClusterForEnv(cluster, env)
}

// ArgoCD returns the argo cd state
func (c *Config) ArgoCD(env string) repository.ArgoCD {
	cluster := c.RepoData.ClusterForEnv(env)
	if cluster == nil {
		return repository.ArgoCD{}
	}

	return cluster.ArgoCD
}

// SetArgoCD sets the argocd to the provided state
func (c *Config) SetArgoCD(argo repository.ArgoCD, env string) {
	cluster := c.RepoData.ClusterForEnv(env)
	if cluster == nil {
		return
	}

	cluster.ArgoCD = argo

	c.RepoData.SetClusterForEnv(cluster, env)
}
