package dryrun

import (
	"context"
	"fmt"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
	"io"
)

type vpcService struct {
	out io.Writer
}

func (v vpcService) CreateVpc(_ context.Context, opts client.CreateVpcOpts) (*client.Vpc, error) {
	fmt.Fprintf(v.out, formatCreate(fmt.Sprintf("VPC with CIDR %s", opts.Cidr)))

	return &client.Vpc{
		ID:                     opts.ID,
		StackName:              toBeGenerated,
		CloudFormationTemplate: []byte(toBeGenerated),
		VpcID:                  toBeGenerated,
		Cidr:                   opts.Cidr,
		PublicSubnets: []client.VpcSubnet{
			{
				ID:               toBeGenerated,
				Cidr:             toBeGenerated,
				AvailabilityZone: toBeGenerated,
			},
		},
		PrivateSubnets: []client.VpcSubnet{
			{
				ID:               toBeGenerated,
				Cidr:             toBeGenerated,
				AvailabilityZone: toBeGenerated,
			},
		},
		DatabaseSubnets: []client.VpcSubnet{
			{
				ID:               toBeGenerated,
				Cidr:             toBeGenerated,
				AvailabilityZone: toBeGenerated,
			},
		},
		DatabaseSubnetsGroupName: toBeGenerated,
	}, nil
}

func (v vpcService) DeleteVpc(_ context.Context, _ client.DeleteVpcOpts) error {
	fmt.Fprintf(v.out, formatDelete("VPC"))

	return nil
}

func (v vpcService) GetVPC(_ context.Context, _ api.ID) (*client.Vpc, error) {
	panic("implement me")
}
