// Package mock provides mocks
package mock

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"

	"github.com/aws/aws-sdk-go/service/servicequotas"
	"github.com/aws/aws-sdk-go/service/servicequotas/servicequotasiface"

	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider/cognitoidentityprovideriface"

	"github.com/aws/aws-sdk-go/service/cloudfront/cloudfrontiface"

	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/aws/aws-sdk-go/service/route53/route53iface"

	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/aws/aws-sdk-go/service/eks/eksiface"
	"github.com/oslokommune/okctl/pkg/api/mock"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	awspkg "github.com/oslokommune/okctl/pkg/credentials/aws"
)

const (
	// DefaultAccessKeyID is a mocked aws access key id
	DefaultAccessKeyID = "ASIAV3ZUEFP6EXAMPLE"
	// DefaultSecretAccessKey is a mocked aws secret access key
	// nolint: gosec
	DefaultSecretAccessKey = "8P+SQvWIuLnKhh8d++jpw0nNmQRBZvNEXAMPLEKEY"
	// DefaultSessionToken is a mocked aws session token
	// nolint: gosec, lll
	DefaultSessionToken = "IQoJb3JpZ2luX2VjEOz////////////////////wEXAMPLEtMSJHMEUCIDoKK3JH9uG\nQE1z0sINr5M4jk+Na8KHDcCYRVjJCZEvOAiEA3OvJGtw1EcViOleS2vhs8VdCKFJQWP\nQrmGdeehM4IC1NtBmUpp2wUE8phUZampKsburEDy0KPkyQDYwT7WZ0wq5VSXDvp75YU\n9HFvlRd8Tx6q6fE8YQcHNVXAkiY9q6d+xo0rKwT38xVqr7ZD0u0iPPkUL64lIZbqBAz\n+scqKmlzm8FDrypNC9Yjc8fPOLn9FX9KSYvKTr4rvx3iSIlTJabIQwj2ICCR/oLxBA=="
	// DefaultExpiration is a mocked aws expiration
	DefaultExpiration = "2019-11-01T20:26:47Z"
	// DefaultPrincipalARN is a mocked principal arn
	DefaultPrincipalARN = "arn:aws:sts::123456789012:assumed-role/admin/user"
	// DefaultRegion is a mocked default region
	DefaultRegion = "eu-west-1"
	// DefaultStackName is a mocked default stack name
	DefaultStackName = "myStack"
	// DefaultCloudFrontDistributionARN is a mocked default
	DefaultCloudFrontDistributionARN = "arn:::::/distribution/FHH78FAKE"
	// DefaultFargateProfilePodExecutionRoleARN is the default name of the pod execution role
	DefaultFargateProfilePodExecutionRoleARN = "arn:aws:iam::123456789012:role/fargatePodExecutionRole-GHEFFAKE"
)

// DefaultCredentials returns a mocked set of aws credentials
func DefaultCredentials() *awspkg.Credentials {
	t, _ := time.Parse(time.RFC3339, DefaultExpiration)

	return &awspkg.Credentials{
		AccessKeyID:     DefaultAccessKeyID,
		SecretAccessKey: DefaultSecretAccessKey,
		SessionToken:    DefaultSessionToken,
		SecurityToken:   DefaultSessionToken,
		PrincipalARN:    DefaultPrincipalARN,
		Expires:         t.Local(),
		Region:          DefaultRegion,
	}
}

// DefaultStsCredentials returns a mocked set of aws sts credentials
func DefaultStsCredentials() *sts.Credentials {
	t, _ := time.Parse(time.RFC3339, DefaultExpiration)

	return &sts.Credentials{
		AccessKeyId:     aws.String(DefaultAccessKeyID),
		Expiration:      aws.Time(t.Local()),
		SecretAccessKey: aws.String(DefaultSecretAccessKey),
		SessionToken:    aws.String(DefaultSessionToken),
	}
}

