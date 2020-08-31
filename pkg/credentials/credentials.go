// Package credentials knows how to keep credentials up to date and
// make the available in various formats
package credentials

import (
	"github.com/oslokommune/okctl/pkg/credentials/aws"
	"github.com/oslokommune/okctl/pkg/credentials/github"
)

// Provider defines how credentials should
// be made available
type Provider interface {
	Aws() aws.Authenticator
	Github() github.Authenticator
}

// provider stores state required for fetching
// AWS credentials
type provider struct {
	aws    aws.Authenticator
	github github.Authenticator
}

// Aws returns an AWS credentials provider
func (p *provider) Aws() aws.Authenticator {
	return p.aws
}

// Github returns a github credentials provider
func (p *provider) Github() github.Authenticator {
	return p.github
}

// New returns a credentials provider
func New(aws aws.Authenticator, github github.Authenticator) Provider {
	return &provider{
		aws:    aws,
		github: github,
	}
}
