package cloud

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/login"
)

type Provider struct {
	Provider v1alpha1.CloudProvider

	login login.Loginer
	creds *sts.Credentials
}

type Services struct {
	cfn cloudformationiface.CloudFormationAPI

	region string
}

func New(region string, l login.Loginer) (*Provider, error) {
	services := &Services{
		region: region,
	}
	p := &Provider{
		Provider: services,
		login:    l,
	}

	sess, err := p.newSession()
	if err != nil {
		return nil, err
	}

	services.cfn = cloudformation.New(sess)

	return p, nil
}

func (p *Services) CloudFormation() cloudformationiface.CloudFormationAPI {
	return p.cfn
}

func (p *Services) Region() string {
	return p.region
}

func (p *Provider) newSession() (*session.Session, error) {
	// Credentials have expired
	if p.creds != nil && time.Since(*p.creds.Expiration) < 0 {
		p.creds = nil
	}

	// No credentials available
	if p.creds == nil {
		creds, err := p.login.Login()
		if err != nil {
			return nil, err
		}

		p.creds = creds
	}

	config := aws.NewConfig().
		WithRegion(p.Provider.Region()).
		WithCredentials(
			credentials.NewStaticCredentials(
				*p.creds.AccessKeyId,
				*p.creds.SecretAccessKey,
				*p.creds.SessionToken,
			),
		)

	sess, err := session.NewSession(config)
	if err != nil {
		return nil, err
	}

	return sess, err
}
