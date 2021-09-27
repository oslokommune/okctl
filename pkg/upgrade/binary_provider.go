package upgrade

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/oslokommune/okctl/pkg/binaries/fetch"
	"github.com/oslokommune/okctl/pkg/binaries/run"
	"github.com/oslokommune/okctl/pkg/binaries/run/okctlupgrade"
	"github.com/sirupsen/logrus"
)

type upgradeBinaryProvider struct {
	repoDir  string
	progress io.Writer
	fetcher  fetch.Provider
	logger   *logrus.Logger

	binaryRunners map[string]*okctlupgrade.BinaryRunner
}

func (p *upgradeBinaryProvider) okctlUpgradeRunner(version string) (*okctlupgrade.BinaryRunner, error) {
	_, ok := p.binaryRunners[version]

	if !ok {
		binaryName := fmt.Sprintf(okctlupgrade.BinaryNameFormat, version)

		binaryPath, err := p.fetcher.Fetch(binaryName, version)
		if err != nil {
			return nil, err
		}

		p.binaryRunners[version] = okctlupgrade.New(p.repoDir, p.progress, p.logger, binaryPath, cmd())
	}

	return p.binaryRunners[version], nil
}

// Using a custom Cmd instead of the default one, as we have to set stdin to os.Stdin, in order to for our upgrades
// to be able to get input from the user
func cmd() run.CmdFn {
	return func(workingDir, path string, env, args []string) *exec.Cmd {
		cmd := run.Cmd()(workingDir, path, env, args)
		cmd.Stdin = os.Stdin

		return cmd
	}
}

func newUpgradeBinaryProvider(
	repoDir string,
	logger *logrus.Logger,
	progress io.Writer,
	fetcher fetch.Provider,
) upgradeBinaryProvider {
	return upgradeBinaryProvider{
		repoDir:       repoDir,
		progress:      progress,
		fetcher:       fetcher,
		logger:        logger,
		binaryRunners: map[string]*okctlupgrade.BinaryRunner{},
	}
}
