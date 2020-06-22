// Package config interacts with all configuration state
package config

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/mitchellh/go-homedir"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/config/application"
	"github.com/oslokommune/okctl/pkg/config/repository"
	"github.com/oslokommune/okctl/pkg/context"
	"github.com/oslokommune/okctl/pkg/storage"
	"github.com/pkg/errors"
	"github.com/sanathkr/go-yaml"
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

	// DefaultRepositoryConfig is the default filename of the okctl repository config
	DefaultRepositoryConfig = ".okctl.yml"
	// DefaultRepositoryConfigName is the default name of the okctl repository config
	DefaultRepositoryConfigName = ".okctl"
	// DefaultRepositoryConfigType is the default type of the okctl repository config
	DefaultRepositoryConfigType = "yml"

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

	homeDir string
	repoDir string
}

// New Config initialises a default okctl configuration
func New() *Config {
	return &Config{
		Context:        context.New(),
		AppDataLoader:  NoopDataLoader,
		RepoDataLoader: NoopDataLoader,
	}
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

// GetRepoOutputDir return the repository output directory,
// where cloud formation stacks, etc., should be written
func (c *Config) GetRepoOutputDir(env string) (string, error) {
	base, err := c.GetRepoDataDir()
	if err != nil {
		return "", err
	}

	return path.Join(base, c.RepoData.OutputDir, env), nil
}

// WriteToOutputDir will write state to the given output dir, env and
// filepath
func (c *Config) WriteToOutputDir(env, filePath string, b []byte) error {
	outDir, err := c.GetRepoOutputDir(env)
	if err != nil {
		return err
	}

	store := storage.NewFileSystemStorage(outDir)

	base, file := path.Split(filePath)

	writer, err := store.Recreate(base, file, 0644)
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

// DeleteFromOutputDir removes everything for a given env from path
func (c *Config) DeleteFromOutputDir(env, filepath string) error {
	outDir, err := c.GetRepoOutputDir(env)
	if err != nil {
		return err
	}

	store := storage.NewFileSystemStorage(outDir)

	return store.RemoveAll(filepath)
}

// WriteClusterConfig stores the cluster config
func (c *Config) WriteClusterConfig(env string, cfg *v1alpha1.ClusterConfig) error {
	b, err := cfg.YAML()
	if err != nil {
		return err
	}

	return c.WriteToOutputDir(env, path.Join("cluster", "config.yml"), b)
}

// DeleteClusterConfig deletes the cluster config for a given environment
func (c *Config) DeleteClusterConfig(env string) error {
	return c.DeleteFromOutputDir(env, path.Join("cluster", "config.yml"))
}

// ClusterConfig loads the cluster configuration for the given environment
func (c *Config) ClusterConfig(env string) (*v1alpha1.ClusterConfig, error) {
	outDir, err := c.GetRepoOutputDir(env)
	if err != nil {
		return nil, err
	}

	store := storage.NewFileSystemStorage(outDir)

	b, err := store.ReadAll(path.Join("cluster", "config.yml"))
	if err != nil {
		return nil, err
	}

	cfg := v1alpha1.NewClusterConfig()

	return cfg, yaml.Unmarshal(b, cfg)
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

// HasCluster returns true if we can find the cluster in the repo data
func (c *Config) HasCluster(env string) bool {
	for _, c := range c.RepoData.Clusters {
		if c.Environment == env {
			return true
		}
	}

	return false
}
