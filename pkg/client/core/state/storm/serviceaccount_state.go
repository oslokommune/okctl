package storm

import (
	"time"

	stormpkg "github.com/asdine/storm/v3"
	"github.com/oslokommune/okctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/oslokommune/okctl/pkg/client"
)

// ServiceAccount contains storm compatible state
type ServiceAccount struct {
	Metadata `storm:"inline"`

	ID        ID
	Name      string `storm:"unique,index"`
	PolicyArn string
	Config    *v1alpha5.ClusterConfig
}

// NewServiceAccount constructs a storm compatible ServiceAccount
func NewServiceAccount(sa *client.ServiceAccount, meta Metadata) *ServiceAccount {
	return &ServiceAccount{
		Metadata:  meta,
		ID:        NewID(sa.ID),
		Name:      sa.Name,
		PolicyArn: sa.PolicyArn,
		Config:    sa.Config,
	}
}

// Convert to a client.ServiceAccount
func (sa *ServiceAccount) Convert() *client.ServiceAccount {
	return &client.ServiceAccount{
		ID:        sa.ID.Convert(),
		Name:      sa.Name,
		PolicyArn: sa.PolicyArn,
		Config:    sa.Config,
	}
}

type serviceAccountState struct {
	node stormpkg.Node
}

func (s *serviceAccountState) SaveServiceAccount(account *client.ServiceAccount) error {
	return s.node.Save(NewServiceAccount(account, NewMetadata()))
}

func (s *serviceAccountState) RemoveServiceAccount(name string) error {
	sa := &ServiceAccount{}

	err := s.node.One("Name", name, sa)
	if err != nil {
		return err
	}

	return s.node.DeleteStruct(sa)
}

func (s *serviceAccountState) GetServiceAccount(name string) (*client.ServiceAccount, error) {
	sa := &ServiceAccount{}

	err := s.node.One("Name", name, sa)
	if err != nil {
		return nil, err
	}

	return sa.Convert(), nil
}

func (s *serviceAccountState) UpdateServiceAccount(account *client.ServiceAccount) error {
	sa := &ServiceAccount{}

	err := s.node.One("Name", account.Name, sa)
	if err != nil {
		return err
	}

	updated := NewServiceAccount(account, sa.Metadata)
	updated.UpdatedAt = time.Now()

	return s.node.Save(updated)
}

// NewServiceAccountState returns an initialised state store
func NewServiceAccountState(node stormpkg.Node) client.ServiceAccountState {
	return &serviceAccountState{
		node: node,
	}
}
