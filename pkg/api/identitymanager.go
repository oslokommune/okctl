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

// IdentityManagerService implements the service layer
type IdentityManagerService interface {
	CreateIdentityPool(ctx context.Context, opts CreateIdentityPoolOpts) (*IdentityPool, error)
	CreateIdentityPoolClient(ctx context.Context, opts CreateIdentityPoolClientOpts) (*IdentityPoolClient, error)
}

// IdentityManagerCloudProvider implements the cloud layer
type IdentityManagerCloudProvider interface {
	CreateIdentityPool(certificateARN string, opts CreateIdentityPoolOpts) (*IdentityPool, error)
	CreateIdentityPoolClient(opts CreateIdentityPoolClientOpts) (*IdentityPoolClient, error)
}
