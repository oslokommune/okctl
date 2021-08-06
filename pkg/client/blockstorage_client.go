package client

import (
	"context"

	"github.com/oslokommune/okctl/pkg/api"
)

// Blockstorage is the content of a blockstorage deployment
type Blockstorage struct {
	Policy         *ManagedPolicy
	ServiceAccount *ServiceAccount
	Chart          *Helm
}

// CreateBlockstorageOpts contains the required inputs
type CreateBlockstorageOpts struct {
	ID api.ID
}

// BlockstorageService is an implementation of the business logic
type BlockstorageService interface {
	CreateBlockstorage(ctx context.Context, opts CreateBlockstorageOpts) (*Blockstorage, error)
	DeleteBlockstorage(ctx context.Context, id api.ID) error
}

// BlockstorageState defines the state layer
type BlockstorageState interface {
	HasBlockstorage() (bool, error)
}
