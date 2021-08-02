package upgrade

import (
	"testing"

	semverPkg "github.com/Masterminds/semver"
	"github.com/stretchr/testify/assert"
)

func TestName(t *testing.T) {
	testCases := []struct {
		name        string
		version     string
		assert      func(t *testing.T, version upgradeBinaryVersion)
		expectError bool
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
			name:        "Should return error for invalid inputr",
			version:     "0.0.1.some.hotfix",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			version, err := newVersion(tc.version)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				tc.assert(t, version)
			}
		})
	}
}
