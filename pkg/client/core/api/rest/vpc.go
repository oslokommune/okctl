package rest

import (
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
)

// TargetVPC is the API route for HTTP requests
const TargetVPC = "vpcs/"

type vpcAPI struct {
	client *client.HTTPClient
}

func (a *vpcAPI) CreateVpc(opts api.CreateVpcOpts) (*api.Vpc, error) {
	into := &api.Vpc{}
	return into, a.client.DoPost(TargetVPC, opts, into)
}

func (a *vpcAPI) DeleteVpc(opts api.DeleteVpcOpts) error {
	return a.client.DoDelete(TargetVPC, opts)
}

// NewVPCAPI returns an initialised API REST client
func NewVPCAPI(client *client.HTTPClient) client.VPCAPI {
	return &vpcAPI{
		client: client,
	}
}
