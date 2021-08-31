package aws

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/config/constant"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/cfn"
	"github.com/oslokommune/okctl/pkg/cfn/components"
)

type containerRepositoryCloudProvider struct {
	provider v1alpha1.CloudProvider
}

func (c *containerRepositoryCloudProvider) CreateContainerRepository(opts *api.CreateContainerRepositoryOpts) (*api.ContainerRepository, error) {
	composition := components.NewECRRepositoryComposer(opts.Name)

	template, err := cfn.New(composition).Build()
	if err != nil {
		return nil, fmt.Errorf(constant.BuildCloudFormationTemplateError, err)
	}

	r := cfn.NewRunner(c.provider)

	err = r.CreateIfNotExists(
		opts.ClusterID.ClusterName,
		opts.StackName,
		template,
		nil,
		defaultTimeOut,
	)
	if err != nil {
		return nil, fmt.Errorf(constant.CreateCloudFormationStackError, err)
	}

	repository := &api.ContainerRepository{
		Name:                   opts.Name,
		ClusterID:              opts.ClusterID,
		StackName:              opts.StackName,
		CloudFormationTemplate: string(template),
	}

	err = r.Outputs(opts.StackName, map[string]cfn.ProcessOutputFn{
		composition.ResourceRepositoryNameOutput(): cfn.String(&repository.Name),
	})
	if err != nil {
		return nil, fmt.Errorf(constant.CollectStackFormationOutputsError, err)
	}

	return repository, nil
}

func (c *containerRepositoryCloudProvider) DeleteContainerRepository(opts *api.DeleteContainerRepositoryOpts) error {
	return cfn.NewRunner(c.provider).Delete(opts.StackName)
}

// NewContainerRepositoryCloudProvider returns an initialised container repository cloud provider
func NewContainerRepositoryCloudProvider(provider v1alpha1.CloudProvider) api.ContainerRepositoryCloudProvider {
	return &containerRepositoryCloudProvider{
		provider: provider,
	}
}
