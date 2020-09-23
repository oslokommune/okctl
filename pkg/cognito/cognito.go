package cognito

import (
	"fmt"
	"path"

	"github.com/aws/aws-sdk-go/service/cloudfront"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/oslokommune/okctl/pkg/api/okctl.io/v1alpha1"
)

const (
	// DefaultCloudFrontHostedZoneID is the default hosted zone id for all cloud front distributions
	DefaultCloudFrontHostedZoneID = "Z2FDTNDATAQYW2"
	// DefaultCloudFrontACMRegion is the default region for ACM certificates used with a cloud front distribution
	DefaultCloudFrontACMRegion = "us-east-1"
)

// Cognito contains all required state for interacting
// with the Cognito API
type Cognito struct {
	provider v1alpha1.CloudProvider
}

// UserPoolDomainInfo contains the retrieved state about
// a cognito user pool domain
type UserPoolDomainInfo struct {
	CloudFrontDomainName string
	UserPoolDomain       string
}

// UserPoolDomainInfo returns information about the cognito user pool domain
func (c *Cognito) UserPoolDomainInfo(domain string) (*UserPoolDomainInfo, error) {
	pd, err := c.provider.CognitoIdentityProvider().DescribeUserPoolDomain(&cognitoidentityprovider.DescribeUserPoolDomainInput{
		Domain: aws.String(domain),
	})
	if err != nil {
		return nil, fmt.Errorf("describing user pool domain: %w", err)
	}

	dist, err := c.provider.CloudFront().GetDistribution(&cloudfront.GetDistributionInput{
		Id: aws.String(path.Base(*pd.DomainDescription.CloudFrontDistribution)),
	})
	if err != nil {
		return nil, err
	}

	return &UserPoolDomainInfo{
		UserPoolDomain:       *pd.DomainDescription.Domain,
		CloudFrontDomainName: *dist.Distribution.DomainName,
	}, nil
}

// New returns an initialised cognito interaction
func New(provider v1alpha1.CloudProvider) *Cognito {
	return &Cognito{
		provider: provider,
	}
}
