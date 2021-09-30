package originalclusterversion_test

import (
	"testing"

	"github.com/oslokommune/okctl/pkg/upgrade/originalclusterversion"

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
		expectedError          string
	}{
		{
			name:                   "Should store cluster tag version",
			existingClusterVersion: "0.0.50",
			expectedSavedVersion:   "0.0.50",
		},
		{
			name:                   "Should return error if cluster tag version is not in semver format",
			existingClusterVersion: "dev",
			expectedError:          "getting cluster state version: parsing version 'dev': Invalid Semantic Version",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// Given
			upgradeState := testutils.MockUpgradeState(tc.existingClusterVersion)
			clusterState := testutils.MockClusterState(tc.existingClusterVersion)

			versioner := originalclusterversion.New(api.ID{ClusterName: "my-cluster"}, upgradeState, clusterState)

			// When
			err := versioner.SaveOriginalClusterVersionFromClusterTagIfNotExists()

			// Then
			if len(tc.expectedError) > 0 {
				assert.Contains(t, tc.expectedError, err.Error())
				return
			}

			require.NoError(t, err)

			savedVersion, err := upgradeState.GetOriginalClusterVersion()
			assert.NoError(t, err)

			assert.Equal(t, tc.expectedSavedVersion, savedVersion.Value)
		})
	}
}
