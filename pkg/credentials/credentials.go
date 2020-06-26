// Package credentials knows how to keep credentials up to date and
// make the available in various formats
package credentials

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/service/sts"
)

// AwsAuthenticator authenticates with AWS
type AwsAuthenticator interface {
	Get() (*sts.Credentials, error)
}

// Provider defines how credentials should
// be made available
type Provider interface {
	AwsEnv() ([]string, error)
	AwsRaw() (*sts.Credentials, error)
}

// provider stores state required for fetching
// AWS credentials
type provider struct {
	aws AwsAuthenticator

	creds *sts.Credentials
}

// New returns an AWS credentials provider
func New(aws AwsAuthenticator) Provider {
	return &provider{
		aws: aws,
	}
}

// AwsEnv returns the AWS credentials as env vars
func (p *provider) AwsEnv() ([]string, error) {
	creds, err := p.AwsRaw()
	if err != nil {
		return nil, err
	}

	return []string{
		fmt.Sprintf("AWS_ACCESS_KEY_ID=%s", *creds.AccessKeyId),
		fmt.Sprintf("AWS_SECRET_ACCESS_KEY=%s", *creds.SecretAccessKey),
		fmt.Sprintf("AWS_SESSION_TOKEN=%s", *creds.SessionToken),
	}, nil
}

// AwsRaw returns the raw credentials
func (p *provider) AwsRaw() (*sts.Credentials, error) {
	// Credentials have expired
	if p.creds != nil && time.Since(*p.creds.Expiration) < 0 {
		p.creds = nil
	}

	// No credentials available
	if p.creds == nil {
		creds, err := p.aws.Get()
		if err != nil {
			return nil, err
		}

		p.creds = creds
	}

	return p.creds, nil
}
