package api

import "context"

// IdentityPool contains the state about the created
// identity management
type IdentityPool struct {
	ID                      ID
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
	ID           ID
	AuthDomain   string
	AuthFQDN     string
	HostedZoneID string
}

// IdentityPoolClient contains the state about a client
type IdentityPoolClient struct {
	ID                      ID
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
	ID          ID
	UserPoolID  string
	Purpose     string
	CallbackURL string
}

// CreateIdentityPoolUserOpts input
type CreateIdentityPoolUserOpts struct {
	ID         ID
	Email      string
	UserPoolID string
}

// DeleteIdentityPoolOpts input
type DeleteIdentityPoolOpts struct {
	ID         ID
	UserPoolID string
	Domain     string
}

// DeleteIdentityPoolClientOpts contains the required inputs
type DeleteIdentityPoolClientOpts struct {
	ID      ID
	Purpose string
}

// IdentityPoolUser state of user
type IdentityPoolUser struct {
	ID                     ID
	Email                  string
	UserPoolID             string
	StackName              string
	CloudFormationTemplate []byte
}

// IdentityManagerService implements the service layer
type IdentityManagerService interface {
	CreateIdentityPool(ctx context.Context, opts CreateIdentityPoolOpts) (*IdentityPool, error)
	CreateIdentityPoolClient(ctx context.Context, opts CreateIdentityPoolClientOpts) (*IdentityPoolClient, error)
	CreateIdentityPoolUser(ctx context.Context, opts CreateIdentityPoolUserOpts) (*IdentityPoolUser, error)
	DeleteIdentityPool(ctx context.Context, opts DeleteIdentityPoolOpts) error
	DeleteIdentityPoolClient(ctx context.Context, opts DeleteIdentityPoolClientOpts) error
}

// IdentityManagerCloudProvider implements the cloud layer
type IdentityManagerCloudProvider interface {
	CreateIdentityPool(certificateARN string, opts CreateIdentityPoolOpts) (*IdentityPool, error)
	CreateIdentityPoolClient(opts CreateIdentityPoolClientOpts) (*IdentityPoolClient, error)
	CreateIdentityPoolUser(opts CreateIdentityPoolUserOpts) (*IdentityPoolUser, error)
	DeleteIdentityPool(opts DeleteIdentityPoolOpts) error
	DeleteIdentityPoolClient(opts DeleteIdentityPoolClientOpts) error
}
