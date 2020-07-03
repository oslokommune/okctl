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
	"github.com/oslokommune/okctl/pkg/config/application"
	"github.com/oslokommune/okctl/pkg/config/repository"
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
	//DefaultLogDir is the default directory name for logs
	DefaultLogDir = "logs"
	//DefaultLogName is the default name of the file to log to
	DefaultLogName = "console.log"
	// DefaultLogDays determines how many days we keep the logs
	DefaultLogDays = 28
	// DefaultLogBackups determines how many backups we will keep
	DefaultLogBackups = 3
	// DefaultLogSizeInMb determines how much storage we will consume
	DefaultLogSizeInMb = 50

	// DefaultRepositoryConfig is the default filename of the okctl repository config
	DefaultRepositoryConfig = ".okctl.yml"
	// DefaultRepositoryConfigName is the default name of the okctl repository config
	DefaultRepositoryConfigName = ".okctl"
	// DefaultRepositoryConfigType is the default type of the okctl repository config
	DefaultRepositoryConfigType = "yml"

	// DefaultClusterConfig is the default filename of the eksctl cluster config
	DefaultClusterConfig = "cluster.yml"
	// DefaultClusterBaseDir is the default directory name of the eksctl cluster config
	DefaultClusterBaseDir = "cluster"

	// DefaultVpcOutputs is the default filename of the vpc outputs information
	DefaultVpcOutputs = "outputs.json"
	// DefaultVpcCloudFormationTemplate is the default filename of the vpc cloud formation template
	DefaultVpcCloudFormationTemplate = "vpc.yml"
	// DefaultVpcBaseDir is the default directory of the vpc resources
	DefaultVpcBaseDir = "vpc"

	// EnvPrefix of environment variables that will be processed by okctl
	EnvPrefix = "OKCTL"
	// EnvHome is the default env var parsed for determining the application home
	EnvHome = "OKCTL_HOME"
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

	AppDataLoader DataLoaderFn
	AppData       *application.Data

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
		AppDataLoader:  NoopDataLoader,
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

// LoadAppData will attempt to load okctl application data
func (c *Config) LoadAppData() error {
	c.AppData = nil

	if c.AppDataLoader == nil {
		c.AppDataLoader = NoopDataLoader
	}

	return c.AppDataLoader(c)
}

// WriteAppData will store the current app data state to disk
func (c *Config) WriteAppData(b []byte) error {
	home, err := c.GetHomeDir()
	if err != nil {
		return err
	}

	store := storage.NewFileSystemStorage(home)

	writer, err := store.Recreate(DefaultDir, DefaultConfig, 0644)
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

// WriteCurrentAppData writes the current app data state
// to disk
func (c *Config) WriteCurrentAppData() error {
	b, err := c.AppData.YAML()
	if err != nil {
		return err
	}

	return c.WriteAppData(b)
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

	writer, err := store.Recreate("", DefaultRepositoryConfig, 0644)
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

// GetAppDataDir will get the directory to where okctl
// application data should be written
func (c *Config) GetAppDataDir() (string, error) {
	home, err := c.GetHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, DefaultDir), nil
}

// GetAppDataPath returns the path to the okctl application
// config path
func (c *Config) GetAppDataPath() (string, error) {
	base, err := c.GetAppDataDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(base, DefaultConfig), nil
}

// GetLogName returns the path to a logfile
func (c *Config) GetLogName() (string, error) {
	base, err := c.GetAppDataDir()
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
