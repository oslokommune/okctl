// Package originalclusterversioner knows how to save the original version of a cluster
package originalclusterversioner

import (
	"errors"
	"fmt"

	"github.com/Masterminds/semver"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/client"
)

// Versioner stores the version of okctl first used to apply a cluster
type Versioner struct {
	clusterID    api.ID
	upgradeState client.UpgradeState
	clusterState client.ClusterState
}

// SaveErrorMessage contains the error messsage to return if saving original cluster version fails
const SaveErrorMessage = "saving original cluster version. Your cluster will not work as expected. Please" +
	" retry and see if this error goes away. If it does, everything is okay. If not, contact the" +
	" developers to address this issue. Error details: %w"

// SaveOriginalClusterVersionIfNotExists saves the original cluster version
func (v Versioner) SaveOriginalClusterVersionIfNotExists(version string) error {
	_, err := v.upgradeState.GetOriginalClusterVersion()
	if err != nil && !errors.Is(err, client.ErrOriginalClusterVersionNotFound) {
		return fmt.Errorf("getting original cluster version: %w", err)
	}

	if errors.Is(err, client.ErrOriginalClusterVersionNotFound) {
		err = v.upgradeState.SaveOriginalClusterVersionIfNotExists(&client.OriginalClusterVersion{
			ID:    v.clusterID,
			Value: version,
		})
		if err != nil {
			return fmt.Errorf("saving version '%s': %w", version, err)
		}
	}

	return nil
}

// SaveOriginalClusterVersionFromClusterTagIfNotExists stores the original cluster version if it doesn't exist, with
// the value from the cluster tag.
//
// When we're sure all users has run this code, i.e. stored original okctl version in state, the caller of this function
// can replace the call with a call to SaveOriginalClusterVersionIfNotExists(version.GetVersionInfo().Version).
//
// Ideally, we want to store just version.GetVersionInfo().Version if it's missing. However, this can cause
// problems. Let's say we release the upgrade functionality in version 0.0.60. A user might have a cluster made
// with version 0.0.50. They then wait until 0.0.70 to download a new version of okctl, and then runs apply
// cluster or upgrade.
//
// If we simply store version.GetVersionInfo().Version, i.e. 0.0.70, okctl upgrade won't download upgrades for
// 0.0.65 - 0.0.70, as these will be considered as too old, when in fact they should run.
//
// To fix this, we use the information from state, i.e. a tag applied to the cluster cloudformation stack, which
// has been present since 0.0.40 (see: https://github.com/oslokommune/okctl/pull/299).
//
// To check if users has run this code, just check all users' cluster's state, and see if
// upgrade/OriginalClusterVersion has been set or not. If it's set, it means this code has been run.
func (v Versioner) SaveOriginalClusterVersionFromClusterTagIfNotExists() error {
	_, err := v.upgradeState.GetOriginalClusterVersion()
	if err != nil && !errors.Is(err, client.ErrOriginalClusterVersionNotFound) {
		return fmt.Errorf("getting original cluster version: %w", err)
	}

	if errors.Is(err, client.ErrOriginalClusterVersionNotFound) {
		clusterStateVersion, err := v.getClusterStateVersion()
		if err != nil {
			return fmt.Errorf("getting cluster state version: %w", err)
		}

		err = v.upgradeState.SaveOriginalClusterVersionIfNotExists(&client.OriginalClusterVersion{
			ID:    v.clusterID,
			Value: clusterStateVersion.String(),
		})
		if err != nil {
			return fmt.Errorf("saving version '%s': %w", clusterStateVersion.String(), err)
		}
	}

	return nil
}

func (v Versioner) getClusterStateVersion() (*semver.Version, error) {
	cluster, err := v.clusterState.GetCluster(v.clusterID.ClusterName)
	if err != nil {
		return nil, fmt.Errorf("getting cluster version tag: %w", err)
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

func (v Versioner) GetOriginalClusterVersion() (string, error) {
	version, err := v.upgradeState.GetOriginalClusterVersion()
	if err != nil && !errors.Is(err, client.ErrOriginalClusterVersionNotFound) {
		return "", fmt.Errorf("getting original cluster version: %w", err)
	}

	if err != nil {
		return "", err
	}

	return version.Value, nil
}

// New returns a Versioner
func New(
	clusterID api.ID,
	upgradeState client.UpgradeState,
	clusterState client.ClusterState,
) Versioner {
	return Versioner{
		clusterID:    clusterID,
		upgradeState: upgradeState,
		clusterState: clusterState,
	}
}
