package upgrade

import (
	"fmt"
	sortPkg "sort"

	semverPkg "github.com/Masterminds/semver"
	"github.com/oslokommune/okctl/pkg/binaries/run/okctlupgrade"
	"github.com/oslokommune/okctl/pkg/config/state"
)

// okctlUpgradeBinary contains metadata for an upgrade that can be run to upgrade okctl to some specific version.
// Note that an okctlUpgradeBinary represents multiple binaries, one for each combination of OS and architecture, see
// comment for field checksums.
type okctlUpgradeBinary struct {
	// fileExtension can be for instance "tar.gz"
	fileExtension string
	// version is the upgrade version
	version upgradeBinaryVersion
	// checksum is a list of checksums, one for every combination of host OS and architecture that exists for this
	// binary, for instance Linux-amd64
	checksums []state.Checksum
	// git tag associated with the binary
	gitTag string
}

func (b okctlUpgradeBinary) String() string {
	return b.BinaryName()
}

// BinaryName returns a string with a version, for instance "okctl-upgrade_0.0.56.some-component"
func (b okctlUpgradeBinary) BinaryName() string {
	return fmt.Sprintf(okctlupgrade.BinaryNameFormat, b.RawVersion())
}

func (b okctlUpgradeBinary) RawVersion() string {
	return b.version.raw
}

func (b okctlUpgradeBinary) SemverVersion() *semverPkg.Version {
	return b.version.semver
}

func (b okctlUpgradeBinary) HotfixVersion() string {
	return b.version.hotfix
}

func (b okctlUpgradeBinary) GitTag() string {
	return b.gitTag
}

func sort(upgradeBinaries []okctlUpgradeBinary) {
	sortPkg.SliceStable(upgradeBinaries, func(i, j int) bool {
		if upgradeBinaries[i].SemverVersion().LessThan(upgradeBinaries[j].SemverVersion()) {
			return true
		}

		if upgradeBinaries[i].SemverVersion().GreaterThan(upgradeBinaries[j].SemverVersion()) {
			return false
		}

		// semvers are equal, order on hotfix
		return upgradeBinaries[i].HotfixVersion() < upgradeBinaries[j].HotfixVersion()
	})
}

func newOkctlUpgradeBinary(version upgradeBinaryVersion, checksums []state.Checksum, gitTag string) okctlUpgradeBinary {
	return okctlUpgradeBinary{
		fileExtension: ".tar.gz",
		version:       version,
		checksums:     checksums,
		gitTag:        gitTag,
	}
}
