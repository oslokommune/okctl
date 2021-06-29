package upgrade

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
	"testing"

	"github.com/oslokommune/okctl/pkg/osarch"

	"github.com/jarcoal/httpmock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"

	"github.com/oslokommune/okctl/pkg/config/state"
	"github.com/oslokommune/okctl/pkg/github"
	"github.com/oslokommune/okctl/pkg/storage"
	"github.com/stretchr/testify/assert"
)

// Vurder: (Egen migrator test?)
// withOriginalOkctlVersion     string
// withOkctlVersion             string
// withAlreadyAppliedMigrations []string
// expectMigrationsToBeRun      []string

// Test cases
// ------------------------------------------------
// Given these releases, ..., then these binaries should be run
// Should run a upgrade
// Should not run already applied migrations
// Should run migrations up to the current okctl version

// Should run Darwin upgrades

// Failure situations
// 		Verify invalid digests in okctl-upgrade.txt
//
// I okctl-upgrade-checksuyms.txt, endre filnavn. Da bør man få feil at ting ikke matcher. Får noe annet unyttig.
// ------------------------------------------------

//nolint:funlen
func TestRunUpgrades(t *testing.T) {
	testCases := []struct {
		name string

		withGithubReleases                []*github.RepositoryRelease
		withGithubReleaseAssetsFromFolder string
		withHost                          state.Host
		expectedBinaryVersionsRun         []string
		expectedErrorContains             string
	}{
		{
			name:                              "Should run a Linux upgrade",
			withGithubReleases:                createGithubReleases(osarch.Linux, osarch.Amd64, []string{"0.0.61"}),
			withGithubReleaseAssetsFromFolder: "0.0.61",
			withHost: state.Host{
				Os:   osarch.Linux,
				Arch: osarch.Amd64,
			},
			expectedBinaryVersionsRun: []string{"0.0.61"},
		},
		{
			name:                              "Should run a Darwin upgrade",
			withGithubReleases:                createGithubReleases(osarch.Darwin, osarch.Amd64, []string{"0.0.61"}),
			withGithubReleaseAssetsFromFolder: "0.0.61",
			withHost: state.Host{
				Os:   osarch.Darwin,
				Arch: osarch.Amd64,
			},
			expectedBinaryVersionsRun: []string{"0.0.61"},
		},
		{
			name:                              "Should detect if binary's digest doesn't match the expected digest",
			withGithubReleases:                createGithubReleases(osarch.Linux, osarch.Amd64, []string{"0.0.61"}),
			withGithubReleaseAssetsFromFolder: "invalid_digest",
			withHost: state.Host{
				Os:   osarch.Linux,
				Arch: osarch.Amd64,
			},
			expectedBinaryVersionsRun: []string{},
			expectedErrorContains: "failed to verify binary signature: verification failed, hash mismatch, " +
				"got: 83bae1d215407ff3715063a621afa9138d2b15392d930e6377ed4a6058fea0ba, " +
				"expected: a3bae1d215407ff3715063a621afa9138d2b15392d930e6377ed4a6058fea0ba",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// Given
			var buffer bytes.Buffer
			var err error

			tmpStore, err := storage.NewTemporaryStorage()
			assert.NoError(t, err)

			repoDir := "my-iac-repo"
			repositoryAbsoluteDir := path.Join(tmpStore.BasePath, repoDir)

			err = tmpStore.MkdirAll(repoDir)
			assert.NoError(t, err)

			upgrader := New(Opts{
				Debug:               false,
				Logger:              logrus.StandardLogger(),
				Out:                 &buffer,
				RepositoryDirectory: repositoryAbsoluteDir,
				GithubService:       NewGithubServiceMock(tc.withGithubReleases),
				ChecksumDownloader:  NewChecksumDownloader(),
				FetcherOpts: FetcherOpts{
					Host:  tc.withHost,
					Store: tmpStore,
				},
			})

			err = mockHTTPResponse(tc.withGithubReleaseAssetsFromFolder, tc.withGithubReleases)
			require.NoError(t, err)

			httpmock.Activate()
			defer httpmock.DeactivateAndReset()

			// When
			err = upgrader.Run()

			// Then
			if len(tc.expectedErrorContains) > 0 {
				//goland:noinspection GoNilness
				assert.Contains(t, err.Error(), tc.expectedErrorContains)
			} else {
				assert.NoError(t, err)
			}

			for _, version := range tc.expectedBinaryVersionsRun {
				expectedFilepath := path.Join(repoDir, fmt.Sprintf(
					"okctl-upgrade_%s_%s_%s_ran_successfully",
					version,
					capitalizeFirst(tc.withHost.Os),
					tc.withHost.Arch),
				)
				upgradeRan, err := tmpStore.Exists(expectedFilepath)
				assert.NoError(t, err)

				assert.True(t, upgradeRan, fmt.Sprintf(
					"the upgrade should have produced the file %s, but no file was found",
					path.Join(tmpStore.BasePath, expectedFilepath)))
			}
		})
	}
}

func createGithubReleases(os string, arch string, versions []string) []*github.RepositoryRelease {
	releases := make([]*github.RepositoryRelease, 0, len(versions))
	os = capitalizeFirst(os)

	for i, version := range versions {
		release := &github.RepositoryRelease{
			TagName: &versions[i],
			Name:    &versions[i],
			Assets: []*github.ReleaseAsset{
				{
					Name:        github.StringPtr("okctl-upgrade-checksums.txt"),
					ContentType: github.StringPtr("text/plain"),
					BrowserDownloadURL: github.StringPtr(fmt.Sprintf(
						"https://github.com/oslokommune/okctl-upgrade/releases/download/%s/okctl-upgrade-checksums.txt", version)),
				},
				{
					Name:        github.StringPtr(fmt.Sprintf("okctl_upgrade-%s_%s_%s.tar.gz", os, arch, version)),
					ContentType: github.StringPtr("application/gzip"),
					BrowserDownloadURL: github.StringPtr(fmt.Sprintf(
						"https://github.com/oslokommune/okctl-upgrade/releases/download/%s/okctl-upgrade_%s_%s_%s.tar.gz", version, version, os, arch)),
				},
			},
		}

		releases = append(releases, release)
	}

	return releases
}

func capitalizeFirst(os string) string {
	return strings.ToUpper(os[0:1]) + strings.ToLower(os[1:])
}

func mockHTTPResponse(folder string, releases []*github.RepositoryRelease) error {
	for _, release := range releases {
		for _, asset := range release.Assets {
			assetFilename := getAssetFilename(*asset.BrowserDownloadURL)

			data, err := readBytesFromFile(fmt.Sprintf(`testdata/%s/%s`, folder, assetFilename))
			if err != nil {
				return err
			}

			responder := httpmock.NewBytesResponder(http.StatusOK, data)
			httpmock.RegisterResponder(http.MethodGet, *asset.BrowserDownloadURL, responder)
		}
	}

	return nil
}

func getAssetFilename(url string) string {
	split := strings.Split(url, "/")
	return split[len(split)-1]
}

func readBytesFromFile(file string) ([]byte, error) {
	//nolint: gosec
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	return b, nil
}
