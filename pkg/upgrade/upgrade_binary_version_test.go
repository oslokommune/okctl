package upgrade

import (
	"testing"

	semverPkg "github.com/Masterminds/semver"
	"github.com/stretchr/testify/assert"
)

//nolint:funlen
func TestParseUpgradeBinaryVersion(t *testing.T) {
	testCases := []struct {
		name                string
		version             string
		assert              func(t *testing.T, version upgradeBinaryVersion)
		expectErrorContains string
	}{
		{
			name:    "Should parse regular version",
			version: "0.0.1",
			assert: func(t *testing.T, actual upgradeBinaryVersion) {
				semver, err := semverPkg.NewVersion("0.0.1")
				assert.NoError(t, err)

				assert.Equal(t, upgradeBinaryVersion{
					raw:    "0.0.1",
					semver: semver,
					hotfix: "",
				}, actual)
			},
		},
		{
			name:    "Should parse hotfix version",
			version: "0.0.1.some-hotfix",
			assert: func(t *testing.T, actual upgradeBinaryVersion) {
				semver, err := semverPkg.NewVersion("0.0.1")
				assert.NoError(t, err)

				assert.Equal(t, upgradeBinaryVersion{
					raw:    "0.0.1.some-hotfix",
					semver: semver,
					hotfix: "some-hotfix",
				}, actual)
			},
		},
		{
			name:                "Should return error for invalid semantic version",
			version:             "0.0.a",
			expectErrorContains: "parsing to semantic version from '0.0.a'",
		},
		{
			name:                "Should return error for hotfix version",
			version:             "0.0.a.some-hotfix",
			expectErrorContains: "parsing to semantic with hotfix version from '0.0.a.some-hotfix'",
		},
		{
			name:                "Should return error for too many dots",
			version:             "0.0.1.some.hotfix",
			expectErrorContains: "not a valid version: 0.0.1.some.hotfix",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			version, err := parseUpgradeBinaryVersion(tc.version)

			if len(tc.expectErrorContains) > 0 {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectErrorContains)
			} else {
				assert.NoError(t, err)
				tc.assert(t, version)
			}
		})
	}
}
