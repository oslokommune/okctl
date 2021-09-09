package testutils

import (
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/client"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type clusterStateMock struct {
	cluster *client.Cluster
}

//goland:noinspection GoUnusedParameter
func (c clusterStateMock) SaveCluster(cluster *client.Cluster) error {
	panic("implement me")
}

//goland:noinspection GoUnusedParameter
func (c clusterStateMock) GetCluster(name string) (*client.Cluster, error) {
	return c.cluster, nil
}

//goland:noinspection GoUnusedParameter
func (c clusterStateMock) RemoveCluster(name string) error {
	panic("implement me")
}

//goland:noinspection GoUnusedParameter
func (c clusterStateMock) HasCluster(name string) (bool, error) {
	panic("implement me")
}

// MockClusterState returns a mocked upgrade state with the given cluster version
func MockClusterState(clusterVersion string) client.ClusterState {
	cluster := createCluster()

	cluster.Config.Metadata.Tags[v1alpha1.OkctlVersionTag] = clusterVersion

	return &clusterStateMock{
		cluster: cluster,
	}
}

func createCluster() *client.Cluster {
	return &client.Cluster{
		ID:   api.ID{},
		Name: "someCluster",
		Config: &v1alpha5.ClusterConfig{
			TypeMeta: metav1.TypeMeta{},
			Metadata: v1alpha5.ClusterMeta{
				Tags: make(map[string]string),
			},
		},
	}
}
