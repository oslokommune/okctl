package upgrade

import (
	"fmt"

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
	// version is the upgrade version, for instance "0.0.56" or "0.0.56_some_hotfix"
	version upgradeBinaryVersion
	// checksum is a list of checksums, one for every combination of host OS and architecture that exists for this
	// binary, for instance Linux-amd64
	checksums []state.Checksum
}

func (b okctlUpgradeBinary) String() string {
	return b.BinaryName()
}

// BinaryName returns a string with a version, for instance "okctl-upgrade_0.0.56"
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

func newOkctlUpgradeBinary(version upgradeBinaryVersion, checksums []state.Checksum) okctlUpgradeBinary {
	return okctlUpgradeBinary{
		fileExtension: ".tar.gz",
		version:       version,
		checksums:     checksums,
	}
}