// DefaultValidCredentials returns a mocked set of valid credentials
func DefaultValidCredentials() *awspkg.Credentials {
	creds := DefaultCredentials()

	creds.Expires = time.Now().Add(1 * time.Hour).Local()

	return creds
}

// DefaultValidStsCredentials returns a mocked set of valid aws sts credentials
func DefaultValidStsCredentials() *sts.Credentials {
	creds := DefaultStsCredentials()

	creds.Expiration = aws.Time(time.Now().Add(1 * time.Hour).Local())

	return creds
}

// SQAPI mocks the service quotas interface
type SQAPI struct {
	servicequotasiface.ServiceQuotasAPI

	GetServiceQuotaFn func(*servicequotas.GetServiceQuotaInput) (*servicequotas.GetServiceQuotaOutput, error)
}

// GetServiceQuota invokes the mocked response
func (a *SQAPI) GetServiceQuota(input *servicequotas.GetServiceQuotaInput) (*servicequotas.GetServiceQuotaOutput, error) {
	return a.GetServiceQuotaFn(input)
}

// STSAPI stores state for mocking out the STS API
type STSAPI struct {
	stsiface.STSAPI

	AssumeRoleWithSAMLFn func(*sts.AssumeRoleWithSAMLInput) (*sts.AssumeRoleWithSAMLOutput, error)
}

// AssumeRoleWithSAML invokes the mocked function
func (a *STSAPI) AssumeRoleWithSAML(in *sts.AssumeRoleWithSAMLInput) (*sts.AssumeRoleWithSAMLOutput, error) {
	return a.AssumeRoleWithSAMLFn(in)
}

// NewGoodSTSAPI returns a mocked sts api that will succeed
func NewGoodSTSAPI() stsiface.STSAPI {
	return &STSAPI{
		AssumeRoleWithSAMLFn: func(input *sts.AssumeRoleWithSAMLInput) (*sts.AssumeRoleWithSAMLOutput, error) {
			return &sts.AssumeRoleWithSAMLOutput{
				AssumedRoleUser: &sts.AssumedRoleUser{
					Arn: aws.String(DefaultPrincipalARN),
				},
				Credentials: DefaultStsCredentials(),
			}, nil
		},
	}
}

// EKSAPI mocks eksiface.EKSAPI
type EKSAPI struct {
	eksiface.EKSAPI

	DescribeClusterFn        func(*eks.DescribeClusterInput) (*eks.DescribeClusterOutput, error)
	DescribeFargateProfileFn func(*eks.DescribeFargateProfileInput) (*eks.DescribeFargateProfileOutput, error)
}

// DescribeCluster mocks describe cluster invocation
func (a *EKSAPI) DescribeCluster(input *eks.DescribeClusterInput) (*eks.DescribeClusterOutput, error) {
	return a.DescribeClusterFn(input)
}

// DescribeFargateProfile mocks the API invocation
func (a *EKSAPI) DescribeFargateProfile(input *eks.DescribeFargateProfileInput) (*eks.DescribeFargateProfileOutput, error) {
	return a.DescribeFargateProfileFn(input)
}

// EC2API provides a mocked structure of the ec2 API
type EC2API struct {
	ec2iface.EC2API

	DescribeSubnetsFn          func(*ec2.DescribeSubnetsInput) (*ec2.DescribeSubnetsOutput, error)
	DescribeAddressesFn        func(*ec2.DescribeAddressesInput) (*ec2.DescribeAddressesOutput, error)
	DescribeInternetGatewaysFn func(*ec2.DescribeInternetGatewaysInput) (*ec2.DescribeInternetGatewaysOutput, error)
	DescribeVpcsFn             func(*ec2.DescribeVpcsInput) (*ec2.DescribeVpcsOutput, error)
}

// DescribeSubnets invokes the mocked response
func (a *EC2API) DescribeSubnets(sub *ec2.DescribeSubnetsInput) (*ec2.DescribeSubnetsOutput, error) {
	return a.DescribeSubnetsFn(sub)
}

