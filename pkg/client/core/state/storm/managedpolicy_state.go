package storm

import (
	"errors"
	"time"

	stormpkg "github.com/asdine/storm/v3"
	"github.com/oslokommune/okctl/pkg/breeze"
	"github.com/oslokommune/okctl/pkg/client"
)

type managedPolicyState struct {
	node breeze.Client
}

// ManagedPolicy contains storm compatible state
type ManagedPolicy struct {
	Metadata `storm:"inline"`

	ID                     ID
	StackName              string `storm:"unique"`
	PolicyARN              string
	CloudFormationTemplate string
}

// NewManagedPolicy returns a storm compatible ManagedPolicy
func NewManagedPolicy(p *client.ManagedPolicy, meta Metadata) *ManagedPolicy {
	return &ManagedPolicy{
		Metadata:               meta,
		ID:                     NewID(p.ID),
		StackName:              p.StackName,
		PolicyARN:              p.PolicyARN,
		CloudFormationTemplate: string(p.CloudFormationTemplate),
	}
}

// Convert to a client.ManagedPolicy
func (p *ManagedPolicy) Convert() *client.ManagedPolicy {
	return &client.ManagedPolicy{
		ID:                     p.ID.Convert(),
		StackName:              p.StackName,
		PolicyARN:              p.PolicyARN,
		CloudFormationTemplate: []byte(p.CloudFormationTemplate),
	}
}

func (m *managedPolicyState) SavePolicy(policy *client.ManagedPolicy) error {
	existing, err := m.getPolicy(policy.StackName)
	if err != nil && !errors.Is(err, stormpkg.ErrNotFound) {
		return err
	}

	if errors.Is(err, stormpkg.ErrNotFound) {
		return m.node.Save(NewManagedPolicy(policy, NewMetadata()))
	}

	existing.Metadata.UpdatedAt = time.Now()

	return m.node.Save(NewManagedPolicy(policy, existing.Metadata))
}

func (m *managedPolicyState) GetPolicy(stackName string) (*client.ManagedPolicy, error) {
	p, err := m.getPolicy(stackName)
	if err != nil {
		return nil, err
	}

	return p.Convert(), nil
}

func (m *managedPolicyState) getPolicy(stackName string) (*ManagedPolicy, error) {
	p := &ManagedPolicy{}

	err := m.node.One("StackName", stackName, p)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (m *managedPolicyState) RemovePolicy(stackName string) error {
	p := &ManagedPolicy{}

	err := m.node.One("StackName", stackName, p)
	if err != nil {
		if errors.Is(err, stormpkg.ErrNotFound) {
			return nil
		}

		return err
	}

	return m.node.DeleteStruct(p)
}

// NewManagedPolicyState returns an initialised managed policy state
func NewManagedPolicyState(node breeze.Client) client.ManagedPolicyState {
	return &managedPolicyState{
		node: node,
	}
}
