package upgrade

import (
	"fmt"
	"io"

	"github.com/Masterminds/semver"
)

func (f filter) get(binaries []okctlUpgradeBinary, alreadyExecuted map[string]bool) ([]okctlUpgradeBinary, error) {
	var err error

	binaries = f.removeAlreadyExecuted(binaries, alreadyExecuted)

	printUpgradesIfDebug(f.debug, f.out, "%d remaining upgrades after removing already executed:", binaries)

	binaries, err = f.removeTooNew(binaries)
	if err != nil {
		return nil, fmt.Errorf("removing too new upgrade binaries: %w", err)
	}

	printUpgradesIfDebug(f.debug, f.out, "%d remaining upgrades after removing too new upgrades:", binaries)

	binaries, err = f.removeTooOld(binaries)
	if err != nil {
		return nil, fmt.Errorf("removing too old upgrade binaries: %w", err)
	}

	printUpgradesIfDebug(f.debug, f.out, "%d remaining upgrades after removing too old upgrades:", binaries)

	return binaries, nil
}

func (f filter) removeAlreadyExecuted(binaries []okctlUpgradeBinary, alreadyExecuted map[string]bool) []okctlUpgradeBinary {
	var notExecuted []okctlUpgradeBinary

	for _, binary := range binaries {
		_, hasBinaryExecuted := alreadyExecuted[binary.RawVersion()]

		if !hasBinaryExecuted {
			notExecuted = append(notExecuted, binary)
		}
	}

	return notExecuted
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

// removeTooOld removes binaries with versions below original clster version.
//
// If we make a cluster with okctl 0.0.62, we should never run upgrade binaries that are meant
// to upgrade older versions of the cluster up to 0.0.62. This is because there is no need to run those
// upgrades, as a fresh okctl cluster already is up-to-date. So if we have an upgrade-binary-0.0.60, which
// means "upgrade cluster and attached resources to support okctl version 0.0.60", we don't need to run this
// upgrade binary.
func (f filter) removeTooOld(binaries []okctlUpgradeBinary) ([]okctlUpgradeBinary, error) {
	var versionIsNewThanOriginalClusterVersion []okctlUpgradeBinary

	originalClusterSemver, err := semver.NewVersion(f.originalClusterVersion)
	if err != nil {
		return nil, fmt.Errorf("could not create semver from original okctl version '%s': %w", f.originalClusterVersion, err)
	}

	for _, binary := range binaries {
		if binary.SemverVersion().GreaterThan(originalClusterSemver) {
			versionIsNewThanOriginalClusterVersion = append(versionIsNewThanOriginalClusterVersion, binary)
		}
	}

	return versionIsNewThanOriginalClusterVersion, nil
}

type filter struct {
	debug                  bool
	out                    io.Writer
	okctlVersion           string
	originalClusterVersion string
}
