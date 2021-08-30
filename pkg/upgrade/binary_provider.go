package upgrade

import (
	"fmt"
	"io"

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

		p.binaryRunners[version] = okctlupgrade.New(p.repoDir, p.progress, p.logger, binaryPath, run.Cmd())
	}

	return p.binaryRunners[version], nil
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
