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

func (m *managedPolicy) CreateExternalDNSPolicy(opts api.CreateExternalDNSPolicyOpts) (*api.ManagedPolicy, error) {
	b := cfn.New(
		components.NewExternalDNSPolicyComposer(opts.Repository, opts.Environment),
	)

	stackName := cfn.NewStackNamer().
		ExternalDNSPolicy(opts.Repository, opts.Environment)

	return m.createPolicy(stackName, opts.Environment, opts.Repository, "ExternalDNSPolicy", b)
}

func (m *managedPolicy) CreateAlbIngressControllerPolicy(opts api.CreateAlbIngressControllerPolicyOpts) (*api.ManagedPolicy, error) {
	b := cfn.New(
		components.NewAlbIngressControllerPolicyComposer(opts.Repository, opts.Environment),
	)

	stackName := cfn.NewStackNamer().
		AlbIngressControllerPolicy(opts.Repository, opts.Environment)

	return m.createPolicy(stackName, opts.Environment, opts.Repository, "AlbIngressControllerPolicy", b)
}

// CreateExternalSecretsPolicy builds and applies a cloud formation template
func (m *managedPolicy) CreateExternalSecretsPolicy(opts api.CreateExternalSecretsPolicyOpts) (*api.ManagedPolicy, error) {
	b := cfn.New(
		components.NewExternalSecretsPolicyComposer(opts.Repository, opts.Environment),
	)

	stackName := cfn.NewStackNamer().
		ExternalSecretsPolicy(opts.Repository, opts.Environment)

	return m.createPolicy(stackName, opts.Environment, opts.Repository, "ExternalSecretsPolicy", b)
}

func (m *managedPolicy) createPolicy(stackName, env, repoName, outputName string, builder cfn.StackBuilder) (*api.ManagedPolicy, error) {
	template, err := builder.Build()
	if err != nil {
		return nil, errors.E(err, "failed to build cloud formation template")
	}

	r := cfn.NewRunner(m.provider)

	err = r.CreateIfNotExists(stackName, template, []string{cfn.CapabilityNamedIam}, defaultTimeOut)
	if err != nil {
		return nil, errors.E(err, "failed to create cloud formation template")
	}

	p := &api.ManagedPolicy{
		StackName:              stackName,
		Repository:             repoName,
		Environment:            env,
		CloudFormationTemplate: template,
	}

	err = r.Outputs(stackName, map[string]cfn.ProcessOutputFn{
		outputName: cfn.String(&p.PolicyARN),
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
