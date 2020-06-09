package config

import (
	"github.com/spf13/viper"
)

const (
	DefaultRepositoryConfig     = ".okctl.yml"
	DefaultRepositoryConfigName = ".okctl"
	DefaultRepositoryConfigType = "yml"
)

type RepoConfig struct {
	Name     string
	Region   string
	BaseDir  string
	Clusters []Cluster
}

type Cluster struct {
	Name string
	AWS  AWS
}

type AWS struct {
	Account int
}

// LoadRepo reads in the configuration of a repository
// from the provided baseDir.
func LoadRepo(baseDir string) (*RepoConfig, error) {
	v := viper.New()
	v.SetConfigName(DefaultRepositoryConfigName)
	v.SetConfigType(DefaultRepositoryConfigType)
	v.AddConfigPath(baseDir)

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	cfg := &RepoConfig{}

	err := v.Unmarshal(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