// DescribeAddresses invokes the mocked response
func (a *EC2API) DescribeAddresses(input *ec2.DescribeAddressesInput) (*ec2.DescribeAddressesOutput, error) {
	return a.DescribeAddressesFn(input)
}

// DescribeInternetGateways invokes the mocked response
func (a *EC2API) DescribeInternetGateways(input *ec2.DescribeInternetGatewaysInput) (*ec2.DescribeInternetGatewaysOutput, error) {
	return a.DescribeInternetGatewaysFn(input)
}

// DescribeVpcs invokes the mocked response
func (a *EC2API) DescribeVpcs(input *ec2.DescribeVpcsInput) (*ec2.DescribeVpcsOutput, error) {
	return a.DescribeVpcsFn(input)
}

// CFAPI provides a mocked structure of the cf API
type CFAPI struct {
	cloudformationiface.CloudFormationAPI

	DescribeStacksFn   []func(*cloudformation.DescribeStacksInput) (*cloudformation.DescribeStacksOutput, error)
	nextDescribeStacks int

	CreateStackFn         func(*cloudformation.CreateStackInput) (*cloudformation.CreateStackOutput, error)
	DeleteStackFn         func(*cloudformation.DeleteStackInput) (*cloudformation.DeleteStackOutput, error)
	DescribeStackEventsFn func(*cloudformation.DescribeStackEventsInput) (*cloudformation.DescribeStackEventsOutput, error)
}

// DescribeStackEvents invokes a mocked response
func (a *CFAPI) DescribeStackEvents(input *cloudformation.DescribeStackEventsInput) (*cloudformation.DescribeStackEventsOutput, error) {
	return a.DescribeStackEventsFn(input)
}

// DeleteStack invokes a mocked a response
func (a *CFAPI) DeleteStack(in *cloudformation.DeleteStackInput) (*cloudformation.DeleteStackOutput, error) {
	return a.DeleteStackFn(in)
}

// CreateStack invokes a mocked response
func (a *CFAPI) CreateStack(in *cloudformation.CreateStackInput) (*cloudformation.CreateStackOutput, error) {
	return a.CreateStackFn(in)
}

// DescribeStacks invokes a mocked response
func (a *CFAPI) DescribeStacks(stack *cloudformation.DescribeStacksInput) (*cloudformation.DescribeStacksOutput, error) {
	i := a.nextDescribeStacks
	a.nextDescribeStacks++

	return a.DescribeStacksFn[i](stack)
}

// R53API mocks the Route53 API
type R53API struct {
	route53iface.Route53API

	ListHostedZonesFn func(*route53.ListHostedZonesInput) (*route53.ListHostedZonesOutput, error)
}

// ListHostedZones invokes a mocked response
func (a *R53API) ListHostedZones(input *route53.ListHostedZonesInput) (*route53.ListHostedZonesOutput, error) {
	return a.ListHostedZonesFn(input)
}

// CFRONTAPI mocks the CloudFront API
type CFRONTAPI struct {
	cloudfrontiface.CloudFrontAPI

	GetDistributionFn func(*cloudfront.GetDistributionInput) (*cloudfront.GetDistributionOutput, error)
}

// GetDistribution invokes a mocked response
func (a *CFRONTAPI) GetDistribution(input *cloudfront.GetDistributionInput) (*cloudfront.GetDistributionOutput, error) {
	return a.GetDistributionFn(input)
}

// CIPAPI mocks the CognitoIdentityProvider API
type CIPAPI struct {
	cognitoidentityprovideriface.CognitoIdentityProviderAPI

	DescribeUserPoolDomainFn func(*cognitoidentityprovider.DescribeUserPoolDomainInput) (*cognitoidentityprovider.DescribeUserPoolDomainOutput, error)
}

