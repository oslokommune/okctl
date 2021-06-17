package upgrade

import (
	"errors"
	"fmt"
	binariesPkg "github.com/oslokommune/okctl/pkg/binaries"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/config/state"
	"github.com/oslokommune/okctl/pkg/github"
	"strings"
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
	githubService    client.GithubService
	binaryService    client.BinaryService
	binariesProvider binariesPkg.Provider
	checksumFetcher  checksumFetcher
}

type upgradeBinary struct {
	filename                 string
	filenameWithoutExtension string
	version                  string
	checksums                []state.Checksum
}

func (u Upgrader) Run() error {
	releases, err := u.githubService.ListReleases("oslokommune", "okctl-upgrade")
	if err != nil {
		return fmt.Errorf("listing github releases: %w", err)
	}

	// TODO take care of linux vs mac, amd64 vs x86 somewhere

	// [x] Last ned liste av releases
	// [ ] Filtrer
	// [ ] Last ned <--- NÅ
	// [ ] Kjør

	// Convert to upgrades
	upgradeBinaries, err := u.parseUpgradeBinaries(releases)
	if err != nil {
		return fmt.Errorf("parsing upgrade binaries: %w", err)
	}

	// Remove too new upgrades (semver compare)
	upgradeBinaries = u.removeTooNewMigrations(upgradeBinaries)

	// Remove upgrades that are earlier than than or equal to original cluster version
	upgradeBinaries = u.removeTooOldUpgrades(upgradeBinaries)

	// Remove already applied upgrades (from state.db)
	upgradeBinaries = u.removeAlreadyAppliedUpgrades(upgradeBinaries)

	// Downlod upgrade binaries
	// Runs binaries. Saves that they have been run. Handles errors.
	//migrator := u.newMigrator(binaryRunner, state, migrations)
	//migrator.Run()

	// Verify checksum

	// Run resulting migrations, and store that they have been run

	for _, ub := range upgradeBinaries {
		URLPattern := fmt.Sprintf(
			"https://github.com/oslokommune/okctl-upgrade/releases/download/%s/okctl-upgrade_%s_#{os}_#{arch}",
			ub.version,
			ub.version,
		)

		binary := state.Binary{
			Name:       ub.filenameWithoutExtension,
			Version:    ub.version,
			BufferSize: "300mb",
			URLPattern: URLPattern,
			Archive: state.Archive{
				Type:   ".tar.gz",
				Target: ub.filenameWithoutExtension,
			},
			Checksums: ub.checksums,
			Preload:   true,
		}

		err = u.binaryService.Add(binary)
		if err != nil {
			return fmt.Errorf("adding binary %s: %w", binary.Id(), err)
		}
	}

	// Download binaries
	err = u.binariesProvider.ReloadBinaries()
	if err != nil {
		return fmt.Errorf("loading binaries: %w", err)
	}

	for _, binary := range upgradeBinaries {
		upgradeBinary := u.binariesProvider.OkctlUpgrade(binary.GitReleaseTag)
		// something something, binary runner or something. See eksctl.go.
		upgradeBinary.Run()
	}

	return nil
}

func (u Upgrader) parseUpgradeBinaries(releases []*github.RepositoryRelease) ([]upgradeBinary, error) {
	var binaries []upgradeBinary

	for _, release := range releases {
		err := u.validateRelease(release)
		if err != nil {
			return nil, fmt.Errorf("validating release: %w", err)
		}

		checksums, err := u.checksumFetcher.getFor(release)
		if err != nil {
			return nil, fmt.Errorf("fetching checksum: %w", err)
		}

		binaries = append(binaries, upgradeBinary{
			filenameWithoutExtension: "",
			version:                  *release.TagName,
			checksums:                checksums,
		})
	}

	return binaries, nil
}

func (u Upgrader) validateRelease(release *github.RepositoryRelease) error {
	if release.Name == nil || len(*release.Name) == 0 {
		return fmt.Errorf("release ID '%d' name must be non-empty", *release.ID)
	}

	if release.TagName == nil {
		return fmt.Errorf("release '%s' tag name must be non-empty", *release.Name)
	}

	if len(release.Assets) < 2 {
		return fmt.Errorf("release '%s' must have at least two assets (binary and checksum) ", *release.Name)
	}

	for _, asset := range release.Assets {
		parts := strings.Split(*asset.BrowserDownloadURL, "/")
		if len(parts) == 0 {
			return fmt.Errorf("expected at least 1 '/' in browser download URL %s", *asset.BrowserDownloadURL)
		}

		filename := parts[len(parts)-1]

		upgradeFile, err := parseOkctlUpgradeFilename(filename)
		if err != nil {
			return fmt.Errorf("cannot parse okctl upgrade filename '%s': %w", filename, err)
		}

		if upgradeFile.version != *release.TagName {
			return fmt.Errorf("expected upgrade file '%s' version to equal tag name '%s'", upgradeFile.version, *release.TagName)
		}
	}

	return nil
}

func NewUpgrader(githubService client.GithubService, binaryService client.BinaryService, binariesProvider binariesPkg.Provider, checksumFetcher checksumFetcher) Upgrader {
	return Upgrader{
		githubService:    githubService,
		binaryService:    binaryService,
		binariesProvider: binariesProvider,
		checksumFetcher:  checksumFetcher,
	}
}
