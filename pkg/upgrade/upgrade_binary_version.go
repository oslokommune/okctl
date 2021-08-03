package upgrade

import (
	"fmt"
	"strings"

	semverPkg "github.com/Masterminds/semver"
)

const (
	dotCountForRegularSemver    = 3
	dotCountForSemverWithHotfix = 4
)

type upgradeBinaryVersion struct {
	// raw is semver + hotfix, for instance "0.0.56_a"
	raw string
	// semver is the semantic version, for instance "0.0.56"
	semver *semverPkg.Version
	// hotfix is the hotfix version, for instance "a"
	hotfix string
}

// newUpgradeBinaryVersion parses a given version and returns an instance of upgradeBinaryVersion or
// an error if unable to parse the version.
//
// Valid input examples:
//
// 0.0.56
//
// 0.0.56.some-hotfix
func newUpgradeBinaryVersion(text string) (upgradeBinaryVersion, error) {
	var semver *semverPkg.Version

	var err error

	hotfix := ""
	parts := strings.Split(text, ".")

	switch {
	case len(parts) == dotCountForRegularSemver:
		semver, err = semverPkg.NewVersion(text)
		if err != nil {
			return upgradeBinaryVersion{}, fmt.Errorf("parsing semantic version: %w", err)
		}
	case len(parts) == dotCountForSemverWithHotfix:
		semverString := strings.Join(parts[0:3], ".")

		semver, err = semverPkg.NewVersion(semverString)
		if err != nil {
			return upgradeBinaryVersion{}, fmt.Errorf("parsing semantic hotfix version: %w", err)
		}

		hotfix = parts[3]
	default:
		return upgradeBinaryVersion{}, fmt.Errorf("not a valid version: %s", text)
	}

	return upgradeBinaryVersion{
		raw:    text,
		semver: semver,
		hotfix: hotfix,
	}, nil
}
