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
	repoStatePaths Paths
	clusterConfig  Paths
	fs             *afero.Afero
	repoState      *repository.Data
}

// SaveCluster knows how to save cluster state
func (c *cluster) SaveCluster(clu *api.Cluster) error {
	for _, cluster := range c.repoState.Clusters {
		if cluster.Environment == clu.Environment {
			return nil
		}
	}

	c.repoState.Clusters = append(c.repoState.Clusters, repository.Cluster{
		Environment: clu.Environment,
		AWS: repository.AWS{
			AccountID: clu.AWSAccountID,
			Cidr:      clu.Cidr,
		},
	})

	err := c.updateRepoStateFile()
	if err != nil {
		return fmt.Errorf("failed to update cluster state: %w", err)
	}

	data, err := clu.Config.YAML()
	if err != nil {
		return err
	}

	err = c.fs.MkdirAll(c.clusterConfig.BaseDir, 0o744)
	if err != nil {
		return err
	}

	return c.fs.WriteFile(path.Join(c.clusterConfig.BaseDir, c.clusterConfig.ConfigFile), data, 0o644)
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

	err := c.updateRepoStateFile()
	if err != nil {
		return fmt.Errorf("failed to delete cluster state: %w", err)
	}

	return c.fs.Remove(path.Join(c.clusterConfig.BaseDir, c.clusterConfig.ConfigFile))
}

func (c *cluster) updateRepoStateFile() error {
	err := c.fs.MkdirAll(c.repoStatePaths.BaseDir, 0o744)
	if err != nil {
		return err
	}

	data, err := c.repoState.YAML()
	if err != nil {
		return err
	}

	return c.fs.WriteFile(path.Join(c.repoStatePaths.BaseDir, c.repoStatePaths.ConfigFile), data, 0o644)
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
func NewClusterStore(repoStatePaths, clusterConfig Paths, fs *afero.Afero, repoState *repository.Data) api.ClusterStore {
	return &cluster{
		repoStatePaths: repoStatePaths,
		clusterConfig:  clusterConfig,
		fs:             fs,
		repoState:      repoState,
	}
}
