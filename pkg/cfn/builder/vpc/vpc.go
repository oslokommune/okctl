// Package vpc knows how to create a VPC cloud formation stack
package vpc

import (
	"fmt"

	"github.com/awslabs/goformation/v4/cloudformation"
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

	template *cloudformation.Template
}

// New returns a VPC builder
func New(name, env, cidrBlock, region string) *Builder {
	return &Builder{
		Name:      name,
		Env:       env,
		CidrBlock: cidrBlock,
		Region:    region,
		template:  cloudformation.NewTemplate(),
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
func (b *Builder) Build() ([]byte, error) {
	err := b.create()
	if err != nil {
		return nil, err
	}

	err = b.collectOutputs()
	if err != nil {
		return nil, err
	}

	err = b.collectResources()
	if err != nil {
		return nil, err
	}

	return b.template.YAML()
}

func (b *Builder) collectResources() error {
	for _, resource := range b.Resources() {
		if _, hasKey := b.template.Resources[resource.Name()]; hasKey {
			return fmt.Errorf("already have resource with name: %s", resource.Name())
		}

		b.template.Resources[resource.Name()] = resource.Resource()
	}

	return nil
}

func (b *Builder) collectOutputs() error {
	for _, output := range b.Outputs() {
		for key, value := range output.NamedOutputs() {
			if _, hasKey := b.template.Outputs[key]; hasKey {
				return fmt.Errorf("already have output with name: %s", key)
			}

			b.template.Outputs[key] = value
		}
	}

	return nil
}

//nolint: funlen
func (b *Builder) create() error {
	cluster := clusterPkg.New(b.Name, b.Env)

	cidr, err := cidrPkg.NewDefault(b.CidrBlock)
	if err != nil {
		return err
	}

	vpc := vpcPkg.New(cluster, cidr.Block)
	igw := internetgateway.New()
	gwa := vpcgatewayattachment.New(vpc, igw)
	b.resources = append(b.resources, vpc, igw, gwa)
	b.outputs = append(b.outputs, vpc)

	subnets, err := subnet.NewDefault(cidr.Block, b.Region, vpc, cluster)
	if err != nil {
		return err
	}

	nats := make([]*natgateway.NatGateway, len(subnets.Public))

	// Public subnets
	prt := routetable.NewPublic(vpc)
	pr := route.NewPublic(gwa, prt, igw)
	b.resources = append(b.resources, prt, pr)

	for i, sub := range subnets.Public {
		// Create one NAT gateway for each public subnet
		e := eip.New(i, gwa)
		ngw := natgateway.New(i, gwa, e, sub)
		nats[i] = ngw

		// Associate the public subnet with the public route table
		assoc := routetableassociation.NewPublic(i, sub, prt)

		b.resources = append(b.resources, sub, assoc, ngw, e)
	}

	// Private subnets
	for i, sub := range subnets.Private {
		// Create a route table for each private subnet and associate
		// it with the subnet. Also add a route to the NAT gateway
		// so the instances can reach the internet
		rt := routetable.NewPrivate(i, vpc)
		r := route.NewPrivate(i, gwa, rt, nats[i%len(subnets.Private)])
		assoc := routetableassociation.NewPrivate(i, sub, rt)

		b.resources = append(b.resources, sub, rt, r, assoc)
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
	return fmt.Sprintf("okctl-vpc-%s-%s", name, env)
}
