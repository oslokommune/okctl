package storm

import (
	"errors"
	"time"

	"github.com/oslokommune/okctl/pkg/cfn"

	stormpkg "github.com/asdine/storm/v3"
	"github.com/oslokommune/okctl/pkg/breeze"
	"github.com/oslokommune/okctl/pkg/client"
)

type vpcState struct {
	node breeze.Client
}

// Vpc contains storm compatible state
type Vpc struct {
	Metadata `storm:"inline"`

	ID                       ID
	StackName                string `storm:"unique"`
	CloudFormationTemplate   string
	VpcID                    string
	Cidr                     string
	PublicSubnets            []VpcSubnet
	PrivateSubnets           []VpcSubnet
	DatabaseSubnets          []VpcSubnet
	DatabaseSubnetsGroupName string
}

// NewVpc returns storm compatible state
func NewVpc(v *client.Vpc, meta Metadata) *Vpc {
	return &Vpc{
		Metadata:               meta,
		ID:                     NewID(v.ID),
		StackName:              v.StackName,
		CloudFormationTemplate: string(v.CloudFormationTemplate),
		VpcID:                  v.VpcID,
		Cidr:                   v.Cidr,
		PublicSubnets: func() (subs []VpcSubnet) {
			for _, s := range v.PublicSubnets {
				subs = append(subs, NewVpcSubnet(s))
			}

			return subs
		}(),
		PrivateSubnets: func() (subs []VpcSubnet) {
			for _, s := range v.PrivateSubnets {
				subs = append(subs, NewVpcSubnet(s))
			}

			return subs
		}(),
		DatabaseSubnets: func() (subs []VpcSubnet) {
			for _, s := range v.DatabaseSubnets {
				subs = append(subs, NewVpcSubnet(s))
			}

			return subs
		}(),
		DatabaseSubnetsGroupName: v.DatabaseSubnetsGroupName,
	}
}

// Convert to client.Vpc
func (v *Vpc) Convert() *client.Vpc {
	return &client.Vpc{
		ID:                     v.ID.Convert(),
		StackName:              v.StackName,
		CloudFormationTemplate: []byte(v.CloudFormationTemplate),
		VpcID:                  v.VpcID,
		Cidr:                   v.Cidr,
		PublicSubnets: func() (subs []client.VpcSubnet) {
			for _, s := range v.PublicSubnets {
				subs = append(subs, s.Convert())
			}

			return subs
		}(),
		PrivateSubnets: func() (subs []client.VpcSubnet) {
			for _, s := range v.PrivateSubnets {
				subs = append(subs, s.Convert())
			}

			return subs
		}(),
		DatabaseSubnets: func() (subs []client.VpcSubnet) {
			for _, s := range v.DatabaseSubnets {
				subs = append(subs, s.Convert())
			}

			return subs
		}(),
		DatabaseSubnetsGroupName: v.DatabaseSubnetsGroupName,
	}
}

// VpcSubnet represents an aws vpc subnet
type VpcSubnet struct {
	ID               string
	Cidr             string
	AvailabilityZone string
}

// NewVpcSubnet returns storm compatible state
func NewVpcSubnet(s client.VpcSubnet) VpcSubnet {
	return VpcSubnet{
		ID:               s.ID,
		Cidr:             s.Cidr,
		AvailabilityZone: s.AvailabilityZone,
	}
}

// Convert to client.VpcSubnet
func (s *VpcSubnet) Convert() client.VpcSubnet {
	return client.VpcSubnet{
		ID:               s.ID,
		Cidr:             s.Cidr,
		AvailabilityZone: s.AvailabilityZone,
	}
}

func (v *vpcState) SaveVpc(vpc *client.Vpc) error {
	existing, err := v.getVpc(vpc.StackName)
	if err != nil && !errors.Is(err, stormpkg.ErrNotFound) {
		return err
	}

	if errors.Is(err, stormpkg.ErrNotFound) {
		return v.node.Save(NewVpc(vpc, NewMetadata()))
	}

	existing.Metadata.UpdatedAt = time.Now()

	return v.node.Save(NewVpc(vpc, existing.Metadata))
}

func (v *vpcState) GetVpc(stackName string) (*client.Vpc, error) {
	vpc, err := v.getVpc(stackName)
	if err != nil {
		return nil, err
	}

	return vpc.Convert(), nil
}

func (v *vpcState) getVpc(stackName string) (*Vpc, error) {
	vpc := &Vpc{}

	err := v.node.One("StackName", stackName, vpc)
	if err != nil {
		return nil, err
	}

	return vpc, nil
}

func (v *vpcState) HasVPC(clusterName string) (bool, error) {
	_, err := v.GetVpc(cfn.NewStackNamer().Vpc(clusterName))
	if err != nil {
		if errors.Is(err, stormpkg.ErrNotFound) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (v *vpcState) RemoveVpc(stackName string) error {
	vpc := &Vpc{}

	err := v.node.One("StackName", stackName, vpc)
	if err != nil {
		if errors.Is(err, stormpkg.ErrNotFound) {
			return nil
		}

		return err
	}

	return v.node.DeleteStruct(vpc)
}

// NewVpcState returns an initialised state client
func NewVpcState(node breeze.Client) client.VPCState {
	return &vpcState{
		node: node,
	}
}
