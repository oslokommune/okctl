package upgrade

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFileName(t *testing.T) {
	testCases := []struct {
		name         string
		withFilename string
		expectFile   okctlUpgradeFile
		expectError  bool
	}{
		{
			name:         "Should parse filename",
			withFilename: "okctl-upgrade_0.0.63_Darwin_amd64.tar.gz",
			expectFile: okctlUpgradeFile{
				filename:  "okctl-upgrade_0.0.63_Darwin_amd64.tar.gz",
				version:   "0.0.63",
				os:        "Darwin",
				arch:      "amd64",
				extension: "tar.gz",
			},
		},
		{
			name:         "Should return error if too few underscores in filename",
			withFilename: "okctl-upgrade_0.0.63Darwin_amd64.tar.gz",
			expectError:  true,
		},
		{
			name:         "Should return error if too many underscores in filename",
			withFilename: "okctl-upgrade_0.0.63_Dar_win_amd64.tar.gz",
			expectError:  true,
		},
		{
			name:         "Should return error if too few dots in filename",
			withFilename: "okctl-upgrade_0.0.63_Darwin_amd64targz",
			expectError:  true,
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			upgradeFile, err := parseOkctlUpgradeFilename(tc.withFilename)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				assert.Equal(t, tc.expectFile.filename, upgradeFile.filename)
				assert.Equal(t, tc.expectFile.version, upgradeFile.version)
				assert.Equal(t, tc.expectFile.os, upgradeFile.os)
				assert.Equal(t, tc.expectFile.arch, upgradeFile.arch)
				assert.Equal(t, tc.expectFile.extension, upgradeFile.extension)
			}
		})
	}
}
