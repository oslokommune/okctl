// Package upgrade knows how to upgrade okctl
package upgrade

import (
	"fmt"
	"io"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/oslokommune/okctl/pkg/binaries/run/okctlupgrade"

	"github.com/oslokommune/okctl/pkg/upgrade/clusterversion"
	"github.com/oslokommune/okctl/pkg/upgrade/originalclusterversion"

	"github.com/oslokommune/okctl/pkg/upgrade/survey"

	"github.com/oslokommune/okctl/pkg/commands"
	"github.com/oslokommune/okctl/pkg/github"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/binaries/fetch"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/config/state"
	"github.com/oslokommune/okctl/pkg/storage"
	"github.com/sirupsen/logrus"
)

// OkctlUpgradeRepo contains the repository for upgrades
const OkctlUpgradeRepo = "okctl-upgrade"

// OkctlUpgradeRepoURL contains URL for the upgrade repository
var OkctlUpgradeRepoURL = fmt.Sprintf("https://github.com/oslokommune/%s", OkctlUpgradeRepo) //nolint:gochecknoglobals

// Run upgrades okctl
func (u Upgrader) Run() error {
	// Fetch
	releases, err := u.githubService.ListReleases(github.DefaultOrg, OkctlUpgradeRepo)
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

		userConfirmedContinue, err := u.runBinaries(upgradeBinaries)
		if err != nil {
			return fmt.Errorf("running upgrade binaries: %w", err)
		}

		if !userConfirmedContinue {
			return nil
		}
	} else {
		_, _ = fmt.Fprintln(u.out, "Did not find any applicable upgrades. Cluster version will be updated regardless.")
	}

	// Update cluster version
	err = u.clusterVersioner.SaveClusterVersion(u.okctlVersion)
	if err != nil {
		return fmt.Errorf(commands.SaveClusterVersionErr, err)
	}

	_, _ = fmt.Fprintf(u.out, "\nUpgrade complete, cluster version is now %s.\n", u.okctlVersion)

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

func (u Upgrader) runBinaries(upgradeBinaries []okctlUpgradeBinary) (bool, error) {
	binaryProvider, err := u.createBinaryProvider(upgradeBinaries)
	if err != nil {
		return false, fmt.Errorf("creating binary provider: %w", err)
	}

	err = u.dryRunBinaries(upgradeBinaries, binaryProvider)
	if err != nil {
		return false, fmt.Errorf("simulating upgrades: %w", err)
	}

	userConfirmedContinue, err := u.surveyor.PromptUser(
		"This will upgrade your okctl cluster, are you sure you want to continue?")
	if err != nil {
		return false, fmt.Errorf("asking user for input: %w", err)
	}

	if !userConfirmedContinue {
		_, _ = fmt.Fprintln(u.out, "Upgrade aborted by user.")
		return false, nil
	}

	err = u.doRunBinaries(upgradeBinaries, binaryProvider)
	if err != nil {
		return false, fmt.Errorf("running upgrades: %w", err)
	}

	return true, nil
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

	binaryProvider := newUpgradeBinaryProvider(u.repositoryDirectory, u.logger, u.out, fetcher, u.binaryEnvironmentVariables)

	return binaryProvider, nil
}

func (u Upgrader) dryRunBinaries(upgradeBinaries []okctlUpgradeBinary, binaryProvider upgradeBinaryProvider) error {
	_, _ = fmt.Fprint(u.out, "Simulating upgrades (we're not doing any actual changes yet, "+
		"just printing what's going to happen)... \n\n")

	for _, binary := range upgradeBinaries {
		// Get
		binaryRunner, err := binaryProvider.okctlUpgradeRunner(binary.RawVersion())
		if err != nil {
			return fmt.Errorf("getting okctl upgrade binary: %w", err)
		}

		_, _ = fmt.Fprintf(u.out, "--- Simulating upgrade: %s ---\n", binary)

		// Run
		_, err = binaryRunner.DryRun(okctlupgrade.Flags{
			Debug: u.debug,
		})
		if err != nil {
			_, _ = fmt.Fprintf(u.out, "--- Upgrade failed: %s ---\n", binary)
			return fmt.Errorf("running upgrade binary %s: %w", binary, err)
		}
	}

	_, _ = fmt.Fprintf(u.out, "\nSimulating upgrades complete.\n\n")

	return nil
}

