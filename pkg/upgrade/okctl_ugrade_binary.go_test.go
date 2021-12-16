package upgrade

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSort(t *testing.T) {
	testCases := []struct {
		name                string
		withUpgradeBinaries []okctlUpgradeBinary
		expect              []okctlUpgradeBinary
	}{
		{
			name:                "Should sort by semver",
			withUpgradeBinaries: createUpgradeBinaries(t, []string{"0.0.3", "0.0.2", "0.0.1"}),
			expect:              createUpgradeBinaries(t, []string{"0.0.1", "0.0.2", "0.0.3"}),
		},
		{
			name:                "Test case sort by semver and hotfix",
			withUpgradeBinaries: createUpgradeBinaries(t, []string{"0.0.2.b", "0.0.3", "0.0.20", "0.0.2", "0.0.2.a", "0.0.1"}),
			expect:              createUpgradeBinaries(t, []string{"0.0.1", "0.0.2", "0.0.2.a", "0.0.2.b", "0.0.3", "0.0.20"}),
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			sort(tc.withUpgradeBinaries)

			expected := make([]string, len(tc.expect))
			actual := make([]string, len(tc.expect))

			for i := 0; i < len(tc.expect); i++ {
				expected[i] = tc.expect[i].RawVersion()
				actual[i] = tc.withUpgradeBinaries[i].RawVersion()
			}

			assert.Equal(t, expected, actual)
		})
	}
}

func createUpgradeBinaries(t *testing.T, versions []string) []okctlUpgradeBinary {
	binaries := make([]okctlUpgradeBinary, 0, len(versions))

	for _, versionString := range versions {
		version, err := parseUpgradeBinaryVersion(versionString)
		assert.NoError(t, err)

		b := newOkctlUpgradeBinary(version, nil, "")
		binaries = append(binaries, b)
	}

	return binaries
}
