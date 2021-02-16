package client

import (
	"context"

	"github.com/oslokommune/okctl/pkg/client/store"

	"github.com/oslokommune/okctl/pkg/api"
)

// Blockstorage is the content of a blockstorage deployment
type Blockstorage struct {
	Policy         *api.ManagedPolicy
	ServiceAccount *api.ServiceAccount
	Chart          *api.Helm
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

// BlockstorageAPI invokes REST API endpoints
type BlockstorageAPI interface {
	CreateBlockstoragePolicy(opts api.CreateBlockstoragePolicy) (*api.ManagedPolicy, error)
	DeleteBlockstoragePolicy(id api.ID) error
	CreateBlockstorageServiceAccount(opts api.CreateBlockstorageServiceAccountOpts) (*api.ServiceAccount, error)
	DeleteBlockstorageServiceAccount(id api.ID) error
	CreateBlockstorageHelmChart(opts api.CreateBlockstorageHelmChartOpts) (*api.Helm, error)
}

// BlockstorageStore is a storage layer implementation
type BlockstorageStore interface {
	SaveBlockstorage(scaler *Blockstorage) (*store.Report, error)
	RemoveBlockstorage(id api.ID) (*store.Report, error)
}

// BlockstorageReport is a report layer
type BlockstorageReport interface {
	ReportCreateBlockstorage(scaler *Blockstorage, report *store.Report) error
	ReportDeleteBlockstorage(report *store.Report) error
}
