package storm

import (
	"errors"
	"time"

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
	CloudFormationTemplates string
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
		CloudFormationTemplates: string(p.CloudFormationTemplates),
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
		CloudFormationTemplates: []byte(p.CloudFormationTemplates),
		Certificate:             p.Certificate.Convert(),
		RecordSetAlias:          p.RecordSetAlias.Convert(),
	}
}

// RecordSetAlias contains storm compatible state
type RecordSetAlias struct {
	AliasDomain            string
	AliasHostedZones       string
	StackName              string
	CloudFormationTemplate string
}

// NewRecordSetAlias returns a storm compatible RecordSetAlias
func NewRecordSetAlias(a *client.RecordSetAlias) *RecordSetAlias {
	return &RecordSetAlias{
		AliasDomain:            a.AliasDomain,
		AliasHostedZones:       a.AliasHostedZones,
		StackName:              a.StackName,
		CloudFormationTemplate: string(a.CloudFormationTemplate),
	}
}

// Convert to client.RecordSetAlias
func (a *RecordSetAlias) Convert() *client.RecordSetAlias {
	return &client.RecordSetAlias{
		AliasDomain:            a.AliasDomain,
		AliasHostedZones:       a.AliasHostedZones,
		StackName:              a.StackName,
		CloudFormationTemplate: []byte(a.CloudFormationTemplate),
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
	CloudFormationTemplates string
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
		CloudFormationTemplates: string(c.CloudFormationTemplates),
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
		CloudFormationTemplates: []byte(c.CloudFormationTemplates),
	}
}

// IdentityPoolUser contains storm compatible state
type IdentityPoolUser struct {
	Metadata `storm:"inline"`

	ID                     ID
	Email                  string
	UserPoolID             string
	StackName              string `storm:"unique,index"`
	CloudFormationTemplate string
}

// NewIdentityPoolUser returns a storm compatible IdentityPoolUser
func NewIdentityPoolUser(u *client.IdentityPoolUser, meta Metadata) *IdentityPoolUser {
	return &IdentityPoolUser{
		Metadata:               meta,
		ID:                     NewID(u.ID),
		Email:                  u.Email,
		UserPoolID:             u.UserPoolID,
		StackName:              u.StackName,
		CloudFormationTemplate: string(u.CloudFormationTemplate),
	}
}

// Convert to client.IdentityPoolUser
func (u *IdentityPoolUser) Convert() *client.IdentityPoolUser {
	return &client.IdentityPoolUser{
		ID:                     u.ID.Convert(),
		Email:                  u.Email,
		UserPoolID:             u.UserPoolID,
		StackName:              u.StackName,
		CloudFormationTemplate: []byte(u.CloudFormationTemplate),
	}
}

func (s *identityManagerState) SaveIdentityPool(pool *client.IdentityPool) error {
	existing, err := s.getIdentityPool(pool.StackName)
	if err != nil && !errors.Is(err, stormpkg.ErrNotFound) {
		return err
	}

	if errors.Is(err, stormpkg.ErrNotFound) {
		return s.node.Save(NewIdentityPool(pool, NewMetadata(), NewMetadata()))
	}

	existing.Metadata.UpdatedAt = time.Now()
	existing.Certificate.Metadata.UpdatedAt = time.Now()

	return s.node.Save(NewIdentityPool(pool, existing.Metadata, existing.Certificate.Metadata))
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
	p, err := s.getIdentityPool(stackName)
	if err != nil {
		return nil, err
	}

	return p.Convert(), nil
}

func (s *identityManagerState) getIdentityPool(stackName string) (*IdentityPool, error) {
	p := &IdentityPool{}

	err := s.node.One("StackName", stackName, p)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (s *identityManagerState) SaveIdentityPoolClient(client *client.IdentityPoolClient) error {
	existing, err := s.getIdentityPoolClient(client.StackName)
	if err != nil && !errors.Is(err, stormpkg.ErrNotFound) {
		return err
	}

	if errors.Is(err, stormpkg.ErrNotFound) {
		return s.node.Save(NewIdentityPoolClient(client, NewMetadata()))
	}

	existing.Metadata.UpdatedAt = time.Now()

	return s.node.Save(NewIdentityPoolClient(client, existing.Metadata))
}

func (s *identityManagerState) GetIdentityPoolClient(stackName string) (*client.IdentityPoolClient, error) {
	p, err := s.getIdentityPoolClient(stackName)
	if err != nil {
		return nil, err
	}

	return p.Convert(), nil
}

func (s *identityManagerState) getIdentityPoolClient(stackName string) (*IdentityPoolClient, error) {
	p := &IdentityPoolClient{}

	err := s.node.One("StackName", stackName, p)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (s *identityManagerState) RemoveIdentityPoolClient(stackName string) error {
	p := &IdentityPoolClient{}

	err := s.node.One("StackName", stackName, p)
	if err != nil {
		if errors.Is(err, stormpkg.ErrNotFound) {
			return nil
		}

		return err
	}

	return s.node.DeleteStruct(p)
}

func (s *identityManagerState) GetIdentityPoolUser(stackName string) (*client.IdentityPoolUser, error) {
	u, err := s.getIdentityPoolUser(stackName)
	if err != nil {
		return nil, err
	}

	return u.Convert(), nil
}

func (s *identityManagerState) getIdentityPoolUser(stackName string) (*IdentityPoolUser, error) {
	u := &IdentityPoolUser{}

	err := s.node.One("StackName", stackName, u)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (s *identityManagerState) RemoveIdentityPoolUser(stackName string) error {
	p := &IdentityPoolUser{}

	err := s.node.One("StackName", stackName, p)
	if err != nil {
		if errors.Is(err, stormpkg.ErrNotFound) {
			return nil
		}

		return err
	}

	return s.node.DeleteStruct(p)
}

func (s *identityManagerState) SaveIdentityPoolUser(user *client.IdentityPoolUser) error {
	existing, err := s.getIdentityPoolUser(user.StackName)
	if err != nil && !errors.Is(err, stormpkg.ErrNotFound) {
		return err
	}

	if errors.Is(err, stormpkg.ErrNotFound) {
		return s.node.Save(NewIdentityPoolUser(user, NewMetadata()))
	}

	existing.Metadata.UpdatedAt = time.Now()

	return s.node.Save(NewIdentityPoolUser(user, existing.Metadata))
}

// NewIdentityManager returns an initialised state
func NewIdentityManager(node stormpkg.Node) client.IdentityManagerState {
	return &identityManagerState{
		node: node,
	}
}
