package aws

import (
	"bytes"
	"context"
	"fmt"

	merrors "github.com/mishudark/errors"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/cfn"
	"github.com/oslokommune/okctl/pkg/cfn/components"
	"github.com/oslokommune/okctl/pkg/cfn/components/securitygroup"
)

type securityGroupCloudProvider struct {
	provider v1alpha1.CloudProvider
}

const (
	defaultSecurityGroupResourceName = "SecurityGroup"
	defaultSecurityGroupIDOutput     = defaultSecurityGroupResourceName + "GroupId"
)

// CreateSecurityGroup knows how to create a Security Group through CloudFormation
func (s *securityGroupCloudProvider) CreateSecurityGroup(_ context.Context, opts api.CreateSecurityGroupOpts) (api.SecurityGroup, error) {
	composer := components.NewSecurityGroupComposer(components.NewSecurityGroupComposerOpts{
		ClusterName:   opts.ClusterID.ClusterName,
		ResourceName:  defaultSecurityGroupResourceName,
		VPCID:         opts.VPCID,
		Name:          opts.Name,
		Description:   opts.Description,
		InboundRules:  opts.InboundRules,
		OutboundRules: opts.OutboundRules,
	})

	b := cfn.New(composer)

	template, err := b.Build()
	if err != nil {
		return api.SecurityGroup{}, fmt.Errorf("building cloud formation template: %w", err)
	}

	r := cfn.NewRunner(s.provider)

	stackName := cfn.NewStackNamer().SecurityGroup(opts.ClusterID.ClusterName, opts.Name)

	err = r.CreateIfNotExists(
		opts.ClusterID.ClusterName,
		stackName,
		template,
		nil,
		postgresTimeOutInMinutes,
	)
	if err != nil {
		return api.SecurityGroup{}, fmt.Errorf("creating cloud formation stack: %w", err)
	}

	securityGroup := api.SecurityGroup{
		InboundRules:  opts.InboundRules,
		OutboundRules: opts.OutboundRules,
	}

	err = r.Outputs(stackName, map[string]cfn.ProcessOutputFn{
		defaultSecurityGroupIDOutput: cfn.String(&securityGroup.ID),
	})
	if err != nil {
		return api.SecurityGroup{}, fmt.Errorf("collecting stack outputs: %w", err)
	}

	return securityGroup, nil
}

// GetSecurityGroup knows how to extract a Security Group from a CloudFormation stack
func (s *securityGroupCloudProvider) GetSecurityGroup(_ context.Context, opts api.GetSecurityGroupOpts) (api.SecurityGroup, error) {
	stackName := cfn.NewStackNamer().SecurityGroup(opts.ClusterName, opts.Name)

	stack, err := cfn.NewRunner(s.provider).Get(stackName)
	if err != nil {
		return api.SecurityGroup{}, merrors.E(err, fmt.Sprintf("getting stack %s", stackName))
	}

	securityGroupID, err := cfn.GetOutput(stack.Outputs, defaultSecurityGroupIDOutput)
	if err != nil {
		return api.SecurityGroup{}, fmt.Errorf("getting security group ID: %w", err)
	}

	return api.SecurityGroup{
		ID: securityGroupID,
	}, nil
}

// DeleteSecurityGroup knows how to delete a Security Group through CloudFormation
func (s *securityGroupCloudProvider) DeleteSecurityGroup(_ context.Context, opts api.DeleteSecurityGroupOpts) error {
	stackName := cfn.NewStackNamer().SecurityGroup(opts.ClusterName, opts.Name)

	return cfn.NewRunner(s.provider).Delete(stackName)
}

// AddRule knows how to add a rule to an existing Security Group using CloudFormation
func (s *securityGroupCloudProvider) AddRule(_ context.Context, opts api.AddRuleOpts) (api.Rule, error) {
	r := cfn.NewRunner(s.provider)

	originalTemplate, err := r.GetTemplate(opts.SecurityGroupStackName)
	if err != nil {
		return api.Rule{}, fmt.Errorf("getting security group stack template: %w", err)
	}

	var patchFunc func([]byte, string, api.Rule) ([]byte, error)

	switch opts.RuleType {
	case api.RuleTypeIngress:
		patchFunc = securitygroup.PatchAppendIngressRule
	case api.RuleTypeEgress:
		patchFunc = securitygroup.PatchAppendEgressRule
	default:
		return api.Rule{}, fmt.Errorf("unknown rule type %s", string(opts.RuleType))
	}

	updatedTemplate, err := patchFunc(originalTemplate, opts.SecurityGroupResourceName, opts.Rule)
	if err != nil {
		return api.Rule{}, fmt.Errorf("applying patch: %w", err)
	}

	if bytes.Equal(updatedTemplate, originalTemplate) {
		return opts.Rule, nil
	}

	err = r.Update(opts.SecurityGroupStackName, updatedTemplate)
	if err != nil {
		return api.Rule{}, fmt.Errorf("updating cloud formation stack: %w", err)
	}

	return opts.Rule, nil
}

// RemoveRule knows how to remove a rule from an existing Security Group using CloudFormation
func (s *securityGroupCloudProvider) RemoveRule(_ context.Context, opts api.RemoveRuleOpts) error {
	r := cfn.NewRunner(s.provider)

	originalTemplate, err := r.GetTemplate(opts.SecurityGroupStackName)
	if err != nil {
		return fmt.Errorf("getting security group stack template: %w", err)
	}

	var patchFunc func([]byte, string, api.Rule) ([]byte, error)

	switch opts.RuleType {
	case api.RuleTypeIngress:
		patchFunc = securitygroup.PatchRemoveIngressRule
	case api.RuleTypeEgress:
		patchFunc = securitygroup.PatchRemoveEgressRule
	default:
		return fmt.Errorf("unknown rule type %s", string(opts.RuleType))
	}

	updatedTemplate, err := patchFunc(originalTemplate, opts.SecurityGroupResourceName, opts.Rule)
	if err != nil {
		return fmt.Errorf("applying patch: %w", err)
	}

	err = r.Update(opts.SecurityGroupStackName, updatedTemplate)
	if err != nil {
		return fmt.Errorf("updating cloud formation stack: %w", err)
	}

	return nil
}

// NewSecurityGroupCloudProvider initializes a new SecurityGroupCloudProvider
func NewSecurityGroupCloudProvider(provider v1alpha1.CloudProvider) api.SecurityGroupCloudProvider {
	return &securityGroupCloudProvider{provider: provider}
}
