// Package components contains functionality for
// creating cloud formation templates
package components

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/cfn/components/certificate"

	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/oslokommune/okctl/pkg/cfn"
	cidrPkg "github.com/oslokommune/okctl/pkg/cfn/components/cidr"
	clusterPkg "github.com/oslokommune/okctl/pkg/cfn/components/cluster"
	"github.com/oslokommune/okctl/pkg/cfn/components/dbsubnetgroup"
	"github.com/oslokommune/okctl/pkg/cfn/components/eip"
	"github.com/oslokommune/okctl/pkg/cfn/components/hostedzone"
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

// MinimalVPCComposer contains the required state for building
// a VPC using cloud formation components
type MinimalVPCComposer struct {
	Name        string
	Environment string
	CidrBlock   string
	Region      string
}

// NewMinimalVPCComposer returns an initialised VPC composer
func NewMinimalVPCComposer(name, env, cidrBlock, region string) *MinimalVPCComposer {
	return &MinimalVPCComposer{
		Name:        name,
		Environment: env,
		CidrBlock:   cidrBlock,
		Region:      region,
	}
}

// Compose constructs the required cloud formation components
// nolint: funlen
func (v *MinimalVPCComposer) Compose() (*cfn.Composition, error) {
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

	var nat *natgateway.NatGateway

	// Public subnets
	prt := routetable.NewPublic(vpc)
	pr := route.NewPublic(gwa, prt, igw)
	composition.Resources = append(composition.Resources, prt, pr)

	for i, sub := range subnets.Public {
		// Create only one NAT gateway
		if nat == nil {
			e := eip.New(i, gwa)
			ngw := natgateway.New(i, gwa, e, sub)
			nat = ngw

			composition.Resources = append(composition.Resources, ngw, e)
		}

		// Associate the public subnet with the public route table
		assoc := routetableassociation.NewPublic(i, sub, prt)

		composition.Resources = append(composition.Resources, sub, assoc)
	}

	// Private subnets
	for i, sub := range subnets.Private {
		// Create a route table for each private subnet and associate
		// it with the subnet. Also add a route to the NAT gateway
		// so the instances can reach the internet
		rt := routetable.NewPrivate(i, vpc)
		r := route.NewPrivate(i, gwa, rt, nat) // Route all egress traffic through one NAT
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

// Ensure that MinimalVPCComposer implements the Composer interface
var _ cfn.Composer = &MinimalVPCComposer{}

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

// AlbIngressControllerPolicyComposer contains state for building
// a managed iam policy compatible with aws-alb-ingress-controller
type AlbIngressControllerPolicyComposer struct {
	Repository  string
	Environment string
}

// NewAlbIngressControllerPolicyComposer returns an initialised alb ingress controller composer
func NewAlbIngressControllerPolicyComposer(repository, env string) *AlbIngressControllerPolicyComposer {
	return &AlbIngressControllerPolicyComposer{
		Repository:  repository,
		Environment: env,
	}
}

// Compose builds the policy and returns the result
func (a *AlbIngressControllerPolicyComposer) Compose() (*cfn.Composition, error) {
	p := a.ManagedPolicy()

	return &cfn.Composition{
		Outputs:   []cfn.StackOutputer{p},
		Resources: []cfn.ResourceNamer{p},
	}, nil
}

// ManagedPolicy creates a managed policy
// nolint: funlen
func (a *AlbIngressControllerPolicyComposer) ManagedPolicy() *managedpolicy.ManagedPolicy {
	policyName := fmt.Sprintf("okctl-%s-%s-AlbIngressControllServiceAccountPolicy", a.Repository, a.Environment)
	policyDesc := "Service account policy for creat ALBs"

	d := &policydocument.PolicyDocument{
		Version: policydocument.Version,
		Statement: []policydocument.StatementEntry{
			{
				Effect: policydocument.EffectTypeAllow,
				Action: []string{
					"acm:DescribeCertificate",
					"acm:ListCertificates",
					"acm:GetCertificate",
				},
				Resource: []string{
					"*",
				},
			},
			{
				Effect: policydocument.EffectTypeAllow,
				Action: []string{
					"ec2:AuthorizeSecurityGroupIngress",
					"ec2:CreateSecurityGroup",
					"ec2:CreateTags",
					"ec2:DeleteTags",
					"ec2:DeleteSecurityGroup",
					"ec2:DescribeAccountAttributes",
					"ec2:DescribeAddresses",
					"ec2:DescribeInstances",
					"ec2:DescribeInstanceStatus",
					"ec2:DescribeInternetGateways",
					"ec2:DescribeNetworkInterfaces",
					"ec2:DescribeSecurityGroups",
					"ec2:DescribeSubnets",
					"ec2:DescribeTags",
					"ec2:DescribeVpcs",
					"ec2:ModifyInstanceAttribute",
					"ec2:ModifyNetworkInterfaceAttribute",
					"ec2:RevokeSecurityGroupIngress",
				},
				Resource: []string{
					"*",
				},
			},
			{
				Effect: policydocument.EffectTypeAllow,
				Action: []string{
					"elasticloadbalancing:AddListenerCertificates",
					"elasticloadbalancing:AddTags",
					"elasticloadbalancing:CreateListener",
					"elasticloadbalancing:CreateLoadBalancer",
					"elasticloadbalancing:CreateRule",
					"elasticloadbalancing:CreateTargetGroup",
					"elasticloadbalancing:DeleteListener",
					"elasticloadbalancing:DeleteLoadBalancer",
					"elasticloadbalancing:DeleteRule",
					"elasticloadbalancing:DeleteTargetGroup",
					"elasticloadbalancing:DeregisterTargets",
					"elasticloadbalancing:DescribeListenerCertificates",
					"elasticloadbalancing:DescribeListeners",
					"elasticloadbalancing:DescribeLoadBalancers",
					"elasticloadbalancing:DescribeLoadBalancerAttributes",
					"elasticloadbalancing:DescribeRules",
					"elasticloadbalancing:DescribeSSLPolicies",
					"elasticloadbalancing:DescribeTags",
					"elasticloadbalancing:DescribeTargetGroups",
					"elasticloadbalancing:DescribeTargetGroupAttributes",
					"elasticloadbalancing:DescribeTargetHealth",
					"elasticloadbalancing:ModifyListener",
					"elasticloadbalancing:ModifyLoadBalancerAttributes",
					"elasticloadbalancing:ModifyRule",
					"elasticloadbalancing:ModifyTargetGroup",
					"elasticloadbalancing:ModifyTargetGroupAttributes",
					"elasticloadbalancing:RegisterTargets",
					"elasticloadbalancing:RemoveListenerCertificates",
					"elasticloadbalancing:RemoveTags",
					"elasticloadbalancing:SetIpAddressType",
					"elasticloadbalancing:SetSecurityGroups",
					"elasticloadbalancing:SetSubnets",
					"elasticloadbalancing:SetWebACL",
				},
				Resource: []string{
					"*",
				},
			},
			{
				Effect: policydocument.EffectTypeAllow,
				Action: []string{
					"iam:CreateServiceLinkedRole",
					"iam:GetServerCertificate",
					"iam:ListServerCertificates",
				},
				Resource: []string{
					"*",
				},
			},
			{
				Effect: policydocument.EffectTypeAllow,
				Action: []string{
					"cognito-idp:DescribeUserPoolClient",
				},
				Resource: []string{
					"*",
				},
			},
			{
				Effect: policydocument.EffectTypeAllow,
				Action: []string{
					"waf-regional:GetWebACLForResource",
					"waf-regional:GetWebACL",
					"waf-regional:AssociateWebACL",
					"waf-regional:DisassociateWebACL",
				},
				Resource: []string{
					"*",
				},
			},
			{
				Effect: policydocument.EffectTypeAllow,
				Action: []string{
					"tag:GetResources",
					"tag:TagResources",
				},
				Resource: []string{
					"*",
				},
			},
			{
				Effect: policydocument.EffectTypeAllow,
				Action: []string{
					"waf:GetWebACL",
				},
				Resource: []string{
					"*",
				},
			},
			{
				Effect: policydocument.EffectTypeAllow,
				Action: []string{
					"wafv2:GetWebACL",
					"wafv2:GetWebACLForResource",
					"wafv2:AssociateWebACL",
					"wafv2:DisassociateWebACL",
				},
				Resource: []string{
					"*",
				},
			},
			{
				Effect: policydocument.EffectTypeAllow,
				Action: []string{
					"shield:DescribeProtection",
					"shield:GetSubscriptionState",
					"shield:DeleteProtection",
					"shield:CreateProtection",
					"shield:DescribeSubscription",
					"shield:ListProtections",
				},
				Resource: []string{
					"*",
				},
			},
		},
	}

	return managedpolicy.New("AlbIngressControllerPolicy", policyName, policyDesc, d)
}

// ExternalDNSPolicyComposer contains state for building
// a managed iam policy compatible with aws-alb-ingress-controller
type ExternalDNSPolicyComposer struct {
	Repository  string
	Environment string
}

// NewExternalDNSPolicyComposer returns an initialised alb ingress controller composer
func NewExternalDNSPolicyComposer(repository, env string) *ExternalDNSPolicyComposer {
	return &ExternalDNSPolicyComposer{
		Repository:  repository,
		Environment: env,
	}
}

// Compose builds the policy and returns the result
func (c *ExternalDNSPolicyComposer) Compose() (*cfn.Composition, error) {
	p := c.ManagedPolicy()

	return &cfn.Composition{
		Outputs:   []cfn.StackOutputer{p},
		Resources: []cfn.ResourceNamer{p},
	}, nil
}

// ManagedPolicy returns the policy
func (c *ExternalDNSPolicyComposer) ManagedPolicy() *managedpolicy.ManagedPolicy {
	policyName := fmt.Sprintf("okctl-%s-%s-ExternalDNSServiceAccountPolicy", c.Repository, c.Environment)
	policyDesc := "Service account policy for creating route53 hostnames"

	d := &policydocument.PolicyDocument{
		Version: policydocument.Version,
		Statement: []policydocument.StatementEntry{
			{
				Effect: policydocument.EffectTypeAllow,
				Action: []string{
					"route53:ChangeResourceRecordSets",
				},
				Resource: []string{
					"arn:aws:route53:::hostedzone/*",
				},
			},
			{
				Effect: policydocument.EffectTypeAllow,
				Action: []string{
					"route53:ListHostedZones",
					"route53:ListResourceRecordSets",
				},
				Resource: []string{
					"*",
				},
			},
		},
	}

	return managedpolicy.New("ExternalDNSPolicy", policyName, policyDesc, d)
}

// HostedZoneComposer contains state for creating a hosted zone
type HostedZoneComposer struct {
	FQDN    string
	Comment string
}

// NewHostedZoneComposer returns an initialised hosted zone composer
func NewHostedZoneComposer(fqdn, comment string) *HostedZoneComposer {
	return &HostedZoneComposer{
		FQDN:    fqdn,
		Comment: comment,
	}
}

// Compose returns the cloud formation components required for building
// the hosted zone
func (h *HostedZoneComposer) Compose() (*cfn.Composition, error) {
	zone := hostedzone.New(h.FQDN, h.Comment)

	return &cfn.Composition{
		Outputs:   []cfn.StackOutputer{zone},
		Resources: []cfn.ResourceNamer{zone},
	}, nil
}

// PublicCertificateComposer stores the state for the composer
type PublicCertificateComposer struct {
	FQDN         string
	HostedZoneID string
}

// NewPublicCertificateComposer returns an initialised composer
func NewPublicCertificateComposer(fqdn, hostedZoneID string) *PublicCertificateComposer {
	return &PublicCertificateComposer{
		FQDN:         fqdn,
		HostedZoneID: hostedZoneID,
	}
}

// Compose returns the resources and outputts for creating a certificate
func (c *PublicCertificateComposer) Compose() (*cfn.Composition, error) {
	cert := certificate.New(c.FQDN, c.HostedZoneID)

	return &cfn.Composition{
		Outputs:   []cfn.StackOutputer{cert},
		Resources: []cfn.ResourceNamer{cert},
	}, nil
}
