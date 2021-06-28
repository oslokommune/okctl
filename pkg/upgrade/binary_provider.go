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

	okctlUpgrades map[string]*okctlupgrade.OkctlUpgrade
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
		okctlUpgrades: map[string]*okctlupgrade.OkctlUpgrade{},
	}
}

func (p *upgradeBinaryProvider) okctlUpgrade(version string) (*okctlupgrade.OkctlUpgrade, error) {
	_, ok := p.okctlUpgrades[version]

	if !ok {
		binaryName := fmt.Sprintf(okctlupgrade.BinaryNameFormat, version)

		binaryPath, err := p.fetcher.Fetch(binaryName, version)
		if err != nil {
			return nil, err
		}

		p.okctlUpgrades[version] = okctlupgrade.New(p.repoDir, p.progress, p.logger, binaryPath, run.Cmd())
	}

	return p.okctlUpgrades[version], nil
}
