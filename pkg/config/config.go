package config

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-yaml/yaml"
	"github.com/google/uuid"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

const (
	DefaultAppDir        = ".okctl"
	DefaultAppConfig     = "conf.yml"
	DefaultAppConfigName = "conf"
	DefaultAppConfigType = "yml"
	DefaultAppEnvPrefix  = "OKCTL"
	DefaultAppEnvHome    = "OKCTL_HOME"
)

type AppConfig struct {
	BaseDir  string
	User     User
	Host     Host
	Binaries []Binary
}

type User struct {
	ID       string
	Username string
}

type Host struct {
	Os   string
	Arch string
}

const (
	OsDarwin = "darwin"
	OsLinux  = "linux"

	ArchAmd64 = "amd64"
)

// host ensures that the operating system and architecture of the host machine
// are supported.
func (h Host) Valid() error {
	return validation.ValidateStruct(&h,
		validation.Field(&h.Arch,
			validation.Required,
			validation.In(
				ArchAmd64,
			),
		),
		validation.Field(&h.Os,
			validation.Required,
			validation.In(
				OsDarwin,
				OsLinux,
			),
		),
	)
}

type Binary struct {
	Name       string
	Version    string
	BufferSize string
	URLPattern string
	Archive    Archive
	Checksums  []Checksum
}

type Archive struct {
	Type   string
	Target string
}

type Checksum struct {
	Os     string
	Arch   string
	Type   string
	Digest string
}

type AppCfgFn func() *AppConfig

// NewDefaultAppCfg is the default configuration
// of the okctl CLI.
func NewDefaultAppCfg() *AppConfig {
	return &AppConfig{
		User: User{
			ID: uuid.New().String(),
		},
		Host: Host{
			Os:   runtime.GOOS,
			Arch: runtime.GOARCH,
		},
		Binaries: []Binary{
			{
				Name:       "eksctl",
				Version:    "0.18.0",
				BufferSize: "100mb",
				URLPattern: "https://github.com/weaveworks/eksctl/releases/download/#{ver}/eksctl_#{os}_#{arch}.tar.gz",
				Archive: Archive{
					Type:   ".tar.gz",
					Target: "eksctl",
				},
				Checksums: []Checksum{
					{
						Os:     "darwin",
						Arch:   "amd64",
						Type:   "sha256",
						Digest: "fc320f3e52adef9f8d06a98f1996801ee3b59d1d74bac11e24123f593875a344",
					},
					{
						Os:     "linux",
						Arch:   "amd64",
						Type:   "sha256",
						Digest: "a8f83394a12051bd6bf539dca7db2005237d36c6b1a67073bcf2070d034356f0",
					},
				},
			},
		},
	}
}

// NewDefaultAppCfgFn provides a function that can be
// invoked to get a default configuration of okctl.
func NewDefaultAppCfgFn() AppCfgFn {
	return NewDefaultAppCfg
}

type AppNotFoundFn func(baseDir, appDir, config string) error

func NewDefaultNotFoundFn() AppNotFoundFn {
	return func(baseDir, appDir, config string) error {
		return fmt.Errorf("failed to load your configuration: %s, run: 'configure' command first", path.Join(baseDir, appDir, config))
	}
}

// GetHomeDir will attempt to determine what location
// we should read the okctl CLI configuration from.
func GetHomeDir() (string, error) {
	homeDir := os.Getenv(DefaultAppEnvHome)

	if len(homeDir) == 0 {
		dir, err := homedir.Dir()
		if err != nil {
			return "", err
		}

		homeDir = dir
	}

	return homeDir, nil
}

// LoadApp will load the okctl CLI configuration.
func LoadApp(baseDir string, defaultAppCfgFn AppCfgFn, notFoundFn AppNotFoundFn) (*AppConfig, error) {
	v := viper.New()

	b, err := yaml.Marshal(defaultAppCfgFn())
	if err != nil {
		return nil, err
	}

	defaultConfig := bytes.NewReader(b)

	v.SetConfigType(DefaultAppConfigType)

	err = v.MergeConfig(defaultConfig)
	if err != nil {
		return nil, err
	}

	v.AddConfigPath(path.Join(baseDir, DefaultAppDir))
	v.SetConfigName(DefaultAppConfigName)

	err = v.MergeInConfig()
	if err != nil {
		switch err.(type) {
		case viper.ConfigFileNotFoundError:
			err := notFoundFn(baseDir, DefaultAppDir, DefaultAppConfig)
			if err != nil {
				return nil, err
			}
		default:
			return nil, err
		}
	}

	v.AutomaticEnv()
	v.SetEnvPrefix(DefaultAppEnvPrefix)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	cfg := &AppConfig{}

	err = v.Unmarshal(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
