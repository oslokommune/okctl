// Package cloud provides access to AWS APIs
package cloud

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	awsCreds "github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/oslokommune/okctl/pkg/api/okctl.io/v1alpha1"
	awsauth "github.com/oslokommune/okctl/pkg/credentials/aws"
)

// Provider stores state required for interacting with the AWS API
type Provider struct {
	Provider v1alpha1.CloudProvider
}

// New returns a new AWS API provider and builds a session from
// the provided authenticator
func New(region string, a awsauth.Authenticator) (*Provider, error) {
	sess, err := NewSession(region, a)
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate with aws: %w", err)
	}

	return NewFromSession(region, sess)
}

// NewFromSession returns a new AWS API provider and builds a session
// from the provided authenticator
func NewFromSession(region string, sess *session.Session) (*Provider, error) {
	services := &Services{
		region: region,
	}
	p := &Provider{
		Provider: services,
	}

	services.cfn = cloudformation.New(sess)
	services.ec2 = ec2.New(sess)

	return p, nil
}

// NewSession returns an AWS session using the provided authenticator
func NewSession(region string, auth awsauth.Authenticator) (*session.Session, error) {
	creds, err := auth.Raw()
	if err != nil {
		return nil, err
	}

	config := aws.NewConfig().
		WithRegion(region).
		WithCredentials(
			awsCreds.NewStaticCredentials(
				creds.AccessKeyID,
				creds.SecretAccessKey,
				creds.SessionToken,
			),
		)

	sess, err := session.NewSession(config)
	if err != nil {
		return nil, err
	}

	return sess, err
}

// Services stores access to the various AWS APIs
type Services struct {
	cfn cloudformationiface.CloudFormationAPI
	ec2 ec2iface.EC2API

	region string
}

// EC2 returns an interface to the AWS EC2 API
func (s *Services) EC2() ec2iface.EC2API {
	return s.ec2
}

// CloudFormation returns an interface to the AWS CloudFormation API
func (s *Services) CloudFormation() cloudformationiface.CloudFormationAPI {
	return s.cfn
}

// Region returns the configured AWS region
func (s *Services) Region() string {
	return s.region
}
