package upgrade

import (
	"errors"
	"fmt"
	"io"

	"github.com/Masterminds/semver"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
)

func (f filter) get(binaries []okctlUpgradeBinary) ([]okctlUpgradeBinary, error) {
	var err error

	binaries, err = f.removeAlreadyExecuted(binaries)
	if err != nil {
		return nil, fmt.Errorf("removing already executed upgrade binaries: %w", err)
	}

	printIfDebug(f.debug, f.out, "%d remaining upgrades after removing already executed:", binaries)

	binaries, err = f.removeTooNew(binaries)
	if err != nil {
		return nil, fmt.Errorf("removing too new upgrade binaries: %w", err)
	}

	printIfDebug(f.debug, f.out, "%d remaining upgrades after removing too new upgrades:", binaries)

	binaries, err = f.removeTooOld(binaries)
	if err != nil {
		return nil, fmt.Errorf("removing too old upgrade binaries: %w", err)
	}

	printIfDebug(f.debug, f.out, "%d remaining upgrades after removing too old upgrades:", binaries)

	return binaries, nil
}

func (f filter) removeAlreadyExecuted(binaries []okctlUpgradeBinary) ([]okctlUpgradeBinary, error) {
	var notExecuted []okctlUpgradeBinary

	for _, binary := range binaries {
		hasBinaryExecuted, err := f.hasBinaryRun(binary)
		if err != nil {
			return nil, fmt.Errorf("checking if binary '%s' has run: %w", binary, err)
		}

		if !hasBinaryExecuted {
			notExecuted = append(notExecuted, binary)
		}
	}

	return notExecuted, nil
}

func (f filter) hasBinaryRun(binary okctlUpgradeBinary) (bool, error) {
	_, err := f.state.GetUpgrade(binary.RawVersion())
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

	okctlSemver, err := semver.NewVersion(f.okctlVersion)
	if err != nil {
		return nil, fmt.Errorf("could not create semver from okctl version '%s': %w", f.okctlVersion, err)
	}

	for _, binary := range binaries {
		if binary.SemverVersion().Equal(okctlSemver) || binary.SemverVersion().LessThan(okctlSemver) {
			versionIsEqualOrLess = append(versionIsEqualOrLess, binary)
		}
	}

	return versionIsEqualOrLess, nil
}

// removeTooOld removes binaries with versions below original okctl version.
//
// If we make a cluster with okctl 0.0.62, we should never run upgrade binaries that are meant
// to upgrade older versions of the cluster up to 0.0.62. This is because there is no need to run those
// upgrades, as a fresh okctl cluster already is up-to-date. So if we have an upgrade-binary-0.0.60, which
// means "upgrade cluster and attached resources to support okctl version 0.0.60", we don't need to run this
// upgrade binary.
func (f filter) removeTooOld(binaries []okctlUpgradeBinary) ([]okctlUpgradeBinary, error) {
	var versionIsNewThanOriginalOkctlVersion []okctlUpgradeBinary

	originalOkctlSemver, err := semver.NewVersion(f.originalOkctlVersion)
	if err != nil {
		return nil, fmt.Errorf("could not create semver from original okctl version '%s': %w", f.originalOkctlVersion, err)
	}

	for _, binary := range binaries {
		if binary.SemverVersion().GreaterThan(originalOkctlSemver) {
			versionIsNewThanOriginalOkctlVersion = append(versionIsNewThanOriginalOkctlVersion, binary)
		}
	}

	return versionIsNewThanOriginalOkctlVersion, nil
}

func (f filter) markAsRun(binary okctlUpgradeBinary) error {
	u := &client.Upgrade{
		ID:      f.clusterID,
		Version: binary.RawVersion(),
	}

	err := f.state.SaveUpgrade(u)
	if err != nil {
		return fmt.Errorf("saving upgrade %s: %w", u.Version, err)
	}

	return nil
}

type filter struct {
	debug                bool
	out                  io.Writer
	state                client.UpgradeState
	clusterID            api.ID
	okctlVersion         string
	originalOkctlVersion string
}
