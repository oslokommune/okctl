// Package userpoolclient provides an implementation of Cognito UserPoolClient
// - https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-cognito-userpoolclient.html
package userpoolclient

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/cfn"

	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/awslabs/goformation/v4/cloudformation/cognito"
)

const (
	refreshTokenValidityDays = 30
)

// UserPoolClient stores the state for a cloud formation
// cognito user pool client
type UserPoolClient struct {
	StoredName  string
	ClientName  string
	CallBackURL string
	ClusterName string
	Purpose     string
	UserPoolID  string
}

// NamedOutputs returns the named outputs
func (c *UserPoolClient) NamedOutputs() map[string]cloudformation.Output {
	return cfn.NewValue(fmt.Sprintf("%sClientID", c.Purpose), c.Ref()).NamedOutputs()
}

// Resource returns the cloud formation resource for a
// cognito user pool client
func (c *UserPoolClient) Resource() cloudformation.Resource {
	return &cognito.UserPoolClient{
		AllowedOAuthFlows: []string{
			// https://auth0.com/docs/flows/authorization-code-flow
			"code",
		},
		AllowedOAuthFlowsUserPoolClient: true,
		AllowedOAuthScopes: []string{
			"email",
			"profile",
			"openid",
		},
		CallbackURLs: []string{
			c.CallBackURL,
		},
		ClientName: c.ClientName,
		// https://tools.ietf.org/html/rfc6749#section-3.1.2
		DefaultRedirectURI:         c.CallBackURL,
		ExplicitAuthFlows:          nil, // ?
		GenerateSecret:             true,
		PreventUserExistenceErrors: "ENABLED",
		RefreshTokenValidity:       refreshTokenValidityDays,
		SupportedIdentityProviders: []string{
			"COGNITO",
		},
		UserPoolId: c.UserPoolID,
	}
}

// Ref returns a cloud formation intrinsic ref to the resource
func (c *UserPoolClient) Ref() string {
	return cloudformation.Ref(c.Name())
}

// Name returns the logical identifier of the resource
func (c *UserPoolClient) Name() string {
	return c.StoredName
}

// New returns an initialised cloud formation creator for
// a cognito user pool client
func New(purpose, clusterName, callBackURL, userPoolID string) *UserPoolClient {
	return &UserPoolClient{
		StoredName:  fmt.Sprintf("UserPoolClient%s", purpose),
		CallBackURL: callBackURL,
		ClusterName: clusterName,
		ClientName:  fmt.Sprintf("okctl-%s-%s", clusterName, purpose),
		UserPoolID:  userPoolID,
		Purpose:     purpose,
	}
}
