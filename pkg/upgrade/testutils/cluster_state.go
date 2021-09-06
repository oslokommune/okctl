package testutils

import "github.com/oslokommune/okctl/pkg/client"

type clusterStateMock struct {
}

//goland:noinspection GoUnusedParameter
func (c clusterStateMock) SaveCluster(cluster *client.Cluster) error {
	panic("implement me")
}

//goland:noinspection GoUnusedParameter
func (c clusterStateMock) GetCluster(name string) (*client.Cluster, error) {
	panic("implement me")
}

//goland:noinspection GoUnusedParameter
func (c clusterStateMock) RemoveCluster(name string) error {
	panic("implement me")
}

//goland:noinspection GoUnusedParameter
func (c clusterStateMock) HasCluster(name string) (bool, error) {
	panic("implement me")
}

// MockClusterState returns a mocked upgrade state
func MockClusterState() client.ClusterState {
	return &clusterStateMock{}
}
