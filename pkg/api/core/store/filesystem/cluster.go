// Package filesystem implements the store layer using the filesystem for persistence
package filesystem

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/client/store"

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
		if cluster.Environment == clu.ID.Environment {
			return nil
		}
	}

	c.repoState.Clusters = append(c.repoState.Clusters, repository.Cluster{
		Name:        clu.ID.ClusterName,
		Environment: clu.ID.Environment,
		AWS: repository.AWS{
			AccountID: clu.ID.AWSAccountID,
			Cidr:      clu.Cidr,
		},
	})

	_, err := store.NewFileSystem(c.repoStatePaths.BaseDir, c.fs).
		StoreStruct(c.repoStatePaths.ConfigFile, c.repoState, store.ToYAML()).
		AlterStore(store.SetBaseDir(c.clusterConfig.BaseDir)).
		StoreStruct(c.clusterConfig.ConfigFile, clu.Config, store.ToYAML()).
		Do()
	if err != nil {
		return fmt.Errorf("failed to store cluster: %w", err)
	}

	return nil
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

	_, err := store.NewFileSystem(c.repoStatePaths.BaseDir, c.fs).
		StoreStruct(c.repoStatePaths.ConfigFile, c.repoState, store.ToYAML()).
		AlterStore(store.SetBaseDir(c.clusterConfig.BaseDir)).
		Remove(c.clusterConfig.ConfigFile).
		Do()
	if err != nil {
		return fmt.Errorf("failed to remove cluster from storage: %w", err)
	}

	return nil
}

// GetCluster knows how to get cluster state
func (c *cluster) GetCluster(env string) (*api.Cluster, error) {
	clusters := c.repoState.Clusters

	for _, cluster := range clusters {
		if cluster.Environment == env {
			return &api.Cluster{
				ID: api.ID{
					Region:       c.repoState.Region,
					AWSAccountID: cluster.AWS.AccountID,
					Environment:  cluster.Environment,
					Repository:   c.repoState.Name,
					ClusterName:  cluster.Name,
				},
				Cidr: cluster.AWS.Cidr,
			}, nil
		}
	}

	return nil, fmt.Errorf("failed to find cluster for env: %s", env)
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
