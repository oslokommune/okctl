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

	"github.com/stretchr/testify/require"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/sebdah/goldie/v2"

	"github.com/oslokommune/okctl/pkg/osarch"

	"github.com/jarcoal/httpmock"
	"github.com/oslokommune/okctl/pkg/config/state"
	"github.com/oslokommune/okctl/pkg/github"
	"github.com/oslokommune/okctl/pkg/storage"
	"github.com/sirupsen/logrus"
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
// x Should run an upgrade
// x Should not run already applied upgrades - custom: Må kjøre upgrade flere ganger. Assert: each binary was run once
// x Should run upgrades up to the current okctl version, but no newer
// ?   See: // "DO: Remove file verification" , should be easier to do verifications.
// x Should not run too old upgrades
// x Should run hot fixes in correct order
// x Should run a hotfix even if it is older than the last applied upgrade.
//     Så hvis upgrade bumper 0.0.65, og det kommer en hotfix 0.0.63_my-hotfix, så skal den fortsatt kjøres. - custom sjekk
// x Må støtte OriginalOkctlVersion

// Failure situations
// x Should detect if release has invalid tag name or assets (must support hot fixes)
// - Should verify digest before running. Trengs ikke, veldig sjelden man laster ned og kjører senere.
// I okctl-upgrade-checksums.txt, endre filnavn. Da bør man få feil at ting ikke matcher. Får noe annet unyttig. - unit test
// Other?

// Refaktorer: Bruker Sprintf med %s-%s-%s mange steder. Bør være ett sted.

// upgrade --dry-run ?
const folderWorking = "working"

