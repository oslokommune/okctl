package aws

import (
	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/api/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/cfn"
	"github.com/oslokommune/okctl/pkg/cfn/components"
)

type managedPolicy struct {
	provider v1alpha1.CloudProvider
}

func (m *managedPolicy) CreateAlbIngressControllerPolicy(opts api.CreateAlbIngressControllerPolicyOpts) (*api.ManagedPolicy, error) {
	b := cfn.New(components.NewAlbIngressControllerPolicyComposer(opts.Repository, opts.Environment))

	template, err := b.Build()
	if err != nil {
		return nil, errors.E(err, "failed to build cloud formation template")
	}

	stackName := cfn.NewStackNamer().AlbIngressControllerPolicy(opts.Repository, opts.Environment)

	r := cfn.NewRunner(m.provider)

	err = r.CreateIfNotExists(stackName, template, []string{cfn.CapabilityNamedIam}, defaultTimeOut)
	if err != nil {
		return nil, errors.E(err, "cloud provider failed to create policy")
	}

	p := &api.ManagedPolicy{
		StackName:              stackName,
		Repository:             opts.Repository,
		Environment:            opts.Environment,
		CloudFormationTemplate: template,
	}

	err = r.Outputs(stackName, map[string]cfn.ProcessOutputFn{
		"AlbIngressControllerPolicy": cfn.String(&p.PolicyARN),
	})
	if err != nil {
		return nil, errors.E(err, "failed to process outputs")
	}

	return p, nil
}

// CreateExternalSecretsPolicy builds and applies a cloud formation template
func (m *managedPolicy) CreateExternalSecretsPolicy(opts api.CreateExternalSecretsPolicyOpts) (*api.ManagedPolicy, error) {
	b := cfn.New(components.NewExternalSecretsPolicyComposer(opts.Repository, opts.Environment))

	template, err := b.Build()
	if err != nil {
		return nil, errors.E(err, "failed to build cloud formation template", errors.Internal)
	}

	stackName := cfn.NewStackNamer().ExternalSecretsPolicy(opts.Repository, opts.Environment)

	r := cfn.NewRunner(m.provider)

	err = r.CreateIfNotExists(stackName, template, []string{cfn.CapabilityNamedIam}, defaultTimeOut)
	if err != nil {
		return nil, errors.E(err, "cloud provider failed to create policy", errors.Unknown)
	}

	p := &api.ManagedPolicy{
		StackName:              stackName,
		Repository:             opts.Repository,
		Environment:            opts.Environment,
		CloudFormationTemplate: template,
	}

	err = r.Outputs(stackName, map[string]cfn.ProcessOutputFn{
		"ExternalSecretsPolicy": cfn.String(&p.PolicyARN),
	})
	if err != nil {
		return nil, errors.E(err, "failed to process outputs")
	}

	return p, nil
}

// NewManagedPolicyCloudProvider returns an initialised cloud provider
func NewManagedPolicyCloudProvider(provider v1alpha1.CloudProvider) api.ManagedPolicyCloudProvider {
	return &managedPolicy{
		provider: provider,
	}
}
