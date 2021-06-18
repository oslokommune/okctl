package upgrade

import (
	"fmt"
	"github.com/oslokommune/okctl/pkg/binaries/fetch"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/config/state"
	"github.com/oslokommune/okctl/pkg/storage"
	"github.com/sirupsen/logrus"
	"io"
)

/*

How to do upgrade - long version

USE DOMAIN DRIVEN DESIGN FOR MIGRATION CALCULATION LOGIC

- [API-call] get all releases from okctl-upgrade repo. ->
{
	"tag_name": "1-0.0.63",
	"assets": [
		"browser_download_url":  "https://github.com/oslokommune/okctl-upgrade/releases/download/1_0.0.64/okctl-upgrade-1-0.0.63"
	],
}

- Fjern de som ikke har asset.state == "uploaded"
	- Kast error dersom fil ikke finnes? Tjaa. Test hva som skjer dersom fil ikke finnes.

- Parse into struct
Migration {
	Index 			int,
	Filename		string,
	Version			SemverVersion

}
	->
		0.0.61
		0.0.62
		0.0.62_a
		0.0.63
		0.0.63_a
		0.0.64
		0.0.65

		Betydning: x-<okctl version> - Migrasjon som kjøres for å komme seg til <okctl version>. Så for å komme seg til
			versjon 0.0.62, må man kjøre 0.0.61 og 0.0.62.

- [okctl-call]
	get current okctl version - 0.0.62
	okctl apply was run with 0.0.62.

- [state-call] find already applied migrations (from state.db)
	if empty
		add all migrations up to and including current okctl versions to applied migrations table:
			0.0.61 (applied_at: nil, comment: Doesn't need to be run.)
			0.0.62
			0.0.62_a

- calculate which migrations to run
	- remove too new migrations
		0.0.61
		0.0.62
		0.0.62_a
	- remove too old migrations
		logic: okctl original version = 0.0.62, current = 0.0.62. if (migrationVersion <= okctlOriginalVersion) continue / skip;
		(empty)
	- remove applied
		(empty)

	since okctl version is 0.0.62, and we have applied all migrations up to 0.0.62, we're done.

- user downloads okctl 0.0.64

- [state-call] find already applied migrations (from state.db)7
			0.0.61
			0.0.62
			0.0.62_a
-

- calculate which migrations to run
	- get list of all migrations
		(see list above)
    - remove too new migrations (all upgrades having version number higher than current okctl version (semver compare))
		0.0.61
		0.0.62
		0.0.62_a
		0.0.63
		0.0.63_a
		0.0.64
	- remove too old migrations
	- remove already applied migrations
		0.0.63
		0.0.63_a
		0.0.64

- [okctl-call binaries-provider] Download upgrade binaries for migrations to run
	->	okctl-upgrade_4-0.0.63_linux_amd64
		okctl-upgrade_5-0.0.63_linux_amd64
		okctl-upgrade_6-0.0.64_linux_amd64

- Verify checksum

- [shell-call] run resulting migrations

	- If exit code 0, update upgrade state: Ran OK.
	- If exit code > 0, update upgrade state: Ran not-OK.
		Write result to state. MigrationResult struct or something.


How to require a specific version:
- Feks: okctl 0.64: apply cluster: sjekkom upgrade state indexOfLatestAppliedMigration >= 5.
	If OK, continue apply cluster.
	If not OK, write error message:
		Some okctl resources are out of date. You need to:
			upgrade your okctl resources by running 'okctl upgrade'. To see what this will do, run
			'okctl upgrade --dry-run'.

			Or, if you don't want to upgrade at this time, download a previous version of okctl and try again.

		Technical details:
			apply cluster requires minimum upgrade version: 5-0.0.63 [or just 5]
*/

type Upgrader struct {
	progress            io.Writer
	logger              *logrus.Logger
	githubService       client.GithubService
	githubReleaseParser GithubReleaseParser
	fetcherOpts         FetcherOpts
}

