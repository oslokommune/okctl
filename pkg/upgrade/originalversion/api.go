// Package originalversion knows how to save the original version of a cluster
package originalversion

import (
	"errors"
	"fmt"

	"github.com/Masterminds/semver"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/client"
)

// Saver stores the version of okctl first used to apply a cluster
type Saver struct {
	clusterID    api.ID
	upgradeState client.UpgradeState
	clusterState client.ClusterState
}

// SaveErrorMessage contains the error messsage to return if saving original version fails
const SaveErrorMessage = "saving original version. Your cluster will not work as expected. Please" +
	" retry and see if this error goes away. If it does, everything is okay. If not, contact the" +
	" developers to address this issue. Error details: %w"

// SaveOriginalClusterVersionIfNotExists stores the version of okctl first used to apply a cluster
// Ideally, we want to store just version.GetVersionInfo().Version if it's missing. However, this can cause
// problems. Let's say we release the upgrade functionality in version 0.0.60. A user might have a cluster made
// with version 0.0.50. They then wait until 0.0.70 to download a new version of okctl, and then runs apply
// cluster or upgrade.
// If we simply store version.GetVersionInfo().Version, i.e. 0.0.70, okctl upgrade won't download upgrades for
// 0.0.65 - 0.0.70, as these will be considered as too old, when in fact they should run.
// To fix this, we use the information from state, i.e. a tag applied to the cluster cloudformation stack, which
// has been present since 0.0.40 (see: https://github.com/oslokommune/okctl/pull/299).
// When we're sure all users has run this code, i.e. stored original okctl version in the correct place,
// we can simplify this logic to just use version.GetVersionInfo().Version and ignore the cluster tag version.
// To check if users has run this code, just check all users' cluster's state, and see if
// upgrade/OriginalClusterVersion has been set or not. If it's set, it means this code has been run.
func (o Saver) SaveOriginalClusterVersionIfNotExists() error {
	_, err := o.upgradeState.GetOriginalClusterVersion()
	if err != nil && !errors.Is(err, client.ErrOriginalClusterVersionNotFound) {
		return fmt.Errorf("getting original okctl version: %w", err)
	}

	if errors.Is(err, client.ErrOriginalClusterVersionNotFound) {
		var versionToSave *semver.Version

		clusterStateVersion, err := o.getClusterStateVersion()
		if err != nil {
			return fmt.Errorf("getting cluster state version: %w", err)
		}

		err = o.upgradeState.SaveOriginalClusterVersionIfNotExists(&client.OriginalClusterVersion{
			ID:    o.clusterID,
			Value: clusterStateVersion.String(),
		})
		if err != nil {
			return fmt.Errorf("saving version '%s': %w", versionToSave.String(), err)
		}
	}

	return nil
}

func (o Saver) getClusterStateVersion() (*semver.Version, error) {
	cluster, err := o.clusterState.GetCluster(o.clusterID.ClusterName)
	if err != nil {
		return nil, fmt.Errorf("getting cluster '%s': %w", o.clusterID.ClusterName, err)
	}

	clusterStateVersionString, ok := cluster.Config.Metadata.Tags[v1alpha1.OkctlVersionTag]
	if !ok {
		return nil, fmt.Errorf("could not find cluster state tag '%s'", v1alpha1.OkctlVersionTag)
	}

	clusterStateVersion, err := semver.NewVersion(clusterStateVersionString)
	if err != nil {
		return nil, fmt.Errorf("parsing version '%s': %w", clusterStateVersion, err)
	}

	return clusterStateVersion, nil
}

// New returns a Saver
func New(
	clusterID api.ID,
	upgradeState client.UpgradeState,
	clusterState client.ClusterState,
) (Saver, error) {
	return Saver{
		clusterID:    clusterID,
		upgradeState: upgradeState,
		clusterState: clusterState,
	}, nil
}
