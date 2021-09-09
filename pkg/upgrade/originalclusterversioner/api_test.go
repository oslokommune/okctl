package originalclusterversioner_test

import (
	"github.com/oslokommune/okctl/pkg/upgrade/originalclusterversioner"
	"testing"

	"github.com/oslokommune/okctl/pkg/upgrade/testutils"

	"github.com/oslokommune/okctl/pkg/api"
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
			upgradeState := testutils.MockUpgradeState(tc.existingClusterVersion)
			clusterState := testutils.MockClusterState(tc.existingClusterVersion)

			versioner := originalclusterversioner.New(
				api.ID{ClusterName: "my-cluster"},
				upgradeState,
				clusterState,
			)

			// When
			err := versioner.SaveOriginalClusterVersionFromClusterTagIfNotExists()

			// Then
			require.NoError(t, err)

			savedVersion, err := upgradeState.GetOriginalClusterVersion()
			assert.NoError(t, err)

			assert.Equal(t, tc.expectedSavedVersion, savedVersion.Value)
		})
	}
}
