// Package components contains functionality for
// creating cloud formation templates
package components

import (
	"fmt"

	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/oslokommune/okctl/pkg/cfn"
	cidrPkg "github.com/oslokommune/okctl/pkg/cfn/components/cidr"
	clusterPkg "github.com/oslokommune/okctl/pkg/cfn/components/cluster"
	"github.com/oslokommune/okctl/pkg/cfn/components/dbsubnetgroup"
	"github.com/oslokommune/okctl/pkg/cfn/components/eip"
	"github.com/oslokommune/okctl/pkg/cfn/components/internetgateway"
	"github.com/oslokommune/okctl/pkg/cfn/components/managedpolicy"
	"github.com/oslokommune/okctl/pkg/cfn/components/natgateway"
	"github.com/oslokommune/okctl/pkg/cfn/components/policydocument"
	"github.com/oslokommune/okctl/pkg/cfn/components/route"
	"github.com/oslokommune/okctl/pkg/cfn/components/routetable"
	"github.com/oslokommune/okctl/pkg/cfn/components/routetableassociation"
	"github.com/oslokommune/okctl/pkg/cfn/components/subnet"
	vpcPkg "github.com/oslokommune/okctl/pkg/cfn/components/vpc"
	"github.com/oslokommune/okctl/pkg/cfn/components/vpcgatewayattachment"
)

// VPCComposer contains the required state for building
// a VPC using cloud formation components
type VPCComposer struct {
	Name        string
	Environment string
	CidrBlock   string
	Region      string
}

// NewVPCComposer returns an initialised VPC composer
func NewVPCComposer(name, env, cidrBlock, region string) *VPCComposer {
	return &VPCComposer{
		Name:        name,
		Environment: env,
		CidrBlock:   cidrBlock,
		Region:      region,
	}
}

// Compose constructs the required cloud formation components
// nolint: funlen
func (v *VPCComposer) Compose() (*cfn.Composition, error) {
	composition := &cfn.Composition{}

	cluster := clusterPkg.New(v.Name, v.Environment)

	cidr, err := cidrPkg.NewDefault(v.CidrBlock)
	if err != nil {
		return nil, err
	}

	vpc := vpcPkg.New(cluster, cidr.Block)
	igw := internetgateway.New()
	gwa := vpcgatewayattachment.New(vpc, igw)
	composition.Resources = append(composition.Resources, vpc, igw, gwa)
	composition.Outputs = append(composition.Outputs, vpc)

	subnets, err := subnet.NewDefault(cidr.Block, v.Region, vpc, cluster)
	if err != nil {
		return nil, err
	}

	nats := make([]*natgateway.NatGateway, len(subnets.Public))

	// Public subnets
	prt := routetable.NewPublic(vpc)
	pr := route.NewPublic(gwa, prt, igw)
	composition.Resources = append(composition.Resources, prt, pr)

	for i, sub := range subnets.Public {
		// Create one NAT gateway for each public subnet
		e := eip.New(i, gwa)
		ngw := natgateway.New(i, gwa, e, sub)
		nats[i] = ngw

		// Associate the public subnet with the public route table
		assoc := routetableassociation.NewPublic(i, sub, prt)

		composition.Resources = append(composition.Resources, sub, assoc, ngw, e)
	}

	// Private subnets
	for i, sub := range subnets.Private {
		// Create a route table for each private subnet and associate
		// it with the subnet. Also add a route to the NAT gateway
		// so the instances can reach the internet
		rt := routetable.NewPrivate(i, vpc)
		r := route.NewPrivate(i, gwa, rt, nats[i%len(subnets.Private)])
		assoc := routetableassociation.NewPrivate(i, sub, rt)

		composition.Resources = append(composition.Resources, sub, rt, r, assoc)
	}

	composition.Outputs = append(composition.Outputs, subnets)

	dbSubnets := make([]cfn.Referencer, len(subnets.Database))

	for i, sub := range subnets.Database {
		dbSubnets[i] = sub

		composition.Resources = append(composition.Resources, sub)
	}

	dsg := dbsubnetgroup.New(dbSubnets)

	composition.Resources = append(composition.Resources, dsg)

	return composition, nil
}

// Ensure that VPCComposer implements the Composer interface
var _ cfn.Composer = &VPCComposer{}

// ExternalSecretsPolicyComposer contains state for building
// a managed iam policy compatible with external-secrets
type ExternalSecretsPolicyComposer struct {
	Repository  string
	Environment string
}

// NewExternalSecretsPolicyComposer returns a managed IAM policy
// that allows: https://github.com/godaddy/kubernetes-external-secrets
// to read SSM parameters and make them available as Kubernetes Secrets
func NewExternalSecretsPolicyComposer(repository, env string) *ExternalSecretsPolicyComposer {
	return &ExternalSecretsPolicyComposer{
		Repository:  repository,
		Environment: env,
	}
}

// Compose returns the cloud formation components required for building
// the policy
func (e *ExternalSecretsPolicyComposer) Compose() (*cfn.Composition, error) {
	p := e.ManagedPolicy()

	return &cfn.Composition{
		Outputs:   []cfn.StackOutputer{p},
		Resources: []cfn.ResourceNamer{p},
	}, nil
}

// ManagedPolicy returns a managed policy
func (e *ExternalSecretsPolicyComposer) ManagedPolicy() *managedpolicy.ManagedPolicy {
	policyName := fmt.Sprintf("okctl-%s-%s-ExternalSecretsServiceAccountPolicy", e.Repository, e.Environment)
	policyDesc := "Service account policy for reading SSM parameters"

	d := &policydocument.PolicyDocument{
		Version: policydocument.Version,
		Statement: []policydocument.StatementEntry{
			{
				Effect: policydocument.EffectTypeAllow,
				Action: []string{
					"ssm:GetParameter",
				},
				Resource: []string{
					ssmParameterARN("*"),
				},
			},
		},
	}

	return managedpolicy.New("ExternalSecretsPolicy", policyName, policyDesc, d)
}

// ssmParameterARN returns a valid resource SSM
// parameter ARN
func ssmParameterARN(resource string) string {
	return cloudformation.Sub(
		fmt.Sprintf(
			"arn:aws:ssm:${%s}:${%s}:parameter/%s",
			policydocument.PseudoParamRegion,
			policydocument.PseudoParamAccountID,
			resource,
		),
	)
}