// DescribeUserPoolDomain invokes a mocked response
func (a *CIPAPI) DescribeUserPoolDomain(input *cognitoidentityprovider.DescribeUserPoolDomainInput) (*cognitoidentityprovider.DescribeUserPoolDomainOutput, error) {
	return a.DescribeUserPoolDomainFn(input)
}

// IAMAPI mocks the IAM API
type IAMAPI struct {
	iamiface.IAMAPI

	AttachRolePolicyFn func(*iam.AttachRolePolicyInput) (*iam.AttachRolePolicyOutput, error)
}

// AttachRolePolicy mocks the invocation
func (a *IAMAPI) AttachRolePolicy(input *iam.AttachRolePolicyInput) (*iam.AttachRolePolicyOutput, error) {
	return a.AttachRolePolicyFn(input)
}

// DescribeStackEventsSuccess sets a success response on the describe event
func (p *CloudProvider) DescribeStackEventsSuccess() *CloudProvider {
	p.CFAPI.DescribeStackEventsFn = func(input *cloudformation.DescribeStackEventsInput) (*cloudformation.DescribeStackEventsOutput, error) {
		return &cloudformation.DescribeStackEventsOutput{
			StackEvents: []*cloudformation.StackEvent{
				{
					ResourceStatus:       aws.String(cloudformation.ResourceStatusCreateFailed),
					ResourceStatusReason: aws.String("something went wrong"),
					ResourceType:         aws.String("ec2"),
				},
			},
		}, nil
	}

	return p
}

// DeleteStackSuccess sets a success response on the mocked DeleteStack function
func (p *CloudProvider) DeleteStackSuccess() *CloudProvider {
	p.CFAPI.DeleteStackFn = func(input *cloudformation.DeleteStackInput) (*cloudformation.DeleteStackOutput, error) {
		return &cloudformation.DeleteStackOutput{}, nil
	}

	return p
}

// CreateStackSuccess sets a success response on the mocked CreateStack function
func (p *CloudProvider) CreateStackSuccess() *CloudProvider {
	p.CFAPI.CreateStackFn = func(input *cloudformation.CreateStackInput) (*cloudformation.CreateStackOutput, error) {
		return &cloudformation.CreateStackOutput{
			StackId: input.StackName,
		}, nil
	}

	return p
}

// DescribeStacksEmpty pushes an empty response onto the describe stacks list
func (p *CloudProvider) DescribeStacksEmpty() *CloudProvider {
	p.CFAPI.DescribeStacksFn = append(p.CFAPI.DescribeStacksFn, func(*cloudformation.DescribeStacksInput) (*cloudformation.DescribeStacksOutput, error) {
		return nil, awserr.New(
			"ValidationError",
			fmt.Sprintf("Stack with id %s does not exist", DefaultStackName),
			fmt.Errorf("something"),
		)
	})

	return p
}

// DescribeStacksResponse pushes a success response onto the describe stacks list
func (p *CloudProvider) DescribeStacksResponse(status string) *CloudProvider {
	p.CFAPI.DescribeStacksFn = append(p.CFAPI.DescribeStacksFn, func(*cloudformation.DescribeStacksInput) (*cloudformation.DescribeStacksOutput, error) {
		return &cloudformation.DescribeStacksOutput{
			Stacks: Stacks(status),
		}, nil
	})

	return p
}

// CloudProvider provides a structure for the mocked CloudProvider
type CloudProvider struct {
	v1alpha1.CloudProvider

	IAMAPI    *IAMAPI
	EC2API    *EC2API
	CFAPI     *CFAPI
	EKSAPI    *EKSAPI
	R53API    *R53API
	CFRONTAPI *CFRONTAPI
	CIPAPI    *CIPAPI
	SQAPI     *SQAPI
}

// IAM returns the mocked IAM API
func (p *CloudProvider) IAM() iamiface.IAMAPI {
	return p.IAMAPI
}

