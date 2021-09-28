package core

import (
	"context"
	"fmt"

	"github.com/oslokommune/okctl/pkg/api"
)

type securityGroupService struct {
	provider api.SecurityGroupCloudProvider
}

// CreateSecurityGroup knows how to use the provider to create a Security Group
func (s *securityGroupService) CreateSecurityGroup(ctx context.Context, opts api.CreateSecurityGroupOpts) (api.SecurityGroup, error) {
	err := opts.Validate()
	if err != nil {
		return api.SecurityGroup{}, fmt.Errorf("validating opts: %w", err)
	}

	return s.provider.CreateSecurityGroup(ctx, opts)
}

// GetSecurityGroup knows how to use the provider to acquire an existing Security Group
func (s *securityGroupService) GetSecurityGroup(ctx context.Context, opts api.GetSecurityGroupOpts) (api.SecurityGroup, error) {
	err := opts.Validate()
	if err != nil {
		return api.SecurityGroup{}, fmt.Errorf("validating opts: %w", err)
	}

	return s.provider.GetSecurityGroup(ctx, opts)
}

// DeleteSecurityGroup knows how to use the provider to delete a Security Group
func (s *securityGroupService) DeleteSecurityGroup(ctx context.Context, opts api.DeleteSecurityGroupOpts) error {
	err := opts.Validate()
	if err != nil {
		return fmt.Errorf("validating opts: %w", err)
	}

	return s.provider.DeleteSecurityGroup(ctx, opts)
}

// AddRule knows how to use the provider to add a rule to an existing Security Group
func (s *securityGroupService) AddRule(ctx context.Context, opts api.AddRuleOpts) (api.Rule, error) {
	err := opts.Validate()
	if err != nil {
		return api.Rule{}, fmt.Errorf("validating opts: %w", err)
	}

	return s.provider.AddRule(ctx, opts)
}

// RemoveRule knows how to use the provider to remove a rule from an existing Security Group
func (s *securityGroupService) RemoveRule(ctx context.Context, opts api.RemoveRuleOpts) error {
	err := opts.Validate()
	if err != nil {
		return fmt.Errorf("validating opts: %w", err)
	}

	return s.provider.RemoveRule(ctx, opts)
}

// NewSecurityGroupService returns an initialized SecurityGroupService
func NewSecurityGroupService(provider api.SecurityGroupCloudProvider) api.SecurityGroupService {
	return &securityGroupService{provider: provider}
}
