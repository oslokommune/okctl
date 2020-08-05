// Package vpc knows how to create a VPC cloud formation stack
package vpc

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/cfn"
	cidrPkg "github.com/oslokommune/okctl/pkg/cfn/components/cidr"
	clusterPkg "github.com/oslokommune/okctl/pkg/cfn/components/cluster"
	"github.com/oslokommune/okctl/pkg/cfn/components/dbsubnetgroup"
	"github.com/oslokommune/okctl/pkg/cfn/components/eip"
	"github.com/oslokommune/okctl/pkg/cfn/components/internetgateway"
	"github.com/oslokommune/okctl/pkg/cfn/components/natgateway"
	"github.com/oslokommune/okctl/pkg/cfn/components/route"
	"github.com/oslokommune/okctl/pkg/cfn/components/routetable"
	"github.com/oslokommune/okctl/pkg/cfn/components/routetableassociation"
	"github.com/oslokommune/okctl/pkg/cfn/components/subnet"
	vpcPkg "github.com/oslokommune/okctl/pkg/cfn/components/vpc"
	"github.com/oslokommune/okctl/pkg/cfn/components/vpcgatewayattachment"
)

// Builder stores state for creating a cloud formation VPC stack
type Builder struct {
	Name      string
	Env       string
	CidrBlock string
	Region    string

	outputs   []cfn.Outputer
	resources []cfn.ResourceNamer
}

// New returns a VPC builder
func New(name, env, cidrBlock, region string) *Builder {
	return &Builder{
		Name:      name,
		Env:       env,
		CidrBlock: cidrBlock,
		Region:    region,
	}
}

// Resources returns all cloud formation resources
func (b *Builder) Resources() []cfn.ResourceNamer {
	return b.resources
}

// Outputs returns all cloud formation outputs
func (b *Builder) Outputs() []cfn.Outputer {
	return b.outputs
}

// StackName returns the name of the stack
func (b *Builder) StackName() string {
	return StackName(b.Name, b.Env)
}

// Build all resources needed for creating a cloud formation VPC
func (b *Builder) Build() error {
	cluster := clusterPkg.New(b.Name, b.Env)

	cidr, err := cidrPkg.NewDefault(b.CidrBlock)
	if err != nil {
		return err
	}

	vpc := vpcPkg.New(cluster, cidr.Block)
	igw := internetgateway.New()
	gwa := vpcgatewayattachment.New(vpc, igw)
	prt := routetable.NewPublic(vpc)
	pr := route.NewPublic(gwa, prt, igw)

	b.resources = append(b.resources, vpc, igw, gwa, prt, pr)
	b.outputs = append(b.outputs, vpc)

	subnets, err := subnet.NewDefault(cidr.Block, b.Region, vpc, cluster)
	if err != nil {
		return err
	}

	for i, sub := range subnets.Private {
		rt := routetable.NewPrivate(i, vpc)
		e := eip.New(i, gwa)
		ngw := natgateway.New(i, gwa, e, sub)
		r := route.NewPrivate(i, gwa, rt, ngw)
		assoc := routetableassociation.NewPrivate(i, sub, rt)

		b.resources = append(b.resources, sub, rt, e, ngw, r, assoc)
	}

	for i, sub := range subnets.Public {
		assoc := routetableassociation.NewPublic(i, sub, prt)

		b.resources = append(b.resources, sub, assoc)
	}

	b.outputs = append(b.outputs, subnets)

	dbSubnets := make([]cfn.Referencer, len(subnets.Database))

	for i, sub := range subnets.Database {
		dbSubnets[i] = sub

		b.resources = append(b.resources, sub)
	}

	dsg := dbsubnetgroup.New(dbSubnets)

	b.resources = append(b.resources, dsg)

	return nil
}

// StackName returns a consistent stack name for a VPC
func StackName(name, env string) string {
	return fmt.Sprintf("%s-%s-okctl-vpc", name, env)
}
