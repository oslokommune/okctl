package config

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/mitchellh/go-homedir"
	"github.com/oslokommune/okctl/pkg/config/application"
	"github.com/oslokommune/okctl/pkg/config/repository"
	"github.com/oslokommune/okctl/pkg/context"
	"github.com/pkg/errors"
)

const (
	DefaultDir        = ".okctl"
	DefaultConfig     = "conf.yml"
	DefaultConfigName = "conf"
	DefaultConfigType = "yml"

	DefaultRepositoryConfig     = ".okctl.yml"
	DefaultRepositoryConfigName = ".okctl"
	DefaultRepositoryConfigType = "yml"

	EnvPrefix = "OKCTL"
	EnvHome   = "OKCTL_HOME"
)

type DataLoaderFn func(*Config) error

func NoopDataLoader(_ *Config) error {
	return nil
}

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

func (c *Config) LoadRepoData() error {
	c.RepoData = nil

	if c.RepoDataLoader == nil {
		c.RepoDataLoader = NoopDataLoader
	}

	return c.RepoDataLoader(c)
}

func (c *Config) LoadAppData() error {
	c.AppData = nil

	if c.AppDataLoader == nil {
		c.AppDataLoader = NoopDataLoader
	}

	return c.AppDataLoader(c)
}

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

func (c *Config) GetRepoDataDir() (string, error) {
	return c.GetRepoDir()
}

func (c *Config) GetRepoDataPath() (string, error) {
	base, err := c.GetRepoDataDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(base, DefaultRepositoryConfig), nil
}

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

func (c *Config) GetAppDataDir() (string, error) {
	home, err := c.GetHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, DefaultDir), nil
}

func (c *Config) GetAppDataPath() (string, error) {
	base, err := c.GetAppDataDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(base, DefaultConfig), nil
}

func (c *Config) GetRepoOutputDir(env string) (string, error) {
	base, err := c.GetRepoDataDir()
	if err != nil {
		return "", err
	}

	return path.Join(base, c.RepoData.OutputDir, env), nil
}

func (c *Config) ClusterName(env string) string {
	return fmt.Sprintf("%s-%s", c.RepoData.Name, env)
}
