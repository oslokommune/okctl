// Package credentials knows how to keep credentials up to date and
// make the available in various formats
package credentials

import (
	"github.com/oslokommune/okctl/pkg/credentials/aws"
)

// Provider defines how credentials should
// be made available
type Provider interface {
	Aws() aws.Authenticator
}

// provider stores state required for fetching
// AWS credentials
type provider struct {
	aws aws.Authenticator
}

// Aws returns an AWS credentials provider
func (p *provider) Aws() aws.Authenticator {
	return p.aws
}

// New returns a credentials provider
func New(aws aws.Authenticator) Provider {
	return &provider{
		aws: aws,
	}
}
