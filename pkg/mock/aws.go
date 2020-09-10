// Package mock provides mocks
package mock

import (
	"encoding/base64"
	"fmt"
	"time"

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
	"github.com/oslokommune/okctl/pkg/api/okctl.io/v1alpha1"
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

	DescribeClusterFn func(*eks.DescribeClusterInput) (*eks.DescribeClusterOutput, error)
}

// DescribeCluster mocks describe cluster invocation
func (a *EKSAPI) DescribeCluster(input *eks.DescribeClusterInput) (*eks.DescribeClusterOutput, error) {
	return a.DescribeClusterFn(input)
}

// EC2API provides a mocked structure of the ec2 API
type EC2API struct {
	ec2iface.EC2API

	DescribeSubnetsFn func(*ec2.DescribeSubnetsInput) (*ec2.DescribeSubnetsOutput, error)
}

// DescribeSubnets invokes the mocked response
func (a *EC2API) DescribeSubnets(sub *ec2.DescribeSubnetsInput) (*ec2.DescribeSubnetsOutput, error) {
	return a.DescribeSubnetsFn(sub)
}

// CFAPI provides a mocked structure of the cf API
type CFAPI struct {
	cloudformationiface.CloudFormationAPI

	DescribeStacksFn   []func(*cloudformation.DescribeStacksInput) (*cloudformation.DescribeStacksOutput, error)
	nextDescribeStacks int

	CreateStackFn func(*cloudformation.CreateStackInput) (*cloudformation.CreateStackOutput, error)
	DeleteStackFn func(*cloudformation.DeleteStackInput) (*cloudformation.DeleteStackOutput, error)
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

	EC2API *EC2API
	CFAPI  *CFAPI
	EKSAPI *EKSAPI
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

// PrincipalARN mocks the principal arn
func (p *CloudProvider) PrincipalARN() string {
	return "arn:::::/someuser"
}

// NewCloudProvider returns a mocked cloud provider with no mocks sets
func NewCloudProvider() *CloudProvider {
	return &CloudProvider{
		EC2API: &EC2API{},
		CFAPI:  &CFAPI{},
		EKSAPI: &EKSAPI{},
	}
}

// NewGoodCloudProvider returns a mocked cloud provider with success set on all
func NewGoodCloudProvider() *CloudProvider {
	return &CloudProvider{
		EC2API: &EC2API{
			DescribeSubnetsFn: func(*ec2.DescribeSubnetsInput) (*ec2.DescribeSubnetsOutput, error) {
				return &ec2.DescribeSubnetsOutput{
					Subnets: Subnets(),
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
		},
	}
}

// NewBadCloudProvider returns a mocked cloud provider with failure set on all
func NewBadCloudProvider() *CloudProvider {
	return &CloudProvider{
		EC2API: &EC2API{
			DescribeSubnetsFn: func(*ec2.DescribeSubnetsInput) (*ec2.DescribeSubnetsOutput, error) {
				return nil, fmt.Errorf("something bad")
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
			StackId:     aws.String(DefaultStackName),
			StackName:   aws.String(DefaultStackName),
			StackStatus: aws.String(status),
		},
	}
}
