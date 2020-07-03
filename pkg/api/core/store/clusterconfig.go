package store

import (
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/storage/state"
	"sigs.k8s.io/yaml"
)

type clusterConfig struct {
	provider state.PersisterProvider
}

// SaveClusterConfig stores the cluster config
func (c *clusterConfig) SaveClusterConfig(config *api.ClusterConfig) error {
	data, err := config.YAML()
	if err != nil {
		return err
	}

	return c.provider.Repository().WriteToDefault("cluster_config", data)
}

// DeleteClusterConfig deletes a cluster config
func (c *clusterConfig) DeleteClusterConfig(env string) error {
	return c.provider.Repository().DeleteDefault("cluster_config")
}

// GetClusterConfig returns a stored cluster config
func (c *clusterConfig) GetClusterConfig(env string) (*api.ClusterConfig, error) {
	data, err := c.provider.Repository().ReadFromDefault("cluster_config")
	if err != nil {
		return nil, err
	}

	cfg := api.NewClusterConfig()

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

// NewClusterConfigStore returns an instantiated cluster config store
func NewClusterConfigStore(provider state.PersisterProvider) api.ClusterConfigStore {
	return &clusterConfig{
		provider: provider,
	}
}
