package client

import (
	"errors"

	"github.com/oslokommune/okctl/pkg/api"
)

// ErrUpgradeNotFound is returned when the specified upgrade is not found in state
var ErrUpgradeNotFound = errors.New("not found")

// Upgrade contains state about an okctl upgrade
type Upgrade struct {
	ID      api.ID
	Version string
}

// UpgradeState updates the state
type UpgradeState interface {
	SaveUpgrade(upgrade *Upgrade) error
	GetUpgrade(version string) (*Upgrade, error)
}
