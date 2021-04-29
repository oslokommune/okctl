// Package components contains functionality for
// creating cloud formation templates
package components

import (
	"fmt"
	"strings"

	"github.com/oslokommune/okctl/pkg/cfn/components/containerrepository"

	"github.com/oslokommune/okctl/pkg/cfn/components/lambdafunction"

	"github.com/oslokommune/okctl/pkg/cfn/components/lambdapermission"

	"github.com/oslokommune/okctl/pkg/cfn/components/s3bucket"

	"github.com/oslokommune/okctl/pkg/cfn/components/securitygroup"

	"github.com/oslokommune/okctl/pkg/cfn/components/secrettargetattachment"

	"github.com/oslokommune/okctl/pkg/cfn/components/secret"

	"github.com/oslokommune/okctl/pkg/cfn/components/vpcendpoint"

	"github.com/oslokommune/okctl/pkg/cfn/components/dbinstance"

	"github.com/oslokommune/okctl/pkg/cfn/components/dbparametergroup"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/cfn/components/role"

	"github.com/oslokommune/okctl/pkg/cfn/components/userpooluser"
	"github.com/oslokommune/okctl/pkg/cfn/components/userpoolusertogroupattachment"

	"github.com/oslokommune/okctl/pkg/cfn/components/userpoolgroup"

	"github.com/oslokommune/okctl/pkg/cfn/components/recordset"

	"github.com/oslokommune/okctl/pkg/cfn/components/aliasrecordset"

	"github.com/oslokommune/okctl/pkg/cfn/components/userpooldomain"

	"github.com/oslokommune/okctl/pkg/cfn/components/userpoolclient"

	"github.com/oslokommune/okctl/pkg/cfn/components/userpool"

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
	ClusterName string
	CidrBlock   string
	Region      string
}

// NewVPCComposer returns an initialised VPC composer
func NewVPCComposer(name, cidrBlock, region string) *VPCComposer {
	return &VPCComposer{
		ClusterName: name,
		CidrBlock:   cidrBlock,
		Region:      region,
	}
}

// Compose constructs the required cloud formation components
// nolint: funlen
func (v *VPCComposer) Compose() (*cfn.Composition, error) {
	composition := &cfn.Composition{}

	cluster := clusterPkg.New(v.ClusterName)

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
	composition.Outputs = append(composition.Outputs, dsg)

	return composition, nil
}

// Ensure that VPCComposer implements the Composer interface
var _ cfn.Composer = &VPCComposer{}

// MinimalVPCComposer contains the required state for building
// a VPC using cloud formation components
type MinimalVPCComposer struct {
	ClusterName string
	CidrBlock   string
	Region      string
}

// NewMinimalVPCComposer returns an initialised VPC composer
func NewMinimalVPCComposer(clusterName, cidrBlock, region string) *MinimalVPCComposer {
	return &MinimalVPCComposer{
		ClusterName: clusterName,
		CidrBlock:   cidrBlock,
		Region:      region,
	}
}

// Compose constructs the required cloud formation components
// nolint: funlen
func (v *MinimalVPCComposer) Compose() (*cfn.Composition, error) {
	composition := &cfn.Composition{}

	cluster := clusterPkg.New(v.ClusterName)

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
	composition.Outputs = append(composition.Outputs, dsg)

	return composition, nil
}

// Ensure that MinimalVPCComposer implements the Composer interface
var _ cfn.Composer = &MinimalVPCComposer{}

// ExternalSecretsPolicyComposer contains state for building
// a managed iam policy compatible with external-secrets
type ExternalSecretsPolicyComposer struct {
	ClusterName string
}

