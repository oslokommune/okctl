package aws

import (
	"fmt"
	"github.com/oslokommune/okctl/pkg/config/constant"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/cfn"
)

type managedPolicy struct {
	provider v1alpha1.CloudProvider
}

func (m *managedPolicy) CreatePolicy(opts api.CreatePolicyOpts) (*api.ManagedPolicy, error) {
	r := cfn.NewRunner(m.provider)

	err := r.CreateIfNotExists(
		opts.ID.ClusterName,
		opts.StackName,
		opts.CloudFormationTemplate,
		[]string{cfn.CapabilityNamedIam},
		defaultTimeOut,
	)
	if err != nil {
		return nil, fmt.Errorf(constant.CreateCloudFormationStackError, err)
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
		return nil, fmt.Errorf(constant.ProcessOutputsError, err)
	}

	return p, nil
}

func (m *managedPolicy) DeletePolicy(opts api.DeletePolicyOpts) error {
	r := cfn.NewRunner(m.provider)

	err := r.Delete(opts.StackName)
	if err != nil {
		return fmt.Errorf(constant.DeletePolicyError, err)
	}

	return nil
}

// NewManagedPolicyCloudProvider returns an initialised cloud provider
func NewManagedPolicyCloudProvider(provider v1alpha1.CloudProvider) api.ManagedPolicyCloudProvider {
	return &managedPolicy{
		provider: provider,
	}
}