func (u Upgrader) doRunBinaries(upgradeBinaries []okctlUpgradeBinary, binaryProvider upgradeBinaryProvider) error {
	for _, binary := range upgradeBinaries {
		// Get
		binaryRunner, err := binaryProvider.okctlUpgradeRunner(binary.RawVersion())
		if err != nil {
			return fmt.Errorf("getting okctl upgrade binary: %w", err)
		}

		_, _ = fmt.Fprintf(u.out, "--- Running upgrade: %s ---\n", binary)

		// Run
		_, err = binaryRunner.Run(okctlupgrade.Flags{
			Debug:   u.debug,
			Confirm: u.autoConfirm,
		})

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
			return fmt.Errorf(commands.SaveClusterVersionErr, err)
		}
	}

	return nil
}

func (u Upgrader) toStateBinaries(upgradeBinaries []okctlUpgradeBinary) []state.Binary {
	binaries := make([]state.Binary, 0, len(upgradeBinaries))

	for _, upgradeBinary := range upgradeBinaries {
		URLPattern := fmt.Sprintf(
			"%s/releases/download/%s/okctl-upgrade_%s_#{os}_#{arch}.tar.gz",
			OkctlUpgradeRepoURL,
			upgradeBinary.GitTag(),
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

// capitalizeFirst converts for instance "liNUX" to "Linux". We use this because we expect GitHub release assets for
// upgrades to be named this way.
func capitalizeFirst(os string) string {
	return strings.ToUpper(os[0:1]) + strings.ToLower(os[1:])
}

// DocumentationURL is the URL to the upgrade documentation
const DocumentationURL = "https://okctl.io/getting-started/upgrading"

func getOriginalClusterVersion(opts Opts) (*client.OriginalClusterVersion, error) {
	hasOriginalClusterVersion, err := opts.OriginalClusterVersioner.OriginalClusterVersionExists()
	if err != nil {
		return nil, fmt.Errorf("checking if original cluster version exists: %w", err)
	}

	// (tag UPGR01) We can remove replace this if-block and above call to OriginalClusterVersionExists with a call to
	// originalClusterVersioner.SaveOriginalClusterVersionIfNotExists
	// when we're sure all users have stored original cluster version into their state. It should be set by
	// apply cluster, not upgrade. We need to have it here in case people run upgrade before apply cluster.
	if !hasOriginalClusterVersion {
		_, _ = fmt.Fprintf(opts.Out, "Okctl needs to initialize parts of the cluster state to support upgrades."+
			" Afterwards you should commit and push changes to git.\nIf you want more details, see %s\n\n",
			DocumentationURL)

		answer, err := opts.Surveyor.PromptUser("Do you want to proceed?")
		if err != nil {
			return nil, fmt.Errorf("prompting user: %w", err)
		}

		if !answer {
			return nil, fmt.Errorf("upgrade aborted by user")
		}

		err = opts.OriginalClusterVersioner.SaveOriginalClusterVersionFromClusterTagIfNotExists()
		if err != nil {
			return nil, fmt.Errorf(originalclusterversion.SaveErrorMessage, err)
		}
	}

	originalClusterVersion, err := opts.State.GetOriginalClusterVersion()
	if err != nil {
		return nil, fmt.Errorf("getting original okctl version: %w", err)
	}

	return originalClusterVersion, nil
}

// FetcherOpts contains data needed to initialize a fetch.Provider
type FetcherOpts struct {
	Host  state.Host
	Store storage.Storer
}

// Opts contains all data needed to create an Upgrader
type Opts struct {
	Debug                      bool
	AutoConfirm                bool
	Logger                     *logrus.Logger
	Out                        io.Writer
	RepositoryDirectory        string
	GithubService              client.GithubService
	ChecksumDownloader         ChecksumHTTPDownloader
	ClusterVersioner           clusterversion.Versioner
	OriginalClusterVersioner   originalclusterversion.Versioner
	Surveyor                   survey.Surveyor
	FetcherOpts                FetcherOpts
	OkctlVersion               string
	State                      client.UpgradeState
	ClusterID                  api.ID
	BinaryEnvironmentVariables map[string]string
}

// Validate validates the given parameters
func (o Opts) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.Logger, validation.Required),
		validation.Field(&o.Out, validation.Required),
		validation.Field(&o.RepositoryDirectory, validation.Required),
		validation.Field(&o.GithubService, validation.Required),
		validation.Field(&o.ChecksumDownloader, validation.Required),
		validation.Field(&o.ClusterVersioner, validation.Required),
		validation.Field(&o.OriginalClusterVersioner, validation.Required),
		validation.Field(&o.Surveyor, validation.Required),
		validation.Field(&o.FetcherOpts, validation.Required),
		validation.Field(&o.OkctlVersion, validation.Required),
		validation.Field(&o.State, validation.Required),
		validation.Field(&o.BinaryEnvironmentVariables, validation.NotNil),
	)
}

