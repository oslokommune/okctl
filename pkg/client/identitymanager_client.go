package client

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"

	"github.com/oslokommune/okctl/pkg/api"
)

// IdentityPool contains the state about the created
// identity management
type IdentityPool struct {
	ID                      api.ID
	UserPoolID              string
	AuthDomain              string
	HostedZoneID            string
	StackName               string
	CloudFormationTemplates []byte
	Certificate             *Certificate
	RecordSetAlias          *RecordSetAlias
}

// RecordSetAlias contains a record set alias
// this should not be here
type RecordSetAlias struct {
	AliasDomain            string
	AliasHostedZones       string
	StackName              string
	CloudFormationTemplate []byte
}

// CreateIdentityPoolOpts contains the required inputs
type CreateIdentityPoolOpts struct {
	ID           api.ID
	AuthDomain   string
	AuthFQDN     string
	HostedZoneID string
}

// IdentityPoolClient contains the state about a client
type IdentityPoolClient struct {
	ID                      api.ID
	UserPoolID              string
	Purpose                 string
	CallbackURL             string
	ClientID                string
	ClientSecret            string
	StackName               string
	CloudFormationTemplates []byte
}

// CreateIdentityPoolClientOpts contains the required inputs
type CreateIdentityPoolClientOpts struct {
	ID          api.ID
	UserPoolID  string
	Purpose     string
	CallbackURL string
}

// CreateIdentityPoolUserOpts input
type CreateIdentityPoolUserOpts struct {
	ID         api.ID
	Email      string
	UserPoolID string
}

// DeleteIdentityPoolOpts input
type DeleteIdentityPoolOpts struct {
	ID         api.ID
	UserPoolID string
	Domain     string
}

// DeleteIdentityPoolClientOpts contains the required inputs
type DeleteIdentityPoolClientOpts struct {
	ID      api.ID
	Purpose string
}

// IdentityPoolUser state of user
type IdentityPoolUser struct {
	ID                     api.ID
	Email                  string
	UserPoolID             string
	StackName              string
	CloudFormationTemplate []byte
}

// Validate the inputs
func (o CreateIdentityPoolUserOpts) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.UserPoolID, validation.Required),
		validation.Field(&o.Email, validation.Required, is.EmailFormat),
	)
}

// IdentityManagerService orchestrates the creation of an identity pool
type IdentityManagerService interface {
	CreateIdentityPool(ctx context.Context, opts CreateIdentityPoolOpts) (*IdentityPool, error)
	DeleteIdentityPool(ctx context.Context, opts api.ID) error
	CreateIdentityPoolClient(ctx context.Context, opts CreateIdentityPoolClientOpts) (*IdentityPoolClient, error)
	DeleteIdentityPoolClient(ctx context.Context, opts DeleteIdentityPoolClientOpts) error
	CreateIdentityPoolUser(ctx context.Context, opts CreateIdentityPoolUserOpts) (*IdentityPoolUser, error)
}

// IdentityManagerAPI invokes the API calls for creating an identity pool
type IdentityManagerAPI interface {
	CreateIdentityPool(opts api.CreateIdentityPoolOpts) (*api.IdentityPool, error)
	DeleteIdentityPool(opts api.DeleteIdentityPoolOpts) error
	CreateIdentityPoolClient(opts api.CreateIdentityPoolClientOpts) (*api.IdentityPoolClient, error)
	DeleteIdentityPoolClient(opts api.DeleteIdentityPoolClientOpts) error
	CreateIdentityPoolUser(opts api.CreateIdentityPoolUserOpts) (*api.IdentityPoolUser, error)
}

// IdentityManagerState implements the state layer
type IdentityManagerState interface {
	SaveIdentityPool(pool *IdentityPool) error
	RemoveIdentityPool(stackName string) error
	GetIdentityPool(stackName string) (*IdentityPool, error)
	SaveIdentityPoolClient(client *IdentityPoolClient) error
	RemoveIdentityPoolClient(stackName string) error
	SaveIdentityPoolUser(user *IdentityPoolUser) error
}
