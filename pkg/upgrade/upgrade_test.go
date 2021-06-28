package upgrade

import (
	"bytes"
	"fmt"
	"github.com/jarcoal/httpmock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
	"testing"

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
// TODO: Failure situations?

// TODO Vurder: Muligens dette bør være test i binaries provider: Should download if not exists...? Spørs hva som må
// gjøres i binaries provider

// Separer domenelogikk og applikasjonslogikk. Tester for applogikk kjører også logikk for domene.
// Use gock for api mocking?

func TestRunMigrations(t *testing.T) {
	testCases := []struct {
		name string

		withUpgradeGithubReleases    []*github.RepositoryRelease
		withOriginalOkctlVersion     string
		withOkctlVersion             string
		withAlreadyAppliedMigrations []string // TODO: or upgrade?
		withHost                     state.Host
		withBinariesFolder           string
		expectMigrationsToBeRun      []string // TODO: or upgrade?
	}{
		{
			name: "Should run an upgrade",
			withUpgradeGithubReleases: createGithubReleases(state.Host{
				Os:   "Linux",
				Arch: "amd64",
			}, "0.0.61"),
			withOriginalOkctlVersion:     "0.0.60",
			withOkctlVersion:             "0.0.60",
			withAlreadyAppliedMigrations: []string{},
			withHost: state.Host{
				Os:   "linux",
				Arch: "amd64",
			},
			withBinariesFolder:      "0.0.61",
			expectMigrationsToBeRun: []string{"0.0.61"},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// Given
			var buffer bytes.Buffer
			var err error
			githubService := NewGithubServiceMock(tc.withUpgradeGithubReleases)

			tmpStore, err := storage.NewTemporaryStorage()
			assert.NoError(t, err)

			repoDir := path.Join(tmpStore.BasePath, "my-iac-repo")
			err = tmpStore.MkdirAll("my-iac-repo")
			assert.NoError(t, err)

			upgrader := New(Opts{
				Debug:               false,
				Logger:              logrus.StandardLogger(),
				Out:                 &buffer,
				RepoDir:             repoDir,
				GithubService:       githubService,
				GithubReleaseParser: NewGithubReleaseParser(NewChecksumDownloader()),
				FetcherOpts: FetcherOpts{
					Host:  tc.withHost,
					Store: tmpStore,
				},
			})

			err = mockHTTPResponse(tc.withBinariesFolder, tc.withUpgradeGithubReleases)
			require.NoError(t, err)

			httpmock.Activate()
			defer httpmock.DeactivateAndReset()

			// When
			err = upgrader.Run()

			// Then
			assert.NoError(t, err)

			// TODO: Assert that the mock binary runner has ran the upgrade executables we expect
		})
	}
}

// TODO: I okctl-upgrade-checksuyms.txt, endre digest. Da funker testen!
// TODO: I okctl-upgrade-checksuyms.txt, endre filnavn. Da bør man få feil at ting ikke matcher. Får noe annet unyttig.

func createGithubReleases(host state.Host, versions ...string) []*github.RepositoryRelease {
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

			fmt.Println("Creating responder for: " + *asset.BrowserDownloadURL)

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
