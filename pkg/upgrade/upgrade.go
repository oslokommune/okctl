// Package upgrade knows how to upgrade okctl
package upgrade

import (
	"fmt"
	"io"
	"sort"

	semverPkg "github.com/Masterminds/semver"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/binaries/fetch"
	"github.com/oslokommune/okctl/pkg/binaries/run/okctlupgrade"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/config/state"
	"github.com/oslokommune/okctl/pkg/storage"
	"github.com/sirupsen/logrus"
)

// okctlUpgradeBinary contains metadata for an upgrade that can be run to upgrade okctl to some specific version.
// Note that an okctlUpgradeBinary represents multiple binaries, one for each combination of OS and architecture, see
// comment for field checksums.
type okctlUpgradeBinary struct {
	// fileExtension can be for instance "tar.gz"
	fileExtension string
	// version is the upgrade version, for instance "0.0.56" or "0.0.56_some_hotfix"
	version upgradeBinaryVersion
	// checksum is a list of checksums, one for every combination of host OS and architecture that exists for this
	// binary, for instance Linux-amd64
	checksums []state.Checksum
}

func (b okctlUpgradeBinary) String() string {
	return b.BinaryName()
}

// BinaryName returns a string with a version, for instance "okctl-upgrade_0.0.56"
func (b okctlUpgradeBinary) BinaryName() string {
	return fmt.Sprintf(okctlupgrade.BinaryNameFormat, b.RawVersion())
}

func (b okctlUpgradeBinary) RawVersion() string {
	return b.version.raw
}

func (b okctlUpgradeBinary) SemverVersion() *semverPkg.Version {
	return b.version.semver
}

func (b okctlUpgradeBinary) HotfixVersion() string {
	return b.version.hotfix
}

func newOkctlUpgradeBinary(version upgradeBinaryVersion, checksums []state.Checksum) okctlUpgradeBinary {
	return okctlUpgradeBinary{
		fileExtension: ".tar.gz",
		version:       version,
		checksums:     checksums,
	}
}

// Run upgrades okctl
func (u Upgrader) Run() error {
	// Fetch
	releases, err := u.githubService.ListReleases("oslokommune", "okctl-upgrade")
	if err != nil {
		return fmt.Errorf("listing github releases: %w", err)
	}

	upgradeBinaries, err := u.githubReleaseParser.toUpgradeBinaries(releases)
	if err != nil {
		return fmt.Errorf("parsing upgrade binaries: %w", err)
	}

	// Filter
	upgradeBinaries, err = u.filter.get(upgradeBinaries)
	if err != nil {
		return fmt.Errorf("filtering upgrade binaries: %w", err)
	}

	// Sort
	u.sort(upgradeBinaries)

	// Run
	err = u.runBinaries(upgradeBinaries)
	if err != nil {
		return fmt.Errorf("running upgrade binaries: %w", err)
	}

	// Update state
	err = u.filter.markAsRun(upgradeBinaries)
	if err != nil {
		return fmt.Errorf("marking upgrades as run: %w", err)
	}

	return nil
}

func (u Upgrader) sort(upgradeBinaries []okctlUpgradeBinary) {
	sort.SliceStable(upgradeBinaries, func(i, j int) bool {
		if upgradeBinaries[i].version.semver.LessThan(upgradeBinaries[j].version.semver) {
			return true
		}

		if upgradeBinaries[i].version.semver.GreaterThan(upgradeBinaries[j].version.semver) {
			return false
		}

		// semvers are equal, order on hotfix
		return upgradeBinaries[i].version.hotfix < upgradeBinaries[j].version.hotfix
	})
}

func (u Upgrader) runBinaries(upgradeBinaries []okctlUpgradeBinary) error {
	binaryProvider, err := u.createBinaryProvider(upgradeBinaries)
	if err != nil {
		return fmt.Errorf("creating binary provider: %w", err)
	}

	for _, binary := range upgradeBinaries {
		binaryRunner, err := binaryProvider.okctlUpgrade(binary.RawVersion())
		if err != nil {
			return fmt.Errorf("getting okctl upgrade binary: %w", err)
		}

		binaryRunner.Debug(u.debug)

		_, _ = fmt.Fprintf(u.out, "--- Running upgrade: %s ---\n", binary)

		_, err = binaryRunner.Run()
		if err != nil {
			return fmt.Errorf("running upgrade binary %s: %w", binary, err)
		}
	}

	return nil
}

func (u Upgrader) createBinaryProvider(upgradeBinaries []okctlUpgradeBinary) (upgradeBinaryProvider, error) {
	binaries := u.toStateBinaries(upgradeBinaries)

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

	if u.debug {
		_, _ = fmt.Fprintf(u.out, "Found %d upgrade(s)\n", len(upgradeBinaries))
	}

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
	ChecksumDownloader   ChecksumDownloader
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
		filter:              newFilter(opts.State, opts.ClusterID, opts.OkctlVersion, opts.OriginalOkctlVersion),
	}
}
