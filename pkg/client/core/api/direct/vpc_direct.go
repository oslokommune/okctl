package direct

import (
	"context"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
)

type vpcAPIDirectClient struct {
	service api.VpcService
}

func (v *vpcAPIDirectClient) CreateVpc(opts api.CreateVpcOpts) (*api.Vpc, error) {
	return v.service.CreateVpc(context.Background(), opts)
}

func (v *vpcAPIDirectClient) DeleteVpc(opts api.DeleteVpcOpts) error {
	return v.service.DeleteVpc(context.Background(), opts)
}

// NewVPCAPI returns an initialised API client with core service
func NewVPCAPI(service api.VpcService) client.VPCAPI {
	return &vpcAPIDirectClient{
		service: service,
	}
}
