package upgrade

import (
	"errors"
	"fmt"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
)

type filter struct {
	state     client.UpgradeState
	clusterID api.ID
}

func (f filter) get(binaries []okctlUpgradeBinary) ([]okctlUpgradeBinary, error) {
	return f.getNotAlreadyExecuted(binaries)
}

func (f filter) getNotAlreadyExecuted(binaries []okctlUpgradeBinary) ([]okctlUpgradeBinary, error) {
	var filtered []okctlUpgradeBinary

	for _, binary := range binaries {
		hasBinaryRun, err := f.hasBinaryRun(binary)
		if err != nil {
			return nil, fmt.Errorf("checking if binary '%s' has run: %w", binary.version, err)
		}

		if !hasBinaryRun {
			filtered = append(filtered, binary)
		}
	}

	return filtered, nil
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

func newFilter(state client.UpgradeState, clusterID api.ID) filter {
	return filter{
		state:     state,
		clusterID: clusterID,
	}
}
