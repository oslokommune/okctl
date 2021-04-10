package core

import (
	"context"

	"github.com/oslokommune/okctl/pkg/cfn"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
)

type vpcService struct {
	api   client.VPCAPI
	state client.VPCState
}

func (s *vpcService) GetVPC(_ context.Context, id api.ID) (*client.Vpc, error) {
	return s.state.GetVpc(cfn.NewStackNamer().Vpc(id.ClusterName))
}

func (s *vpcService) CreateVpc(_ context.Context, opts client.CreateVpcOpts) (*client.Vpc, error) {
	vpc, err := s.api.CreateVpc(api.CreateVpcOpts{
		ID:      opts.ID,
		Cidr:    opts.Cidr,
		Minimal: opts.Minimal,
	})
	if err != nil {
		return nil, err
	}

	v := &client.Vpc{
		ID:                     vpc.ID,
		StackName:              vpc.StackName,
		CloudFormationTemplate: vpc.CloudFormationTemplate,
		VpcID:                  vpc.VpcID,
		Cidr:                   vpc.Cidr,
		PublicSubnets: func() (subs []client.VpcSubnet) {
			for _, sub := range vpc.PublicSubnets {
				subs = append(subs, client.VpcSubnet{
					ID:               sub.ID,
					Cidr:             sub.Cidr,
					AvailabilityZone: sub.AvailabilityZone,
				})
			}

			return subs
		}(),
		PrivateSubnets: func() (subs []client.VpcSubnet) {
			for _, sub := range vpc.PrivateSubnets {
				subs = append(subs, client.VpcSubnet{
					ID:               sub.ID,
					Cidr:             sub.Cidr,
					AvailabilityZone: sub.AvailabilityZone,
				})
			}

			return subs
		}(),
		DatabaseSubnets: func() (subs []client.VpcSubnet) {
			for _, sub := range vpc.DatabaseSubnets {
				subs = append(subs, client.VpcSubnet{
					ID:               sub.ID,
					Cidr:             sub.Cidr,
					AvailabilityZone: sub.AvailabilityZone,
				})
			}

			return subs
		}(),
		DatabaseSubnetsGroupName: vpc.DatabaseSubnetsGroupName,
	}

	err = s.state.SaveVpc(v)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func (s *vpcService) DeleteVpc(_ context.Context, opts client.DeleteVpcOpts) error {
	err := s.api.DeleteVpc(api.DeleteVpcOpts{
		ID: opts.ID,
	})
	if err != nil {
		return err
	}

	err = s.state.RemoveVpc(cfn.NewStackNamer().Vpc(opts.ID.ClusterName))
	if err != nil {
		return err
	}

	return nil
}

// NewVPCService returns an initialised VPC service
func NewVPCService(api client.VPCAPI, state client.VPCState) client.VPCService {
	return &vpcService{
		api:   api,
		state: state,
	}
}
