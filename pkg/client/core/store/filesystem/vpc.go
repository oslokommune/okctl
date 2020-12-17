package filesystem

import (
	"github.com/oslokommune/okctl/pkg/client"

	"github.com/oslokommune/okctl/pkg/client/store"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/spf13/afero"
)

type vpcStore struct {
	paths Paths
	fs    *afero.Afero
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

func (s *vpcStore) SaveVpc(vpc *api.Vpc) (*store.Report, error) {
	report, err := store.NewFileSystem(s.paths.BaseDir, s.fs).
		StoreStruct(s.paths.OutputFile, storeFromDef(vpc), store.ToJSON()).
		StoreBytes(s.paths.CloudFormationFile, vpc.CloudFormationTemplate).
		Do()
	if err != nil {
		return nil, err
	}

	return report, nil
}

func (s *vpcStore) DeleteVpc(_ api.ID) (*store.Report, error) {
	report, err := store.NewFileSystem(s.paths.BaseDir, s.fs).
		Remove(s.paths.OutputFile).
		Remove(s.paths.CloudFormationFile).
		Do()
	if err != nil {
		return nil, err
	}

	_, _ = store.NewFileSystem(s.paths.BaseDir, s.fs).
		Remove("").
		Do()

	return report, nil
}

func (s *vpcStore) GetVpc(_ api.ID) (*api.Vpc, error) {
	vpcOutputs := &Vpc{}

	var template []byte

	_, err := store.NewFileSystem(s.paths.BaseDir, s.fs).
		GetStruct(s.paths.OutputFile, vpcOutputs, store.FromJSON()).
		GetBytes(s.paths.CloudFormationFile, func(_ string, data []byte) {
			template = data
		}).
		Do()
	if err != nil {
		return nil, err
	}

	ret := defFromStore(vpcOutputs)
	ret.CloudFormationTemplate = template

	return ret, nil
}

// NewVpcStore returns an instantiated vpc store
func NewVpcStore(paths Paths, fs *afero.Afero) client.VPCStore {
	return &vpcStore{
		paths: paths,
		fs:    fs,
	}
}