// NewExternalSecretsPolicyComposer returns a managed IAM policy
// that allows: https://github.com/external-secrets/kubernetes-external-secrets
// to read SSM parameters and make them available as Kubernetes Secrets
func NewExternalSecretsPolicyComposer(clusterName string) *ExternalSecretsPolicyComposer {
	return &ExternalSecretsPolicyComposer{
		ClusterName: clusterName,
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
	policyName := fmt.Sprintf("okctl-%s-ExternalSecretsServiceAccountPolicy", e.ClusterName)
	policyDesc := "Service account policy for reading SSM parameters and ASM secrets"

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
			{
				Effect: policydocument.EffectTypeAllow,
				Action: []string{
					"secretsmanager:GetResourcePolicy",
					"secretsmanager:GetSecretValue",
					"secretsmanager:DescribeSecret",
					"secretsmanager:ListSecretVersionIds",
				},
				Resource: []string{
					asmParameterARN("*"),
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

// asmParameterARN returns a valid resource ASM
// parameter ARN
// arn:aws:secretsmanager:eu-west-1:932360772598:secret:*
func asmParameterARN(resource string) string {
	return cloudformation.Sub(
		fmt.Sprintf(
			"arn:aws:secretsmanager:${%s}:${%s}:secret:%s",
			policydocument.PseudoParamRegion,
			policydocument.PseudoParamAccountID,
			resource,
		),
	)
}

// AlbIngressControllerPolicyComposer contains state for building
// a managed iam policy compatible with aws-alb-ingress-controller
type AlbIngressControllerPolicyComposer struct {
	ClusterName string
}

// NewAlbIngressControllerPolicyComposer returns an initialised alb ingress controller composer
func NewAlbIngressControllerPolicyComposer(clusterName string) *AlbIngressControllerPolicyComposer {
	return &AlbIngressControllerPolicyComposer{
		ClusterName: clusterName,
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
	policyName := fmt.Sprintf("okctl-%s-AlbIngressControllServiceAccountPolicy", a.ClusterName)
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
	ClusterName string
}

// NewExternalDNSPolicyComposer returns an initialised alb ingress controller composer
func NewExternalDNSPolicyComposer(clusterName string) *ExternalDNSPolicyComposer {
	return &ExternalDNSPolicyComposer{
		ClusterName: clusterName,
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
	policyName := fmt.Sprintf("okctl-%s-ExternalDNSServiceAccountPolicy", c.ClusterName)
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

// Compose returns the resources and outputs for creating a certificate
func (c *PublicCertificateComposer) Compose() (*cfn.Composition, error) {
	cert := certificate.New(c.FQDN, c.HostedZoneID)

	return &cfn.Composition{
		Outputs:   []cfn.StackOutputer{cert},
		Resources: []cfn.ResourceNamer{cert},
	}, nil
}

// UserPool contains all state for building
// a cognito user pool cloud formation template
type UserPool struct {
	ClusterName    string
	CertificateARN string
	Domain         string
	HostedZoneID   string
}

// Compose returns the resources and outputs
func (u *UserPool) Compose() (*cfn.Composition, error) {
	composition := &cfn.Composition{}

	userPool := userpool.New(u.ClusterName)
	// Cognito User Pool Domain requires an A record on the base domain, for some or another
	// reason. Frequently we don't have this, so we create a placeholder record.
	placeholder := recordset.New("PlaceHolder", "1.1.1.1", RootDomain(u.Domain), u.HostedZoneID)
	upDomain := userpooldomain.New(u.Domain, u.CertificateARN, userPool, placeholder)
	group := userpoolgroup.New("admins", "Default admin group", userPool)

	composition.Resources = append(composition.Resources, userPool, upDomain, placeholder, group)
	composition.Outputs = append(composition.Outputs, userPool)

	return composition, nil
}

// UserPoolUser output of command
type UserPoolUser struct {
	UserPoolID string
	Email      string
}

// Compose userpool user and admin group attachment
func (u *UserPoolUser) Compose() (*cfn.Composition, error) {
	composition := &cfn.Composition{}

	userPoolUser := userpooluser.New(u.Email, "User pool user", u.UserPoolID)
	attachment := userpoolusertogroupattachment.New(userPoolUser, u.Email, "admins", u.UserPoolID)

	composition.Resources = append(composition.Resources, userPoolUser, attachment)
	composition.Outputs = append(composition.Outputs, userPoolUser)

	return composition, nil
}

// NewUserPoolUser add a new user into a userpool
func NewUserPoolUser(email, userpoolid string) *UserPoolUser {
	return &UserPoolUser{
		Email:      email,
		UserPoolID: userpoolid,
	}
}

// RootDomain extract the root domain
func RootDomain(domain string) string {
	if len(domain) == 0 || !strings.Contains(domain, ".") {
		return domain
	}

	parts := strings.Split(domain, ".")
	if len(parts) == 2 { // nolint: gomnd
		return domain
	}

	return strings.Join(parts[1:], ".")
}

// NewUserPool returns an initialised composer
// for creating a cognito user pool with clients
func NewUserPool(clusterName, domain, hostedZoneID, certificateARN string) *UserPool {
	return &UserPool{
		ClusterName:    clusterName,
		Domain:         domain,
		CertificateARN: certificateARN,
		HostedZoneID:   hostedZoneID,
	}
}

// UserPoolClient contains state for building a
// a cognito user pool client cloud formation template
type UserPoolClient struct {
	ClusterName string
	Purpose     string
	CallbackURL string
	UserPoolID  string
}

// Compose returns outputs and resources for a cloud formation stack
func (c *UserPoolClient) Compose() (*cfn.Composition, error) {
	composition := &cfn.Composition{}

	upc := userpoolclient.New(c.Purpose, c.ClusterName, c.CallbackURL, c.UserPoolID)

	composition.Resources = append(composition.Resources, upc)
	composition.Outputs = append(composition.Outputs, upc)

	return composition, nil
}

// NewUserPoolClient returns an initialised composer for
// creating a cognito user pool client
func NewUserPoolClient(purpose, clusterName, callbackURL, userPoolID string) *UserPoolClient {
	return &UserPoolClient{
		ClusterName: clusterName,
		Purpose:     purpose,
		CallbackURL: callbackURL,
		UserPoolID:  userPoolID,
	}
}

// AliasRecordSet contains the state required for
// building an alias record set
type AliasRecordSet struct {
	Name              string
	AliasDNS          string
	AliasHostedZoneID string
	Domain            string
	HostedZoneID      string
}

// Compose returns the cloud formation outputs and resources
func (s *AliasRecordSet) Compose() (*cfn.Composition, error) {
	composition := &cfn.Composition{}

	composition.Resources = append(composition.Resources, aliasrecordset.New(
		"Auth",
		s.AliasDNS,
		s.AliasHostedZoneID,
		s.Domain,
		s.HostedZoneID,
	))

	return composition, nil
}

// NewAliasRecordSet returns an initialised composer
func NewAliasRecordSet(name, aliasDNS, aliasHostedZoneID, domain, hostedZoneID string) *AliasRecordSet {
	return &AliasRecordSet{
		Name:              name,
		AliasDNS:          aliasDNS,
		AliasHostedZoneID: aliasHostedZoneID,
		Domain:            domain,
		HostedZoneID:      hostedZoneID,
	}
}

// AWSLoadBalancerControllerComposer contains state for building
// a managed iam policy compatible with aws-load-balancer-controller
type AWSLoadBalancerControllerComposer struct {
	ClusterName string
}

// NewAWSLoadBalancerControllerComposer returns an initialised aws load balancer controller composer
func NewAWSLoadBalancerControllerComposer(clusterName string) *AWSLoadBalancerControllerComposer {
	return &AWSLoadBalancerControllerComposer{
		ClusterName: clusterName,
	}
}

// Compose builds the policy and returns the result
func (a *AWSLoadBalancerControllerComposer) Compose() (*cfn.Composition, error) {
	p := a.ManagedPolicy()

	return &cfn.Composition{
		Outputs:   []cfn.StackOutputer{p},
		Resources: []cfn.ResourceNamer{p},
	}, nil
}

// ManagedPolicy creates a managed policy
// nolint: funlen
func (a *AWSLoadBalancerControllerComposer) ManagedPolicy() *managedpolicy.ManagedPolicy {
	policyName := fmt.Sprintf("okctl-%s-AWSLoadBalancerControllerServiceAccountPolicy", a.ClusterName)
	policyDesc := "Service account policy for creating AWS load balancers"

	d := &policydocument.PolicyDocument{
		Version: policydocument.Version,
		Statement: []policydocument.StatementEntry{
			{
				Effect: policydocument.EffectTypeAllow,
				Action: []string{
					"iam:CreateServiceLinkedRole",
					"ec2:DescribeAccountAttributes",
					"ec2:DescribeAddresses",
					"ec2:DescribeInternetGateways",
					"ec2:DescribeVpcs",
					"ec2:DescribeSubnets",
					"ec2:DescribeSecurityGroups",
					"ec2:DescribeInstances",
					"ec2:DescribeNetworkInterfaces",
					"ec2:DescribeTags",
					"elasticloadbalancing:DescribeLoadBalancers",
					"elasticloadbalancing:DescribeLoadBalancerAttributes",
					"elasticloadbalancing:DescribeListeners",
					"elasticloadbalancing:DescribeListenerCertificates",
					"elasticloadbalancing:DescribeSSLPolicies",
					"elasticloadbalancing:DescribeRules",
					"elasticloadbalancing:DescribeTargetGroups",
					"elasticloadbalancing:DescribeTargetGroupAttributes",
					"elasticloadbalancing:DescribeTargetHealth",
					"elasticloadbalancing:DescribeTags",
				},
				Resource: []string{
					"*",
				},
			},
			{
				Effect: policydocument.EffectTypeAllow,
				Action: []string{
					"cognito-idp:DescribeUserPoolClient",
					"acm:ListCertificates",
					"acm:DescribeCertificate",
					"iam:ListServerCertificates",
					"iam:GetServerCertificate",
					"waf-regional:GetWebACL",
					"waf-regional:GetWebACLForResource",
					"waf-regional:AssociateWebACL",
					"waf-regional:DisassociateWebACL",
					"wafv2:GetWebACL",
					"wafv2:GetWebACLForResource",
					"wafv2:AssociateWebACL",
					"wafv2:DisassociateWebACL",
					"shield:GetSubscriptionState",
					"shield:DescribeProtection",
					"shield:CreateProtection",
					"shield:DeleteProtection",
				},
				Resource: []string{
					"*",
				},
			},
			{
				Effect: policydocument.EffectTypeAllow,
				Action: []string{
					"ec2:AuthorizeSecurityGroupIngress",
					"ec2:RevokeSecurityGroupIngress",
				},
				Resource: []string{
					"*",
				},
			},
			{
				Effect: policydocument.EffectTypeAllow,
				Action: []string{
					"ec2:CreateSecurityGroup",
				},
				Resource: []string{
					"*",
				},
			},
			{
				Effect: policydocument.EffectTypeAllow,
				Action: []string{
					"ec2:CreateTags",
				},
				Resource: []string{
					"arn:aws:ec2:*:*:security-group/*",
				},
				Condition: map[policydocument.ConditionOperatorType]map[string]string{
					policydocument.ConditionOperatorTypeStringEquals: {
						"ec2:CreateAction": "CreateSecurityGroup",
					},
					policydocument.ConditionOperatorTypeNull: {
						"aws:RequestTag/elbv2.k8s.aws/cluster": "false",
					},
				},
			},
			{
				Effect: policydocument.EffectTypeAllow,
				Action: []string{
					"ec2:CreateTags",
					"ec2:DeleteTags",
				},
				Resource: []string{
					"arn:aws:ec2:*:*:security-group/*",
				},
				Condition: map[policydocument.ConditionOperatorType]map[string]string{
					policydocument.ConditionOperatorTypeNull: {
						"aws:RequestTag/elbv2.k8s.aws/cluster":  "true",
						"aws:ResourceTag/elbv2.k8s.aws/cluster": "false",
					},
				},
			},
			{
				Effect: policydocument.EffectTypeAllow,
				Action: []string{
					"ec2:AuthorizeSecurityGroupIngress",
					"ec2:RevokeSecurityGroupIngress",
					"ec2:DeleteSecurityGroup",
				},
				Resource: []string{
					"*",
				},
				Condition: map[policydocument.ConditionOperatorType]map[string]string{
					policydocument.ConditionOperatorTypeNull: {
						"aws:ResourceTag/elbv2.k8s.aws/cluster": "false",
					},
				},
			},
			{
				Effect: policydocument.EffectTypeAllow,
				Action: []string{
					"elasticloadbalancing:CreateLoadBalancer",
					"elasticloadbalancing:CreateTargetGroup",
				},
				Resource: []string{
					"*",
				},
				Condition: map[policydocument.ConditionOperatorType]map[string]string{
					policydocument.ConditionOperatorTypeNull: {
						"aws:RequestTag/elbv2.k8s.aws/cluster": "false",
					},
				},
			},
			{
				Effect: policydocument.EffectTypeAllow,
				Action: []string{
					"elasticloadbalancing:CreateListener",
					"elasticloadbalancing:DeleteListener",
					"elasticloadbalancing:CreateRule",
					"elasticloadbalancing:DeleteRule",
				},
				Resource: []string{
					"*",
				},
			},
			{
				Effect: policydocument.EffectTypeAllow,
				Action: []string{
					"elasticloadbalancing:AddTags",
					"elasticloadbalancing:RemoveTags",
				},
				Resource: []string{
					"arn:aws:elasticloadbalancing:*:*:targetgroup/*/*",
					"arn:aws:elasticloadbalancing:*:*:loadbalancer/net/*/*",
					"arn:aws:elasticloadbalancing:*:*:loadbalancer/app/*/*",
				},
				Condition: map[policydocument.ConditionOperatorType]map[string]string{
					policydocument.ConditionOperatorTypeNull: {
						"aws:RequestTag/elbv2.k8s.aws/cluster":  "true",
						"aws:ResourceTag/elbv2.k8s.aws/cluster": "false",
					},
				},
			},
			{
				Effect: policydocument.EffectTypeAllow,
				Action: []string{
					"elasticloadbalancing:ModifyLoadBalancerAttributes",
					"elasticloadbalancing:SetIpAddressType",
					"elasticloadbalancing:SetSecurityGroups",
					"elasticloadbalancing:SetSubnets",
					"elasticloadbalancing:DeleteLoadBalancer",
					"elasticloadbalancing:ModifyTargetGroup",
					"elasticloadbalancing:ModifyTargetGroupAttributes",
					"elasticloadbalancing:DeleteTargetGroup",
				},
				Resource: []string{
					"*",
				},
				Condition: map[policydocument.ConditionOperatorType]map[string]string{
					policydocument.ConditionOperatorTypeNull: {
						"aws:ResourceTag/elbv2.k8s.aws/cluster": "false",
					},
				},
			},
			{
				Effect: policydocument.EffectTypeAllow,
				Action: []string{
					"elasticloadbalancing:RegisterTargets",
					"elasticloadbalancing:DeregisterTargets",
				},
				Resource: []string{
					"arn:aws:elasticloadbalancing:*:*:targetgroup/*/*",
				},
			},
			{
				Effect: policydocument.EffectTypeAllow,
				Action: []string{
					"elasticloadbalancing:SetWebAcl",
					"elasticloadbalancing:ModifyListener",
					"elasticloadbalancing:AddListenerCertificates",
					"elasticloadbalancing:RemoveListenerCertificates",
					"elasticloadbalancing:ModifyRule",
				},
				Resource: []string{
					"*",
				},
			},
		},
	}

	return managedpolicy.New("AWSLoadBalancerControllerPolicy", policyName, policyDesc, d)
}

// AutoscalerPolicyComposer contains state for building
// a managed iam policy compatible with cluster autoscaler
type AutoscalerPolicyComposer struct {
	ClusterName string
}

// NewAutoscalerPolicyComposer returns an initialised cluster autoscaler composer
func NewAutoscalerPolicyComposer(clusterName string) *AutoscalerPolicyComposer {
	return &AutoscalerPolicyComposer{
		ClusterName: clusterName,
	}
}

// Compose builds the policy and returns the result
func (c *AutoscalerPolicyComposer) Compose() (*cfn.Composition, error) {
	p := c.ManagedPolicy()

	return &cfn.Composition{
		Outputs:   []cfn.StackOutputer{p},
		Resources: []cfn.ResourceNamer{p},
	}, nil
}

// ManagedPolicy returns the policy
func (c *AutoscalerPolicyComposer) ManagedPolicy() *managedpolicy.ManagedPolicy {
	policyName := fmt.Sprintf("okctl-%s-AutoscalerServiceAccountPolicy", c.ClusterName)
	policyDesc := "Service account policy for automatically scaling the cluster nodegroup"

	d := &policydocument.PolicyDocument{
		Version: policydocument.Version,
		Statement: []policydocument.StatementEntry{
			{
				Effect: policydocument.EffectTypeAllow,
				Action: []string{
					"autoscaling:DescribeAutoScalingGroups",
					"autoscaling:DescribeAutoScalingInstances",
					"autoscaling:DescribeLaunchConfigurations",
					"autoscaling:SetDesiredCapacity",
					"autoscaling:TerminateInstanceInAutoScalingGroup",
					"autoscaling:DescribeTags",
				},
				Resource: []string{
					"*",
				},
			},
			{
				Effect: policydocument.EffectTypeAllow,
				Action: []string{
					"ec2:DescribeLaunchTemplateVersions",
				},
				Resource: []string{
					"*",
				},
			},
		},
	}

	return managedpolicy.New("AutoscalerPolicy", policyName, policyDesc, d)
}

// BlockstoragePolicyComposer contains state for building
// a managed iam policy compatible with ebs csi blockstorage driver
type BlockstoragePolicyComposer struct {
	ClusterName string
}

// NewBlockstoragePolicyComposer returns an initialised ebs csi blockstorage driver composer
func NewBlockstoragePolicyComposer(clusterName string) *BlockstoragePolicyComposer {
	return &BlockstoragePolicyComposer{
		ClusterName: clusterName,
	}
}

// Compose builds the policy and returns the result
func (c *BlockstoragePolicyComposer) Compose() (*cfn.Composition, error) {
	p := c.ManagedPolicy()

	return &cfn.Composition{
		Outputs:   []cfn.StackOutputer{p},
		Resources: []cfn.ResourceNamer{p},
	}, nil
}

// ManagedPolicy returns the policy
func (c *BlockstoragePolicyComposer) ManagedPolicy() *managedpolicy.ManagedPolicy {
	policyName := fmt.Sprintf("okctl-%s-BlockstorageServiceAccountPolicy", c.ClusterName)
	policyDesc := "Service account policy for provisioning persistent volume claims"

	d := &policydocument.PolicyDocument{
		Version: policydocument.Version,
		Statement: []policydocument.StatementEntry{
			{
				Effect: policydocument.EffectTypeAllow,
				Action: []string{
					"ec2:AttachVolume",
					"ec2:CreateSnapshot",
					"ec2:CreateTags",
					"ec2:CreateVolume",
					"ec2:DeleteSnapshot",
					"ec2:DeleteTags",
					"ec2:DeleteVolume",
					"ec2:DescribeAvailabilityZones",
					"ec2:DescribeInstances",
					"ec2:DescribeSnapshots",
					"ec2:DescribeTags",
					"ec2:DescribeVolumes",
					"ec2:DescribeVolumesModifications",
					"ec2:DetachVolume",
					"ec2:ModifyVolume",
				},
				Resource: []string{
					"*",
				},
			},
		},
	}

	return managedpolicy.New("BlockstoragePolicy", policyName, policyDesc, d)
}

// CloudwatchDatasourcePolicyComposer contains state for building
// a managed iam policy that allows grafana to read cloudwatch metrics
// and logs
// - https://grafana.com/docs/grafana/latest/datasources/cloudwatch/#iam-policies
type CloudwatchDatasourcePolicyComposer struct {
	ClusterName string
}

// NewCloudwatchDatasourcePolicyComposer returns an initialised cloudwatch datasource policy composer
func NewCloudwatchDatasourcePolicyComposer(clusterName string) *CloudwatchDatasourcePolicyComposer {
	return &CloudwatchDatasourcePolicyComposer{
		ClusterName: clusterName,
	}
}

// Compose builds the policy and returns the result
func (c *CloudwatchDatasourcePolicyComposer) Compose() (*cfn.Composition, error) {
	p := c.ManagedPolicy()

	return &cfn.Composition{
		Outputs:   []cfn.StackOutputer{p},
		Resources: []cfn.ResourceNamer{p},
	}, nil
}

// ManagedPolicy returns the policy
func (c *CloudwatchDatasourcePolicyComposer) ManagedPolicy() *managedpolicy.ManagedPolicy {
	policyName := fmt.Sprintf("okctl-%s-CloudwatchDatasourceServiceAccountPolicy", c.ClusterName)
	policyDesc := "Service account policy for reading cloudwatch metrics and logs from grafana"

	d := &policydocument.PolicyDocument{
		Version: policydocument.Version,
		Statement: []policydocument.StatementEntry{
			{
				Effect: policydocument.EffectTypeAllow,
				Action: []string{
					"cloudwatch:DescribeAlarmsForMetric",
					"cloudwatch:DescribeAlarmHistory",
					"cloudwatch:DescribeAlarms",
					"cloudwatch:ListMetrics",
					"cloudwatch:GetMetricStatistics",
					"cloudwatch:GetMetricData",
				},
				Resource: []string{
					"*",
				},
			},
			{
				Effect: policydocument.EffectTypeAllow,
				Action: []string{
					"logs:DescribeLogGroups",
					"logs:GetLogGroupFields",
					"logs:StartQuery",
					"logs:StopQuery",
					"logs:GetQueryResults",
					"logs:GetLogEvents",
				},
				Resource: []string{
					"*",
				},
			},
			{
				Effect: policydocument.EffectTypeAllow,
				Action: []string{
					"ec2:DescribeTags",
					"ec2:DescribeInstances",
					"ec2:DescribeRegions",
				},
				Resource: []string{
					"*",
				},
			},
			{
				Effect: policydocument.EffectTypeAllow,
				Action: []string{
					"tag:GetResources",
				},
				Resource: []string{
					"*",
				},
			},
		},
	}

	return managedpolicy.New("CloudwatchDatasourcePolicy", policyName, policyDesc, d)
}

// FargateCloudwatchPolicyComposer contains state for building
// a managed iam policy that allows the fargate pod execution role
// to send logs to cloudwatch
// - https://docs.aws.amazon.com/eks/latest/userguide/fargate-logging.html
// - https://github.com/aws-samples/amazon-eks-fluent-logging-examples/blob/mainline/examples/fargate/cloudwatchlogs/permissions.json
type FargateCloudwatchPolicyComposer struct {
	ClusterName string
}

// NewFargateCloudwatchPolicyComposer returns an initialised cloudwatch datasource policy composer
func NewFargateCloudwatchPolicyComposer(clusterName string) *FargateCloudwatchPolicyComposer {
	return &FargateCloudwatchPolicyComposer{
		ClusterName: clusterName,
	}
}

// Compose builds the policy and returns the result
func (c *FargateCloudwatchPolicyComposer) Compose() (*cfn.Composition, error) {
	p := c.ManagedPolicy()

	return &cfn.Composition{
		Outputs:   []cfn.StackOutputer{p},
		Resources: []cfn.ResourceNamer{p},
	}, nil
}

// ManagedPolicy returns the policy
func (c *FargateCloudwatchPolicyComposer) ManagedPolicy() *managedpolicy.ManagedPolicy {
	policyName := fmt.Sprintf("okctl-%s-FargateCloudwatchPolicy", c.ClusterName)
	policyDesc := "Service account policy for reading cloudwatch metrics and logs from grafana"

	d := &policydocument.PolicyDocument{
		Version: policydocument.Version,
		Statement: []policydocument.StatementEntry{
			{
				Effect: policydocument.EffectTypeAllow,
				Action: []string{
					"logs:CreateLogStream",
					"logs:CreateLogGroup",
					"logs:DescribeLogStreams",
					"logs:PutLogEvents",
				},
				Resource: []string{
					"*",
				},
			},
		},
	}

	return managedpolicy.New("FargateCloudwatchPolicy", policyName, policyDesc, d)
}

// RDSPostgresComposerOpts contains the required
// inputs
type RDSPostgresComposerOpts struct {
	ApplicationDBName string
	AWSAccountID      string
	ClusterName       string
	DBSubnetGroupName string
	UserName          string
	VpcID             string
	VPCDBSubnetIDs    []string
	VPCDBSubnetCIDRs  []string
	Bucket            string
	Key               string
}

// RDSPostgresComposer contains state for building
// an RDS Postgres database
// - https://aws.amazon.com/rds/postgresql/
type RDSPostgresComposer struct {
	ApplicationDBName string
	AWSAccountID      string
	ClusterName       string
	DBSubnetGroupName string
	UserName          string
	VpcID             string
	VPCDBSubnetIDs    []string
	VPCDBSubnetCIDRs  []string
	Bucket            string
	Key               string
}

// NewRDSPostgresComposer returns an initialised RDS postgres composer
func NewRDSPostgresComposer(opts RDSPostgresComposerOpts) *RDSPostgresComposer {
	return &RDSPostgresComposer{
		ApplicationDBName: opts.ApplicationDBName,
		AWSAccountID:      opts.AWSAccountID,
		ClusterName:       opts.ClusterName,
		DBSubnetGroupName: opts.DBSubnetGroupName,
		UserName:          opts.UserName,
		VpcID:             opts.VpcID,
		VPCDBSubnetIDs:    opts.VPCDBSubnetIDs,
		VPCDBSubnetCIDRs:  opts.VPCDBSubnetCIDRs,
		Bucket:            opts.Bucket,
		Key:               opts.Key,
	}
}

const (
	// We use the policy document described here:
	// - https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/USER_Monitoring.OS.html
	amazonRDSEnhancedMonitoringRole = "arn:aws:iam::aws:policy/service-role/AmazonRDSEnhancedMonitoringRole"
	awsLambdaBasicExecutionRole     = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
)

// CloudFormationResourceName returns the CF logical name
func (c *RDSPostgresComposer) CloudFormationResourceName(resource string) string {
	return resource
}

// PhysicalName returns the name of the resource
func (c *RDSPostgresComposer) PhysicalName(resource string) string {
	return fmt.Sprintf("%s%s%s", c.ApplicationDBName, c.ClusterName, resource)
}

// AdminSecretFriendlyName returns the friendly name of the secret
func (c *RDSPostgresComposer) AdminSecretFriendlyName() string {
	return fmt.Sprintf("/%s/%s/postgres_admin", c.ApplicationDBName, c.ClusterName)
}

// Compose builds the policy and returns the result
// nolint: funlen
func (c *RDSPostgresComposer) Compose() (*cfn.Composition, error) {
	monitoringRole := role.New(
		c.PhysicalName("RDSPGMonitoringRole"),
		c.CloudFormationResourceName("RDSPGMonitoringRole"),
		v1alpha1.PermissionsBoundaryARN(c.AWSAccountID),
		[]string{amazonRDSEnhancedMonitoringRole},
		policydocument.PolicyDocument{
			Version: policydocument.Version,
			Statement: []policydocument.StatementEntry{
				{
					Effect: policydocument.EffectTypeAllow,
					Action: []string{
						"sts:AssumeRole",
					},
					Principal: &policydocument.Principal{
						Service: []string{
							"monitoring.rds.amazonaws.com",
						},
					},
				},
			},
		},
		nil,
	)

	params := map[string]string{
		"shared_preload_libraries":   "pg_stat_statements",
		"pg_stat_statements.max":     "10000",
		"pg_stat_statements.track":   "all",
		"log_min_duration_statement": "1000",
		"log_duration":               "on",
		"log_statement":              "all",
		"max_connections":            "100",
	}
	parameterGroup := dbparametergroup.New(
		c.CloudFormationResourceName("RDSPGParamGroup"),
		params,
	)

	admin := secret.NewRDSInstanceSecret(
		c.CloudFormationResourceName("RDSInstanceAdmin"),
		c.AdminSecretFriendlyName(),
		c.UserName,
	)

	lambdaSG := securitygroup.NewLambdaFunctionOutgoing(
		c.PhysicalName("RDSPGOutgoingSG"),
		c.CloudFormationResourceName("RDSPGOutgoingSG"),
		c.VpcID,
		c.VPCDBSubnetCIDRs,
	)

	secretsManagerIncoming := securitygroup.NewRDSPGSMVPCEndpointIncoming(
		c.PhysicalName("RDSPGSMVPCEndpointSG"),
		c.CloudFormationResourceName("RDSPGSMVPCEndpointSG"),
		c.VpcID,
		lambdaSG,
	)

	sme := vpcendpoint.NewSecretsManager(
		c.CloudFormationResourceName("RDSPGSMVPCEndpoint"),
		secretsManagerIncoming,
		c.VpcID,
		c.VPCDBSubnetIDs,
	)

	outgoing := securitygroup.NewPostgresOutgoing(
		c.PhysicalName("RDSPostgresOutgoing"),
		c.CloudFormationResourceName("RDSPostgresOutgoing"),
		c.VpcID,
		c.VPCDBSubnetCIDRs,
	)

	incoming := securitygroup.NewPostgresIncoming(
		c.PhysicalName("RDSPostgresIncoming"),
		c.CloudFormationResourceName("RDSPostgresIncoming"),
		c.VpcID,
		outgoing,
		lambdaSG,
	)

	postgres := dbinstance.New(
		c.CloudFormationResourceName("RDSPostgres"),
		fmt.Sprintf("%s-%s", c.ClusterName, c.ApplicationDBName),
		c.ApplicationDBName,
		c.DBSubnetGroupName,
		parameterGroup,
		monitoringRole,
		admin,
		incoming,
	)

	attachment := secrettargetattachment.NewRDSDBInstance(
		c.CloudFormationResourceName("SecretTargetAttachment"),
		admin,
		postgres,
	)

	// Based on the following content:
	// https://docs.aws.amazon.com/secretsmanager/latest/userguide/rotating-secrets-required-permissions.html
	lambdaRole := role.New(
		c.PhysicalName("RDSPGRotaterRole"),
		c.CloudFormationResourceName("RDSPGRotaterRole"),
		v1alpha1.PermissionsBoundaryARN(c.AWSAccountID),
		[]string{awsLambdaBasicExecutionRole},
		policydocument.PolicyDocument{
			Version: policydocument.Version,
			Statement: []policydocument.StatementEntry{
				{
					Effect: policydocument.EffectTypeAllow,
					Action: []string{
						"sts:AssumeRole",
					},
					Principal: &policydocument.Principal{
						Service: []string{"lambda.amazonaws.com"},
					},
				},
			},
		},
		map[string]interface{}{
			c.PhysicalName("RDSPGRotaterPolicy"): policydocument.PolicyDocument{
				Version: policydocument.Version,
				Statement: []policydocument.StatementEntry{
					{
						Effect: policydocument.EffectTypeAllow,
						Action: []string{
							"secretsmanager:GetRandomPassword",
						},
						Resource: []string{
							"*",
						},
					},
					{
						Effect: policydocument.EffectTypeAllow,
						Action: []string{
							"ec2:CreateNetworkInterface",
							"ec2:DeleteNetworkInterface",
							"ec2:DescribeNetworkInterfaces",
						},
						Resource: []string{
							"*",
						},
					},
				},
			},
		},
	)

	lambdaFunction := lambdafunction.NewRotateLambda(
		c.CloudFormationResourceName("RDSPGRotater"),
		c.Bucket,
		c.Key,
		lambdaRole,
		lambdaSG,
		c.VPCDBSubnetIDs,
		sme,
	)

	lambdaManagedPolicy := managedpolicy.New(
		c.CloudFormationResourceName("RDSPGRotaterPolicy"),
		c.PhysicalName("RDSPGRotaterPolicy"),
		c.CloudFormationResourceName("RDSPGRotaterPolicy"),
		policydocument.PolicyDocument{
			Version: policydocument.Version,
			Statement: []policydocument.StatementEntry{
				{
					Effect: policydocument.EffectTypeAllow,
					Action: []string{
						"secretsmanager:DescribeSecret",
						"secretsmanager:GetSecretValue",
						"secretsmanager:PutSecretValue",
						"secretsmanager:UpdateSecretVersionStage",
					},
					Resource: []string{
						"*",
					},
					Condition: map[policydocument.ConditionOperatorType]map[string]string{
						policydocument.ConditionOperatorTypeStringEquals: {
							"secretsmanager:resource/AllowRotationLambdaArn": cloudformation.GetAtt(c.CloudFormationResourceName("RDSPGRotater"), "Arn"),
						},
					},
				},
			},
		},
	)

	lambdaPermission := lambdapermission.NewRotateLambdaPermission(
		c.CloudFormationResourceName("RDSPGRotaterPermission"),
		lambdaFunction,
	)

	return &cfn.Composition{
		Outputs: []cfn.StackOutputer{
			postgres,
			admin,
			outgoing,
			lambdaManagedPolicy,
			lambdaRole,
			lambdaFunction,
		},
		Resources: []cfn.ResourceNamer{
			monitoringRole,
			parameterGroup,
			admin,
			outgoing,
			incoming,
			postgres,
			attachment,
			lambdaSG,
			lambdaRole,
			lambdaFunction,
			lambdaPermission,
			secretsManagerIncoming,
			sme,
			lambdaManagedPolicy,
		},
	}, nil
}

// S3BucketComposer contains the state required for creating
// the AWS S3 bucket
type S3BucketComposer struct {
	BucketName  string
	ClusterName string
}

// ResourceBucketNameOutput returns the name of the resource
func (s *S3BucketComposer) ResourceBucketNameOutput() string {
	return fmt.Sprintf("%s%sS3Bucket", s.BucketName, strings.ReplaceAll(s.ClusterName, "-", ""))
}

// Compose returns the outputs and resources
func (s *S3BucketComposer) Compose() (*cfn.Composition, error) {
	b := s3bucket.New(
		s.ResourceBucketNameOutput(),
		s.BucketName,
	)

	return &cfn.Composition{
		Outputs:   []cfn.StackOutputer{b},
		Resources: []cfn.ResourceNamer{b},
	}, nil
}

// NewS3BucketComposer returns an initialised AWS S3 bucket composer
func NewS3BucketComposer(bucketName, clusterName string) *S3BucketComposer {
	return &S3BucketComposer{
		BucketName:  bucketName,
		ClusterName: clusterName,
	}
}

// ECRRepositoryComposer contains state required for creating a ECR repository
type ECRRepositoryComposer struct {
	ImageName string
}

// ResourceRepositoryNameOutput returns the name of the resource
func (e *ECRRepositoryComposer) ResourceRepositoryNameOutput() string {
	return fmt.Sprintf("%sECRRepository", e.ImageName)
}

// Compose returns the outputs and resources
func (e *ECRRepositoryComposer) Compose() (*cfn.Composition, error) {
	repository := containerrepository.New(e.ImageName)

	return &cfn.Composition{
		Outputs:   []cfn.StackOutputer{repository},
		Resources: []cfn.ResourceNamer{repository},
	}, nil
}

// NewECRRepositoryComposer returns an initialized ECR composer
func NewECRRepositoryComposer(imageName string) *ECRRepositoryComposer {
	return &ECRRepositoryComposer{
		ImageName: imageName,
	}
}
