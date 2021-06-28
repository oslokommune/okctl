// Package upgrade knows how to upgrade okctl
package upgrade

import (
	"fmt"
	"io"

	"github.com/oslokommune/okctl/pkg/binaries/fetch"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/config/state"
	"github.com/oslokommune/okctl/pkg/storage"
	"github.com/sirupsen/logrus"
)

type okctlUpgradeBinary struct {
	name          string
	fileExtension string
	version       string
	checksums     []state.Checksum
}

// Run upgrades okctl
func (u Upgrader) Run() error {
	releases, err := u.GithubService.ListReleases("oslokommune", "okctl-upgrade")
	if err != nil {
		return fmt.Errorf("listing github releases: %w", err)
	}

	// Convert to upgrades
	upgradeBinaries, err := u.GithubReleaseParser.toUpgradeBinaries(releases)
	if err != nil {
		return fmt.Errorf("parsing upgrade binaries: %w", err)
	}

	// DO: Filter

	binaries := make([]state.Binary, 0, len(upgradeBinaries))

	if u.Debug {
		_, _ = fmt.Fprintf(u.Out, "Found %d upgrade(s)\n", len(upgradeBinaries))
	}

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
		u.Out,
		u.Logger,
		true,
		u.FetcherOpts.Host,
		binaries,
		u.FetcherOpts.Store,
	)
	if err != nil {
		return fmt.Errorf("creating upgrade binaries fetcher: %w", err)
	}

	binaryProvider := newUpgradeBinaryProvider(u.RepositoryDirectory, u.Logger, u.Out, fetcher)

	for _, binary := range upgradeBinaries {
		upgradeBinary, err := binaryProvider.OkctlUpgrade(binary.version)
		if err != nil {
			return fmt.Errorf("getting okctl upgrade binary: %w", err)
		}

		upgradeBinary.Debug(u.Debug)

		if u.Debug {
			_, _ = fmt.Fprintf(u.Out, "--- Running upgrade: %s ---\n", binary.version)
		}

		_, err = upgradeBinary.Run()
		if err != nil {
			return fmt.Errorf("running upgrade binary %s: %w", binary.version, err)
		}
	}

	// DO: Store that upgrades have been run.
	// DO: Consider letting upgrades edit state.db instead of okctl upgrade.

	return nil
}

// FetcherOpts contains data needed to initialize a fetch.Provider
type FetcherOpts struct {
	Host  state.Host
	Store storage.Storer
}

// Opts contains all data needed to create an Upgrader
type Opts struct {
	Debug               bool
	Logger              *logrus.Logger
	Out                 io.Writer
	RepositoryDirectory string
	GithubService       client.GithubService
	GithubReleaseParser GithubReleaseParser
	FetcherOpts         FetcherOpts
}

// Upgrader knows how to upgrade okctl
type Upgrader struct {
	Opts
}

// New returns a new Upgrader
func New(opts Opts) Upgrader {
	return Upgrader{
		Opts: opts,
	}
}
