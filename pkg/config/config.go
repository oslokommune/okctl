// Package config interacts with all configuration state
package config

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/oslokommune/okctl/pkg/paths"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/logging"

	"github.com/rancher/k3d/v3/cmd/util"

	"github.com/oslokommune/okctl/pkg/config/constant"

	"github.com/oslokommune/okctl/pkg/client/store"

	"github.com/go-git/go-git/v5"
	"github.com/mitchellh/go-homedir"
	"github.com/oslokommune/okctl/pkg/config/state"
	"github.com/oslokommune/okctl/pkg/context"
	"github.com/pkg/errors"
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

	RepoDataLoader DataLoaderFn

	Declaration *v1alpha1.Cluster

	Destination   string
	ServerURL     string
	ServerBaseURL string

	homeDir string
	repoDir string
}

// New Config initialises a default okctl configuration
func New() *Config {
	port, err := util.GetFreePort()
	if err != nil {
		port = 8085
	}

	dest := fmt.Sprintf("127.0.0.1:%d", port)

	return &Config{
		Context:        context.New(),
		UserDataLoader: NoopDataLoader,
		UserState:      &state.User{},
		RepoDataLoader: NoopDataLoader,
		Destination:    dest,
		ServerURL:      fmt.Sprintf("http://%s/v1/", dest),
		ServerBaseURL:  fmt.Sprintf("http://%s/", dest),
	}
}

// EnableFileLog for writing logs to a file
func (c *Config) EnableFileLog() error {
	logFile, err := c.GetFullLogFilePath(constant.DefaultLogName)
	if err != nil {
		return err
	}

	return logging.AddLogFileHook(c.Logger, logFile)
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
		StoreStruct(constant.DefaultConfig, c.UserState, store.ToYAML()).
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

	repoDir, err := paths.GetAbsoluteRepositoryRootDirectory()
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

	return filepath.Join(base, constant.DefaultRepositoryStateFile), nil
}

// GetHomeDir will return the value of OKCTL_HOME. Will default to the user's home directory. I.e.: /home/user/ in unix
// based systems
func (c *Config) GetHomeDir() (string, error) {
	if len(c.homeDir) != 0 {
		return c.homeDir, nil
	}

	homeDir := os.Getenv(constant.EnvHome)

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
// application data should be written. I.e.: /home/user/.okctl in unix based systems
func (c *Config) GetUserDataDir() (string, error) {
	home, err := c.GetHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, constant.DefaultDir), nil
}

// GetUserDataPath returns the path to the okctl application
// config path. I.e.: /home/user/.okctl/conf.yml
func (c *Config) GetUserDataPath() (string, error) {
	base, err := c.GetUserDataDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(base, constant.DefaultConfig), nil
}

// GetFullLogFilePath returns the full path to a logfile
func (c *Config) GetFullLogFilePath(logFileName string) (string, error) {
	base, err := c.GetUserDataDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(base, constant.DefaultLogDir, logFileName), nil
}

// GetRepoOutputDir return the repository output directory,
// where cloud formation stacks, etc., should be written
func (c *Config) GetRepoOutputDir() (string, error) {
	base, err := c.GetRepoStateDir()
	if err != nil {
		return "", err
	}

	return path.Join(base, c.Declaration.Github.OutputPath, c.Declaration.Metadata.Name), nil
}

// GetRepoApplicationsOutputDir returns the directory where application
// resources are stored
func (c *Config) GetRepoApplicationsOutputDir() (string, error) {
	base, err := c.GetRepoStateDir()
	if err != nil {
		return "", err
	}

	return path.Join(base, c.Declaration.Github.OutputPath, paths.DefaultApplicationsOutputDir), nil
}
