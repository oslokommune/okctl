package upgrade

import (
	"bytes"
	"fmt"
	"github.com/sebdah/goldie/v2"
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
// Should run hotfixes in order

// Should verify digest before running
// Lots of failure situations
//
// I okctl-upgrade-checksums.txt, endre filnavn. Da bør man få feil at ting ikke matcher. Får noe annet unyttig.
// ------------------------------------------------

// Valg: La binæren selv oppdatere state at den er kjørt. Men bruke en pakke fra okctl. Da kan okctl selv bruke denne
// pakken.

// upgrade --dry-run ?

//nolint:funlen
func TestRunUpgrades(t *testing.T) {
	testCases := []struct {
		name string

		withGithubReleases                []*github.RepositoryRelease
		withGithubReleaseAssetsFromFolder string
		withHost                          state.Host
		expectBinaryVersionsRun           []string
		expectErrorContains               string
		expectedStdOutGolden              bool
	}{
		{
			name:                    "Should run zero upgrades",
			withGithubReleases:      []*github.RepositoryRelease{},
			withHost:                state.Host{Os: osarch.Linux, Arch: osarch.Amd64},
			expectBinaryVersionsRun: []string{},
		},
		{
			name:                              "Should run a Linux upgrade",
			withGithubReleases:                createGithubReleases([]string{osarch.Linux, osarch.Darwin}, osarch.Amd64, []string{"0.0.61"}),
			withGithubReleaseAssetsFromFolder: "working",
			withHost:                          state.Host{Os: osarch.Linux, Arch: osarch.Amd64},
			expectBinaryVersionsRun:           []string{"0.0.61"},
		},
		{
			name:                              "Should run a Darwin upgrade",
			withGithubReleases:                createGithubReleases([]string{osarch.Linux, osarch.Darwin}, osarch.Amd64, []string{"0.0.61"}),
			withGithubReleaseAssetsFromFolder: "working",
			withHost:                          state.Host{Os: osarch.Darwin, Arch: osarch.Amd64},
			expectBinaryVersionsRun:           []string{"0.0.61"},
		},
		{
			name:                              "Should detect if binary's digest doesn't match the expected digest",
			withGithubReleases:                createGithubReleases([]string{osarch.Linux, osarch.Darwin}, osarch.Amd64, []string{"0.0.61"}),
			withGithubReleaseAssetsFromFolder: "invalid_digest",
			withHost:                          state.Host{Os: osarch.Linux, Arch: osarch.Amd64},
			expectBinaryVersionsRun:           []string{},
			expectErrorContains: "failed to verify binary signature: verification failed, hash mismatch, " +
				"got: 83bae1d215407ff3715063a621afa9138d2b15392d930e6377ed4a6058fea0ba, " +
				"expected: a3bae1d215407ff3715063a621afa9138d2b15392d930e6377ed4a6058fea0ba",
		},
		{
			name:                              "Should print upgrade's stdout to stdout",
			withGithubReleases:                createGithubReleases([]string{osarch.Linux, osarch.Darwin}, osarch.Amd64, []string{"0.0.61"}),
			withGithubReleaseAssetsFromFolder: "working",
			withHost:                          state.Host{Os: osarch.Linux, Arch: osarch.Amd64},
			expectedStdOutGolden:              true,
			expectBinaryVersionsRun:           []string{"0.0.61"},
		},
		{
			name:                              "Should return exit status if upgrade crashes",
			withGithubReleases:                createGithubReleases([]string{osarch.Linux, osarch.Darwin}, osarch.Amd64, []string{"0.0.58"}),
			withGithubReleaseAssetsFromFolder: "upgrade_crashes",
			withHost:                          state.Host{Os: osarch.Linux, Arch: osarch.Amd64},
			expectBinaryVersionsRun:           []string{},
			expectErrorContains:               "exit status 1",
		},
		{
			name:                              "Should run two upgrades",
			withGithubReleases:                createGithubReleases([]string{osarch.Linux, osarch.Darwin}, osarch.Amd64, []string{"0.0.61", "0.0.62"}),
			withGithubReleaseAssetsFromFolder: "working",
			withHost:                          state.Host{Os: osarch.Linux, Arch: osarch.Amd64},
			expectBinaryVersionsRun:           []string{"0.0.61", "0.0.62"},
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
			if len(tc.expectErrorContains) > 0 {
				//goland:noinspection GoNilness
				assert.Contains(t, err.Error(), tc.expectErrorContains)
			} else {
				assert.NoError(t, err)
			}

			for _, version := range tc.expectBinaryVersionsRun {
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

			if tc.expectedStdOutGolden {
				g := goldie.New(t)
				g.Assert(t, tc.name, buffer.Bytes())
			}

			t.Log(buffer.String())
		})
	}
}

//nolint:unparam
func createGithubReleases(oses []string, arch string, versions []string) []*github.RepositoryRelease {
	releases := make([]*github.RepositoryRelease, 0, len(versions))

	for i, version := range versions {
		assets := make([]*github.ReleaseAsset, 0, len(oses)+1)

		for _, os := range oses {
			os = capitalizeFirst(os)

			asset := &github.ReleaseAsset{
				Name:        github.StringPtr(fmt.Sprintf("okctl_upgrade-%s_%s_%s.tar.gz", os, arch, version)),
				ContentType: github.StringPtr("application/gzip"),
				BrowserDownloadURL: github.StringPtr(fmt.Sprintf(
					"https://github.com/oslokommune/okctl-upgrade/releases/download/%s/okctl-upgrade_%s_%s_%s.tar.gz", version, version, os, arch)),
			}

			assets = append(assets, asset)
		}

		assets = append(assets, &github.ReleaseAsset{
			Name:        github.StringPtr("okctl-upgrade-checksums.txt"),
			ContentType: github.StringPtr("text/plain"),
			BrowserDownloadURL: github.StringPtr(fmt.Sprintf(
				"https://github.com/oslokommune/okctl-upgrade/releases/download/%s/okctl-upgrade-checksums.txt", version)),
		})

		release := &github.RepositoryRelease{
			TagName: &versions[i],
			Name:    &versions[i],
			Assets:  assets,
		}

		releases = append(releases, release)
	}

	return releases
}

func capitalizeFirst(os string) string {
	return strings.ToUpper(os[0:1]) + strings.ToLower(os[1:])
}

func mockHTTPResponse(baseFolder string, releases []*github.RepositoryRelease) error {
	for _, release := range releases {
		versionFolder := *release.TagName

		for _, asset := range release.Assets {
			assetFilename := getAssetFilename(*asset.BrowserDownloadURL)

			data, err := readBytesFromFile(fmt.Sprintf(`testdata/%s/%s/%s`, baseFolder, versionFolder, assetFilename))
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
