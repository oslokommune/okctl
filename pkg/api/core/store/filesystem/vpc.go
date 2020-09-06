package filesystem

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/client/store"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/spf13/afero"
)

type vpc struct {
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

// SaveVpc stores a vpc
func (v *vpc) SaveVpc(vpc *api.Vpc) error {
	_, err := store.NewFileSystem(v.baseDir, v.fs).
		StoreStruct(v.stackOutputsFileName, storeFromDef(vpc), store.ToJSON()).
		StoreBytes(v.cloudFormationFileName, vpc.CloudFormationTemplate).
		Do()
	if err != nil {
		return fmt.Errorf("failed to store vpc: %w", err)
	}

	return nil
}

// DeleteVpc removes a vpc from storage
func (v *vpc) DeleteVpc(_, _ string) error {
	_, err := store.NewFileSystem(v.baseDir, v.fs).
		Remove(v.stackOutputsFileName).
		Remove(v.cloudFormationFileName).
		Do()
	if err != nil {
		return fmt.Errorf("failed to delete vpc: %w", err)
	}

	return nil
}

// GetVpc returns a vpc from storage
func (v *vpc) GetVpc() (*api.Vpc, error) {
	vpcOutputs := &Vpc{}

	var template []byte

	callback := func(_ string, data []byte) error {
		template = data
		return nil
	}

	_, err := store.NewFileSystem(v.baseDir, v.fs).
		GetStruct(v.stackOutputsFileName, vpcOutputs, store.FromJSON()).
		GetBytes(v.cloudFormationFileName, callback).
		Do()
	if err != nil {
		return nil, fmt.Errorf("failed to get vpc: %w", err)
	}

	ret := defFromStore(vpcOutputs)
	ret.CloudFormationTemplate = template

	return ret, nil
}

// NewVpcStore returns an instantiated vpc store
func NewVpcStore(stackOutputsFileName, cloudFormationFileName, baseDir string, fs *afero.Afero) api.VpcStore {
	return &vpc{
		stackOutputsFileName:   stackOutputsFileName,
		cloudFormationFileName: cloudFormationFileName,
		baseDir:                baseDir,
		fs:                     fs,
	}
}
