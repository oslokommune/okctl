package filesystem

import (
	"encoding/json"
	"path"

	"github.com/mishudark/errors"
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
	data, err := json.Marshal(storeFromDef(vpc))
	if err != nil {
		return errors.E(err, "failed to create json of vpc state")
	}

	err = v.fs.MkdirAll(v.baseDir, 0o744)
	if err != nil {
		return err
	}

	err = v.fs.WriteFile(path.Join(v.baseDir, v.stackOutputsFileName), data, 0o644)
	if err != nil {
		return err
	}

	return v.fs.WriteFile(path.Join(v.baseDir, v.cloudFormationFileName), vpc.CloudFormationTemplate, 0o644)
}

// DeleteVpc removes a vpc from storage
func (v *vpc) DeleteVpc(env, repoName string) error {
	err := v.fs.Remove(path.Join(v.baseDir, v.stackOutputsFileName))
	if err != nil {
		return err
	}

	return v.fs.Remove(path.Join(v.baseDir, v.cloudFormationFileName))
}

// GetVpc returns a vpc from storage
func (v *vpc) GetVpc() (*api.Vpc, error) {
	data, err := v.fs.ReadFile(path.Join(v.baseDir, v.stackOutputsFileName))
	if err != nil {
		return nil, errors.E(err, "failed to read vpc state")
	}

	vpcOutputs := &Vpc{}

	err = json.Unmarshal(data, vpcOutputs)
	if err != nil {
		return nil, errors.E(err, "failed to unmarshal vpc outputs")
	}

	ret := defFromStore(vpcOutputs)

	template, err := v.fs.ReadFile(path.Join(v.baseDir, v.cloudFormationFileName))
	if err != nil {
		return nil, errors.E(err, "failed to read vpc cloud formation template")
	}

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
