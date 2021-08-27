package originalversion_test

import (
	"testing"

	"github.com/oslokommune/okctl/pkg/upgrade/originalversion"
	"github.com/oslokommune/okctl/pkg/upgrade/testutils"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//nolint:funlen
func TestOriginalVersionSaver(t *testing.T) {
	testCases := []struct {
		name                   string
		existingClusterVersion string
		expectedSavedVersion   string
	}{
		{
			name:                   "Should store cluster tag version",
			existingClusterVersion: "0.0.50",
			expectedSavedVersion:   "0.0.50",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// Given
			upgradeState := testutils.MockUpgradeState()
			clusterState := mockClusterState(tc.existingClusterVersion)

			saver, err := originalversion.New(
				api.ID{ClusterName: "my-cluster"},
				upgradeState,
				clusterState,
			)
			require.NoError(t, err)

			// When
			err = saver.SaveOriginalClusterVersionIfNotExists()

			// Then
			require.NoError(t, err)

			savedVersion, err := upgradeState.GetOriginalClusterVersion()
			assert.NoError(t, err)

			assert.Equal(t, tc.expectedSavedVersion, savedVersion.Value)
		})
	}
}

type clusterStateMock struct {
	clusterVersion string
}

//goland:noinspection GoUnusedParameter
func (c clusterStateMock) SaveCluster(cluster *client.Cluster) error {
	panic("implement me")
}

//goland:noinspection GoUnusedParameter
func (c clusterStateMock) GetCluster(name string) (*client.Cluster, error) {
	return &client.Cluster{
		Config: &v1alpha5.ClusterConfig{
			Metadata: v1alpha5.ClusterMeta{
				Tags: map[string]string{
					v1alpha1.OkctlVersionTag: c.clusterVersion,
				},
			},
		},
	}, nil
}

//goland:noinspection GoUnusedParameter
func (c clusterStateMock) RemoveCluster(name string) error {
	panic("implement me")
}

//goland:noinspection GoUnusedParameter
func (c clusterStateMock) HasCluster(name string) (bool, error) {
	panic("implement me")
}

func mockClusterState(clusterVersion string) client.ClusterState {
	return &clusterStateMock{
		clusterVersion: clusterVersion,
	}
}
