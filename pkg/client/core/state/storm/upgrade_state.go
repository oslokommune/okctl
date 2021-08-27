package storm

import (
	"errors"
	"time"

	"github.com/oslokommune/okctl/pkg/version"

	stormpkg "github.com/asdine/storm/v3"
	"github.com/oslokommune/okctl/pkg/breeze"
	"github.com/oslokommune/okctl/pkg/client"
)

//
// Upgrades
//

type upgradesState struct {
	node breeze.Client
}

// Upgrades contains state about upgrades
type Upgrades struct {
	Metadata `storm:"inline"`

	ID ID
}

// Upgrade contains state about an upgrade
type Upgrade struct {
	Metadata `storm:"inline"`

	ID      ID
	Version string
}

// NewUpgrade returns storm compatible state
func newUpgrade(u *client.Upgrade, meta Metadata) *Upgrade {
	return &Upgrade{
		Metadata: meta,
		ID:       NewID(u.ID),
		Version:  u.Version,
	}
}

// Convert to client.Upgrade
func (u *Upgrade) Convert() *client.Upgrade {
	return &client.Upgrade{
		ID:      u.ID.Convert(),
		Version: u.Version,
	}
}

// SaveUpgrade saves the upgrade
func (u *upgradesState) SaveUpgrade(upgrade *client.Upgrade) error {
	existing, err := u.getUpgrade(upgrade.Version)
	if err != nil && !errors.Is(err, stormpkg.ErrNotFound) {
		return err
	}

	if errors.Is(err, stormpkg.ErrNotFound) {
		return u.node.Save(newUpgrade(upgrade, NewMetadata()))
	}

	existing.Metadata.UpdatedAt = time.Now()

	return u.node.Save(newUpgrade(upgrade, existing.Metadata))
}

func (u *upgradesState) getUpgrade(version string) (*Upgrade, error) {
	upgrade := &Upgrade{}

	err := u.node.One("Version", version, upgrade)
	if err != nil {
		return nil, err
	}

	return upgrade, nil
}

// GetUpgrades returns all upgrades, or an error
func (u *upgradesState) GetUpgrades() ([]*client.Upgrade, error) {
	var upgrades []*Upgrade

	err := u.node.All(&upgrades)
	if err != nil {
		return nil, err
	}

	clientUpgrades := make([]*client.Upgrade, len(upgrades))
	for i, upgrade := range upgrades {
		clientUpgrades[i] = upgrade.Convert()
	}

	return clientUpgrades, nil
}

//
// Original version
//

const originalVersionValue = "OriginalOkctlVersion"

// OriginalOkctlVersion contains state about the original installed version of okctl
type OriginalOkctlVersion struct {
	Metadata `storm:"inline"`

	ID    ID
	Value string
	Key   string
}

func newOriginalOkctlVersion(o *client.OriginalOkctlVersion, meta Metadata) *OriginalOkctlVersion {
	return &OriginalOkctlVersion{
		Metadata: meta,
		ID:       NewID(o.ID),
		Value:    o.Value,
		Key:      originalVersionValue,
	}
}

// Convert to client.Upgrade
func (o *OriginalOkctlVersion) Convert() *client.OriginalOkctlVersion {
	return &client.OriginalOkctlVersion{
		ID:    o.ID.Convert(),
		Value: o.Value,
	}
}

// SaveOriginalOkctlVersionIfNotExists saves the original version if it hasn't been saved before
func (u *upgradesState) SaveOriginalOkctlVersionIfNotExists(originalOkctlVersion *client.OriginalOkctlVersion) error {
	_, err := u.getOriginalOkctlVersion()
	if err != nil && !errors.Is(err, stormpkg.ErrNotFound) {
		return err
	}

	if errors.Is(err, stormpkg.ErrNotFound) {
		return u.node.Save(newOriginalOkctlVersion(originalOkctlVersion, NewMetadata()))
	}

	return nil
}

// GetOriginalOkctlVersion returns the original version, or an error if it doesn't exist
func (u *upgradesState) GetOriginalOkctlVersion() (*client.OriginalOkctlVersion, error) {
	originalVersion, err := u.getOriginalOkctlVersion()
	if err != nil && !errors.Is(err, stormpkg.ErrNotFound) {
		return nil, err
	}

	if errors.Is(err, stormpkg.ErrNotFound) {
		return nil, client.ErrOriginalOkctlVersionNotFound
	}

	return originalVersion.Convert(), nil
}

func (u *upgradesState) getOriginalOkctlVersion() (*OriginalOkctlVersion, error) {
	originalOkctlVersion := &OriginalOkctlVersion{}

	err := u.node.One("Key", originalVersionValue, originalOkctlVersion)
	if err != nil {
		return nil, err
	}

	return originalOkctlVersion, nil
}

//
// Cluster version
//

const clusterVersionValue = "ClusterVersion"

// ClusterVersion contains state about the cluster version
type ClusterVersion struct {
	Metadata `storm:"inline"`

	ID    ID
	Value version.Info
	Key   string
}

// Convert to client.Upgrade
func (o *ClusterVersion) Convert() *client.ClusterVersion {
	return &client.ClusterVersion{
		ID:    o.ID.Convert(),
		Value: o.Value,
	}
}

func (u *upgradesState) GetClusterVersionInfo() (*client.ClusterVersion, error) {
	clusterVersion, err := u.getClusterVersion()
	if err != nil && !errors.Is(err, stormpkg.ErrNotFound) {
		return nil, err
	}

	if errors.Is(err, stormpkg.ErrNotFound) {
		return nil, client.ErrClusterVersionNotFound
	}

	return clusterVersion.Convert(), nil
}

func (u *upgradesState) getClusterVersion() (*ClusterVersion, error) {
	clusterVersion := &ClusterVersion{}

	err := u.node.One("Key", clusterVersionValue, clusterVersion)
	if err != nil {
		return nil, err
	}

	return clusterVersion, nil
}

func newClusterVersion(o *client.ClusterVersion, meta Metadata) *ClusterVersion {
	return &ClusterVersion{
		Metadata: meta,
		ID:       NewID(o.ID),
		Value:    o.Value,
		Key:      clusterVersionValue,
	}
}

func (u *upgradesState) SaveClusterVersionInfo(version *client.ClusterVersion) error {
	return u.node.Save(newClusterVersion(version, NewMetadata()))
}

// NewUpgradesState returns an initialised state client
func NewUpgradesState(node breeze.Client) client.UpgradeState {
	return &upgradesState{
		node: node,
	}
}
