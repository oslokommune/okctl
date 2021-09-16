package upgrade_test

import (
	"bytes"
	"fmt"
	"path"
	"regexp"
	"strings"
	"testing"

	"github.com/oslokommune/okctl/pkg/upgrade/clusterversion"
	"github.com/oslokommune/okctl/pkg/upgrade/originalclusterversion"

	"github.com/oslokommune/okctl/pkg/client"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/require"

	"github.com/oslokommune/okctl/pkg/upgrade/testutils"

	"github.com/jarcoal/httpmock"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/config/state"
	"github.com/oslokommune/okctl/pkg/github"
	"github.com/oslokommune/okctl/pkg/storage"
	"github.com/oslokommune/okctl/pkg/upgrade"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

const (
	folderWorking  = "working"
	folderCrashing = "upgrade_crashes"
	linux          = "linux"
	darwin         = "darwin"
	amd64          = "amd64"
)

type TestCase struct {
	name                               string
	withDebug                          bool
	withOkctlVersion                   string
	withOriginalClusterVersion         string
	withClusterVersion                 string
	withGithubReleases                 []*github.RepositoryRelease
	withGithubReleaseAssetsFromFolder  string
	withHost                           state.Host
	withUserAnswers                    []bool
	withTestRun                        func(t *testing.T, tc TestCase, defaultOpts DefaultTestOpts)
	expectBinaryVersionsRunOnce        []string
	expectedClusterVersionAfterUpgrade string
	expectErrorContains                string
}

//nolint:funlen
func TestRunUpgrades(t *testing.T) {
	testCases := []TestCase{
		{
			name:                              "Should detect if binary's digest doesn't match the expected digest",
			withOkctlVersion:                  "0.0.61",
			withOriginalClusterVersion:        "0.0.50",
			withGithubReleases:                createGithubReleases([]string{linux, darwin}, amd64, []string{"0.0.61"}),
			withGithubReleaseAssetsFromFolder: "invalid_digest",
			withHost:                          state.Host{Os: linux, Arch: amd64},
			expectBinaryVersionsRunOnce:       []string{},
			expectErrorContains: "failed to verify binary signature: verification failed, hash mismatch, " +
				"got: 83bae1d215407ff3715063a621afa9138d2b15392d930e6377ed4a6058fea0ba, " +
				"expected: a3bae1d215407ff3715063a621afa9138d2b15392d930e6377ed4a6058fea0ba",
		},
		{
			name:                              "Should print upgrade's stdout to stdout",
			withOkctlVersion:                  "0.0.61",
			withOriginalClusterVersion:        "0.0.50",
			withGithubReleases:                createGithubReleases([]string{linux, darwin}, amd64, []string{"0.0.61"}),
			withGithubReleaseAssetsFromFolder: folderWorking,
			withHost:                          state.Host{Os: linux, Arch: amd64},
			expectBinaryVersionsRunOnce:       []string{"0.0.61"},
		},
		{
			name:                              "Should return exit status if upgrade crashes",
			withOkctlVersion:                  "0.0.58",
			withOriginalClusterVersion:        "0.0.50",
			withGithubReleases:                createGithubReleases([]string{linux, darwin}, amd64, []string{"0.0.58"}),
			withGithubReleaseAssetsFromFolder: folderCrashing,
			withHost:                          state.Host{Os: linux, Arch: amd64},
			expectBinaryVersionsRunOnce:       []string{},
			expectErrorContains:               "exit status 1",
		},
		{
			name:                               "Should run zero upgrades",
			withOkctlVersion:                   "0.0.60",
			withOriginalClusterVersion:         "0.0.50",
			withClusterVersion:                 "0.0.52",
			withGithubReleases:                 []*github.RepositoryRelease{},
			withHost:                           state.Host{Os: linux, Arch: amd64},
			expectBinaryVersionsRunOnce:        []string{},
			expectedClusterVersionAfterUpgrade: "0.0.52",
		},
		{
			name:                              "Should run a Linux upgrade",
			withOkctlVersion:                  "0.0.61",
			withOriginalClusterVersion:        "0.0.50",
			withGithubReleases:                createGithubReleases([]string{linux, darwin}, amd64, []string{"0.0.61"}),
			withGithubReleaseAssetsFromFolder: folderWorking,
			withHost:                          state.Host{Os: linux, Arch: amd64},
			expectBinaryVersionsRunOnce:       []string{"0.0.61"},
		},
		{
			// Bumping cluster version happens after running every upgrade binary, to that upgrade's version. In this
			// test, that would be 0.0.61. Additionally, when all upgrades have completed successfully, we have chosen
			// to bump the cluster version to the current version of okctl. The reason for that, is that we want that
			// cluster to be in the same state as if creating a brand-new cluster with that okctl version.
			//
			// If creating a new cluster, the cluster version would be 0.0.70 (not 0.0.61). Therefore, we do a final
			// bump of cluster version to the current version of okctl, which is what this test verifies.
			//
			// The exception to this, is if we have run zero upgrades. Since absolutely no changes have been made,
			// we don't care about bumping cluster version either. This behavior is verified in another test.
			name:                              "Should bump cluster version to the current okctl version after upgrading",
			withOkctlVersion:                  "0.0.70",
			withOriginalClusterVersion:        "0.0.50",
			withGithubReleases:                createGithubReleases([]string{linux, darwin}, amd64, []string{"0.0.61"}),
			withGithubReleaseAssetsFromFolder: folderWorking,
			withHost:                          state.Host{Os: linux, Arch: amd64},
			expectBinaryVersionsRunOnce:       []string{"0.0.61"},
		},
		{
			name:                               "Should not bump cluster version if user aborts upgrading",
			withOkctlVersion:                   "0.0.70",
			withOriginalClusterVersion:         "0.0.50",
			withGithubReleases:                 createGithubReleases([]string{linux, darwin}, amd64, []string{"0.0.61"}),
			withGithubReleaseAssetsFromFolder:  folderWorking,
			withHost:                           state.Host{Os: linux, Arch: amd64},
			withUserAnswers:                    []bool{true, false},
			expectBinaryVersionsRunOnce:        []string{},
			expectedClusterVersionAfterUpgrade: "0.0.50",
		},
		{
			name:                              "Should run a Darwin upgrade",
			withOkctlVersion:                  "0.0.61",
			withOriginalClusterVersion:        "0.0.50",
			withGithubReleases:                createGithubReleases([]string{linux, darwin}, amd64, []string{"0.0.61"}),
			withGithubReleaseAssetsFromFolder: folderWorking,
			withHost:                          state.Host{Os: darwin, Arch: amd64},
			expectBinaryVersionsRunOnce:       []string{"0.0.61"},
		},
		{
			name:                              "Should return error if okctl version is below cluster version",
			withOkctlVersion:                  "0.0.60",
			withOriginalClusterVersion:        "0.0.50",
			withGithubReleases:                createGithubReleases([]string{linux, darwin}, amd64, []string{"0.0.61"}),
			withGithubReleaseAssetsFromFolder: folderWorking,
			withHost:                          state.Host{Os: darwin, Arch: amd64},
			expectBinaryVersionsRunOnce:       []string{},
			withTestRun: func(t *testing.T, tc TestCase, defaultOpts DefaultTestOpts) {
				err := defaultOpts.ClusterVersioner.SaveClusterVersion("0.0.70")
				require.NoError(t, err)

				_, err = upgrade.New(defaultOpts.Opts)

				assert.Error(t, err)
				assert.Contains(t, err.Error(), "okctl binary version 0.0.60 cannot be less than cluster"+
					" version 0.0.70. Get okctl version 0.0.70 or later and try again")
			},
		},
		{
			// In the future, when upgrade doesn't need to store original cluster version anymore, we should remove
			// this functionality (and this test). See comment in function upgrade.New. (tag UPGR01)
			name:                              "Should save original cluster version if it doesn't exist",
			withOkctlVersion:                  "0.0.61",
			withOriginalClusterVersion:        "0.0.50",
			withGithubReleases:                createGithubReleases([]string{linux, darwin}, amd64, []string{"0.0.61"}),
			withGithubReleaseAssetsFromFolder: folderWorking,
			withHost:                          state.Host{Os: darwin, Arch: amd64},
			expectBinaryVersionsRunOnce:       []string{"0.0.61"},
			withTestRun: func(t *testing.T, tc TestCase, defaultOpts DefaultTestOpts) {
				_, err := defaultOpts.OriginalClusterVersioner.GetOriginalClusterVersion()
				require.ErrorIs(t, client.ErrOriginalClusterVersionNotFound, err)

				_, err = upgrade.New(defaultOpts.Opts)
				require.NoError(t, err)

				originalClusterVersion, err := defaultOpts.OriginalClusterVersioner.GetOriginalClusterVersion()
				assert.NoError(t, err)

				assert.Equal(t, "0.0.50", originalClusterVersion)
			},
		},
		{
			name:                              "Should run multiple upgrades",
			withOkctlVersion:                  "0.0.64",
			withOriginalClusterVersion:        "0.0.50",
			withGithubReleases:                createGithubReleases([]string{linux, darwin}, amd64, []string{"0.0.61", "0.0.62", "0.0.64"}),
			withGithubReleaseAssetsFromFolder: folderWorking,
			withHost:                          state.Host{Os: linux, Arch: amd64},
			expectBinaryVersionsRunOnce:       []string{"0.0.61", "0.0.62", "0.0.64"},
		},
		{
			name:                              "Should run upgrades once",
			withOkctlVersion:                  "0.0.64",
			withOriginalClusterVersion:        "0.0.50",
			withGithubReleases:                createGithubReleases([]string{linux, darwin}, amd64, []string{"0.0.61", "0.0.62", "0.0.64"}),
			withGithubReleaseAssetsFromFolder: folderWorking,
			withHost:                          state.Host{Os: linux, Arch: amd64},
			expectBinaryVersionsRunOnce:       []string{"0.0.61", "0.0.62", "0.0.64"},
			withTestRun: func(t *testing.T, tc TestCase, defaultOpts DefaultTestOpts) {
				// Given
				mockHTTPResponsesForGithubReleases(t, folderWorking, tc.withGithubReleases)
				defer httpmock.DeactivateAndReset()

				stdOutBuffer := new(bytes.Buffer)
				defaultOpts.setStdOut(stdOutBuffer)

				upgrader, err := upgrade.New(defaultOpts.Opts)
				require.NoError(t, err)

				// When running first time
				err = upgrader.Run()
				assert.NoError(t, err)

				// Then
				t.Log(defaultOpts.StdOutBuffer.String())

				assert.NoError(t, err)

				originaltestName := tc.name
				tc.name = originaltestName + "_run1"

				doAsserts(t, tc, defaultOpts)

				// Given
				stdOutBuffer = new(bytes.Buffer)
				defaultOpts.setStdOut(stdOutBuffer)

				upgrader, err = upgrade.New(defaultOpts.Opts)
				require.NoError(t, err)

				// When running second time
				err = upgrader.Run()

				// Then
				t.Log(defaultOpts.StdOutBuffer.String())

				assert.NoError(t, err)

				tc.name = originaltestName + "_run2"
				tc.expectBinaryVersionsRunOnce = []string{}

				doAsserts(t, tc, defaultOpts)
			},
		},
		{
			name:                              "Should run upgrades with version up to and including current okctl version, but no newer",
			withOkctlVersion:                  "0.0.63",
			withOriginalClusterVersion:        "0.0.50",
			withGithubReleases:                createGithubReleases([]string{linux, darwin}, amd64, []string{"0.0.61", "0.0.62", "0.0.63", "0.0.64"}),
			withGithubReleaseAssetsFromFolder: folderWorking,
			withHost:                          state.Host{Os: linux, Arch: amd64},
			expectBinaryVersionsRunOnce:       []string{"0.0.61", "0.0.62", "0.0.63"},
		},
		{
			name:                              "Should not run upgrades that are older than the first installed okctl version",
			withOkctlVersion:                  "0.0.64",
			withOriginalClusterVersion:        "0.0.62",
			withGithubReleases:                createGithubReleases([]string{linux, darwin}, amd64, []string{"0.0.61", "0.0.62", "0.0.63", "0.0.64"}),
			withGithubReleaseAssetsFromFolder: folderWorking,
			withHost:                          state.Host{Os: linux, Arch: amd64},
			expectBinaryVersionsRunOnce:       []string{"0.0.63", "0.0.64"},
		},
		{
			name:                              "Should print correct debug output",
			withDebug:                         true,
			withOkctlVersion:                  "0.0.63",
			withOriginalClusterVersion:        "0.0.61",
			withGithubReleases:                createGithubReleases([]string{linux, darwin}, amd64, []string{"0.0.61", "0.0.62", "0.0.63", "0.0.64"}),
			withGithubReleaseAssetsFromFolder: folderWorking,
			withHost:                          state.Host{Os: linux, Arch: amd64},
			expectBinaryVersionsRunOnce:       []string{"0.0.62", "0.0.63"},
		},
		{
			name:                       "Should run upgrade hot fixes, and in correct order",
			withOkctlVersion:           "0.0.63",
			withOriginalClusterVersion: "0.0.50",
			withGithubReleases: createGithubReleases([]string{linux, darwin}, amd64,
				[]string{"0.0.63.a", "0.0.62", "0.0.62.b", "0.0.61", "0.0.62.a", "0.0.64.a", "0.0.63"}),
			withGithubReleaseAssetsFromFolder: folderWorking,
			withHost:                          state.Host{Os: linux, Arch: amd64},
			expectBinaryVersionsRunOnce:       []string{"0.0.61", "0.0.62", "0.0.62.a", "0.0.62.b", "0.0.63", "0.0.63.a"},
		},
		{
			name:                       "Should not run upgrades, including hot fixes, that are older than the first installed okctl version",
			withOkctlVersion:           "0.0.63",
			withOriginalClusterVersion: "0.0.62",
			withGithubReleases: createGithubReleases([]string{linux, darwin}, amd64,
				[]string{"0.0.62", "0.0.62.a", "0.0.62.b", "0.0.63", "0.0.63.a", "0.0.64.a"}),
			withGithubReleaseAssetsFromFolder: folderWorking,
			withHost:                          state.Host{Os: linux, Arch: amd64},
			expectBinaryVersionsRunOnce:       []string{"0.0.63", "0.0.63.a"},
		},
		{
			// Explanation: The whole point of hotfixes are that you may be in the situation where okctl already has
			// upgraded you to version 0.0.63. But then we discover that we made an error in the upgrade for 0.0.62.
			// So we make a hotfix "0.0.62.a" which will be run even though you already upgraded to 0.0.63.
			name:                              "Should run a hotfix even if it is older than the last applied upgrade",
			withOkctlVersion:                  "0.0.63",
			withOriginalClusterVersion:        "0.0.50",
			withGithubReleaseAssetsFromFolder: folderWorking,
			withHost:                          state.Host{Os: linux, Arch: amd64},
			withTestRun: func(t *testing.T, tc TestCase, defaultOpts DefaultTestOpts) {
				// Given settings for first upgrade
				githubReleaseVersions := []string{"0.0.61", "0.0.62", "0.0.63"}
				githubReleases := createGithubReleases([]string{linux, darwin}, amd64, githubReleaseVersions)

				defaultOpts.GithubService = newGithubServiceMock(githubReleases)

				mockHTTPResponsesForGithubReleases(t, folderWorking, githubReleases)
				defer httpmock.DeactivateAndReset()

				stdOutBuffer := new(bytes.Buffer)
				defaultOpts.setStdOut(stdOutBuffer)

				upgrader, err := upgrade.New(defaultOpts.Opts)
				require.NoError(t, err)

				// When running upgrade first time
				err = upgrader.Run()

				// Then
				assert.NoError(t, err)
				t.Log(stdOutBuffer.String())

				originaltestName := tc.name

				tc.expectBinaryVersionsRunOnce = githubReleaseVersions
				tc.name = originaltestName + "_run1"

				doAsserts(t, tc, defaultOpts)

				httpmock.DeactivateAndReset()

				// Given settings for second upgrade
				githubReleaseVersions = []string{"0.0.61", "0.0.62", "0.0.62.a", "0.0.62.b", "0.0.63", "0.0.63.a", "0.0.64"}
				githubReleases = createGithubReleases([]string{linux, darwin}, amd64, githubReleaseVersions)

				defaultOpts.GithubService = newGithubServiceMock(githubReleases)

				mockHTTPResponsesForGithubReleases(t, folderWorking, githubReleases)

				stdOutBuffer = new(bytes.Buffer)
				defaultOpts.setStdOut(stdOutBuffer)

				upgrader, err = upgrade.New(defaultOpts.Opts)
				require.NoError(t, err)

				// When running upgrade second time
				err = upgrader.Run()

				// Then
				t.Log(stdOutBuffer.String())

				assert.NoError(t, err)

				tc.name = originaltestName + "_run2"
				tc.expectBinaryVersionsRunOnce = []string{"0.0.62.a", "0.0.62.b", "0.0.63.a"}

				doAsserts(t, tc, defaultOpts)
			},
		},
		{
			name:                       "Should return error if github release isn't valid",
			withOkctlVersion:           "0.0.61",
			withOriginalClusterVersion: "0.0.50",
			withGithubReleases: []*github.RepositoryRelease{
				{
					ID: github.Int64Ptr(123),
				},
			},
			withGithubReleaseAssetsFromFolder: folderWorking,
			withHost:                          state.Host{Os: linux, Arch: amd64},
			withTestRun: func(t *testing.T, tc TestCase, defaultOpts DefaultTestOpts) {
				upgrader, err := upgrade.New(defaultOpts.Opts)
				require.NoError(t, err)

				err = upgrader.Run()

				assert.Error(t, err)
				assert.Contains(t, err.Error(),
					"release 'ID: 123' must have at least 2 assets (binary and checksum); "+
						"name: cannot be blank; tag_name: cannot be blank.",
				)
			},
		},
		{
			name:                       "Should return error if github release assets don't include checksum",
			withOkctlVersion:           "0.0.64",
			withOriginalClusterVersion: "0.0.50",
			withGithubReleases: []*github.RepositoryRelease{
				{
					ID:      github.Int64Ptr(123),
					TagName: github.StringPtr("0.0.61"),
					Name:    github.StringPtr("0.0.61"),
					Assets: []*github.ReleaseAsset{
						createGihubReleaseAssetBinary(linux, amd64, "0.0.61"),
						createGihubReleaseAssetBinary(darwin, amd64, "0.0.61"),
					},
				},
			},
			withGithubReleaseAssetsFromFolder: folderWorking,
			withHost:                          state.Host{Os: linux, Arch: amd64},
			withTestRun: func(t *testing.T, tc TestCase, defaultOpts DefaultTestOpts) {
				upgrader, err := upgrade.New(defaultOpts.Opts)
				require.NoError(t, err)

				err = upgrader.Run()

				assert.Error(t, err)
				assert.Contains(t, err.Error(),
					"parsing upgrade binaries: validating release: could not find checksum asset for "+
						"release 0.0.61 (assets: okctl_upgrade-linux_amd64_0.0.61.tar.gz,"+
						"okctl_upgrade-darwin_amd64_0.0.61.tar.gz)",
				)
			},
		},
		{
			name:                       "Should return error if release version does not match release download URL version",
			withOkctlVersion:           "0.0.64",
			withOriginalClusterVersion: "0.0.50",
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
								"0.0.61", "0.0.61", linux, amd64)),
						},
						{
							Name: github.StringPtr("okctl_upgrade-darwin_amd64_0.0.61.tar.gz"),
							BrowserDownloadURL: github.StringPtr(fmt.Sprintf(
								"https://github.com/oslokommune/okctl-upgrade/releases/download/%s/okctl-upgrade_%s_%s_%s.tar.gz",
								"0.0.61", "0.0.61", darwin, amd64)),
						},
					},
				},
			},
			withGithubReleaseAssetsFromFolder: folderWorking,
			withHost:                          state.Host{Os: linux, Arch: amd64},
			withTestRun: func(t *testing.T, tc TestCase, defaultOpts DefaultTestOpts) {
				upgrader, err := upgrade.New(defaultOpts.Opts)
				require.NoError(t, err)

				err = upgrader.Run()

				assert.Error(t, err)
				assert.Contains(t, err.Error(),
					"parsing upgrade binaries: validating release: could not find checksum asset for "+
						"release 0.0.61 (assets: okctl_upgrade-linux_amd64_0.0.61.tar.gz,"+
						"okctl_upgrade-darwin_amd64_0.0.61.tar.gz)",
				)
			},
		},
		{
			// In case we accidentally publish an upgrade binary that crashes, the solution is to delete it and replace
			// with a hotfix.
			name:                       "Should support replacing an erroneous upgrade binary with a hotfix",
			withOkctlVersion:           "0.0.65",
			withOriginalClusterVersion: "0.0.50",
			withHost:                   state.Host{Os: linux, Arch: amd64},
			withTestRun: func(t *testing.T, tc TestCase, defaultOpts DefaultTestOpts) {
				// Given configuration for first run
				githubReleaseVersions := []string{"0.0.61", "0.0.62", "0.0.63"}
				githubReleases := createGithubReleases([]string{linux, darwin}, amd64, githubReleaseVersions)

				defaultOpts.GithubService = newGithubServiceMock(githubReleases)

				mockHTTPResponseForGithubRelease(t, path.Join(folderWorking, "0.0.61"), githubReleases[0])
				mockHTTPResponseForGithubRelease(t, path.Join(folderCrashing, "0.0.62"), githubReleases[1])
				mockHTTPResponseForGithubRelease(t, path.Join(folderWorking, "0.0.63"), githubReleases[2])
				httpmock.Activate()
				defer httpmock.DeactivateAndReset()

				stdOutBuffer := new(bytes.Buffer)
				defaultOpts.setStdOut(stdOutBuffer)

				upgrader, err := upgrade.New(defaultOpts.Opts)
				require.NoError(t, err)

				// When 0.0.61 runs OK and 0.0.62 fails
				err = upgrader.Run()

				// Then
				t.Log(defaultOpts.StdOutBuffer.String())

				assert.Error(t, err)
				assert.Contains(t, err.Error(), "It will crash")

				originaltestName := tc.name

				// Even if 0.0.62 crashed, we still expect the binary to have been run.
				tc.expectBinaryVersionsRunOnce = []string{"0.0.61", "0.0.62"}

				// 0.0.61 completed successfully, so we expect the cluster version to be 0.0.61.
				tc.expectedClusterVersionAfterUpgrade = "0.0.61"

				tc.name = originaltestName + "_run1"

				doAsserts(t, tc, defaultOpts)

				httpmock.DeactivateAndReset()

				// Given configuration for second run
				githubReleaseVersions = []string{"0.0.61", "0.0.62.a", "0.0.63"}
				githubReleases = createGithubReleases([]string{linux, darwin}, amd64, githubReleaseVersions)

				defaultOpts.GithubService = newGithubServiceMock(githubReleases)

				mockHTTPResponseForGithubRelease(t, path.Join(folderWorking, "0.0.61"), githubReleases[0])
				mockHTTPResponseForGithubRelease(t, path.Join(folderWorking, "0.0.62.a"), githubReleases[1])
				mockHTTPResponseForGithubRelease(t, path.Join(folderWorking, "0.0.63"), githubReleases[2])
				httpmock.Activate()
				defer httpmock.DeactivateAndReset()

				stdOutBuffer = new(bytes.Buffer)
				defaultOpts.setStdOut(stdOutBuffer)

				upgrader, err = upgrade.New(defaultOpts.Opts)
				require.NoError(t, err)

				// When
				err = upgrader.Run()

				// Then
				t.Log(defaultOpts.StdOutBuffer.String())

				assert.NoError(t, err)

				tc.name = originaltestName + "_run2"
				tc.expectBinaryVersionsRunOnce = []string{"0.0.62.a", "0.0.63"}

				// All upgrades complteted successfully, so we expect cluster version to be the current okctl version
				tc.expectedClusterVersionAfterUpgrade = tc.withOkctlVersion

				doAsserts(t, tc, defaultOpts)
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

			if len(tc.withClusterVersion) == 0 {
				tc.withClusterVersion = tc.withOriginalClusterVersion
			}

			upgradeState := testutils.MockUpgradeState(tc.withClusterVersion)
			clusterState := testutils.MockClusterState(tc.withOriginalClusterVersion)

			clusterVersioner := clusterversion.New(stdOutBuffer, api.ID{}, upgradeState)
			originalClusterVersioner := originalclusterversion.New(api.ID{}, upgradeState, clusterState)

			surveyor := testutils.NewAutoAnsweringSurveyor(tc.withUserAnswers)

			defaultOpts := DefaultTestOpts{
				Opts: upgrade.Opts{
					Debug:                    tc.withDebug,
					Logger:                   logrus.StandardLogger(),
					Out:                      stdOutBuffer,
					RepositoryDirectory:      repositoryAbsoluteDir,
					GithubService:            newGithubServiceMock(tc.withGithubReleases),
					ChecksumDownloader:       upgrade.NewChecksumDownloader(),
					ClusterVersioner:         clusterVersioner,
					OriginalClusterVersioner: originalClusterVersioner,
					Surveyor:                 surveyor,
					FetcherOpts: upgrade.FetcherOpts{
						Host:  tc.withHost,
						Store: tmpStore,
					},
					OkctlVersion: tc.withOkctlVersion,
					State:        upgradeState,
					ClusterID:    api.ID{},
				},
				StdOutBuffer: stdOutBuffer,
			}

			// When
			if tc.withTestRun != nil {
				tc.withTestRun(t, tc, defaultOpts)
				return
			}

			// Or when
			mockHTTPResponsesForGithubReleases(t, tc.withGithubReleaseAssetsFromFolder, tc.withGithubReleases)
			defer httpmock.DeactivateAndReset()

			upgrader, err := upgrade.New(defaultOpts.Opts)
			require.NoError(t, err)

			err = upgrader.Run()

			t.Log(stdOutBuffer.String())

			// Then
			if len(tc.expectErrorContains) > 0 {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectErrorContains)
				return
			}

			assert.NoError(t, err)

			doAsserts(t, tc, defaultOpts)
		})
	}
}

