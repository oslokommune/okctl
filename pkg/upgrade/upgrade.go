// Package upgrade knows how to upgrade okctl
package upgrade

import (
	"fmt"
	"io"
	"strings"

	"github.com/oslokommune/okctl/pkg/upgrade/clusterversioner"
	"github.com/oslokommune/okctl/pkg/upgrade/originalclusterversioner"

	"github.com/AlecAivazis/survey/v2"
	"github.com/oslokommune/okctl/pkg/commands"
	"github.com/oslokommune/okctl/pkg/github"

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
	releases, err := u.githubService.ListReleases(github.DefaultOrg, "okctl-upgrade")
	if err != nil {
		return fmt.Errorf("listing github releases: %w", err)
	}

	upgradeBinaries, err := u.githubReleaseParser.ToUpgradeBinaries(releases)
	if err != nil {
		return fmt.Errorf("parsing upgrade binaries: %w", err)
	}

	printUpgradesIfDebug(u.debug, u.out, "Found %d upgrade(s):", upgradeBinaries)

	// Filter
	alreadyExecuted, err := u.getAlreadyExecutedBinaries()
	if err != nil {
		return fmt.Errorf("getting already executed binaries: %w", err)
	}

	upgradeBinaries, err = u.filter.get(upgradeBinaries, alreadyExecuted)
	if err != nil {
		return fmt.Errorf("filtering upgrade binaries: %w", err)
	}

	// Sort, i.e. determine execution order
	sort(upgradeBinaries)

	// Run
	if len(upgradeBinaries) > 0 {
		printUpgrades(u.out, "Found %d applicable upgrade(s):", upgradeBinaries)
	} else {
		_, _ = fmt.Fprintln(u.out, "Did not find any applicable upgrades.")
		return nil
	}

	err = u.runBinaries(upgradeBinaries)
	if err != nil {
		return fmt.Errorf("running upgrade binaries: %w", err)
	}

	// Update cluster version
	err = u.clusterVersioner.SaveClusterVersion(u.okctlVersion)
	if err != nil {
		return fmt.Errorf(commands.SaveClusterVersionError, err)
	}

	_, _ = fmt.Fprintf(u.out, "\nUpgrade complete! Cluster version is now %s."+
		" Remember to commit and push changes with git.\n", u.okctlVersion)

	return nil
}

func (u Upgrader) getAlreadyExecutedBinaries() (map[string]bool, error) {
	alreadyExecutedSlice, err := u.state.GetUpgrades()
	if err != nil {
		return nil, fmt.Errorf("getting upgrades: %w", err)
	}

	alreadyExecuted := make(map[string]bool)

	for _, upgrade := range alreadyExecutedSlice {
		alreadyExecuted[upgrade.Version] = true
	}

	return alreadyExecuted, nil
}

func (u Upgrader) runBinaries(upgradeBinaries []okctlUpgradeBinary) error {
	binaryProvider, err := u.createBinaryProvider(upgradeBinaries)
	if err != nil {
		return fmt.Errorf("creating binary provider: %w", err)
	}

	err = u.dryRunBinaries(upgradeBinaries, binaryProvider)
	if err != nil {
		return fmt.Errorf("simulating upgrades: %w", err)
	}

	doContinue, err := u.askUserIfReady()
	if err != nil {
		return fmt.Errorf("asking user for input: %w", err)
	}

	if !doContinue {
		_, _ = fmt.Fprintln(u.out, "User aborted.")
		return nil
	}

	err = u.doRunBinaries(upgradeBinaries, binaryProvider)
	if err != nil {
		return fmt.Errorf("running upgrades: %w", err)
	}

	return nil
}

