// Package config interacts with all configuration state
package config

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/oslokommune/okctl/pkg/client/store"

	"github.com/go-git/go-git/v5"
	"github.com/mitchellh/go-homedir"
	"github.com/oslokommune/okctl/pkg/api/core"
	"github.com/oslokommune/okctl/pkg/config/state"
	"github.com/oslokommune/okctl/pkg/context"
	"github.com/oslokommune/okctl/pkg/rotatefilehook"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// nolint: golint
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

	DefaultRepositoryStateFile = ".okctl.yml"

	DefaultEKSKubernetesVersion = "1.17"

	DefaultArgoCDNamespace = "argocd"

	DefaultClusterConfig         = "cluster.yml"
	DefaultClusterKubeConfig     = "kubeconfig"
	DefaultClusterAwsConfig      = "aws-config"
	DefaultClusterAwsCredentials = "aws-credentials"
	DefaultClusterBaseDir        = "cluster"

	DefaultHelmBaseDir          = "helm"
	DefaultHelmRegistryConfig   = "registry.json"
	DefaultHelmRepositoryConfig = "repositories.yaml"
	DefaultHelmRepositoryCache  = "repository"
	DefaultHelmPluginsDirectory = "plugins"

	DefaultVpcOutputs                = "vpc-outputs.json"
	DefaultVpcCloudFormationTemplate = "vpc-cf.yml"
	DefaultVpcBaseDir                = "vpc"

	DefaultAWSLoadBalancerControllerBaseDir         = "aws-load-balancer-controller"
	DefaultAlbIngressControllerBaseDir              = "alb-ingress-controller"
	DefaultAliasBaseDir                             = "alias"
	DefaultAliasCloudFormationTemplate              = "alias-cf.yaml"
	DefaultArgoCDBaseDir                            = "argocd"
	DefaultArgoOutputsFile                          = "argocd-outputs.json"
	DefaultCertificateBaseDir                       = "certificates"
	DefaultCertificateCloudFormationTemplate        = "certificate-cf.yml"
	DefaultCertificateOutputsFile                   = "certificate-outputs.json"
	DefaultDomainBaseDir                            = "domains"
	DefaultDomainCloudFormationTemplate             = "domains-cf.yml"
	DefaultDomainOutputsFile                        = "domains-outputs.json"
	DefaultExternalDNSBaseDir                       = "external-dns"
	DefaultExternalSecretsBaseDir                   = "external-secrets"
	DefaultHelmChartFile                            = "helm-chart.json"
	DefaultHelmOutputsFile                          = "helm-outputs.json"
	DefaultHelmReleaseFile                          = "helm-release.json"
	DefaultIdentityPoolBaseDir                      = "identitypool"
	DefaultIdentityPoolCloudFormationTemplate       = "identitypool-cf.yaml"
	DefaultIdentityPoolOutputsFile                  = "identitypool-outputs.json"
	DefaultIdentityPoolClientsBaseDir               = "clients"
	DefaultIdentityPoolClientCloudFormationTemplate = "ipc-cf.yaml"
	DefaultIdentityPoolClientOutputsFile            = "ipc-outputs.json"
	DefaultIdentityPoolUsersBaseDir                 = "users"
	DefaultIdentityPoolUserOutputsFile              = "ipu-outputs.json"
	DefaultIdentityPoolUserCloudFormationTemplate   = "ipu-cf.yaml"
	DefaultKubeOutputsFile                          = "kube-outputs.json"
	DefaultParameterBaseDir                         = "parameters"
	DefaultParameterOutputsFile                     = "parameter-outputs.json"
	DefaultPolicyCloudFormationTemplateFile         = "policy-cf.yml"
	DefaultPolicyOutputFile                         = "policy-outputs.json"
	DefaultServiceAccountConfigFile                 = "service-account-config.yml"
	DefaultServiceAccountOutputsFile                = "service-account-outputs.json"

	// EnvPrefix of environment variables that will be processed by okctl
	EnvPrefix = "OKCTL"
	// EnvHome is the default env var parsed for determining the application home
	EnvHome = "OKCTL_HOME"

	// DefaultApplicationOverlayBaseDir is where the directory where overlay files reside
	DefaultApplicationOverlayBaseDir = "base"
	// DefaultApplicationDir is where the application overlay files reside
	DefaultApplicationDir = "applications"

	// DefaultKeyringServiceName is the name of the keyring or encrypted file used to store client secrets
	DefaultKeyringServiceName = "okctlService"

	// DefaultRequiredEpis number of elastic ips required for cluster creation
	DefaultRequiredEpis = 3
	// DefaultRequiredVpcs number of vpc(s) required for cluster creation
	DefaultRequiredVpcs = 1
	// DefaultRequiredIgws number of internet gateways required for cluster creation
	DefaultRequiredIgws = 1
	// DefaultRequiredEpisTestCluster number of elastic ips required for testcluster creation
	DefaultRequiredEpisTestCluster = 1
	// DefaultRequiredVpcsTestCluster number of vpc(s) required for testcluster creation
	DefaultRequiredVpcsTestCluster = 1
	// DefaultRequiredIgwsTestCluster number of internet gateways required for testcluster creation
	DefaultRequiredIgwsTestCluster = 1

	DefaultNameserverRecordTTL = 300
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
	UserState      *state.User

	RepoDataLoader   DataLoaderFn
	RepoState        *state.Repository
	RepoStateWithEnv state.RepositoryStateWithEnv

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
		UserState:      &state.User{},
		RepoDataLoader: NoopDataLoader,
		RepoState:      &state.Repository{},
		Destination:    dest,
		ServerURL:      fmt.Sprintf("http://%s/v1/", dest),
	}
}