// Upgrader knows how to upgrade okctl
type Upgrader struct {
	debug                      bool
	autoConfirm                bool
	logger                     *logrus.Logger
	out                        io.Writer
	clusterID                  api.ID
	state                      client.UpgradeState
	repositoryDirectory        string
	githubService              client.GithubService
	githubReleaseParser        GithubReleaseParser
	clusterVersioner           clusterversion.Versioner
	surveyor                   survey.Surveyor
	okctlVersion               string
	fetcherOpts                FetcherOpts
	filter                     filter
	binaryEnvironmentVariables map[string]string
}

// New returns a new Upgrader, or an error if initialization fails
func New(opts Opts) (Upgrader, error) {
	err := opts.Validate()
	if err != nil {
		return Upgrader{}, err
	}

	err = opts.ClusterVersioner.ValidateBinaryVersionNotLessThanClusterVersion(opts.OkctlVersion)
	if err != nil {
		return Upgrader{}, fmt.Errorf(commands.ValidateBinaryVsClusterVersionErr, err)
	}

	originalClusterVersion, err := getOriginalClusterVersion(opts)
	if err != nil {
		return Upgrader{}, err
	}

	fetcherOpts := FetcherOpts{
		Host: state.Host{
			// Github release URL expects OS to be Linux, not linux
			Os:   capitalizeFirst(opts.FetcherOpts.Host.Os),
			Arch: opts.FetcherOpts.Host.Arch,
		},
		Store: opts.FetcherOpts.Store,
	}

	return Upgrader{
		debug:               opts.Debug,
		autoConfirm:         opts.AutoConfirm,
		logger:              opts.Logger,
		out:                 opts.Out,
		clusterID:           opts.ClusterID,
		state:               opts.State,
		repositoryDirectory: opts.RepositoryDirectory,
		githubService:       opts.GithubService,
		githubReleaseParser: NewGithubReleaseParser(opts.ChecksumDownloader),
		clusterVersioner:    opts.ClusterVersioner,
		surveyor:            opts.Surveyor,
		okctlVersion:        opts.OkctlVersion,
		fetcherOpts:         fetcherOpts,
		filter: filter{
			debug:                  opts.Debug,
			out:                    opts.Out,
			okctlVersion:           opts.OkctlVersion,
			originalClusterVersion: originalClusterVersion.Value,
		},
		binaryEnvironmentVariables: opts.BinaryEnvironmentVariables,
	}, nil
}
