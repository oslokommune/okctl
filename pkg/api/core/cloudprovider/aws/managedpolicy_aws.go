package aws

import (
	"fmt"

	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/cfn"
	"github.com/oslokommune/okctl/pkg/cfn/components"
)

type managedPolicy struct {
	provider v1alpha1.CloudProvider
}

func (m *managedPolicy) DeleteBlockstoragePolicy(id api.ID) error {
	return m.deletePolicy(cfn.NewStackNamer().BlockstoragePolicy(id.Repository, id.Environment))
}

func (m *managedPolicy) DeleteAutoscalerPolicy(id api.ID) error {
	return m.deletePolicy(cfn.NewStackNamer().AutoscalerPolicy(id.Repository, id.Environment))
}

func (m *managedPolicy) DeleteExternalSecretsPolicy(id api.ID) error {
	return m.deletePolicy(cfn.NewStackNamer().ExternalSecretsPolicy(id.Repository, id.Environment))
}

func (m *managedPolicy) DeleteAlbIngressControllerPolicy(id api.ID) error {
	return m.deletePolicy(cfn.NewStackNamer().AlbIngressControllerPolicy(id.Repository, id.Environment))
}

func (m *managedPolicy) DeleteAWSLoadBalancerControllerPolicy(id api.ID) error {
	return m.deletePolicy(cfn.NewStackNamer().AWSLoadBalancerControllerPolicy(id.Repository, id.Environment))
}

func (m *managedPolicy) DeleteExternalDNSPolicy(id api.ID) error {
	return m.deletePolicy(cfn.NewStackNamer().ExternalDNSPolicy(id.Repository, id.Environment))
}

func (m *managedPolicy) deletePolicy(stackName string) error {
	r := cfn.NewRunner(m.provider)

	err := r.Delete(stackName)
	if err != nil {
		return fmt.Errorf("deleting policy: %w", err)
	}

	return nil
}

func (m *managedPolicy) CreateBlockstoragePolicy(opts api.CreateBlockstoragePolicy) (*api.ManagedPolicy, error) {
	b := cfn.New(
		components.NewBlockstoragePolicyComposer(opts.ID.Repository, opts.ID.Environment),
	)

	stackName := cfn.NewStackNamer().
		BlockstoragePolicy(opts.ID.Repository, opts.ID.Environment)

	return m.createPolicy(stackName, opts.ID, "BlockstoragePolicy", b)
}

func (m *managedPolicy) CreateAutoscalerPolicy(opts api.CreateAutoscalerPolicy) (*api.ManagedPolicy, error) {
	b := cfn.New(
		components.NewAutoscalerPolicyComposer(opts.ID.Repository, opts.ID.Environment),
	)

	stackName := cfn.NewStackNamer().
		AutoscalerPolicy(opts.ID.Repository, opts.ID.Environment)

	return m.createPolicy(stackName, opts.ID, "AutoscalerPolicy", b)
}

func (m *managedPolicy) CreateExternalDNSPolicy(opts api.CreateExternalDNSPolicyOpts) (*api.ManagedPolicy, error) {
	b := cfn.New(
		components.NewExternalDNSPolicyComposer(opts.ID.Repository, opts.ID.Environment),
	)

	stackName := cfn.NewStackNamer().
		ExternalDNSPolicy(opts.ID.Repository, opts.ID.Environment)

	return m.createPolicy(stackName, opts.ID, "ExternalDNSPolicy", b)
}

func (m *managedPolicy) CreateAlbIngressControllerPolicy(opts api.CreateAlbIngressControllerPolicyOpts) (*api.ManagedPolicy, error) {
	b := cfn.New(
		components.NewAlbIngressControllerPolicyComposer(opts.ID.Repository, opts.ID.Environment),
	)

	stackName := cfn.NewStackNamer().
		AlbIngressControllerPolicy(opts.ID.Repository, opts.ID.Environment)

	return m.createPolicy(stackName, opts.ID, "AlbIngressControllerPolicy", b)
}

func (m *managedPolicy) CreateAWSLoadBalancerControllerPolicy(opts api.CreateAWSLoadBalancerControllerPolicyOpts) (*api.ManagedPolicy, error) {
	b := cfn.New(
		components.NewAWSLoadBalancerControllerComposer(opts.ID.Repository, opts.ID.Environment),
	)

	stackName := cfn.NewStackNamer().
		AWSLoadBalancerControllerPolicy(opts.ID.Repository, opts.ID.Environment)

	return m.createPolicy(stackName, opts.ID, "AWSLoadBalancerControllerPolicy", b)
}

// CreateExternalSecretsPolicy builds and applies a cloud formation template
func (m *managedPolicy) CreateExternalSecretsPolicy(opts api.CreateExternalSecretsPolicyOpts) (*api.ManagedPolicy, error) {
	b := cfn.New(
		components.NewExternalSecretsPolicyComposer(opts.ID.Repository, opts.ID.Environment),
	)

	stackName := cfn.NewStackNamer().
		ExternalSecretsPolicy(opts.ID.Repository, opts.ID.Environment)

	return m.createPolicy(stackName, opts.ID, "ExternalSecretsPolicy", b)
}

func (m *managedPolicy) createPolicy(stackName string, id api.ID, outputName string, builder cfn.StackBuilder) (*api.ManagedPolicy, error) {
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
		ID:                     id,
		StackName:              stackName,
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
