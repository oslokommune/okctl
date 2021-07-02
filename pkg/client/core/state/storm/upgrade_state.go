package storm

import (
	"errors"
	"time"

	stormpkg "github.com/asdine/storm/v3"
	"github.com/oslokommune/okctl/pkg/breeze"
	"github.com/oslokommune/okctl/pkg/client"
)

type upgradeState struct {
	node breeze.Client
}

// Upgrade contains state about okctl upgrades
type Upgrade struct {
	Metadata `storm:"inline"`
	Name     string `storm:"unique"`

	ID      ID
	Version string
}

// NewUpgrade returns storm compatible state
func newUpgrade(u *client.Upgrade, meta Metadata) *Upgrade {
	return &Upgrade{
		Metadata: meta,
		Name:     "upgrade",
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
func (u *upgradeState) SaveUpgrade(upgrade *client.Upgrade) error {
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

// GetUpgrade returns the upgrade with the given version, or an error if it doesn't exist
func (u *upgradeState) GetUpgrade(version string) (*client.Upgrade, error) {
	r, err := u.getUpgrade(version)
	if err != nil && !errors.Is(err, stormpkg.ErrNotFound) {
		return nil, err
	}

	if errors.Is(err, stormpkg.ErrNotFound) {
		return nil, client.ErrUpgradeNotFound
	}

	return r.Convert(), nil
}

func (u *upgradeState) getUpgrade(version string) (*Upgrade, error) {
	upgrade := &Upgrade{}

	err := u.node.One("Version", version, upgrade)
	if err != nil {
		return nil, err
	}

	return upgrade, nil
}

// NewUpgradeState returns an initialised state client
func NewUpgradeState(node breeze.Client) client.UpgradeState {
	return &upgradeState{
		node: node,
	}
}
