package upgrade

/*
import (
	"bytes"
	"fmt"
	"testing"

	"github.com/oslokommune/okctl/pkg/binaries/fetch"
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
	var releases []*github.RepositoryRelease

	createGithubReleases := func(versions ...string) []*github.RepositoryRelease {
		for _, version := range versions {

			release := &github.RepositoryRelease{
				TagName: &version,
				Name:    &version,
				Assets: []*github.ReleaseAsset{
					{
						Name:        github.StringPtr(fmt.Sprintf("okctl_upgrade-%s_Darwin_amd64.tar.gz", version)),
						ContentType: github.StringPtr("application/gzip"),
						BrowserDownloadURL: github.StringPtr(fmt.Sprintf(
							"https://github.com/oslokommune/okctl-upgrade/releases/download/%s/okctl_upgrade-%s_Darwin_amd64.tar.gz", version)),
					},
				},
			}

			releases = append(releases, release)
		}
	}

	testCases := []struct {
		name string

		withUpgradeGithubReleases    []*github.RepositoryRelease
		withOriginalOkctlVersion     string
		withOkctlVersion             string
		withAlreadyAppliedMigrations []string // TODO: or upgrade?
		expectMigrationsToBeRun      []string // TODO: or upgrade?
	}{
		{
			name:                         "Should run an upgrade",
			withUpgradeGithubReleases:    createGithubReleases("0.0.63"),
			withOriginalOkctlVersion:     "0.0.62",
			withOkctlVersion:             "0.0.62",
			withAlreadyAppliedMigrations: []string{},
			expectMigrationsToBeRun:      []string{"0.0.63"},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// Given
			var buffer bytes.Buffer
			githubServiceMock := NewGithubServiceMock(tc.withUpgradeGithubReleases)

			upgrader := New(Opts{
				Debug:               false,
				Logger:              nil,
				Out:                 &buffer,
				RepoDir:             "",
				GithubService:       nil,
				GithubReleaseParser: GithubReleaseParser{},
				FetcherOpts: FetcherOpts{
					Host: state.Host{
						Os:   "linux",
						Arch: "amd64",
					},
					Store: storage.NewEphemeralStorage(),
				},
			})

			// When
			err := upgrader.Run()

			// Then
			assert.NoError(t, err)

			// TODO: Assert that the mock binary runner has ran the upgrade executables we expect
		})
	}
}
*/