//nolint:funlen
func TestRunUpgrades(t *testing.T) {
	testCases := []struct {
		name                              string
		withOkctlVersion                  string
		withOriginalOkctlVersion          string
		withGithubReleases                []*github.RepositoryRelease
		withGithubReleaseAssetsFromFolder string
		withHost                          state.Host
		withTestRun                       func(t *testing.T, defaultOpts Opts)
		expectBinaryVersionsRunOnce       []string
		expectErrorContains               string
	}{
		{
			name:                              "Should detect if binary's digest doesn't match the expected digest",
			withOkctlVersion:                  "0.0.61",
			withOriginalOkctlVersion:          "0.0.50",
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
			withOkctlVersion:                  "0.0.61",
			withOriginalOkctlVersion:          "0.0.50",
			withGithubReleases:                createGithubReleases([]string{osarch.Linux, osarch.Darwin}, osarch.Amd64, []string{"0.0.61"}),
			withGithubReleaseAssetsFromFolder: folderWorking,
			withHost:                          state.Host{Os: osarch.Linux, Arch: osarch.Amd64},
			expectBinaryVersionsRunOnce:       []string{"0.0.61"},
		},
		{
			name:                              "Should return exit status if upgrade crashes",
			withOkctlVersion:                  "0.0.58",
			withOriginalOkctlVersion:          "0.0.50",
			withGithubReleases:                createGithubReleases([]string{osarch.Linux, osarch.Darwin}, osarch.Amd64, []string{"0.0.58"}),
			withGithubReleaseAssetsFromFolder: "upgrade_crashes",
			withHost:                          state.Host{Os: osarch.Linux, Arch: osarch.Amd64},
			expectBinaryVersionsRunOnce:       []string{},
			expectErrorContains:               "exit status 1",
		},
		{
			name:                        "Should run zero upgrades",
			withOkctlVersion:            "0.0.60",
			withOriginalOkctlVersion:    "0.0.50",
			withGithubReleases:          []*github.RepositoryRelease{},
			withHost:                    state.Host{Os: osarch.Linux, Arch: osarch.Amd64},
			expectBinaryVersionsRunOnce: []string{},
		},
		{
			name:                              "Should run a Linux upgrade",
			withOkctlVersion:                  "0.0.61",
			withOriginalOkctlVersion:          "0.0.50",
			withGithubReleases:                createGithubReleases([]string{osarch.Linux, osarch.Darwin}, osarch.Amd64, []string{"0.0.61"}),
			withGithubReleaseAssetsFromFolder: folderWorking,
			withHost:                          state.Host{Os: osarch.Linux, Arch: osarch.Amd64},
			expectBinaryVersionsRunOnce:       []string{"0.0.61"},
		},
		{
			name:                              "Should run a Darwin upgrade",
			withOkctlVersion:                  "0.0.61",
			withOriginalOkctlVersion:          "0.0.50",
			withGithubReleases:                createGithubReleases([]string{osarch.Linux, osarch.Darwin}, osarch.Amd64, []string{"0.0.61"}),
			withGithubReleaseAssetsFromFolder: folderWorking,
			withHost:                          state.Host{Os: osarch.Darwin, Arch: osarch.Amd64},
			expectBinaryVersionsRunOnce:       []string{"0.0.61"},
		},
		{
			name:                              "Should run multiple upgrades",
			withOkctlVersion:                  "0.0.64",
			withOriginalOkctlVersion:          "0.0.50",
			withGithubReleases:                createGithubReleases([]string{osarch.Linux, osarch.Darwin}, osarch.Amd64, []string{"0.0.61", "0.0.62", "0.0.64"}),
			withGithubReleaseAssetsFromFolder: folderWorking,
			withHost:                          state.Host{Os: osarch.Linux, Arch: osarch.Amd64},
			expectBinaryVersionsRunOnce:       []string{"0.0.61", "0.0.62", "0.0.64"},
		},
		{
			name:                              "Should run upgrades once",
			withOkctlVersion:                  "0.0.64",
			withOriginalOkctlVersion:          "0.0.50",
			withGithubReleaseAssetsFromFolder: folderWorking,
			withHost:                          state.Host{Os: osarch.Linux, Arch: osarch.Amd64},
			withTestRun: func(t *testing.T, defaultOpts Opts) {
				// Given
				mockHTTPResponse(t, folderWorking, createGithubReleases([]string{osarch.Linux, osarch.Darwin}, osarch.Amd64, []string{"0.0.61", "0.0.62", "0.0.64"}))
				defer httpmock.DeactivateAndReset()

				// When
				upgrader := New(defaultOpts)

				// Then
				err := upgrader.Run()
				assert.NoError(t, err)

				// When
				upgrader = New(defaultOpts)

				// Then
				err = upgrader.Run()
				assert.NoError(t, err)
			},
			expectBinaryVersionsRunOnce: []string{"0.0.61", "0.0.62", "0.0.64"},
		},
		{
			name:                              "Should run upgrades with version up to and including current okctl version, but no newer",
			withOkctlVersion:                  "0.0.63",
			withOriginalOkctlVersion:          "0.0.50",
			withGithubReleases:                createGithubReleases([]string{osarch.Linux, osarch.Darwin}, osarch.Amd64, []string{"0.0.61", "0.0.62", "0.0.63", "0.0.64"}),
			withGithubReleaseAssetsFromFolder: folderWorking,
			withHost:                          state.Host{Os: osarch.Linux, Arch: osarch.Amd64},
			expectBinaryVersionsRunOnce:       []string{"0.0.61", "0.0.62", "0.0.63"},
		},
		{

			name:                              "Should not run upgrades that are older than the first installed okctl version",
			withOkctlVersion:                  "0.0.64",
			withOriginalOkctlVersion:          "0.0.62",
			withGithubReleases:                createGithubReleases([]string{osarch.Linux, osarch.Darwin}, osarch.Amd64, []string{"0.0.61", "0.0.62", "0.0.63", "0.0.64"}),
			withGithubReleaseAssetsFromFolder: folderWorking,
			withHost:                          state.Host{Os: osarch.Linux, Arch: osarch.Amd64},
			expectBinaryVersionsRunOnce:       []string{"0.0.63", "0.0.64"},
		},
		{
			name:                     "Should run upgrade hot fixes, and in correct order",
			withOkctlVersion:         "0.0.63",
			withOriginalOkctlVersion: "0.0.50",
			withGithubReleases: createGithubReleases([]string{osarch.Linux, osarch.Darwin}, osarch.Amd64,
				[]string{"0.0.63.a", "0.0.62", "0.0.62.b", "0.0.61", "0.0.62.a", "0.0.64.a", "0.0.63"}),
			withGithubReleaseAssetsFromFolder: folderWorking,
			withHost:                          state.Host{Os: osarch.Linux, Arch: osarch.Amd64},
			expectBinaryVersionsRunOnce:       []string{"0.0.61", "0.0.62", "0.0.62.a", "0.0.62.b", "0.0.63", "0.0.63.a"},
		},
		{
			name:                     "Should not run upgrades, including hot fixes, that are older than the first installed okctl version",
			withOkctlVersion:         "0.0.63",
			withOriginalOkctlVersion: "0.0.62",
			withGithubReleases: createGithubReleases([]string{osarch.Linux, osarch.Darwin}, osarch.Amd64,
				[]string{"0.0.62", "0.0.62.a", "0.0.62.b", "0.0.63", "0.0.63.a", "0.0.64.a"}),
			withGithubReleaseAssetsFromFolder: folderWorking,
			withHost:                          state.Host{Os: osarch.Linux, Arch: osarch.Amd64},
			expectBinaryVersionsRunOnce:       []string{"0.0.63", "0.0.63.a"},
		},
		{
			// Explanation: The whole point of hotfixes are that you may be in the situation where okctl already has
			// upgraded you to version 0.0.63. But then we discover that we made an error in the upgrade for 0.0.62.
			// So we make a hotfix "0.0.62.a" which will be run even though you already upgraded to 0.0.63.
			name:                              "Should run a hotfix even if it is older than the last applied upgrade",
			withOkctlVersion:                  "0.0.63",
			withOriginalOkctlVersion:          "0.0.50",
			withGithubReleaseAssetsFromFolder: folderWorking,
			withHost:                          state.Host{Os: osarch.Linux, Arch: osarch.Amd64},
			withTestRun: func(t *testing.T, defaultOpts Opts) {
				// Given settings for first upgrade
				githubReleaseVersions := []string{"0.0.61", "0.0.62", "0.0.63"}
				githubReleases := createGithubReleases([]string{osarch.Linux, osarch.Darwin}, osarch.Amd64, githubReleaseVersions)
				mockHTTPResponse(t, folderWorking, githubReleases)
				defaultOpts.GithubService = NewGithubServiceMock(githubReleases)

				stdOutBuffer := new(bytes.Buffer)
				defaultOpts.Out = stdOutBuffer

				upgrader := New(defaultOpts)

				// When running upgrade first time
				err := upgrader.Run()

				// Then
				t.Log(stdOutBuffer.String())

				assert.NoError(t, err)

				expectedUpgradesRun := getExpectedUpgradesRun(githubReleaseVersions, state.Host{Os: osarch.Linux, Arch: osarch.Amd64})
				upgradesRun := getActualUpgradesRun(stdOutBuffer)
				assert.Equal(t, expectedUpgradesRun, upgradesRun, "Unexpected upgrades were run")

				g := goldie.New(t)
				g.Assert(t, "Should run a hotfix even if it is older than the last applied upgrade_run1", stdOutBuffer.Bytes())

				httpmock.DeactivateAndReset()

				// Given settings for second upgrade
				githubReleaseVersions = []string{"0.0.61", "0.0.62", "0.0.62.a", "0.0.62.b", "0.0.63", "0.0.63.a", "0.0.64"}
				githubReleases = createGithubReleases([]string{osarch.Linux, osarch.Darwin}, osarch.Amd64, githubReleaseVersions)
				mockHTTPResponse(t, folderWorking, githubReleases)
				defer httpmock.DeactivateAndReset()
				defaultOpts.GithubService = NewGithubServiceMock(githubReleases)

				stdOutBuffer = new(bytes.Buffer)
				defaultOpts.Out = stdOutBuffer

				upgrader = New(defaultOpts)

				// When running upgrade second time
				err = upgrader.Run()

				// Then
				t.Log(stdOutBuffer.String())

				assert.NoError(t, err)

				expectedUpgradesRun2 := getExpectedUpgradesRun([]string{"0.0.62.a", "0.0.62.b", "0.0.63.a"}, state.Host{Os: osarch.Linux, Arch: osarch.Amd64})
				upgradesRun2 := getActualUpgradesRun(stdOutBuffer)
				assert.Equal(t, expectedUpgradesRun2, upgradesRun2, "Unexpected upgrades were run")

				g.Assert(t, "Should run a hotfix even if it is older than the last applied upgrade_run2", stdOutBuffer.Bytes())
			},
		},
		{
			name:                     "Should return error if github release isn't valid",
			withOkctlVersion:         "0.0.61",
			withOriginalOkctlVersion: "0.0.50",
			withGithubReleases: []*github.RepositoryRelease{
				{
					ID: github.Int64Ptr(123),
				},
			},
			withGithubReleaseAssetsFromFolder: folderWorking,
			withHost:                          state.Host{Os: osarch.Linux, Arch: osarch.Amd64},
			withTestRun: func(t *testing.T, defaultOpts Opts) {
				upgrader := New(defaultOpts)

				err := upgrader.Run()

				assert.Error(t, err)
				assert.Contains(t, err.Error(),
					"release 'ID: 123' must have at least 2 assets (binary and checksum); "+
						"name: cannot be blank; tag_name: cannot be blank.",
				)
			},
		},
		{
			name:                     "Should return error if github release assets don't include checksum",
			withOkctlVersion:         "0.0.64",
			withOriginalOkctlVersion: "0.0.50",
			withGithubReleases: []*github.RepositoryRelease{
				{
					ID:      github.Int64Ptr(123),
					TagName: github.StringPtr("0.0.61"),
					Name:    github.StringPtr("0.0.61"),
					Assets: []*github.ReleaseAsset{
						createGihubReleaseAssetBinary(osarch.Linux, osarch.Amd64, "0.0.61"),
						createGihubReleaseAssetBinary(osarch.Darwin, osarch.Amd64, "0.0.61"),
					},
				},
			},
			withGithubReleaseAssetsFromFolder: folderWorking,
			withHost:                          state.Host{Os: osarch.Linux, Arch: osarch.Amd64},
			withTestRun: func(t *testing.T, defaultOpts Opts) {
				upgrader := New(defaultOpts)

				err := upgrader.Run()

				assert.Error(t, err)
				assert.Contains(t, err.Error(),
					"parsing upgrade binaries: validating release: could not find checksum asset for "+
						"release 0.0.61 (assets: okctl_upgrade-linux_amd64_0.0.61.tar.gz,"+
						"okctl_upgrade-darwin_amd64_0.0.61.tar.gz)",
				)
			},
		},
		{
			name:                     "Should return error if release version does not match release download URL version",
			withOkctlVersion:         "0.0.64",
			withOriginalOkctlVersion: "0.0.50",
			withGithubReleases: []*github.RepositoryRelease{
				{
					ID:      github.Int64Ptr(123),
					TagName: github.StringPtr("0.0.61"),
					Name:    github.StringPtr("0.0.61"),
					Assets: []*github.ReleaseAsset{
						{
							Name: github.StringPtr("okctl_upgrade-linux_amd64_0.0.61.tar.gz"),
							BrowserDownloadURL: github.StringPtr(fmt.Sprintf(
								"https://github.com/oslokommune/okctl-upgrade/releases/download/%s/okctl-upgrade_%s_%s_%s.tar.gz",
								"0.0.61", "0.0.61", osarch.Linux, osarch.Amd64)),
						},
						{
							Name: github.StringPtr("okctl_upgrade-darwin_amd64_0.0.61.tar.gz"),
							BrowserDownloadURL: github.StringPtr(fmt.Sprintf(
								"https://github.com/oslokommune/okctl-upgrade/releases/download/%s/okctl-upgrade_%s_%s_%s.tar.gz",
								"0.0.61", "0.0.61", osarch.Darwin, osarch.Amd64)),
						},
					},
				},
			},
			withGithubReleaseAssetsFromFolder: folderWorking,
			withHost:                          state.Host{Os: osarch.Linux, Arch: osarch.Amd64},
			withTestRun: func(t *testing.T, defaultOpts Opts) {
				upgrader := New(defaultOpts)

				err := upgrader.Run()

				assert.Error(t, err)
				assert.Contains(t, err.Error(),
					"parsing upgrade binaries: validating release: could not find checksum asset for "+
						"release 0.0.61 (assets: okctl_upgrade-linux_amd64_0.0.61.tar.gz,"+
						"okctl_upgrade-darwin_amd64_0.0.61.tar.gz)",
				)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// Given
			var err error

			stdOutBuffer := new(bytes.Buffer)

			tmpStore, err := storage.NewTemporaryStorage()
			assert.NoError(t, err)

			repoDir := "my-iac-repo"
			repositoryAbsoluteDir := path.Join(tmpStore.BasePath, repoDir)

			err = tmpStore.MkdirAll(repoDir)
			assert.NoError(t, err)

			defaultOpts := Opts{
				Debug:               false,
				Logger:              logrus.StandardLogger(),
				Out:                 stdOutBuffer,
				RepositoryDirectory: repositoryAbsoluteDir,
				GithubService:       NewGithubServiceMock(tc.withGithubReleases),
				ChecksumDownloader:  NewChecksumDownloader(),
				FetcherOpts: FetcherOpts{
					Host:  tc.withHost,
					Store: tmpStore,
				},
				OkctlVersion:         tc.withOkctlVersion,
				OriginalOkctlVersion: tc.withOriginalOkctlVersion,
				State:                mockUpgradeState(),
				ClusterID:            api.ID{},
			}

			// When
			if tc.withTestRun != nil {
				tc.withTestRun(t, defaultOpts)
				return
			}

			mockHTTPResponse(t, tc.withGithubReleaseAssetsFromFolder, tc.withGithubReleases)
			defer httpmock.DeactivateAndReset()

			upgrader := New(defaultOpts)
			err = upgrader.Run()

			t.Log(stdOutBuffer.String())

			// Then
			if len(tc.expectErrorContains) > 0 {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectErrorContains)
				return
			}

			assert.NoError(t, err)

			expectedUpgradesRun := getExpectedUpgradesRun(tc.expectBinaryVersionsRunOnce, tc.withHost)
			upgradesRun := getActualUpgradesRun(stdOutBuffer)
			assert.Equal(t, expectedUpgradesRun, upgradesRun, "Unexpected upgrades were run")

			g := goldie.New(t)
			g.Assert(t, tc.name, stdOutBuffer.Bytes())
		})
	}
}

