package upgrade

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"regexp"
	"strings"
	"testing"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/sebdah/goldie/v2"

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
// withAlreadyAppliedMigrations []string // HM kanskje ikke
// expectMigrationsToBeRun      []string

// Test cases
// ------------------------------------------------
// x Given these releases, ..., then these binaries should be run
// x Should run a upgrade
// x Should not run already applied upgrades - custom: Må kjøre upgrade flere ganger. Assert: each binary was run once
// -> Should run upgrades up to the current okctl version, but no newer - sjekk binaries that ran
//    See: // "DO: Remove file verification" , should be easier to do verifications.
// Should run hot fixes in correct order
// Should detect if release has invalid tag name or assets (must support hot fixes)
// Should not run too old upgrades - sjekk binaries that ran
// Should run hotfixes in order - sjekk binaries that ran
// Should run a hotfix even if it is older than the last applied upgrade.
//   Så hvis upgrade bumper 0.0.65, og det kommer en hotfix 0.0.63_my-hotfix, så skal den fortsatt kjøres. - custom sjekk

// Should verify digest before running
// Lots of failure situations

// I okctl-upgrade-checksums.txt, endre filnavn. Da bør man få feil at ting ikke matcher. Får noe annet unyttig.

// Refaktorer: Bruker Sprintf med %s-%s-%s mange steder. Bør være ett sted.

// upgrade --dry-run ?

//nolint:funlen
func TestRunUpgrades(t *testing.T) {
	testCases := []struct {
		name string

		withGithubReleases                []*github.RepositoryRelease
		withGithubReleaseAssetsFromFolder string
		withHost                          state.Host
		// withOriginalOkctlVersion          string // DO
		withOkctlVersion            string
		withTestRun                 func(t *testing.T, defaultOpts Opts) error
		expectBinaryVersionsRunOnce []string
		expectErrorContains         string
	}{
		{
			name:                              "Should detect if binary's digest doesn't match the expected digest",
			withGithubReleases:                createGithubReleases([]string{osarch.Linux, osarch.Darwin}, osarch.Amd64, []string{"0.0.61"}),
			withGithubReleaseAssetsFromFolder: "invalid_digest",
			withHost:                          state.Host{Os: osarch.Linux, Arch: osarch.Amd64},
			expectBinaryVersionsRunOnce:       []string{},
			expectErrorContains: "failed to verify binary signature: verification failed, hash mismatch, " +
				"got: 83bae1d215407ff3715063a621afa9138d2b15392d930e6377ed4a6058fea0ba, " +
				"expected: a3bae1d215407ff3715063a621afa9138d2b15392d930e6377ed4a6058fea0ba",
		},
		{
			name:                              "Should print upgrade's stdout to stdout",
			withGithubReleases:                createGithubReleases([]string{osarch.Linux, osarch.Darwin}, osarch.Amd64, []string{"0.0.61"}),
			withGithubReleaseAssetsFromFolder: "working",
			withHost:                          state.Host{Os: osarch.Linux, Arch: osarch.Amd64},
			expectBinaryVersionsRunOnce:       []string{"0.0.61"},
		},
		{
			name:                              "Should return exit status if upgrade crashes",
			withGithubReleases:                createGithubReleases([]string{osarch.Linux, osarch.Darwin}, osarch.Amd64, []string{"0.0.58"}),
			withGithubReleaseAssetsFromFolder: "upgrade_crashes",
			withHost:                          state.Host{Os: osarch.Linux, Arch: osarch.Amd64},
			expectBinaryVersionsRunOnce:       []string{},
			expectErrorContains:               "exit status 1",
		},
		{
			name:                        "Should run zero upgrades",
			withGithubReleases:          []*github.RepositoryRelease{},
			withHost:                    state.Host{Os: osarch.Linux, Arch: osarch.Amd64},
			expectBinaryVersionsRunOnce: []string{},
		},
		{
			name:                              "Should run a Linux upgrade",
			withGithubReleases:                createGithubReleases([]string{osarch.Linux, osarch.Darwin}, osarch.Amd64, []string{"0.0.61"}),
			withGithubReleaseAssetsFromFolder: "working",
			withHost:                          state.Host{Os: osarch.Linux, Arch: osarch.Amd64},
			expectBinaryVersionsRunOnce:       []string{"0.0.61"},
		},
		{
			name:                              "Should run a Darwin upgrade",
			withGithubReleases:                createGithubReleases([]string{osarch.Linux, osarch.Darwin}, osarch.Amd64, []string{"0.0.61"}),
			withGithubReleaseAssetsFromFolder: "working",
			withHost:                          state.Host{Os: osarch.Darwin, Arch: osarch.Amd64},
			expectBinaryVersionsRunOnce:       []string{"0.0.61"},
		},
		{
			name:                              "Should run multiple upgrades",
			withGithubReleases:                createGithubReleases([]string{osarch.Linux, osarch.Darwin}, osarch.Amd64, []string{"0.0.61", "0.0.62", "0.0.64"}),
			withGithubReleaseAssetsFromFolder: "working",
			withHost:                          state.Host{Os: osarch.Linux, Arch: osarch.Amd64},
			expectBinaryVersionsRunOnce:       []string{"0.0.61", "0.0.62", "0.0.64"},
		},
		{
			name:                              "Should not run already run upgrades",
			withGithubReleases:                createGithubReleases([]string{osarch.Linux, osarch.Darwin}, osarch.Amd64, []string{"0.0.61", "0.0.62", "0.0.64"}),
			withGithubReleaseAssetsFromFolder: "working",
			withHost:                          state.Host{Os: osarch.Linux, Arch: osarch.Amd64},
			withTestRun: func(t *testing.T, defaultOpts Opts) error {
				upgrader := New(defaultOpts)

				err := upgrader.Run()
				assert.NoError(t, err)

				upgrader = New(defaultOpts)

				err = upgrader.Run()
				assert.NoError(t, err)

				return nil
			},
			expectBinaryVersionsRunOnce: []string{"0.0.61", "0.0.62", "0.0.64"},
		},
		{
			name:                              "Should run upgrades with version up to and including current okctl version, but no newer",
			withGithubReleases:                createGithubReleases([]string{osarch.Linux, osarch.Darwin}, osarch.Amd64, []string{"0.0.61", "0.0.62", "0.0.64"}),
			withGithubReleaseAssetsFromFolder: "working",
			withHost:                          state.Host{Os: osarch.Linux, Arch: osarch.Amd64},
			withOkctlVersion:                  "0.0.62",
			expectBinaryVersionsRunOnce:       []string{"0.0.61", "0.0.62"},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// Given
			var stdOutBuffer bytes.Buffer
			var err error

			tmpStore, err := storage.NewTemporaryStorage()
			assert.NoError(t, err)

			repoDir := "my-iac-repo"
			repositoryAbsoluteDir := path.Join(tmpStore.BasePath, repoDir)

			err = tmpStore.MkdirAll(repoDir)
			assert.NoError(t, err)

			err = mockHTTPResponse(tc.withGithubReleaseAssetsFromFolder, tc.withGithubReleases)
			require.NoError(t, err)

			httpmock.Activate()
			defer httpmock.DeactivateAndReset()

			defaultOpts := Opts{
				Debug:               false,
				Logger:              logrus.StandardLogger(),
				Out:                 &stdOutBuffer,
				RepositoryDirectory: repositoryAbsoluteDir,
				GithubService:       NewGithubServiceMock(tc.withGithubReleases),
				ChecksumDownloader:  NewChecksumDownloader(),
				FetcherOpts: FetcherOpts{
					Host:  tc.withHost,
					Store: tmpStore,
				},
				OkctlVersion: "TODO", // DO
				State:        mockUpgradeState(),
				ClusterID:    api.ID{},
			}

			// When
			if tc.withTestRun == nil {
				upgrader := New(defaultOpts)
				err = upgrader.Run()
			} else {
				err = tc.withTestRun(t, defaultOpts)
			}

			// Then
			if len(tc.expectErrorContains) > 0 {
				//goland:noinspection GoNilness
				assert.Contains(t, err.Error(), tc.expectErrorContains)
			} else {
				assert.NoError(t, err)
			}

			expectedUpgradesRun := getExpectedUpgradesRun(tc.expectBinaryVersionsRunOnce, tc.withHost)
			upgradesRun := getActualUpgradesRun(stdOutBuffer)
			assert.Equal(t, expectedUpgradesRun, upgradesRun, "Unexpected upgrades were run")

			g := goldie.New(t)
			g.Assert(t, tc.name, stdOutBuffer.Bytes())

			t.Log(stdOutBuffer.String())
		})
	}
}

