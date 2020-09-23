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
	Environment string
	Repository  string
	Purpose     string
	UserPool    cfn.NameReferencer
}

// NamedOutputs returns the named outputs
func (c *UserPoolClient) NamedOutputs() map[string]map[string]interface{} {
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
		ReadAttributes: []string{
			"email",
			"openid",
		},
		RefreshTokenValidity:       refreshTokenValidityDays,
		SupportedIdentityProviders: nil,
		UserPoolId:                 c.UserPool.Ref(),
		AWSCloudFormationDependsOn: []string{
			c.UserPool.Name(),
		},
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
func New(purpose, environment, repository, callBackURL string, userPool cfn.NameReferencer) *UserPoolClient {
	return &UserPoolClient{
		StoredName:  fmt.Sprintf("UserPoolClient%s", purpose),
		CallBackURL: callBackURL,
		Environment: environment,
		Repository:  repository,
		ClientName:  fmt.Sprintf("okctl-%s-%s-%s", environment, repository, purpose),
		UserPool:    userPool,
		Purpose:     purpose,
	}
}
