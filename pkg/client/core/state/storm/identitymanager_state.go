package storm

import (
	"errors"

	stormpkg "github.com/asdine/storm/v3"
	"github.com/oslokommune/okctl/pkg/client"
)

type identityManagerState struct {
	node stormpkg.Node
}

// IdentityPool contains storm compatible state
type IdentityPool struct {
	Metadata `storm:"inline"`

	ID                      ID
	UserPoolID              string
	AuthDomain              string
	HostedZoneID            string
	StackName               string `storm:"unique"`
	CloudFormationTemplates []byte
	Certificate             *Certificate
	RecordSetAlias          *RecordSetAlias
}

// NewIdentityPool returns a storm compatible IdentityPool
func NewIdentityPool(p *client.IdentityPool, meta, certMeta Metadata) *IdentityPool {
	return &IdentityPool{
		Metadata:                meta,
		ID:                      NewID(p.ID),
		UserPoolID:              p.UserPoolID,
		AuthDomain:              p.AuthDomain,
		HostedZoneID:            p.HostedZoneID,
		StackName:               p.StackName,
		CloudFormationTemplates: p.CloudFormationTemplates,
		Certificate:             NewCertificate(p.Certificate, certMeta),
		RecordSetAlias:          NewRecordSetAlias(p.RecordSetAlias),
	}
}

// Convert to client.IdentityPool
func (p *IdentityPool) Convert() *client.IdentityPool {
	return &client.IdentityPool{
		ID:                      p.ID.Convert(),
		UserPoolID:              p.UserPoolID,
		AuthDomain:              p.AuthDomain,
		HostedZoneID:            p.HostedZoneID,
		StackName:               p.StackName,
		CloudFormationTemplates: p.CloudFormationTemplates,
		Certificate:             p.Certificate.Convert(),
		RecordSetAlias:          p.RecordSetAlias.Convert(),
	}
}

// RecordSetAlias contains storm compatible state
type RecordSetAlias struct {
	AliasDomain            string
	AliasHostedZones       string
	StackName              string
	CloudFormationTemplate []byte
}

// NewRecordSetAlias returns a storm compatible RecordSetAlias
func NewRecordSetAlias(a *client.RecordSetAlias) *RecordSetAlias {
	return &RecordSetAlias{
		AliasDomain:            a.AliasDomain,
		AliasHostedZones:       a.AliasHostedZones,
		StackName:              a.StackName,
		CloudFormationTemplate: a.CloudFormationTemplate,
	}
}

// Convert to client.RecordSetAlias
func (a *RecordSetAlias) Convert() *client.RecordSetAlias {
	return &client.RecordSetAlias{
		AliasDomain:            a.AliasDomain,
		AliasHostedZones:       a.AliasHostedZones,
		StackName:              a.StackName,
		CloudFormationTemplate: a.CloudFormationTemplate,
	}
}

// IdentityPoolClient contains storm compatible state
type IdentityPoolClient struct {
	Metadata `storm:"inline"`

	ID                      ID
	UserPoolID              string
	Purpose                 string
	CallbackURL             string
	ClientID                string
	ClientSecret            string
	StackName               string `storm:"unique,index"`
	CloudFormationTemplates []byte
}

// NewIdentityPoolClient returns a storm compatible IdentityPoolClient
func NewIdentityPoolClient(c *client.IdentityPoolClient, meta Metadata) *IdentityPoolClient {
	return &IdentityPoolClient{
		Metadata:                meta,
		ID:                      NewID(c.ID),
		UserPoolID:              c.UserPoolID,
		Purpose:                 c.Purpose,
		CallbackURL:             c.CallbackURL,
		ClientID:                c.ClientID,
		ClientSecret:            c.ClientSecret,
		StackName:               c.StackName,
		CloudFormationTemplates: c.CloudFormationTemplates,
	}
}

// Convert to IdentityPoolClient
func (c *IdentityPoolClient) Convert() *client.IdentityPoolClient {
	return &client.IdentityPoolClient{
		ID:                      c.ID.Convert(),
		UserPoolID:              c.UserPoolID,
		Purpose:                 c.Purpose,
		CallbackURL:             c.CallbackURL,
		ClientID:                c.ClientID,
		ClientSecret:            c.ClientSecret,
		StackName:               c.StackName,
		CloudFormationTemplates: c.CloudFormationTemplates,
	}
}

// IdentityPoolUser contains storm compatible state
type IdentityPoolUser struct {
	Metadata `storm:"inline"`

	ID                     ID
	Email                  string
	UserPoolID             string
	StackName              string `storm:"unique,index"`
	CloudFormationTemplate []byte
}

// NewIdentityPoolUser returns a storm compatible IdentityPoolUser
func NewIdentityPoolUser(u *client.IdentityPoolUser, meta Metadata) *IdentityPoolUser {
	return &IdentityPoolUser{
		Metadata:               meta,
		ID:                     NewID(u.ID),
		Email:                  u.Email,
		UserPoolID:             u.UserPoolID,
		StackName:              u.StackName,
		CloudFormationTemplate: u.CloudFormationTemplate,
	}
}

// Convert to client.IdentityPoolUser
func (u *IdentityPoolUser) Convert() *client.IdentityPoolUser {
	return &client.IdentityPoolUser{
		ID:                     u.ID.Convert(),
		Email:                  u.Email,
		UserPoolID:             u.UserPoolID,
		StackName:              u.StackName,
		CloudFormationTemplate: u.CloudFormationTemplate,
	}
}

func (s *identityManagerState) SaveIdentityPool(pool *client.IdentityPool) error {
	return s.node.Save(NewIdentityPool(pool, NewMetadata(), NewMetadata()))
}

func (s *identityManagerState) RemoveIdentityPool(stackName string) error {
	p := &IdentityPool{}

	err := s.node.One("StackName", stackName, p)
	if err != nil {
		if errors.Is(err, stormpkg.ErrNotFound) {
			return nil
		}

		return err
	}

	return s.node.DeleteStruct(p)
}

func (s *identityManagerState) GetIdentityPool(stackName string) (*client.IdentityPool, error) {
	p := &IdentityPool{}

	err := s.node.One("StackName", stackName, p)
	if err != nil {
		return nil, err
	}

	return p.Convert(), nil
}

func (s *identityManagerState) SaveIdentityPoolClient(client *client.IdentityPoolClient) error {
	return s.node.Save(NewIdentityPoolClient(client, NewMetadata()))
}

func (s *identityManagerState) RemoveIdentityPoolClient(stackName string) error {
	p := &IdentityPoolClient{}

	err := s.node.One("StackName", stackName, p)
	if err != nil {
		return err
	}

	return s.node.DeleteStruct(p)
}

func (s *identityManagerState) SaveIdentityPoolUser(user *client.IdentityPoolUser) error {
	return s.node.Save(NewIdentityPoolUser(user, NewMetadata()))
}

// NewIdentityManager returns an initialised state
func NewIdentityManager(node stormpkg.Node) client.IdentityManagerState {
	return &identityManagerState{
		node: node,
	}
}