//
//type okctlUpgradeBinary struct {
//	filenameWithoutExtension string
//	version                  string
//	checksums                []state.Checksum
//}

// upgrade is an okctl upgrade. A example file name for an upgrade is 'okctl-upgrade_0.0.63_Darwin_amd64.tar.gz'.
//
// The semantics of the meaning of the file name and how upgrades are supposed to be run, is explained at
// https://github.com/oslokommune/okctl-upgrade
type okctlUpgradeBinary struct {
	name          string
	fileExtension string
	version       string
	checksums     []state.Checksum
}

func (u Upgrader) Run() error {
	releases, err := u.githubService.ListReleases("oslokommune", "okctl-upgrade")
	if err != nil {
		return fmt.Errorf("listing github releases: %w", err)
	}

	// [x] Last ned liste av releases
	// [ ] Filtrer
	// [ ] Last ned <--- NÅ
	// [ ] Kjør

	// Convert to upgrades
	upgradeBinaries, err := u.githubReleaseParser.toUpgradeBinaries(releases)
	if err != nil {
		return fmt.Errorf("parsing upgrade binaries: %w", err)
	}

	// TODO consider parsing to okctlUpgrade, that contains UpgradeBinary.

	// TODO implement filtering after testing actual functionality
	//// Remove too new upgrades (semver compare)
	//upgradeBinaries = u.removeTooNewMigrations(upgradeBinaries)
	//
	//// Remove upgrades that are earlier than than or equal to original cluster version
	//upgradeBinaries = u.removeTooOldUpgrades(upgradeBinaries)
	//
	//// Remove already applied upgrades (from state.db)
	//upgradeBinaries = u.removeAlreadyAppliedUpgrades(upgradeBinaries)

	// Downlod upgrade binaries
	// Runs binaries. Saves that they have been run. Handles errors.
	//migrator := u.newMigrator(binaryRunner, state, migrations)
	//migrator.Run()

	// Verify checksum

	// Run resulting migrations, and store that they have been run

	var binaries []state.Binary

	for _, upgradeBinary := range upgradeBinaries {
		URLPattern := fmt.Sprintf(
			"https://github.com/oslokommune/okctl-upgrade/releases/download/%s/okctl-upgrade_%s_#{os}_#{arch}.tar.gz",
			upgradeBinary.version,
			upgradeBinary.version,
		)

		binary := state.Binary{
			Name:       upgradeBinary.name,
			Version:    upgradeBinary.version,
			BufferSize: "300mb",
			URLPattern: URLPattern,
			Archive: state.Archive{
				Type:   upgradeBinary.fileExtension,
				Target: upgradeBinary.name,
			},
			Checksums: upgradeBinary.checksums,
		}

		binaries = append(binaries, binary)
	}

	// Download binaries
	fetcher, err := fetch.New(
		u.progress,
		u.logger,
		true,
		u.fetcherOpts.Host,
		binaries,
		u.fetcherOpts.Store,
	)
	if err != nil {
		return fmt.Errorf("creating upgrade binaries fetcher: %w", err)
	}

	binaryProvider := newUpgradeBinaryProvider(u.logger, u.progress, fetcher)

	// TODO verify that binary for current os and arch is run

	for _, binary := range upgradeBinaries {
		upgradeBinary, err := binaryProvider.OkctlUpgrade(binary.version)
		if err != nil {
			return fmt.Errorf("getting okctl upgrade binary: %w", err)
		}

		// something something, binary runner or something. See eksctl.go.
		upgradeBinary.Run()
	}

	return nil
}

type FetcherOpts struct {
	Host  state.Host
	Store storage.Storer
}

func NewUpgrader(
	logger *logrus.Logger,
	progress io.Writer,
	githubService client.GithubService,
	githubReleaseParser GithubReleaseParser,
	fetcherOpts FetcherOpts,
) Upgrader {
	return Upgrader{
		progress:            progress,
		logger:              logger,
		githubService:       githubService,
		githubReleaseParser: githubReleaseParser,
		fetcherOpts:         fetcherOpts,
	}
}
