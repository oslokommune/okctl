// Package filesystem implements the store layer using the filesystem for persistence
package filesystem

import (
	"fmt"
	"path"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/config/repository"
	"github.com/spf13/afero"
)

type cluster struct {
	repoConfigFileName string
	repoBaseDir        string
	fs                 *afero.Afero
	repoState          *repository.Data
}

// SaveCluster knows how to save cluster state
func (c *cluster) SaveCluster(clu *api.Cluster) error {
	c.repoState.Clusters = append(c.repoState.Clusters, repository.Cluster{
		Environment: clu.Environment,
		AWS: repository.AWS{
			AccountID: clu.AWSAccountID,
			Cidr:      clu.Cidr,
		},
	})

	return c.updateConfigFile()
}

// DeleteCluster knows how to delete cluster state
func (c *cluster) DeleteCluster(env string) error {
	s := c.repoState

	for i, curr := range s.Clusters {
		if curr.Environment == env {
			s.Clusters[i] = s.Clusters[len(s.Clusters)-1]
			s.Clusters[len(s.Clusters)-1] = repository.Cluster{}
			s.Clusters = s.Clusters[:len(s.Clusters)-1]
		}
	}

	return c.updateConfigFile()
}

func (c *cluster) updateConfigFile() error {
	err := c.fs.MkdirAll(c.repoBaseDir, 0744)
	if err != nil {
		return err
	}

	data, err := c.repoState.YAML()
	if err != nil {
		return err
	}

	return c.fs.WriteFile(path.Join(c.repoBaseDir, c.repoConfigFileName), data, 0644)
}

// GetCluster knows how to get cluster state
func (c *cluster) GetCluster(env string) (*api.Cluster, error) {
	clusters := c.repoState.Clusters

	for _, cluster := range clusters {
		if cluster.Environment == env {
			return &api.Cluster{
				Environment:  cluster.Environment,
				AWSAccountID: cluster.AWS.AccountID,
				Cidr:         cluster.AWS.Cidr,
			}, nil
		}
	}

	return nil, fmt.Errorf("failed to find cluster configuration for env: %s", env)
}

// NewClusterStore returns a store for cluster
func NewClusterStore(repoConfigFileName, repoBaseDir string, fs *afero.Afero, repoState *repository.Data) api.ClusterStore {
	return &cluster{
		repoConfigFileName: repoConfigFileName,
		repoBaseDir:        repoBaseDir,
		fs:                 fs,
		repoState:          repoState,
	}
}
