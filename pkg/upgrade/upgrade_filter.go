package upgrade

import (
	"errors"
	"fmt"

	"github.com/Masterminds/semver"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
)

type filter struct {
	state        client.UpgradeState
	clusterID    api.ID
	okctlVersion string
}

func (f filter) get(binaries []okctlUpgradeBinary) ([]okctlUpgradeBinary, error) {
	var err error

	binaries, err = f.removeAlreadyExecuted(binaries)
	if err != nil {
		return nil, fmt.Errorf("removing already executed upgrade binaries: %w", err)
	}

	binaries, err = f.removeTooNew(binaries)
	if err != nil {
		return nil, fmt.Errorf("removing too new upgrade binaries: %w", err)
	}

	return binaries, nil
}

func (f filter) removeAlreadyExecuted(binaries []okctlUpgradeBinary) ([]okctlUpgradeBinary, error) {
	var notExecuted []okctlUpgradeBinary

	for _, binary := range binaries {
		hasBinaryExecuted, err := f.hasBinaryRun(binary)
		if err != nil {
			return nil, fmt.Errorf("checking if binary '%s' has run: %w", binary.version, err)
		}

		if !hasBinaryExecuted {
			notExecuted = append(notExecuted, binary)
		}
	}

	return notExecuted, nil
}

func (f filter) hasBinaryRun(binary okctlUpgradeBinary) (bool, error) {
	_, err := f.state.GetUpgrade(binary.version)
	if err != nil {
		if !errors.Is(err, client.ErrUpgradeNotFound) {
			return false, err
		}

		if errors.Is(err, client.ErrUpgradeNotFound) {
			return false, nil
		}
	}

	return true, nil
}

// removeTooNew removes binaries that are too new for the current okctl version. For instance, if okctl is on version
// 0.0.63, and there exist upgrade binaries for 0.0.62, 0.0.63 and 0.0.64, we only want to run binaries 0.0.62 and
// 0.0.63.
func (f filter) removeTooNew(binaries []okctlUpgradeBinary) ([]okctlUpgradeBinary, error) {
	var versionIsEqualOrLess []okctlUpgradeBinary

	for _, binary := range binaries {
		binarySemver, err := semver.NewVersion(binary.version)
		if err != nil {
			return nil, fmt.Errorf("could not create semver from upgrade binary version '%s': %w", binary.version, err)
		}

		okctlSemver, err := semver.NewVersion(f.okctlVersion)
		if err != nil {
			return nil, fmt.Errorf("could not create semver from okctl version '%s': %w", f.okctlVersion, err)
		}

		if binarySemver.Equal(okctlSemver) || binarySemver.LessThan(okctlSemver) {
			versionIsEqualOrLess = append(versionIsEqualOrLess, binary)
		}
	}

	return versionIsEqualOrLess, nil
}

func (f filter) markAsRun(binaries []okctlUpgradeBinary) error {
	for _, binary := range binaries {
		u := &client.Upgrade{
			ID:      f.clusterID,
			Version: binary.version,
		}

		err := f.state.SaveUpgrade(u)
		if err != nil {
			return fmt.Errorf("saving upgrade %s: %w", u.Version, err)
		}
	}

	return nil
}

func newFilter(state client.UpgradeState, clusterID api.ID, okctlVersion string) filter {
	return filter{
		state:        state,
		clusterID:    clusterID,
		okctlVersion: okctlVersion,
	}
}
