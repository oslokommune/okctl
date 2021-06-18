package upgrade

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
/*
func TestRunMigrations(t *testing.T) {
	testCases := []struct {
		name string

		withUpgradeGithubReleases    []*github.RepositoryRelease
		withOriginalOkctlVersion     string
		withOkctlVersion             string
		withAlreadyAppliedMigrations []string // TODO: or upgrade?
		expectMigrationsToBeRun      []string // TODO: or upgrade?
	}{
		{
			name: "Should run an upgrade",
			withUpgradeGithubReleases: []*github.RepositoryRelease{
				{
					TagName: github.StringPtr("0.0.63"),
					Name:    github.StringPtr("0.0.63"),
					Assets: []*github.ReleaseAsset{
						{
							Name:               github.StringPtr("okctl_upgrade-0.0.63_Darwin_amd64.tar.gz"),
							ContentType:        github.StringPtr("application/gzip"),
							BrowserDownloadURL: github.StringPtr("https://github.com/oslokommune/okctl-okctlUpgrade/releases/download/0.0.63/okctl_upgrade-0.0.63_Darwin_amd64.tar.gz"),
						},
					},
				},
			},
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
			githubServiceMock := NewGithubServiceMock(tc.withUpgradeGithubReleases)
			upgrader := NewUpgrader(githubServiceMock, nil)

			// When
			err := upgrader.Run()

			// Then
			assert.NoError(t, err)

			// TODO: Assert that the mock binary runner has ran the upgrade executables we expect
		})
	}
}
*/
