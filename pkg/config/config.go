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
	"github.com/oslokommune/okctl/pkg/config/application"
	"github.com/oslokommune/okctl/pkg/config/repository"
	"github.com/oslokommune/okctl/pkg/context"
	"github.com/oslokommune/okctl/pkg/storage"
	"github.com/pkg/errors"
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

// ClusterName returns a consistent cluster name
func (c *Config) ClusterName(env string) string {
	return fmt.Sprintf("%s-%s", c.RepoData.Name, env)
}
