package upgrade

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"

	"github.com/oslokommune/okctl/pkg/config/state"
	"github.com/oslokommune/okctl/pkg/github"
	"github.com/oslokommune/okctl/pkg/storage"
	"github.com/stretchr/testify/assert"
)

// Test cases:
// Given these releases, ..., then these binaries should be run
// Should run a upgrade
// Should not run already applied migrations
// Should run migrations up to the current okctl version
// Failure situations
// Verify invalid digests in okctl-upgrade.txt
//
// I okctl-upgrade-checksuyms.txt, endre filnavn. Da bør man få feil at ting ikke matcher. Får noe annet unyttig.

// Vurder: (Egen migrator test?)
// withOriginalOkctlVersion     string
// withOkctlVersion             string
// withAlreadyAppliedMigrations []string
// expectMigrationsToBeRun      []string

func TestRunUpgrades(t *testing.T) {
	testCases := []struct {
		name string

		withGithubReleases                []*github.RepositoryRelease
		withGithubReleaseAssetsFromFolder string
		withHost                          state.Host
	}{
		{
			name: "Should run an upgrade",
			withGithubReleases: createGithubReleases(
				state.Host{Os: "Linux", Arch: "amd64"},
				[]string{"0.0.61"}),
			withGithubReleaseAssetsFromFolder: "0.0.61",
			withHost: state.Host{
				Os:   "linux",
				Arch: "amd64",
			},
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
			assert.NoError(t, err)

			upgradeRan, err := tmpStore.Exists(path.Join(repoDir, "okctl_upgrade_0.0.61_ran_successfully"))
			assert.NoError(t, err)

			assert.True(t, upgradeRan, "the upgrade should have produced a file, but no file was found")
			fmt.Println(tmpStore.BasePath)
		})
	}
}

func createGithubReleases(host state.Host, versions []string) []*github.RepositoryRelease {
	releases := make([]*github.RepositoryRelease, 0, len(versions))

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
					Name:        github.StringPtr(fmt.Sprintf("okctl_upgrade-%s_%s_%s.tar.gz", host.Os, host.Arch, version)),
					ContentType: github.StringPtr("application/gzip"),
					BrowserDownloadURL: github.StringPtr(fmt.Sprintf(
						"https://github.com/oslokommune/okctl-upgrade/releases/download/%s/okctl-upgrade_%s_%s_%s.tar.gz", version, version, host.Os, host.Arch)),
				},
			},
		}

		releases = append(releases, release)
	}

	return releases
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