// ServiceQuotas returns the mocked SQ API
func (p *CloudProvider) ServiceQuotas() servicequotasiface.ServiceQuotasAPI {
	return p.SQAPI
}

// CloudFormation returns the mocked CF API
func (p *CloudProvider) CloudFormation() cloudformationiface.CloudFormationAPI {
	return p.CFAPI
}

// EC2 returns the mocked EC2 API
func (p *CloudProvider) EC2() ec2iface.EC2API {
	return p.EC2API
}

// EKS returns the mocked EKS API
func (p *CloudProvider) EKS() eksiface.EKSAPI {
	return p.EKSAPI
}

// Route53 returns the mocked Route53 API
func (p *CloudProvider) Route53() route53iface.Route53API {
	return p.R53API
}

// CognitoIdentityProvider returns the mocked CognitoIdentityProvider API
func (p *CloudProvider) CognitoIdentityProvider() cognitoidentityprovideriface.CognitoIdentityProviderAPI {
	return p.CIPAPI
}

// CloudFront returns the mocked CloudFront API
func (p *CloudProvider) CloudFront() cloudfrontiface.CloudFrontAPI {
	return p.CFRONTAPI
}

// PrincipalARN mocks the principal arn
func (p *CloudProvider) PrincipalARN() string {
	return "arn:::::/someuser"
}

// NewCloudProvider returns a mocked cloud provider with no mocks sets
func NewCloudProvider() *CloudProvider {
	return &CloudProvider{
		EC2API:    &EC2API{},
		CFAPI:     &CFAPI{},
		EKSAPI:    &EKSAPI{},
		R53API:    &R53API{},
		CFRONTAPI: &CFRONTAPI{},
		CIPAPI:    &CIPAPI{},
	}
}