func (u Upgrader) doRunBinaries(upgradeBinaries []okctlUpgradeBinary, binaryProvider upgradeBinaryProvider) error {
	for _, binary := range upgradeBinaries {
		// Get
		binaryRunner, err := binaryProvider.okctlUpgradeRunner(binary.RawVersion())
		if err != nil {
			return fmt.Errorf("getting okctl upgrade binary: %w", err)
		}

		binaryRunner.SetDebug(u.debug)

		_, _ = fmt.Fprintf(u.out, "--- Running upgrade: %s ---\n", binary)

		// Run
		_, err = binaryRunner.Run(true)
		if err != nil {
			_, _ = fmt.Fprintf(u.out, "--- Upgrade failed: %s ---\n", binary)
			return fmt.Errorf("running upgrade binary %s: %w", binary, err)
		}

		// Mark as run
		err = u.markAsRun(binary)
		if err != nil {
			return fmt.Errorf("marking upgrades as run: %w", err)
		}

		// Update cluster version
		// Note that we don't save the hotfix version, only the semver version. This is because we use the cluster
		// versioner later to validate that the current okctl version against the cluster version. The okctl version
		// will always be a semver version without any hotfix in it, so it wouldn't be possible to store a hotfix
		// version here.
		//
		// Also, a hotfix version is supposed to just fix any errors being made in an upgrade binary, so
		// the effect of an upgrade binary version 0.0.10 plus a hotfix binary 0.0.10.a, should be as if running one
		// working upgrade binary with version 0.0.10.
		err = u.clusterVersioner.SaveClusterVersion(binary.SemverVersion().String())
		if err != nil {
			return fmt.Errorf(commands.SaveClusterVersionError, err)
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

func (u Upgrader) dryRunBinaries(upgradeBinaries []okctlUpgradeBinary, binaryProvider upgradeBinaryProvider) error {
	_, _ = fmt.Fprint(u.out, "Simulating upgrades...\n\n")

	for _, binary := range upgradeBinaries {
		// Get
		binaryRunner, err := binaryProvider.okctlUpgradeRunner(binary.RawVersion())
		if err != nil {
			return fmt.Errorf("getting okctl upgrade binary: %w", err)
		}

		binaryRunner.SetDebug(u.debug)

		_, _ = fmt.Fprintf(u.out, "--- Simulating upgrade: %s ---\n", binary)

		// Run
		_, err = binaryRunner.Run(false)
		if err != nil {
			_, _ = fmt.Fprintf(u.out, "--- Upgrade failed: %s ---\n", binary)
			return fmt.Errorf("running upgrade binary %s: %w", binary, err)
		}
	}

	_, _ = fmt.Fprintf(u.out, "\nSimulating upgrades complete.\n\n")

	return nil
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

func (u Upgrader) markAsRun(binary okctlUpgradeBinary) error {
	clientUpgrade := &client.Upgrade{
		ID:      u.clusterID,
		Version: binary.RawVersion(),
	}

	err := u.state.SaveUpgrade(clientUpgrade)
	if err != nil {
		return fmt.Errorf("saving upgrade %s: %w", clientUpgrade.Version, err)
	}

	return nil
}

func (u Upgrader) askUserIfReady() (bool, error) {
	if u.autoConfirmPrompt {
		return true, nil
	}

	doContinue := false
	prompt := &survey.Confirm{
		Message: "This will upgrade your okctl cluster, are you sure you want to continue?",
	}

	err := survey.AskOne(prompt, &doContinue)
	if err != nil {
		return false, err
	}

	_, _ = fmt.Fprintln(u.out, "")

	return doContinue, nil
}

func printUpgradesIfDebug(debug bool, out io.Writer, text string, upgradeBinaries []okctlUpgradeBinary) {
	if debug {
		printUpgrades(out, text, upgradeBinaries)
	}
}

func printUpgrades(out io.Writer, text string, upgradeBinaries []okctlUpgradeBinary) {
	binaries := make([]string, 0)
	for _, binary := range upgradeBinaries {
		binaries = append(binaries, binary.RawVersion())
	}

	joinedBinariesTxt := strings.Join(binaries, ", ")

	_, _ = fmt.Fprintf(out, text+"\n", len(upgradeBinaries))
	_, _ = fmt.Fprintln(out, joinedBinariesTxt)
	_, _ = fmt.Fprintln(out, "")
}

// FetcherOpts contains data needed to initialize a fetch.Provider
type FetcherOpts struct {
	Host  state.Host
	Store storage.Storer
}

// Opts contains all data needed to create an Upgrader
type Opts struct {
	Debug                    bool
	Logger                   *logrus.Logger
	Out                      io.Writer
	AutoConfirmPrompt        bool
	RepositoryDirectory      string
	GithubService            client.GithubService
	ChecksumDownloader       ChecksumHTTPDownloader
	ClusterVersioner         clusterversioner.Versioner
	OriginalClusterVersioner originalclusterversioner.Versioner
	FetcherOpts              FetcherOpts
	OkctlVersion             string
	State                    client.UpgradeState
	ClusterID                api.ID
}

// Upgrader knows how to upgrade okctl
type Upgrader struct {
	debug               bool
	logger              *logrus.Logger
	out                 io.Writer
	autoConfirmPrompt   bool
	clusterID           api.ID
	state               client.UpgradeState
	repositoryDirectory string
	githubService       client.GithubService
	githubReleaseParser GithubReleaseParser
	clusterVersioner    clusterversioner.Versioner
	okctlVersion        string
	fetcherOpts         FetcherOpts
	filter              filter
}

// New returns a new Upgrader, or an error if initialization fails
func New(opts Opts) (Upgrader, error) {
	err := opts.ClusterVersioner.ValidateBinaryVsClusterVersion(opts.OkctlVersion)
	if err != nil {
		return Upgrader{}, fmt.Errorf(commands.ValidateBinaryVsClusterVersionError, err)
	}

	// We can remove this call to SaveOriginalClusterVersionIfNotExists when we're sure all users have stored original
	// cluster version into their state. It should be set by apply cluster, not upgrade. We need to have it here
	// in case people run upgrade before apply cluster.
	err = opts.OriginalClusterVersioner.SaveOriginalClusterVersionFromClusterTagIfNotExists()
	if err != nil {
		return Upgrader{}, fmt.Errorf(originalclusterversioner.SaveErrorMessage, err)
	}

	originalClusterVersion, err := opts.State.GetOriginalClusterVersion()
	if err != nil {
		return Upgrader{}, fmt.Errorf("getting original okctl version: %w", err)
	}

	return Upgrader{
		debug:               opts.Debug,
		logger:              opts.Logger,
		out:                 opts.Out,
		autoConfirmPrompt:   opts.AutoConfirmPrompt,
		clusterID:           opts.ClusterID,
		state:               opts.State,
		repositoryDirectory: opts.RepositoryDirectory,
		githubService:       opts.GithubService,
		githubReleaseParser: NewGithubReleaseParser(opts.ChecksumDownloader),
		clusterVersioner:    opts.ClusterVersioner,
		okctlVersion:        opts.OkctlVersion,
		fetcherOpts:         opts.FetcherOpts,
		filter: filter{
			debug:                  opts.Debug,
			out:                    opts.Out,
			okctlVersion:           opts.OkctlVersion,
			originalClusterVersion: originalClusterVersion.Value,
		},
	}, nil
}
