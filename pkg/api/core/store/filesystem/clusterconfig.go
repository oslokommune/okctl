package filesystem

import (
	"path"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/spf13/afero"
	"sigs.k8s.io/yaml"
)

type clusterConfig struct {
	clusterConfigFileName string
	baseDir               string
	fs                    *afero.Afero
}

// SaveClusterConfig stores the cluster config
func (c *clusterConfig) SaveClusterConfig(config *api.ClusterConfig) error {
	data, err := config.YAML()
	if err != nil {
		return err
	}

	err = c.fs.MkdirAll(c.baseDir, 0o744)
	if err != nil {
		return err
	}

	return c.fs.WriteFile(path.Join(c.baseDir, c.clusterConfigFileName), data, 0o644)
}

// DeleteClusterConfig deletes a cluster config
func (c *clusterConfig) DeleteClusterConfig(env string) error {
	return c.fs.Remove(path.Join(c.baseDir, c.clusterConfigFileName))
}

// GetClusterConfig returns a stored cluster config
func (c *clusterConfig) GetClusterConfig(env string) (*api.ClusterConfig, error) {
	data, err := c.fs.ReadFile(path.Join(c.baseDir, c.clusterConfigFileName))
	if err != nil {
		return nil, err
	}

	cfg := &api.ClusterConfig{}

	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

// NewClusterConfigStore returns an instantiated cluster config store
func NewClusterConfigStore(clusterConfigFileName, baseDir string, fs *afero.Afero) api.ClusterConfigStore {
	return &clusterConfig{
		clusterConfigFileName: clusterConfigFileName,
		baseDir:               baseDir,
		fs:                    fs,
	}
}
