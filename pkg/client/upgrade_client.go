package client

import (
	"errors"

	"github.com/oslokommune/okctl/pkg/version"

	"github.com/oslokommune/okctl/pkg/api"
)

// ErrUpgradeNotFound is returned when an upgrade is not found in state
var ErrUpgradeNotFound = errors.New("not found")

// ErrOriginalOkctlVersionNotFound is returned when the original okctl version is not found in state
var ErrOriginalOkctlVersionNotFound = errors.New("not found")

// ErrClusterVersionNotFound is returned when the cluster version is not found in state
var ErrClusterVersionNotFound = errors.New("not found")

// Upgrade contains state about an okctl upgrade
type Upgrade struct {
	ID      api.ID
	Version string
}

// OriginalOkctlVersion contains state about the original okctl version installed
type OriginalOkctlVersion struct {
	ID    api.ID
	Value string
}

// ClusterVersion contains state about the cluster version
type ClusterVersion struct {
	ID    api.ID
	Value version.Info
}

// UpgradeState updates the state
type UpgradeState interface {
	SaveUpgrade(upgrade *Upgrade) error
	GetUpgrades() ([]*Upgrade, error)
	SaveOriginalOkctlVersionIfNotExists(originalOkctlVersion *OriginalOkctlVersion) error
	GetOriginalOkctlVersion() (*OriginalOkctlVersion, error)
	SaveClusterVersionInfo(version *ClusterVersion) error
	GetClusterVersionInfo() (*ClusterVersion, error)
}
