package cognito

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
)

// Cognito contains all required state for interacting
// with the Cognito API
type Cognito struct {
	provider   v1alpha1.CloudProvider
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

	return &UserPoolDomainInfo{
		UserPoolDomain:       *pd.DomainDescription.Domain,
		CloudFrontDomainName: *pd.DomainDescription.CloudFrontDistribution,
	}, nil
}

// UserPoolClientSecret returns the client secret for a user pool client
func (c *Cognito) UserPoolClientSecret(clientID, userPoolID string) (string, error) {
	out, err := c.provider.CognitoIdentityProvider().DescribeUserPoolClient(&cognitoidentityprovider.DescribeUserPoolClientInput{
		ClientId:   aws.String(clientID),
		UserPoolId: aws.String(userPoolID),
	})
	if err != nil {
		return "", fmt.Errorf("describing user pool client: %w", err)
	}

	return *out.UserPoolClient.ClientSecret, nil
}

// EnableMFA for a user pool
func (c *Cognito) EnableMFA(userPoolID string) error {
	_, err := c.provider.CognitoIdentityProvider().SetUserPoolMfaConfig(&cognitoidentityprovider.SetUserPoolMfaConfigInput{
		MfaConfiguration: aws.String("ON"),
		SoftwareTokenMfaConfiguration: &cognitoidentityprovider.SoftwareTokenMfaConfigType{
			Enabled: aws.Bool(true),
		},
		UserPoolId: aws.String(userPoolID),
	})
	if err != nil {
		return fmt.Errorf("enabling totp mfa: %w", err)
	}

	return nil
}

// DeleteAuthDomain dissociate auth domain with user pool, to allow pool deletion
func (c *Cognito) DeleteAuthDomain(userPoolID, domain string) error {
	_, err := c.provider.CognitoIdentityProvider().DeleteUserPoolDomain(&cognitoidentityprovider.DeleteUserPoolDomainInput{
		Domain:     &domain,
		UserPoolId: &userPoolID,
	})
	if err != nil {
		return err
	}

	return nil
}

// DeleteUserPool removes the userpool
func (c *Cognito) DeleteUserPool(userPoolID string) error {
	_, err := c.provider.CognitoIdentityProvider().DeleteUserPool(&cognitoidentityprovider.DeleteUserPoolInput{
		UserPoolId: &userPoolID,
	})
	if err != nil {
		return err
	}

	return nil
}

// New returns an initialised cognito interaction
func New(provider v1alpha1.CloudProvider) *Cognito {
	return &Cognito{
		provider: provider,
	}
}
