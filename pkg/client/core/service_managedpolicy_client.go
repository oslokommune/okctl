package core

import (
	"context"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
)

type managedPolicyService struct {
	service api.ManagedPolicyService
	state   client.ManagedPolicyState
}

func (m *managedPolicyService) CreatePolicy(context context.Context, opts client.CreatePolicyOpts) (*client.ManagedPolicy, error) {
	p, err := m.service.CreatePolicy(context, api.CreatePolicyOpts{
		ID:                     opts.ID,
		StackName:              opts.StackName,
		PolicyOutputName:       opts.PolicyOutputName,
		CloudFormationTemplate: opts.CloudFormationTemplate,
	})
	if err != nil {
		return nil, err
	}

	policy := &client.ManagedPolicy{
		ID:                     p.ID,
		StackName:              p.StackName,
		PolicyARN:              p.PolicyARN,
		CloudFormationTemplate: p.CloudFormationTemplate,
	}

	err = m.state.SavePolicy(policy)
	if err != nil {
		return nil, err
	}

	return policy, nil
}

func (m *managedPolicyService) DeletePolicy(context context.Context, opts client.DeletePolicyOpts) error {
	err := m.service.DeletePolicy(context, api.DeletePolicyOpts{
		ID:        opts.ID,
		StackName: opts.StackName,
	})
	if err != nil {
		return err
	}

	return m.state.RemovePolicy(opts.StackName)
}

// NewManagedPolicyService returns an initialised service
func NewManagedPolicyService(
	service api.ManagedPolicyService,
	state client.ManagedPolicyState,
) client.ManagedPolicyService {
	return &managedPolicyService{
		service: service,
		state:   state,
	}
}
