package filesystem

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/client"

	"github.com/oslokommune/okctl/pkg/client/store"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/spf13/afero"
)

type vpcStore struct {
	stackOutputsFileName   string
	cloudFormationFileName string
	baseDir                string
	fs                     *afero.Afero
}

// Vpc represents the information that is stored
// about a vpc
type Vpc struct {
	ID             api.ID
	StackName      string
	VpcID          string
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
		ID:             vpc.ID,
		StackName:      vpc.StackName,
		VpcID:          vpc.VpcID,
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
		ID:             vpc.ID,
		StackName:      vpc.StackName,
		VpcID:          vpc.VpcID,
		PublicSubnets:  pub,
		PrivateSubnets: pri,
	}
}

func (s *vpcStore) SaveVpc(vpc *api.Vpc) error {
	_, err := store.NewFileSystem(s.baseDir, s.fs).
		StoreStruct(s.stackOutputsFileName, storeFromDef(vpc), store.ToJSON()).
		StoreBytes(s.cloudFormationFileName, vpc.CloudFormationTemplate).
		Do()
	if err != nil {
		return fmt.Errorf("failed to store vpc: %w", err)
	}

	return nil
}

func (s *vpcStore) DeleteVpc(_ api.ID) error {
	_, err := store.NewFileSystem(s.baseDir, s.fs).
		Remove(s.stackOutputsFileName).
		Remove(s.cloudFormationFileName).
		Do()
	if err != nil {
		return fmt.Errorf("failed to delete vpc: %w", err)
	}

	return nil
}

func (s *vpcStore) GetVpc(_ api.ID) (*api.Vpc, error) {
	vpcOutputs := &Vpc{}

	var template []byte

	callback := func(_ string, data []byte) error {
		template = data
		return nil
	}

	_, err := store.NewFileSystem(s.baseDir, s.fs).
		GetStruct(s.stackOutputsFileName, vpcOutputs, store.FromJSON()).
		GetBytes(s.cloudFormationFileName, callback).
		Do()
	if err != nil {
		return nil, fmt.Errorf("failed to get vpc: %w", err)
	}

	ret := defFromStore(vpcOutputs)
	ret.CloudFormationTemplate = template

	return ret, nil
}

// NewVpcStore returns an instantiated vpc store
func NewVpcStore(stackOutputsFileName, cloudFormationFileName, baseDir string, fs *afero.Afero) client.VPCStore {
	return &vpcStore{
		stackOutputsFileName:   stackOutputsFileName,
		cloudFormationFileName: cloudFormationFileName,
		baseDir:                baseDir,
		fs:                     fs,
	}
}