func getExpectedUpgradesRun(expectBinaryVersionsRunOnce []string, withHost state.Host) []string {
	expectedUpgradesRun := make([]string, 0, len(expectBinaryVersionsRunOnce))

	for _, binaryVersion := range expectBinaryVersionsRunOnce {
		expectedUpgradesRun = append(expectedUpgradesRun,
			fmt.Sprintf("okctl-upgrade_%s_%s_%s",
				binaryVersion,
				capitalizeFirst(withHost.Os),
				withHost.Arch,
			),
		)
	}

	return expectedUpgradesRun
}

func getActualUpgradesRun(stdOutBuffer bytes.Buffer) []string {
	stdOut := stdOutBuffer.String()
	re := regexp.MustCompile(`This is upgrade file for (okctl-upgrade.*)`)
	found := re.FindAllStringSubmatch(stdOut, -1)
	upgradesRun := make([]string, 0)

	for _, match := range found {
		upgradesRun = append(upgradesRun, match[1])
	}

	return upgradesRun
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

// capitalizeFirst converts for instance "linux" to "Linux". We use this because we expect github release assets for
// upgrades to be named this way.
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

type upgradeStateMock struct {
	upgrades map[string]*client.Upgrade
}

func (u upgradeStateMock) SaveUpgrade(upgrade *client.Upgrade) error {
	u.upgrades[upgrade.Version] = upgrade
	return nil
}

func (u upgradeStateMock) GetUpgrade(version string) (*client.Upgrade, error) {
	upgrade, ok := u.upgrades[version]
	if !ok {
		return nil, client.ErrUpgradeNotFound
	}

	return upgrade, nil
}

func mockUpgradeState() client.UpgradeState {
	return &upgradeStateMock{
		upgrades: make(map[string]*client.Upgrade),
	}
}
