package credentials

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/oslokommune/okctl/pkg/credentials/login"
)

type Provider interface {
	AsEnv() ([]string, error)
	Raw() (*sts.Credentials, error)
}

type ErrorProvider struct {
	Err error
}

func (p *ErrorProvider) AsEnv() ([]string, error) {
	return nil, p.Err
}

func (p *ErrorProvider) Raw() (*sts.Credentials, error) {
	return nil, p.Err
}

func NewErrorProvider() *ErrorProvider {
	return &ErrorProvider{
		fmt.Errorf("this is an error provider, cli is probably misconfigured"),
	}
}

type AWSProvider struct {
	Login login.Loginer

	creds *sts.Credentials
}

func New(l login.Loginer) *AWSProvider {
	return &AWSProvider{
		Login: l,
	}
}

func (p *AWSProvider) AsEnv() ([]string, error) {
	creds, err := p.Raw()
	if err != nil {
		return nil, err
	}

	return []string{
		fmt.Sprintf("AWS_ACCESS_KEY_ID=%s", *creds.AccessKeyId),
		fmt.Sprintf("AWS_SECRET_ACCESS_KEY=%s", *creds.SecretAccessKey),
		fmt.Sprintf("AWS_SESSION_TOKEN=%s", *creds.SessionToken),
	}, nil
}

func (p *AWSProvider) Raw() (*sts.Credentials, error) {
	// Credentials have expired
	if p.creds != nil && time.Since(*p.creds.Expiration) < 0 {
		p.creds = nil
	}

	// No credentials available
	if p.creds == nil {
		creds, err := p.Login.Login()
		if err != nil {
			return nil, err
		}

		p.creds = creds
	}

	return p.creds, nil
}
