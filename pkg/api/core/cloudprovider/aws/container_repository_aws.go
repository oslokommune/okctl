package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecr"

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
		return nil, fmt.Errorf("building the cloud formation template: %w", err)
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
		return nil, fmt.Errorf("creating cloud formation stack: %w", err)
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
		return nil, fmt.Errorf("collecting stack outputs: %w", err)
	}

	return repository, nil
}

func (c *containerRepositoryCloudProvider) DeleteContainerRepository(opts *api.DeleteContainerRepositoryOpts) error {
	return cfn.NewRunner(c.provider).Delete(opts.StackName)
}

func (c *containerRepositoryCloudProvider) EmptyContainerRepository(ctx context.Context, opts api.EmptyContainerRepositoryOpts) error {
	result, err := c.provider.ECR().ListImagesWithContext(ctx, &ecr.ListImagesInput{
		RepositoryName: aws.String(opts.Name),
	})
	if err != nil {
		return fmt.Errorf("acquiring image IDs: %w", err)
	}

	_, err = c.provider.ECR().BatchDeleteImageWithContext(ctx, &ecr.BatchDeleteImageInput{
		ImageIds:       result.ImageIds,
		RepositoryName: aws.String(opts.Name),
	})
	if err != nil {
		return fmt.Errorf("deleting images: %w", err)
	}

	return nil
}

// NewContainerRepositoryCloudProvider returns an initialised container repository cloud provider
func NewContainerRepositoryCloudProvider(provider v1alpha1.CloudProvider) api.ContainerRepositoryCloudProvider {
	return &containerRepositoryCloudProvider{
		provider: provider,
	}
}