// EnableFileLog for writing logs to a file
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
		Levels: []logrus.Level{
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel,
			logrus.WarnLevel,
			logrus.InfoLevel,
			logrus.DebugLevel,
			logrus.TraceLevel,
		},
		Formatter: &logrus.JSONFormatter{
			TimestampFormat: time.RFC822,
		},
	})
	if err != nil {
		return fmt.Errorf("initialising the file rotate hook: %v", err)
	}

	c.Logger.AddHook(rotateFileHook)

	return nil
}

// SetFormat sets the response type when encoding
func (c *Config) SetFormat(responseType core.EncodeResponseType) {
	c.format = responseType
}

// Format returns the encode response type
func (c *Config) Format() core.EncodeResponseType {
	return c.format
}

// LoadRepoData will attempt to load repository data
func (c *Config) LoadRepoData() error {
	if c.RepoDataLoader == nil {
		c.RepoDataLoader = NoopDataLoader
	}

	return c.RepoDataLoader(c)
}

// LoadUserData will attempt to load okctl application data
func (c *Config) LoadUserData() error {
	if c.UserDataLoader == nil {
		c.UserDataLoader = NoopDataLoader
	}

	return c.UserDataLoader(c)
}

// WriteCurrentUserData writes the current app data state
// to disk
func (c *Config) WriteCurrentUserData() error {
	userDir, err := c.GetUserDataDir()
	if err != nil {
		return err
	}

	_, err = store.NewFileSystem(userDir, c.FileSystem).
		StoreStruct(DefaultConfig, c.UserState, store.ToYAML()).
		Do()
	if err != nil {
		return err
	}

	return nil
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

// GetRepoStateDir will return the directory where repo data should be read/written
func (c *Config) GetRepoStateDir() (string, error) {
	return c.GetRepoDir()
}

// GetRepoStatePath will return the filename where repo data should be read/written
func (c *Config) GetRepoStatePath() (string, error) {
	base, err := c.GetRepoStateDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(base, DefaultRepositoryStateFile), nil
}

// WriteCurrentRepoState will write current repo state to disk
func (c *Config) WriteCurrentRepoState() error {
	repoDir, err := c.GetRepoDir()
	if err != nil {
		return err
	}

	_, err = store.NewFileSystem(repoDir, c.FileSystem).
		StoreStruct(DefaultRepositoryStateFile, c.RepoState, store.ToYAML()).
		Do()
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
	base, err := c.GetRepoStateDir()
	if err != nil {
		return "", err
	}

	return path.Join(base, c.RepoState.Metadata.OutputDir, env), nil
}

// GetRepoApplicationBaseDir returns the directory where application
// resources are stored
func (c *Config) GetRepoApplicationBaseDir() (string, error) {
	base, err := c.GetRepoStateDir()
	if err != nil {
		return "", err
	}

	return path.Join(base, c.RepoState.Metadata.OutputDir, DefaultApplicationOverlayBaseDir, DefaultApplicationDir), nil
}
