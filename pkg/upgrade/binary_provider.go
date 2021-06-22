package upgrade

import (
	"fmt"
	"github.com/oslokommune/okctl/pkg/binaries/fetch"
	"github.com/oslokommune/okctl/pkg/binaries/run"
	"github.com/oslokommune/okctl/pkg/binaries/run/okctlupgrade"
	"github.com/sirupsen/logrus"
	"io"
)

type upgradeBinaryProvider struct {
	repoDir  string
	progress io.Writer
	fetcher  fetch.Provider
	logger   *logrus.Logger

	okctlUpgrade map[string]*okctlupgrade.OkctlUpgrade
}

func newUpgradeBinaryProvider(
	repoDir string,
	logger *logrus.Logger,
	progress io.Writer,
	fetcher fetch.Provider,
) upgradeBinaryProvider {
	return upgradeBinaryProvider{
		repoDir:      repoDir,
		progress:     progress,
		fetcher:      fetcher,
		logger:       logger,
		okctlUpgrade: map[string]*okctlupgrade.OkctlUpgrade{},
	}
}

// OkctlUpgrade returns an okctl upgrade wrapper. The given version is the version of the okctl upgrade to run.
func (p *upgradeBinaryProvider) OkctlUpgrade(version string) (*okctlupgrade.OkctlUpgrade, error) {
	_, ok := p.okctlUpgrade[version]

	if !ok {
		binaryName := fmt.Sprintf(okctlupgrade.BinaryNameFormat, version)

		binaryPath, err := p.fetcher.Fetch(binaryName, version)
		if err != nil {
			return nil, err
		}

		p.okctlUpgrade[version] = okctlupgrade.New(p.repoDir, p.progress, p.logger, binaryPath, run.Cmd())
	}

	return p.okctlUpgrade[version], nil
}
