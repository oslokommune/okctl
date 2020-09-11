package state

import (
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/client/store"
	"github.com/oslokommune/okctl/pkg/config/state"
)

type vpcState struct {
	state state.Vpcer
}

func (s *vpcState) DeleteVpc(_ api.ID) (*store.Report, error) {
	return s.state.DeleteVPC()
}

func (s *vpcState) SaveVpc(vpc *api.Vpc) (*store.Report, error) {
	v := s.state.GetVPC()

	v.VpcID = vpc.VpcID
	v.CIDR = vpc.Cidr
	v.Subnets = map[string][]*state.VPCSubnet{
		state.SubnetTypePublic: func() (subnets []*state.VPCSubnet) {
			for _, sub := range vpc.PublicSubnets {
				subnets = append(subnets, &state.VPCSubnet{
					CIDR:             sub.Cidr,
					AvailabilityZone: sub.AvailabilityZone,
				})
			}

			return subnets
		}(),
		state.SubnetTypePrivate: func() (subnets []*state.VPCSubnet) {
			for _, sub := range vpc.PrivateSubnets {
				subnets = append(subnets, &state.VPCSubnet{
					CIDR:             sub.Cidr,
					AvailabilityZone: sub.AvailabilityZone,
				})
			}

			return subnets
		}(),
	}

	return s.state.SaveVPC(v)
}

// NewVpcState returns an initialised vpc state handler
func NewVpcState(state state.Vpcer) client.VPCState {
	return &vpcState{
		state: state,
	}
}
