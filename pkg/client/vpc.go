package client

import (
	"context"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client/store"
)

// We are shadowing the api interfaces for now, but
// this is probably not sustainable.

// VPCService orchestrates the creation of a vpc
type VPCService interface {
	CreateVpc(ctx context.Context, opts api.CreateVpcOpts) (*api.Vpc, error)
	DeleteVpc(ctx context.Context, opts api.DeleteVpcOpts) error
}

// VPCAPI invokes the API calls for creating a vpc
type VPCAPI interface {
	CreateVpc(opts api.CreateVpcOpts) (*api.Vpc, error)
	DeleteVpc(opts api.DeleteVpcOpts) error
}

// VPCStore stores the data
type VPCStore interface {
	SaveVpc(vpc *api.Vpc) (*store.Report, error)
	DeleteVpc(id api.ID) (*store.Report, error)
	GetVpc(id api.ID) (*api.Vpc, error)
}

// VPCReport summaries the creation of a VPC
type VPCReport interface {
	ReportCreateVPC(vpc *api.Vpc, reports []*store.Report) error
}

// VPCState implement the state layer
type VPCState interface {
	SaveVpc(vpc *api.Vpc) (*store.Report, error)
	DeleteVpc(id api.ID) (*store.Report, error)
}
