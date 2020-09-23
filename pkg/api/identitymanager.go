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
	Clients                 []*IdentityClient
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

// IdentityClient contains the state about a client
type IdentityClient struct {
	Purpose      string
	CallbackURL  string
	ClientID     string
	ClientSecret string
}

// CreateIdentityPoolOpts contains the required inputs
type CreateIdentityPoolOpts struct {
	ID           ID
	AuthDomain   string
	AuthFQDN     string
	HostedZoneID string
	Clients      []IdentityPoolClientOpts
}

// IdentityPoolClientOpts contains the required inputs
type IdentityPoolClientOpts struct {
	Purpose     string
	CallbackURL string
}

// IdentityManagerService implements the service layer
type IdentityManagerService interface {
	CreateIdentityPool(ctx context.Context, opts CreateIdentityPoolOpts) (*IdentityPool, error)
}

// IdentityManagerCloudProvider implements the cloud layer
type IdentityManagerCloudProvider interface {
	CreateIdentityPool(certificateARN string, opts CreateIdentityPoolOpts) (*IdentityPool, error)
}