func doAsserts(t *testing.T, tc TestCase, defaultOpts DefaultTestOpts) {
	// Upgrades
	expectedUpgradesRun := getExpectedUpgradesRun(tc.expectBinaryVersionsRunOnce, tc.withHost)
	upgradesRun := getActualUpgradesRun(defaultOpts.StdOutBuffer)
	assert.Equal(t, expectedUpgradesRun, upgradesRun, tc.name+": Unexpected upgrades were run")

	// Cluster version
	clusterVersion, err := defaultOpts.ClusterVersioner.GetClusterVersion()
	require.NoError(t, err, tc.name)

	if len(tc.expectedClusterVersionAfterUpgrade) > 0 {
		assert.Equal(t, tc.expectedClusterVersionAfterUpgrade, clusterVersion, tc.name)
	} else {
		assert.Equal(t, tc.withOkctlVersion, clusterVersion, tc.name)
	}

	// Original version
	originalClusterVersion, err := defaultOpts.OriginalClusterVersioner.GetOriginalClusterVersion()
	require.NoError(t, err, tc.name)

	assert.Equal(t, originalClusterVersion, originalClusterVersion, tc.name)

	// Goldie
	g := goldie.New(t)
	t.Log(tc.name)
	g.Assert(t, tc.name, defaultOpts.StdOutBuffer.Bytes())
}

