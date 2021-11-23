package clusterversion_test

import (
	"bytes"
	"github.com/oslokommune/okctl/pkg/client/mock"
	"github.com/oslokommune/okctl/pkg/upgrade/clusterversion"
	"github.com/oslokommune/okctl/pkg/upgrade/testutils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateBinaryEqualsClusterVersion(t *testing.T) {
	testCases := []struct {
		name                string
		withBinaryVersion   string
		withClusterVersion  string
		expectErrorContains string
	}{
		{
			name:                "Should fail when binary version is less than cluster version",
			withBinaryVersion:   "0.0.49",
			withClusterVersion:  "0.0.50",
			expectErrorContains: "okctl binary version must be equal to cluster version 0.0.50, but was 0.0.49",
		},
		{
			name:               "Should validate without error when binary version is equal to cluster version",
			withBinaryVersion:  "0.0.50",
			withClusterVersion: "0.0.50",
		},
		{
			name:                "Should fail when binary version is greater than cluster version",
			withBinaryVersion:   "0.0.51",
			withClusterVersion:  "0.0.50",
			expectErrorContains: "okctl binary version must be equal to cluster version 0.0.50, but was 0.0.51",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			var out bytes.Buffer
			clusterID := mock.Cluster().ID
			upgradeState := testutils.MockUpgradeState("0.0.50")
			versioner := clusterversion.New(&out, clusterID, upgradeState)

			// When
			err := versioner.ValidateBinaryEqualsClusterVersion(tc.withBinaryVersion)

			// Then
			if len(tc.expectErrorContains) > 0 {
				assert.Contains(t, err.Error(), tc.expectErrorContains)
			} else {
				assert.NoError(t, err)
			}
		})
	}

}
