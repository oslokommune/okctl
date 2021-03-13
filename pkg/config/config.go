// Package config interacts with all configuration state
package config

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/oslokommune/okctl/pkg/config/constant"

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
	DefaultComponentBaseDir           = "component"
	DefaultPostgresBaseDir            = "postgres"
	DefaultPostgresOutputFile         = "postgres-output.json"
	DefaultPostgresCloudFormationFile = "postgres-cf.yaml"
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
		MaxSize:    constant.DefaultLogSizeInMb,
		MaxBackups: constant.DefaultLogBackups,
		MaxAge:     constant.DefaultLogDays,
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

	return filepath.Join(base, constant.DefaultRepositoryStateFile), nil
}

// WriteCurrentRepoState will write current repo state to disk
func (c *Config) WriteCurrentRepoState() error {
	repoDir, err := c.GetRepoDir()
	if err != nil {
		return err
	}

	_, err = store.NewFileSystem(repoDir, c.FileSystem).
		StoreStruct(constant.DefaultRepositoryStateFile, c.RepoState, store.ToYAML()).
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
// application data should be written
func (c *Config) GetUserDataDir() (string, error) {
	home, err := c.GetHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, constant.DefaultDir), nil
}

// GetUserDataPath returns the path to the okctl application
// config path
func (c *Config) GetUserDataPath() (string, error) {
	base, err := c.GetUserDataDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(base, constant.DefaultConfig), nil
}

// GetLogName returns the path to a logfile
func (c *Config) GetLogName() (string, error) {
	base, err := c.GetUserDataDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(base, constant.DefaultLogDir, constant.DefaultLogName), nil
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

// GetRepoApplicatiosOutputDir returns the directory where application
// resources are stored
func (c *Config) GetRepoApplicatiosOutputDir() (string, error) {
	base, err := c.GetRepoStateDir()
	if err != nil {
		return "", err
	}

	return path.Join(base, c.RepoState.Metadata.OutputDir, constant.DefaultApplicationsOutputDir), nil
}
