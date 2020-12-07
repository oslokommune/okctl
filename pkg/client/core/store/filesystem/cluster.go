// Package filesystem implements the store layer using the filesystem for persistence
package filesystem

import (
	"github.com/oslokommune/okctl/pkg/client"

	"github.com/oslokommune/okctl/pkg/client/store"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/spf13/afero"
)

type clusterStore struct {
	paths Paths
	fs    *afero.Afero
}

func (s *clusterStore) SaveCluster(c *api.Cluster) (*store.Report, error) {
	report, err := store.NewFileSystem(s.paths.BaseDir, s.fs).
		StoreStruct(s.paths.ConfigFile, c.Config, store.ToYAML()).
		Do()
	if err != nil {
		return nil, err
	}

	return report, nil
}

// DeleteCluster knows how to delete clusterStore state
func (s *clusterStore) DeleteCluster(_ api.ID) (*store.Report, error) {
	report, err := store.NewFileSystem(s.paths.BaseDir, s.fs).
		Remove(s.paths.ConfigFile).
		Remove("").
		Do()
	if err != nil {
		return nil, err
	}

	return report, nil
}

// GetCluster knows how to get clusterStore state
func (s *clusterStore) GetCluster(id api.ID) (*api.Cluster, error) {
	cluster := &api.Cluster{
		ID: id,
	}

	_, err := store.NewFileSystem(s.paths.BaseDir, s.fs).
		GetStruct(s.paths.ConfigFile, cluster, store.FromYAML()).
		Do()
	if err != nil {
		return nil, err
	}

	return cluster, nil
}

// NewClusterStore returns a store for clusterStore
func NewClusterStore(paths Paths, fs *afero.Afero) client.ClusterStore {
	return &clusterStore{
		paths: paths,
		fs:    fs,
	}
}
