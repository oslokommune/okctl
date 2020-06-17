package cloud

import (
	"github.com/aws/aws-sdk-go/aws"
	awsCreds "github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/credentials"
)

type Provider struct {
	Provider    v1alpha1.CloudProvider
	Credentials credentials.Provider
}

type Services struct {
	cfn cloudformationiface.CloudFormationAPI

	region string
}

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

	return p, nil
}

func (p *Services) CloudFormation() cloudformationiface.CloudFormationAPI {
	return p.cfn
}

func (p *Services) Region() string {
	return p.region
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