type DefaultTestOpts struct {
	upgrade.Opts
	StdOutBuffer *bytes.Buffer
}

func (o *DefaultTestOpts) setStdOut(stdOut *bytes.Buffer) {
	o.StdOutBuffer = stdOut
	o.Out = stdOut
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
	upgradesRun := make(map[string]bool)

	actualUpgradesRun := make([]string, 0)

	for _, match := range found {
		// match[1] is the captured regex group, e.g.: okctl-upgrade-0.0.64-amd64-linux
		_, ok := upgradesRun[match[1]]
		if !ok {
			// First match is the simulation
			upgradesRun[match[1]] = true
		} else {
			// Second match is the actual run
			actualUpgradesRun = append(actualUpgradesRun, match[1])
		}
	}

	return actualUpgradesRun
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
			Name:        github.StringPtr(upgrade.ChecksumsTxt),
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

	return releases
}

func createGihubReleaseAssetBinary(os, arch, version string) *github.ReleaseAsset {
	return &github.ReleaseAsset{
		Name:        github.StringPtr(fmt.Sprintf("okctl_upgrade-%s_%s_%s.tar.gz", os, arch, version)),
		ContentType: github.StringPtr("application/gzip"),
		BrowserDownloadURL: github.StringPtr(fmt.Sprintf(
			"https://github.com/oslokommune/okctl-upgrade/releases/download/%s/okctl-upgrade_%s_%s_%s.tar.gz", version, version, os, arch)),
	}
}

// capitalizeFirst converts for instance "linux" to "Linux". We use this because we expect github release assets for
// upgrades to be named this way.
func capitalizeFirst(os string) string {
	return strings.ToUpper(os[0:1]) + strings.ToLower(os[1:])
}
