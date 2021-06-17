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
	progress io.Writer
	fetcher  fetch.Provider
	logger   *logrus.Logger

	okctlUpgrade map[string]*okctlupgrade.OkctlUpgrade
}

func newUpgradeBinaryProvider(
	logger *logrus.Logger,
	progress io.Writer,
	fetcher fetch.Provider,
) upgradeBinaryProvider {
	return upgradeBinaryProvider{
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

		p.okctlUpgrade[version] = okctlupgrade.New(p.logger, p.progress, binaryPath, run.Cmd())
	}

	return p.okctlUpgrade[version], nil
}
