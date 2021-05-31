package upgrade

import (
	"fmt"
	"github.com/oslokommune/okctl/pkg/config/state"
	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/oslokommune/okctl/pkg/version"
)

/*
How to do upgrade

USE DOMAIN DRIVEN DESIGN FOR MIGRATION CALCULATION LOGIC

- [API-call] get all releases from okctl-upgrade repo. ->
{
	"tag_name": "1-0.0.63",
	"assets": [
		"browser_download_url":  "https://github.com/oslokommune/okctl-upgrade/releases/download/1_0.0.64/okctl-upgrade-1-0.0.63"
	],
	"body" // TODO nice to have, add later
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
		1-0.0.61
		2-0.0.62
		3-0.0.62
		4-0.0.63
		5-0.0.63
		6-0.0.64
		7-0.0.65

		Betydning: x-<okctl version> - Migrasjon som kjøres for å komme seg til <okctl version>. Så for å komme seg til
			versjon 0.0.62, må man kjøre 1-0.0.61 og 2-0.0.62.

- [okctl-call]
	get current okctl version - 0.0.62

- [state-call] find already applied migrations (from state.db)
	if empty
		add all migrations up to and including current okctl versions to applied migrations table:
			1-0.0.61 (applied_at: nil, comment: Doesn't need to be run.)
			2-0.0.62
			3-0.0.62

- since okctl version is 0.0.62, and we have applied all migrations up to 0.0.62, we're done.

- user upgrades to 0.0.64

- [state-call] find already applied migrations (from state.db)
			1-0.0.61
			2-0.0.62
			3-0.0.62

- calculate which migrations to run
	- parse migration number from all tags
		-> 1, 2, 3, 4, 5, 6, 7
	- remove all upgrades that have migration numbers that already have been applied
		- > 4, 5, 6, 7
	- remove all upgrades have version number higher than current okctl version (semver compare)
		-> 4, 5, 6

- [okctl-call binaries-provider] Download migration binaries for migrations to run
	->	okctl-upgrade_4-0.0.63_linux_amd64
		okctl-upgrade_5-0.0.63_linux_amd64
		okctl-upgrade_6-0.0.64_linux_amd64

- Verify checksum

- [shell-call] run resulting migrations

	- If exit code 0, update migration state: Ran OK.
	- If exit code > 0, update migration state: Ran not-OK.
		Write result to state. MigrationResult struct or something.


How to require a specific version:
- Feks: okctl 0.64: apply cluster: sjekkom migration state indexOfLatestAppliedMigration >= 5.
	If OK, continue apply cluster.
	If not OK, write error message:
		Some okctl resources are out of date. You need to:
			upgrade your okctl resources by running 'okctl upgrade'. To see what this will do, run
			'okctl upgrade --dry-run'.

			Or, if you don't want to upgrade at this time, download a previous version of okctl and try again.

		Technical details:
			apply cluster requires minimum migration version: 5-0.0.63 [or just 5]
*/

func Run(o *okctl.Okctl) error {
	// TODO Figure out how to add x migrations to binariesprovider

	o.BinariesProvider.OkctlMigration()

	migration := "1" + version.String()
	o.BinariesProvider.OkctlMigration(migration)

	// Somehow calculate that these are the migrations we want. Might want to change use a Migration struct instead of
	// GithubRelease at some point.
	migrations := []Migration{
		{
			Version:            "4",
			OkctlTargetVersion: "0.0.64",
		},
	}

	upgradeBinaries := make([]UpgradeBinary, len(migrations))

	for _, migration := range migrations {
		upgradeBinaries = append(upgradeBinaries, UpgradeBinary{
			Migration:     migration,
			GitReleaseTag: fmt.Sprintf("%s_%s", migration.Version, migration.OkctlTargetVersion),
		})
	}

	for _, binary := range upgradeBinaries {
		URLPattern := fmt.Sprintf(
			"https://github.com/oslokommune/okctl-upgrade/releases/download/%s/okctl-upgrade_%s_#{os}_#{arch}",
			binary.GitReleaseTag,
			binary.GitReleaseTag,
		)

		o.Binaries().Add(state.Binary{ // For this to work, we need to call fetcher.prepareAndLoad again
			Name:       binary.Migration.Name,
			Version:    binary.GitReleaseTag,
			BufferSize: "300mb",
			//URLPattern: "https://github.com/oslokommune/okctl-upgrade/releases/download/3_0.0.64/okctl-upgrade_4_0.0.63_#{os}_#{arch}",
			URLPattern: URLPattern,
			Archive:    state.Archive{},
			Checksums:  nil,
			Preload:    false,
		})

		upgradeBinary := o.BinariesProvider.OkctlMigration(binary.GitReleaseTag)
		upgradeBinary.Run()
	}

	return nil
}

type Migration struct {
	Name               string
	Version            string
	OkctlTargetVersion string
}

type UpgradeBinary struct {
	Migration     Migration
	GitReleaseTag string
}
