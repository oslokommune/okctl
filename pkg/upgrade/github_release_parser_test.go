package upgrade_test

import (
	"fmt"
	"testing"

	"github.com/oslokommune/okctl/pkg/upgrade"

	"github.com/oslokommune/okctl/pkg/github"
	"github.com/oslokommune/okctl/pkg/osarch"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//nolint:funlen
func TestGithubReleaseParser(t *testing.T) {
	testCases := []struct {
		name         string
		releases     []*github.RepositoryRelease
		checksumFile string
		expectError  string
	}{
		{
			name:         "Should work when everything is OK",
			releases:     createGithubReleases([]string{osarch.Linux, osarch.Darwin}, osarch.Amd64, []string{"0.0.61"}),
			checksumFile: fmt.Sprintf("%s/0.0.61/%s", folderWorking, upgrade.ChecksumsTxt),
		},
		{
			name: "Should return unknown release name in error message for at least 2 assets",
			releases: []*github.RepositoryRelease{
				{
					TagName: github.StringPtr("0.0.61"),
					Assets: []*github.ReleaseAsset{
						{},
					},
				},
			},
			expectError: "release '(unknown)' must have at least 2 assets (binary and checksum)",
		},
		{
			name: "Should return release ID in error message for at least 2 assets",
			releases: []*github.RepositoryRelease{
				{
					ID:      github.Int64Ptr(1),
					TagName: github.StringPtr("0.0.61"),
					Assets: []*github.ReleaseAsset{
						{},
					},
				},
			},
			expectError: "release 'ID: 1' must have at least 2 assets (binary and checksum)",
		},
		{
			name: "Should return release name in error message for at least 2 assets",
			releases: []*github.RepositoryRelease{
				{
					ID:      github.Int64Ptr(1),
					TagName: github.StringPtr("0.0.61"),
					Name:    github.StringPtr("0.0.61"),
					Assets: []*github.ReleaseAsset{
						{},
					},
				},
			},
			expectError: "release '0.0.61' must have at least 2 assets (binary and checksum)",
		},
		{
			name:         "Should return error about invalid substrings in checksum file",
			releases:     createGithubReleases([]string{osarch.Linux, osarch.Darwin}, osarch.Amd64, []string{"0.0.61"}),
			checksumFile: "github_release_parser/okctl-upgrade-checksums-substring-error.txt",
			expectError: "expected 2 substrings when splitting digest line on whitespace, got 3 in string " +
				"'716c5b3f517c197bdb748a55983b7a6de9045a66c759ee3f10863d19bbf90a61  " +
				"okctl-upgrade_0.0.61_Darwin_amd64.tar.gz an_extra_substring'",
		},
		{
			name:         "Should return error about invalid filename in checksum file",
			releases:     createGithubReleases([]string{osarch.Linux, osarch.Darwin}, osarch.Amd64, []string{"0.0.61"}),
			checksumFile: "github_release_parser/okctl-upgrade-checksums-filename-error.txt",
			expectError: "expected 4 substrings when splitting on underscore (_), got 5 in string " +
				"'okctl-upgrade_0.0.61_Darwin_amd64_invalidstuff.tar.gz'",
		},
		{
			name:         "Should return error about too many dots in checksum file",
			releases:     createGithubReleases([]string{osarch.Linux, osarch.Darwin}, osarch.Amd64, []string{"0.0.61"}),
			checksumFile: "github_release_parser/okctl-upgrade-checksums-osarch-error.txt",
			expectError: "expected at least 2 substrings when splitting on dot (.), got 1 in string " +
				"'okctl-upgrade_0.0.61_Darwin_amd64'",
		},
		{
			name: "Should return release name in error message for at least 2 assets",
			releases: []*github.RepositoryRelease{
				{
					ID:      github.Int64Ptr(1),
					TagName: github.StringPtr("0.0.61"),
					Name:    github.StringPtr("0.0.61"),
					Assets: []*github.ReleaseAsset{
						createGihubReleaseAssetBinary(osarch.Linux, osarch.Amd64, "0.0.55"),
						createGihubReleaseAssetBinary(osarch.Linux, osarch.Amd64, "0.0.55"),
					},
				},
			},
			expectError: "expected okctl upgrade binary version '0.0.55' to equal release upgrade version (i.e. tag) '0.0.61'",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// Given
			var err error

			checksumFileContents := []byte("")
			if len(tc.checksumFile) > 0 {
				checksumFileContents, err = readBytesFromFile(fmt.Sprintf("testdata/%s", tc.checksumFile))
				require.NoError(t, err)
			}

			parser := upgrade.NewGithubReleaseParser(NewMockChecksumDownloader(checksumFileContents))

			// When
			_, err = parser.ToUpgradeBinaries(tc.releases)

			// Then
			if len(tc.expectError) > 0 {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

type MockChecksumDownloader struct {
	checksumFileContents []byte
}

//goland:noinspection GoUnusedParameter
func (m MockChecksumDownloader) Download(checksumAsset *github.ReleaseAsset) ([]byte, error) {
	return m.checksumFileContents, nil
}

func NewMockChecksumDownloader(checksumDownloadResult []byte) upgrade.ChecksumDownloader {
	return MockChecksumDownloader{
		checksumFileContents: checksumDownloadResult,
	}
}
