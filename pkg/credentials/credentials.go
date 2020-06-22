// Package credentials knows how to keep credentials up to date and
// make the available in various formats
package credentials

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/oslokommune/okctl/pkg/credentials/login"
)

// Provider defines how credentials should
// be made available
type Provider interface {
	AsEnv() ([]string, error)
	Raw() (*sts.Credentials, error)
}

// ErrorProvider stores an error
type ErrorProvider struct {
	Err error
}

// AsEnv returns the stored error
func (p *ErrorProvider) AsEnv() ([]string, error) {
	return nil, p.Err
}

// Raw returns the stored error
func (p *ErrorProvider) Raw() (*sts.Credentials, error) {
	return nil, p.Err
}

// NewErrorProvider returns a provider that simply errors
func NewErrorProvider() *ErrorProvider {
	return &ErrorProvider{
		fmt.Errorf("this is an error provider, cli is probably misconfigured"),
	}
}

// AWSProvider stores state required for fetching
// AWS credentials
type AWSProvider struct {
	Login login.Loginer

	creds *sts.Credentials
}

// New returns an AWS credentials provider
func New(l login.Loginer) *AWSProvider {
	return &AWSProvider{
		Login: l,
	}
}

// AsEnv returns the AWS credentials as env vars
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

// Raw returns the raw credentials
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
