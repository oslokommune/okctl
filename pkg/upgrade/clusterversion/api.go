// Package clusterversion manages the cluster version
package clusterversion

import (
	"errors"
	"fmt"
	"io"

	"github.com/Masterminds/semver"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/version"
)

// ValidateBinaryVsClusterVersion returns an error if binary version is less than cluster version
func (c ClusterVersioner) ValidateBinaryVsClusterVersion(binaryVersionString string) error {
	binaryVersion, err := semver.NewVersion(binaryVersionString)
	if err != nil {
		return fmt.Errorf("parsing binary version to semver from '%s': %w", binaryVersionString, err)
	}

	clusterVersionInfo, err := c.upgradeState.GetClusterVersionInfo()
	if errors.Is(err, client.ErrClusterVersionNotFound) {
		// This means we haven't stored the cluster version yet. In this case we don't return an error, as we don't
		// expect it to be stored yet.
		return nil
	}

	if err != nil {
		return fmt.Errorf(": %w", err)
	}

	clusterVersion, err := semver.NewVersion(clusterVersionInfo.Value.Version)
	if err != nil {
		return fmt.Errorf("parsing cluster verion to semver from '%s': %w", binaryVersionString, err)
	}

	return c.validateBinaryVsClusterVersion(binaryVersion, clusterVersion)
}

func (c ClusterVersioner) validateBinaryVsClusterVersion(binaryVersion *semver.Version, clusterVersion *semver.Version) error {
	if binaryVersion.LessThan(clusterVersion) {
		return fmt.Errorf("okctl binary version '%s' cannot be less than cluster version '%s'. This can cause"+
			" incompatibility errors. Get a newer version of okctl and try again",
			binaryVersion.String(), clusterVersion.String())
	}

	return nil
}

// SaveClusterVersion saves the provided version
func (c ClusterVersioner) SaveClusterVersion(version version.Info) error {
	err := c.upgradeState.SaveClusterVersionInfo(&client.ClusterVersion{
		ID:    c.clusterID,
		Value: version,
	})
	if err != nil {
		return fmt.Errorf("saving cluster version: %w", err)
	}

	_, _ = fmt.Fprintf(c.out, "Cluster version is now: %s\n", version.Version)

	return nil
}

// ClusterVersioner knows how to save and get cluster version
type ClusterVersioner struct {
	out          io.Writer
	clusterID    api.ID
	upgradeState client.UpgradeState
}

// New returns a new ClusterVersioner
func New(out io.Writer, clusterID api.ID, upgradeState client.UpgradeState) ClusterVersioner {
	return ClusterVersioner{
		out:          out,
		clusterID:    clusterID,
		upgradeState: upgradeState,
	}
}
