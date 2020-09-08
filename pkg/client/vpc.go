package client

import "github.com/oslokommune/okctl/pkg/api"

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
	api.VpcStore
}
