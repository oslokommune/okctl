package core

import (
	"context"
	"fmt"

	"github.com/oslokommune/okctl/pkg/cfn"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
)

type vpcService struct {
	service api.VpcService
	state   client.VPCState
}

func (s *vpcService) GetVPC(_ context.Context, id api.ID) (*client.Vpc, error) {
	vpc, err := s.state.GetVpc(cfn.NewStackNamer().Vpc(id.ClusterName))
	if err != nil {
		return nil, fmt.Errorf("getting vpc: %w", err)
	}

	return vpc, nil
}

func (s *vpcService) CreateVpc(context context.Context, opts client.CreateVpcOpts) (*client.Vpc, error) {
	vpc, err := s.service.CreateVpc(context, api.CreateVpcOpts{
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

func (s *vpcService) DeleteVpc(context context.Context, opts client.DeleteVpcOpts) error {
	err := s.service.DeleteVpc(context, api.DeleteVpcOpts{
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
func NewVPCService(service api.VpcService, state client.VPCState) client.VPCService {
	return &vpcService{
		service: service,
		state:   state,
	}
}
