package config

import (
	"path"
	"runtime"
	"strings"

	"github.com/spf13/viper"
)

const EnvPrefix = "OKCTL"

// LoadUserConfiguration will load the configuration related
// to interacting with a repository.
func LoadUserConfiguration(baseDir string) (*UserConfig, error) {
	v := viper.New()
	v.SetConfigName("user")
	v.SetConfigType("yml")
	v.SetEnvPrefix(EnvPrefix)
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AddConfigPath(path.Join(baseDir, ".okctl"))

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	cfg := &UserConfig{}

	err := v.Unmarshal(cfg)
	if err != nil {
		return nil, err
	}

	cfg.Host = Host{
		Os:   runtime.GOOS,
		Arch: runtime.GOARCH,
	}

	return cfg, nil
}

type UserConfig struct {
	Host
	Binaries []Binary
}

type Host struct {
	Os   string
	Arch string
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
