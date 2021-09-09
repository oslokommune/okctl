// Package clusterversion manages the cluster version
package clusterversioner

import (
	"errors"
	"fmt"
	"io"

	"github.com/Masterminds/semver"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
)

// ValidateBinaryVsClusterVersion returns an error if binary version is less than cluster version.
func (v Versioner) ValidateBinaryVsClusterVersion(binaryVersionString string) error {
	binaryVersion, err := semver.NewVersion(binaryVersionString)
	if err != nil {
		return fmt.Errorf("parsing binary version to semver from '%s': %w", binaryVersionString, err)
	}

	clusterVersionInfo, err := v.upgradeState.GetClusterVersion()
	if errors.Is(err, client.ErrClusterVersionNotFound) {
		// This means we haven't stored the cluster version yet. In this case we don't return an error, as we don't
		// expect it to be stored yet.
		return nil
	}

	if err != nil {
		return fmt.Errorf(": %w", err)
	}

	clusterVersion, err := semver.NewVersion(clusterVersionInfo.Value)
	if err != nil {
		return fmt.Errorf("parsing cluster verion to semver from '%s': %w", binaryVersionString, err)
	}

	return v.validateBinaryVsClusterVersion(binaryVersion, clusterVersion)
}

func (v Versioner) validateBinaryVsClusterVersion(binaryVersion *semver.Version, clusterVersion *semver.Version) error {
	if binaryVersion.LessThan(clusterVersion) {
		return fmt.Errorf("okctl binary version %s cannot be less than cluster version %s."+
			" Get okctl version %s or later and try again",
			binaryVersion.String(), clusterVersion.String(), clusterVersion.String())
	}

	return nil
}

// SaveClusterVersion saves the provided version
func (v Versioner) SaveClusterVersion(version string) error {
	err := v.upgradeState.SaveClusterVersion(&client.ClusterVersion{
		ID:    v.clusterID,
		Value: version,
	})
	if err != nil {
		return fmt.Errorf("saving cluster version: %w", err)
	}

	return nil
}

// SaveClusterVersionIfNotExists saves the current cluster version
func (v Versioner) SaveClusterVersionIfNotExists(version string) error {
	_, err := v.upgradeState.GetClusterVersion()
	if err != nil && !errors.Is(err, client.ErrClusterVersionNotFound) {
		return fmt.Errorf("getting cluster version: %w", err)
	}

	if errors.Is(err, client.ErrClusterVersionNotFound) {
		var versionToSave *semver.Version

		err = v.upgradeState.SaveClusterVersion(&client.ClusterVersion{
			ID:    v.clusterID,
			Value: version,
		})
		if err != nil {
			return fmt.Errorf("saving version '%s': %w", versionToSave.String(), err)
		}
	}

	return nil
}

// SaveClusterVersionFromOriginalClusterVersionIfNotExists gets the cluster's original version, and saves it as the
// current cluster version.
//
// When we're sure all users have set cluster version in state, the caller of this function can replace the call
// with a call to SaveClusterVersionIfNotExists(version.GetVersionInfo().Version).
//
// Ideally, we would like to store just version.GetVersionInfo().Version, if it hasn't been stored before.
// Because what we want to achieve is to store current version of okctl into the current version of the
// cluster when a new cluster is made.
//
// However, we're in a transition phase, where users' clusters don't have stored cluster version yet. If we
// simply store version.GetVersionInfo().Version, say 0.0.60, and the original cluster value is 0.0.50,
// this means the upgrade logic believes the cluster version is 0.0.60, when in fact it is 0.0.50. The
// upgrade logic will therefore not run upgrades for 0.0.50 to 0.0.60, when in fact they should be run.
//
// The solution is to store the original cluster version instead. The upgrade logic will bump the cluster
// version to the current version when run.
//
// To check if users has run this code, just check all users' cluster's state, and see if
// upgrade/ClusterVersion has been set or not. If it's set, it means this code has been run.
func (v Versioner) SaveClusterVersionFromOriginalClusterVersionIfNotExists() error {
	_, err := v.upgradeState.GetClusterVersion()
	if err != nil && !errors.Is(err, client.ErrClusterVersionNotFound) {
		return fmt.Errorf("getting cluster version: %w", err)
	}

	if errors.Is(err, client.ErrClusterVersionNotFound) {
		var versionToSave *semver.Version

		originalClusterVersion, err := v.upgradeState.GetOriginalClusterVersion()
		if err != nil {
			return fmt.Errorf("getting original cluster version: %w", err)
		}

		err = v.upgradeState.SaveClusterVersion(&client.ClusterVersion{
			ID:    v.clusterID,
			Value: originalClusterVersion.Value,
		})
		if err != nil {
			return fmt.Errorf("saving version '%s': %w", versionToSave.String(), err)
		}
	}

	return nil
}

// GetClusterVersion returns the current cluster version
func (v Versioner) GetClusterVersion() (string, error) {
	clusterVersion, err := v.upgradeState.GetClusterVersion()
	if err != nil {
		return "", fmt.Errorf("getting cluster version: %w", err)
	}

	return clusterVersion.Value, nil
}

// Versioner knows how to enforce correct version of the okctl binary versus the cluster version.
// The intention is that we want to enforce that no users of a cluster are trying to run 'upgrade' or 'apply cluster'
// with an outdated version of the okctl binary, that is, a version that is older than the cluster version.
// The cluster version should be set to the current version whenever we run 'upgrade' or 'apply cluster'.
type Versioner struct {
	out          io.Writer
	clusterID    api.ID
	upgradeState client.UpgradeState
}

// New returns a new Versioner
func New(out io.Writer, clusterID api.ID, upgradeState client.UpgradeState) Versioner {
	return Versioner{
		out:          out,
		clusterID:    clusterID,
		upgradeState: upgradeState,
	}
}
