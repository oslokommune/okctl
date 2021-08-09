// Package upgrade knows how to upgrade okctl
package upgrade

import (
	"fmt"
	"io"
	"strings"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/binaries/fetch"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/config/state"
	"github.com/oslokommune/okctl/pkg/storage"
	"github.com/sirupsen/logrus"
)

// Run upgrades okctl
func (u Upgrader) Run() error {
	// Fetch
	releases, err := u.githubService.ListReleases("oslokommune", "okctl-upgrade")
	if err != nil {
		return fmt.Errorf("listing github releases: %w", err)
	}

	upgradeBinaries, err := u.githubReleaseParser.ToUpgradeBinaries(releases)
	if err != nil {
		return fmt.Errorf("parsing upgrade binaries: %w", err)
	}

	printIfDebug(u.debug, u.out, "Found %d upgrade(s):", upgradeBinaries)

	// Filter
	upgradeBinaries, err = u.filter.get(upgradeBinaries)
	if err != nil {
		return fmt.Errorf("filtering upgrade binaries: %w", err)
	}

	// Sort, i.e. determine execution order
	sort(upgradeBinaries)
	printIfDebug(u.debug, u.out, "Found %d applicable upgrade(s):", upgradeBinaries)

	// Run
	err = u.runBinaries(upgradeBinaries)
	if err != nil {
		return fmt.Errorf("running upgrade binaries: %w", err)
	}

	return nil
}

func (u Upgrader) runBinaries(upgradeBinaries []okctlUpgradeBinary) error {
	binaryProvider, err := u.createBinaryProvider(upgradeBinaries)
	if err != nil {
		return fmt.Errorf("creating binary provider: %w", err)
	}

	for _, binary := range upgradeBinaries {
		// Get
		binaryRunner, err := binaryProvider.okctlUpgrade(binary.RawVersion())
		if err != nil {
			return fmt.Errorf("getting okctl upgrade binary: %w", err)
		}

		binaryRunner.SetDebug(u.debug)

		_, _ = fmt.Fprintf(u.out, "--- Running upgrade: %s ---\n", binary)

		_, err = binaryRunner.Run()
		if err != nil {
			_, _ = fmt.Fprintf(u.out, "--- Upgrade failed: %s ---\n", binary)
			return fmt.Errorf("running upgrade binary %s: %w", binary, err)
		}

		err = u.filter.markAsRun(binary)
		if err != nil {
			return fmt.Errorf("marking upgrades as run: %w", err)
		}
	}

	return nil
}

func (u Upgrader) createBinaryProvider(upgradeBinaries []okctlUpgradeBinary) (upgradeBinaryProvider, error) {
	binaries := u.toStateBinaries(upgradeBinaries)

	// fetch.New downloads binaries that are missing
	fetcher, err := fetch.New(
		u.out,
		u.logger,
		true,
		u.fetcherOpts.Host,
		binaries,
		u.fetcherOpts.Store,
	)
	if err != nil {
		return upgradeBinaryProvider{}, fmt.Errorf("creating upgrade binaries fetcher: %w", err)
	}

	binaryProvider := newUpgradeBinaryProvider(u.repositoryDirectory, u.logger, u.out, fetcher)

	return binaryProvider, nil
}

func (u Upgrader) toStateBinaries(upgradeBinaries []okctlUpgradeBinary) []state.Binary {
	binaries := make([]state.Binary, 0, len(upgradeBinaries))

	for _, upgradeBinary := range upgradeBinaries {
		URLPattern := fmt.Sprintf(
			"https://github.com/oslokommune/okctl-upgrade/releases/download/%s/okctl-upgrade_%s_#{os}_#{arch}.tar.gz",
			upgradeBinary.RawVersion(),
			upgradeBinary.RawVersion(),
		)

		binary := state.Binary{
			Name:       upgradeBinary.BinaryName(),
			Version:    upgradeBinary.RawVersion(),
			BufferSize: "300mb",
			URLPattern: URLPattern,
			Archive: state.Archive{
				Type:   upgradeBinary.fileExtension,
				Target: upgradeBinary.BinaryName(),
			},
			Checksums: upgradeBinary.checksums,
		}

		binaries = append(binaries, binary)
	}

	return binaries
}

func printIfDebug(debug bool, out io.Writer, text string, upgradeBinaries []okctlUpgradeBinary) {
	if debug {
		binaries := make([]string, 0)
		for _, binary := range upgradeBinaries {
			binaries = append(binaries, binary.RawVersion())
		}

		joinedBinariesTxt := strings.Join(binaries, ", ")

		_, _ = fmt.Fprintf(out, text+"\n", len(upgradeBinaries))
		_, _ = fmt.Fprintln(out, joinedBinariesTxt)
		_, _ = fmt.Fprintln(out, "")
	}
}

// FetcherOpts contains data needed to initialize a fetch.Provider
type FetcherOpts struct {
	Host  state.Host
	Store storage.Storer
}

// Opts contains all data needed to create an Upgrader
type Opts struct {
	Debug                bool
	Logger               *logrus.Logger
	Out                  io.Writer
	RepositoryDirectory  string
	GithubService        client.GithubService
	ChecksumDownloader   ChecksumHTTPDownloader
	FetcherOpts          FetcherOpts
	OkctlVersion         string
	OriginalOkctlVersion string
	State                client.UpgradeState
	ClusterID            api.ID
}

// Upgrader knows how to upgrade okctl
type Upgrader struct {
	debug               bool
	logger              *logrus.Logger
	out                 io.Writer
	repositoryDirectory string
	githubService       client.GithubService
	githubReleaseParser GithubReleaseParser
	fetcherOpts         FetcherOpts
	filter              filter
}

// New returns a new Upgrader
func New(opts Opts) Upgrader {
	return Upgrader{
		debug:               opts.Debug,
		logger:              opts.Logger,
		out:                 opts.Out,
		repositoryDirectory: opts.RepositoryDirectory,
		githubService:       opts.GithubService,
		githubReleaseParser: NewGithubReleaseParser(opts.ChecksumDownloader),
		fetcherOpts:         opts.FetcherOpts,
		filter: filter{
			debug:                opts.Debug,
			out:                  opts.Out,
			state:                opts.State,
			clusterID:            opts.ClusterID,
			okctlVersion:         opts.OkctlVersion,
			originalOkctlVersion: opts.OriginalOkctlVersion,
		},
	}
}
