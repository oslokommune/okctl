package aws

import (
	"context"
	"fmt"

	"github.com/oslokommune/okctl/pkg/version"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/cfn"
)

type managedPolicy struct {
	provider  v1alpha1.CloudProvider
	versioner version.Versioner
}

func (m *managedPolicy) CreatePolicy(ctx context.Context, opts api.CreatePolicyOpts) (*api.ManagedPolicy, error) {
	versionInfo, err := m.versioner.GetVersionInfo(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting version info: %w", err)
	}

	r := cfn.NewRunner(m.provider)

	err = r.CreateIfNotExists(
		versionInfo,
		opts.ID.ClusterName,
		opts.StackName,
		opts.CloudFormationTemplate,
		[]string{cfn.CapabilityNamedIam},
		defaultTimeOut,
	)
	if err != nil {
		return nil, fmt.Errorf("creating cloud formation stack: %w", err)
	}

	p := &api.ManagedPolicy{
		ID:                     opts.ID,
		StackName:              opts.StackName,
		CloudFormationTemplate: opts.CloudFormationTemplate,
	}

	err = r.Outputs(opts.StackName, map[string]cfn.ProcessOutputFn{
		opts.PolicyOutputName: cfn.String(&p.PolicyARN),
	})
	if err != nil {
		return nil, fmt.Errorf("processing outputs: %w", err)
	}

	return p, nil
}

func (m *managedPolicy) DeletePolicy(opts api.DeletePolicyOpts) error {
	r := cfn.NewRunner(m.provider)

	err := r.Delete(opts.StackName)
	if err != nil {
		return fmt.Errorf("deleting policy: %w", err)
	}

	return nil
}

// NewManagedPolicyCloudProvider returns an initialised cloud provider
func NewManagedPolicyCloudProvider(provider v1alpha1.CloudProvider, versioner version.Versioner) api.ManagedPolicyCloudProvider {
	return &managedPolicy{
		provider:  provider,
		versioner: versioner,
	}
}