func mockHTTPResponse(t *testing.T, folder string, releases []*github.RepositoryRelease) {
	err := registerHTTPResponse(
		folder,
		releases,
	)
	require.NoError(t, err)

	httpmock.Activate()
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

func getActualUpgradesRun(stdOutBuffer *bytes.Buffer) []string {
	stdOut := stdOutBuffer.String()
	re := regexp.MustCompile(`This is upgrade file for (okctl-upgrade.*)`)
	found := re.FindAllStringSubmatch(stdOut, -1)
	upgradesRun := make([]string, 0)

	for _, match := range found {
		upgradesRun = append(upgradesRun, match[1])
	}

	return upgradesRun
}

func createGihubReleaseAssetBinary(os, arch, version string) *github.ReleaseAsset {
	return &github.ReleaseAsset{
		Name:        github.StringPtr(fmt.Sprintf("okctl_upgrade-%s_%s_%s.tar.gz", os, arch, version)),
		ContentType: github.StringPtr("application/gzip"),
		BrowserDownloadURL: github.StringPtr(fmt.Sprintf(
			"https://github.com/oslokommune/okctl-upgrade/releases/download/%s/okctl-upgrade_%s_%s_%s.tar.gz", version, version, os, arch)),
	}
}

//nolint:unparam
func createGithubReleases(oses []string, arch string, versions []string) []*github.RepositoryRelease {
	releases := make([]*github.RepositoryRelease, 0, len(versions))

	for i, version := range versions {
		assets := make([]*github.ReleaseAsset, 0, len(oses)+1)

		for _, os := range oses {
			os = capitalizeFirst(os)

			asset := createGihubReleaseAssetBinary(os, arch, version)
			assets = append(assets, asset)
		}

		assets = append(assets, &github.ReleaseAsset{
			Name:        github.StringPtr(upgradeChecksumsTxt),
			ContentType: github.StringPtr("text/plain"),
			BrowserDownloadURL: github.StringPtr(fmt.Sprintf(
				"https://github.com/oslokommune/okctl-upgrade/releases/download/%s/okctl-upgrade-checksums.txt", version)),
		})

		release := &github.RepositoryRelease{
			ID:      github.Int64Ptr(int64(i + 1)),
			TagName: &versions[i],
			Name:    &versions[i],
			Assets:  assets,
		}

		releases = append(releases, release)
	}

	// reverse slice to make sure sorting of upgrade binaries work
	for i, j := 0, len(releases)-1; i < j; i, j = i+1, j-1 {
		releases[i], releases[j] = releases[j], releases[i]
	}

	return releases
}

// capitalizeFirst converts for instance "linux" to "Linux". We use this because we expect github release assets for
// upgrades to be named this way.
func capitalizeFirst(os string) string {
	return strings.ToUpper(os[0:1]) + strings.ToLower(os[1:])
}

func registerHTTPResponse(baseFolder string, releases []*github.RepositoryRelease) error {
	for _, release := range releases {
		versionFolder := *release.TagName

		for _, asset := range release.Assets {
			assetFilename := getAssetFilename(*asset.BrowserDownloadURL)

			data, err := readBytesFromFile(fmt.Sprintf("testdata/%s/%s/%s", baseFolder, versionFolder, assetFilename))
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

//goland:noinspection GoUnusedParameter
func (u upgradeStateMock) SaveOriginalOkctlVersionIfNotExists(originalOkctlVersion *client.OriginalOkctlVersion) error {
	panic("implement me")
}

func (u upgradeStateMock) GetOriginalOkctlVersion() (*client.OriginalOkctlVersion, error) {
	panic("implement me")
}

func mockUpgradeState() client.UpgradeState {
	return &upgradeStateMock{
		upgrades: make(map[string]*client.Upgrade),
	}
}
