package storm

import (
	"errors"

	stormpkg "github.com/asdine/storm/v3"
	"github.com/oslokommune/okctl/pkg/client"
)

type managedPolicyState struct {
	node stormpkg.Node
}

// ManagedPolicy contains storm compatible state
type ManagedPolicy struct {
	Metadata `storm:"inline"`

	ID                     ID
	StackName              string `storm:"unique"`
	PolicyARN              string
	CloudFormationTemplate []byte
}

// NewManagedPolicy returns a storm compatible ManagedPolicy
func NewManagedPolicy(p *client.ManagedPolicy, meta Metadata) *ManagedPolicy {
	return &ManagedPolicy{
		Metadata:               meta,
		ID:                     NewID(p.ID),
		StackName:              p.StackName,
		PolicyARN:              p.PolicyARN,
		CloudFormationTemplate: p.CloudFormationTemplate,
	}
}

// Convert to a client.ManagedPolicy
func (p *ManagedPolicy) Convert() *client.ManagedPolicy {
	return &client.ManagedPolicy{
		ID:                     p.ID.Convert(),
		StackName:              p.StackName,
		PolicyARN:              p.PolicyARN,
		CloudFormationTemplate: p.CloudFormationTemplate,
	}
}

func (m *managedPolicyState) SavePolicy(policy *client.ManagedPolicy) error {
	return m.node.Save(NewManagedPolicy(policy, NewMetadata()))
}

func (m *managedPolicyState) GetPolicy(stackName string) (*client.ManagedPolicy, error) {
	p := &ManagedPolicy{}

	err := m.node.One("StackName", stackName, p)
	if err != nil {
		return nil, err
	}

	return p.Convert(), nil
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
func NewManagedPolicyState(node stormpkg.Node) client.ManagedPolicyState {
	return &managedPolicyState{
		node: node,
	}
}