// NewGoodCloudProvider returns a mocked cloud provider with success set on all
// nolint: funlen
func NewGoodCloudProvider() *CloudProvider {
	return &CloudProvider{
		EC2API: &EC2API{
			DescribeSubnetsFn: func(*ec2.DescribeSubnetsInput) (*ec2.DescribeSubnetsOutput, error) {
				return &ec2.DescribeSubnetsOutput{
					Subnets: Subnets(),
				}, nil
			},
			DescribeAddressesFn: func(*ec2.DescribeAddressesInput) (*ec2.DescribeAddressesOutput, error) {
				return &ec2.DescribeAddressesOutput{
					Addresses: []*ec2.Address{
						{},
					},
				}, nil
			},
			DescribeInternetGatewaysFn: func(*ec2.DescribeInternetGatewaysInput) (*ec2.DescribeInternetGatewaysOutput, error) {
				return &ec2.DescribeInternetGatewaysOutput{
					InternetGateways: []*ec2.InternetGateway{
						{},
					},
				}, nil
			},
			DescribeVpcsFn: func(*ec2.DescribeVpcsInput) (*ec2.DescribeVpcsOutput, error) {
				return &ec2.DescribeVpcsOutput{
					Vpcs: []*ec2.Vpc{
						{},
					},
				}, nil
			},
		},
		EKSAPI: &EKSAPI{
			DescribeClusterFn: func(*eks.DescribeClusterInput) (*eks.DescribeClusterOutput, error) {
				return &eks.DescribeClusterOutput{
					Cluster: &eks.Cluster{
						Arn: aws.String("arn:::something"),
						CertificateAuthority: &eks.Certificate{
							Data: aws.String(base64.StdEncoding.EncodeToString([]byte("something"))),
						},
						Endpoint: aws.String("https://something"),
						Name:     aws.String(mock.DefaultClusterName),
						Status:   aws.String(eks.ClusterStatusActive),
					},
				}, nil
			},
			DescribeFargateProfileFn: func(*eks.DescribeFargateProfileInput) (*eks.DescribeFargateProfileOutput, error) {
				return &eks.DescribeFargateProfileOutput{
					FargateProfile: &eks.FargateProfile{
						PodExecutionRoleArn: aws.String(DefaultFargateProfilePodExecutionRoleARN),
					},
				}, nil
			},
		},
		R53API: &R53API{
			ListHostedZonesFn: func(*route53.ListHostedZonesInput) (*route53.ListHostedZonesOutput, error) {
				return &route53.ListHostedZonesOutput{
					HostedZones: []*route53.HostedZone{
						{
							Config: &route53.HostedZoneConfig{
								PrivateZone: aws.Bool(false),
							},
							Id:   aws.String("/hostedzone/AABBCCDD"),
							Name: aws.String("test.oslo.systems."),
						},
					},
				}, nil
			},
		},
		CFRONTAPI: &CFRONTAPI{
			GetDistributionFn: func(*cloudfront.GetDistributionInput) (*cloudfront.GetDistributionOutput, error) {
				return &cloudfront.GetDistributionOutput{
					Distribution: &cloudfront.Distribution{
						ARN:        aws.String(DefaultCloudFrontDistributionARN),
						DomainName: aws.String("cloudfront-us-east-1.something.aws.com"),
					},
				}, nil
			},
		},
		CIPAPI: &CIPAPI{
			DescribeUserPoolDomainFn: func(*cognitoidentityprovider.DescribeUserPoolDomainInput) (*cognitoidentityprovider.DescribeUserPoolDomainOutput, error) {
				return &cognitoidentityprovider.DescribeUserPoolDomainOutput{
					DomainDescription: &cognitoidentityprovider.DomainDescriptionType{
						CloudFrontDistribution: aws.String("cloudfront-us-east-1.something.aws.com"),
						Domain:                 aws.String("auth.oslo.systems"),
					},
				}, nil
			},
		},
		SQAPI: &SQAPI{
			GetServiceQuotaFn: func(*servicequotas.GetServiceQuotaInput) (*servicequotas.GetServiceQuotaOutput, error) {
				return &servicequotas.GetServiceQuotaOutput{
					Quota: &servicequotas.ServiceQuota{
						Value: aws.Float64(3), // nolint: gomnd
					},
				}, nil
			},
		},
		IAMAPI: &IAMAPI{
			AttachRolePolicyFn: func(*iam.AttachRolePolicyInput) (*iam.AttachRolePolicyOutput, error) {
				return &iam.AttachRolePolicyOutput{}, nil
			},
		},
	}
}

var errBad = fmt.Errorf("something bad")

// NewBadCloudProvider returns a mocked cloud provider with failure set on all
func NewBadCloudProvider() *CloudProvider {
	return &CloudProvider{
		EC2API: &EC2API{
			DescribeSubnetsFn: func(*ec2.DescribeSubnetsInput) (*ec2.DescribeSubnetsOutput, error) {
				return nil, errBad
			},
		},
		EKSAPI: &EKSAPI{
			DescribeFargateProfileFn: func(*eks.DescribeFargateProfileInput) (*eks.DescribeFargateProfileOutput, error) {
				return nil, errBad
			},
		},
		IAMAPI: &IAMAPI{
			AttachRolePolicyFn: func(*iam.AttachRolePolicyInput) (*iam.AttachRolePolicyOutput, error) {
				return nil, errBad
			},
		},
	}
}

// Subnets returns a valid subnet describe response
func Subnets() []*ec2.Subnet {
	return []*ec2.Subnet{
		{
			AvailabilityZone: aws.String("eu-west-1a"),
			CidrBlock:        aws.String("192.168.0.0/24"),
			SubnetId:         aws.String("subnet-0bb1c79de3EXAMPLE"),
		},
	}
}

// Stacks returns a valid stack describe response
func Stacks(status string) []*cloudformation.Stack {
	return []*cloudformation.Stack{
		{
			StackId:           aws.String(DefaultStackName),
			StackName:         aws.String(DefaultStackName),
			StackStatus:       aws.String(status),
			StackStatusReason: aws.String("something"),
		},
	}
}
