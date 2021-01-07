package noop

import "github.com/oslokommune/okctl/pkg/api"

type clusterStore struct{}

func (c *clusterStore) SaveCluster(*api.Cluster) error {
	return nil
}

func (c *clusterStore) DeleteCluster(api.ID) error {
	return nil
}

func (c *clusterStore) GetCluster(api.ID) (*api.Cluster, error) {
	return nil, nil
}

// NewClusterStore returns a no operation cluster store
func NewClusterStore() api.ClusterStore {
	return &clusterStore{}
}
