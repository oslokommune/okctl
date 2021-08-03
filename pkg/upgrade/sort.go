package upgrade

import sortPkg "sort"

func sort(upgradeBinaries []okctlUpgradeBinary) {
	sortPkg.SliceStable(upgradeBinaries, func(i, j int) bool {
		if upgradeBinaries[i].version.semver.LessThan(upgradeBinaries[j].version.semver) {
			return true
		}

		if upgradeBinaries[i].version.semver.GreaterThan(upgradeBinaries[j].version.semver) {
			return false
		}

		// semvers are equal, order on hotfix
		return upgradeBinaries[i].version.hotfix < upgradeBinaries[j].version.hotfix
	})
}
