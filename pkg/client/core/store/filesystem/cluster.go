// Package filesystem implements the store layer using the filesystem for persistence
package filesystem

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/client"

	"github.com/oslokommune/okctl/pkg/client/store"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/config/repository"
	"github.com/spf13/afero"
)

type clusterStore struct {
	repoStatePaths Paths
	clusterConfig  Paths
	fs             *afero.Afero
	repoState      *repository.Data
}

func (s *clusterStore) SaveCluster(c *api.Cluster) (*store.Report, error) {
	if _, ok := s.repoState.Clusters[c.ID.Environment]; ok {
		return nil, nil
	}

	s.repoState.Clusters[c.ID.Environment] = &repository.Cluster{
		Name:         c.ID.ClusterName,
		Environment:  c.ID.Environment,
		AWSAccountID: c.ID.AWSAccountID,
		VPC: &repository.VPC{
			VpcID: "",
			CIDR:  c.Cidr,
		},
	}

	report, err := store.NewFileSystem(s.repoStatePaths.BaseDir, s.fs).
		StoreStruct(s.repoStatePaths.ConfigFile, s.repoState, store.ToYAML()).
		AlterStore(store.SetBaseDir(s.clusterConfig.BaseDir)).
		StoreStruct(s.clusterConfig.ConfigFile, c.Config, store.ToYAML()).
		Do()
	if err != nil {
		return nil, fmt.Errorf("failed to store cluster: %w", err)
	}

	return report, nil
}

// DeleteCluster knows how to delete clusterStore state
func (s *clusterStore) DeleteCluster(id api.ID) (*store.Report, error) {
	delete(s.repoState.Clusters, id.Environment)

	report, err := store.NewFileSystem(s.repoStatePaths.BaseDir, s.fs).
		StoreStruct(s.repoStatePaths.ConfigFile, s.repoState, store.ToYAML()).
		AlterStore(store.SetBaseDir(s.clusterConfig.BaseDir)).
		Remove(s.clusterConfig.ConfigFile).
		Do()
	if err != nil {
		return nil, fmt.Errorf("failed to remove cluster from storage: %w", err)
	}

	return report, nil
}

// GetCluster knows how to get clusterStore state
func (s *clusterStore) GetCluster(id api.ID) (*api.Cluster, error) {
	if c, ok := s.repoState.Clusters[id.Environment]; ok {
		return &api.Cluster{
			ID:   id,
			Cidr: c.VPC.CIDR,
		}, nil
	}

	return nil, fmt.Errorf("failed to find cluster %s", id.ClusterName)
}

// NewClusterStore returns a store for clusterStore
func NewClusterStore(repoStatePaths, clusterConfig Paths, fs *afero.Afero, repoState *repository.Data) client.ClusterStore {
	return &clusterStore{
		repoStatePaths: repoStatePaths,
		clusterConfig:  clusterConfig,
		fs:             fs,
		repoState:      repoState,
	}
}
