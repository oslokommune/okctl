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

	binaryRunners        map[string]*okctlupgrade.BinaryRunner
	environmentVariables map[string]string
}

func (p *upgradeBinaryProvider) okctlUpgradeRunner(version string) (*okctlupgrade.BinaryRunner, error) {
	_, ok := p.binaryRunners[version]

	if !ok {
		binaryName := fmt.Sprintf(okctlupgrade.BinaryNameFormat, version)

		binaryPath, err := p.fetcher.Fetch(binaryName, version)
		if err != nil {
			return nil, err
		}

		envs := toSlice(p.environmentVariables)

		p.binaryRunners[version] = okctlupgrade.New(p.repoDir, p.progress, p.logger, envs, binaryPath, cmd())
	}

	return p.binaryRunners[version], nil
}

func toSlice(m map[string]string) []string {
	result := make([]string, len(m))
	index := 0

	for key, val := range m {
		result[index] = fmt.Sprintf("%s=%s", key, val)

		index++
	}

	return result
}

func cmd() run.CmdFn {
	return func(workingDir, path string, env, args []string) *exec.Cmd {
		cmd := run.Cmd()(workingDir, path, env, args)
		cmd.Stdin = os.Stdin // Enables user input when running a binary

		return cmd
	}
}

func newUpgradeBinaryProvider(
	repoDir string,
	logger *logrus.Logger,
	progress io.Writer,
	fetcher fetch.Provider,
	binaryEnvironmentVariables map[string]string,
) upgradeBinaryProvider {
	return upgradeBinaryProvider{
		repoDir:              repoDir,
		progress:             progress,
		fetcher:              fetcher,
		logger:               logger,
		binaryRunners:        map[string]*okctlupgrade.BinaryRunner{},
		environmentVariables: binaryEnvironmentVariables,
	}
}
