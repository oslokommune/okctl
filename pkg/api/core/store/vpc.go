package store

import (
	"encoding/json"

	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/storage/state"
)

type vpc struct {
	provider state.PersisterProvider
}

// Vpc represents the information that is stored
// about a vpc
type Vpc struct {
	StackName      string
	ID             string
	PublicSubnets  []VpcSubnet
	PrivateSubnets []VpcSubnet
}

// VpcSubnet represents the information that is stored
// about a vpc subnet
type VpcSubnet struct {
	ID               string
	Cidr             string
	AvailabilityZone string
}

func storeFromDef(vpc *api.Vpc) Vpc {
	pub := make([]VpcSubnet, len(vpc.PublicSubnets))

	for i, p := range vpc.PublicSubnets {
		pub[i] = VpcSubnet{
			ID:               p.ID,
			Cidr:             p.Cidr,
			AvailabilityZone: p.AvailabilityZone,
		}
	}

	pri := make([]VpcSubnet, len(vpc.PrivateSubnets))

	for i, p := range vpc.PrivateSubnets {
		pri[i] = VpcSubnet{
			ID:               p.ID,
			Cidr:             p.Cidr,
			AvailabilityZone: p.AvailabilityZone,
		}
	}

	vpcState := Vpc{
		StackName:      vpc.StackName,
		ID:             vpc.ID,
		PublicSubnets:  pub,
		PrivateSubnets: pri,
	}

	return vpcState
}

func defFromStore(vpc *Vpc) *api.Vpc {
	pub := make([]api.VpcSubnet, len(vpc.PublicSubnets))

	for i, p := range vpc.PublicSubnets {
		pub[i] = api.VpcSubnet{
			ID:               p.ID,
			Cidr:             p.Cidr,
			AvailabilityZone: p.AvailabilityZone,
		}
	}

	pri := make([]api.VpcSubnet, len(vpc.PrivateSubnets))

	for i, p := range vpc.PrivateSubnets {
		pri[i] = api.VpcSubnet{
			ID:               p.ID,
			Cidr:             p.Cidr,
			AvailabilityZone: p.AvailabilityZone,
		}
	}

	return &api.Vpc{
		StackName:      vpc.StackName,
		ID:             vpc.ID,
		PublicSubnets:  pub,
		PrivateSubnets: pri,
	}
}

// SaveVpc stores a vpc
func (v *vpc) SaveVpc(vpc *api.Vpc) error {
	err := v.provider.Repository().WriteToDefault("vpc_cloud_formation", vpc.CloudFormationTemplate)
	if err != nil {
		return err
	}

	data, err := json.Marshal(storeFromDef(vpc))
	if err != nil {
		return errors.E(err, "failed to create json of vpc state")
	}

	return v.provider.Repository().WriteToDefault("vpc_outputs", data)
}

// DeleteVpc removes a vpc from storage
func (v *vpc) DeleteVpc(env, repoName string) error {
	err := v.provider.Repository().DeleteDefault("vpc_cloud_formation")
	if err != nil {
		return err
	}

	return v.provider.Repository().DeleteDefault("vpc_outputs")
}

// GetVpc returns a vpc from storage
func (v *vpc) GetVpc() (*api.Vpc, error) {
	data, err := v.provider.Repository().ReadFromDefault("vpc_outputs")
	if err != nil {
		return nil, errors.E(err, "failed to read vpc state")
	}

	vpcStored := &Vpc{}

	err = json.Unmarshal(data, vpcStored)
	if err != nil {
		return nil, errors.E(err, "failed to unmarshal vpc outputs")
	}

	ret := defFromStore(vpcStored)

	template, err := v.provider.Repository().ReadFromDefault("vpc_cloud_formation")
	if err != nil {
		return nil, errors.E(err, "failed to read vpc cloud formation template")
	}

	ret.CloudFormationTemplate = template

	return ret, nil
}

// NewVpcStore returns an instantiated vpc store
func NewVpcStore(provider state.PersisterProvider) api.VpcStore {
	return &vpc{
		provider: provider,
	}
}
