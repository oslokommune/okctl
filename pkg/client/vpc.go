package client

import (
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client/store"
)

// We are shadowing the api interfaces for now, but
// this is probably not sustainable.

// VPCService orchestrates the creation of a vpc
type VPCService interface {
	api.VpcService
}

// VPCAPI invokes the API calls for creating a vpc
type VPCAPI interface {
	api.VpcCloudProvider
}

// VPCStore stores the data
type VPCStore interface {
	SaveVpc(vpc *api.Vpc) (*store.Report, error)
	DeleteVpc(id api.ID) (*store.Report, error)
	GetVpc(id api.ID) (*api.Vpc, error)
}
