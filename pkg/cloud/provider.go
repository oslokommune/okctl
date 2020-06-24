// Package cloud provides access to AWS APIs
package cloud

import (
	"github.com/aws/aws-sdk-go/aws"
	awsCreds "github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/oslokommune/okctl/pkg/api/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/credentials"
)

// Provider stores state required for interacting with the AWS API
type Provider struct {
	Provider    v1alpha1.CloudProvider
	Credentials credentials.Provider
}

// New returns a new AWS API provider
func New(region string, c credentials.Provider) (*Provider, error) {
	services := &Services{
		region: region,
	}
	p := &Provider{
		Provider:    services,
		Credentials: c,
	}

	sess, err := p.newSession()
	if err != nil {
		return nil, err
	}

	services.cfn = cloudformation.New(sess)
	services.ec2 = ec2.New(sess)

	return p, nil
}

func (p *Provider) newSession() (*session.Session, error) {
	creds, err := p.Credentials.Raw()
	if err != nil {
		return nil, err
	}

	config := aws.NewConfig().
		WithRegion(p.Provider.Region()).
		WithCredentials(
			awsCreds.NewStaticCredentials(
				*creds.AccessKeyId,
				*creds.SecretAccessKey,
				*creds.SessionToken,
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
