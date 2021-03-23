// Package containerrepository knows how to create an AWS ECRepository
// cloud formation resource
package containerrepository

import (
	"fmt"

	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/awslabs/goformation/v4/cloudformation/ecr"
	"github.com/oslokommune/okctl/pkg/cfn"
)

type (
	TagMutabilityConfiguration string
	ImageScanConfiguration     bool
)

const (
	TagMutabilityConfigurationMutable   = "MUTABLE"
	TagMutabilityConfigurationImmutable = "IMMUTABLE"

	ImageScanConfigurationOn  = true
	ImageScanConfigurationOff = false
)

type ContainerRepository struct {
	StoredName string
}

// TODO: Lifecycle policy
func (c *ContainerRepository) Resource() cloudformation.Resource {
	return &ecr.Repository{
		ImageScanningConfiguration: ImageScanConfigurationOn,
		ImageTagMutability:         TagMutabilityConfigurationImmutable,
		RepositoryName:             c.StoredName, // TODO: should be something else?
		// TODO: Is this set at a sentralized place?
		// Tags:                                 nil,
		// AWSCloudFormationDeletionPolicy:      "",
		// AWSCloudFormationUpdateReplacePolicy: "",
		// AWSCloudFormationDependsOn:           nil,
		// AWSCloudFormationMetadata:            nil,
		// AWSCloudFormationCondition:           "",
	}
}

func (c *ContainerRepository) Name() string {
	return c.StoredName
}

func (c *ContainerRepository) Ref() string {
	return cloudformation.Ref(c.Name())
}

func (c *ContainerRepository) NamedOutputs() map[string]cloudformation.Output {
	return cfn.NewValue(c.Name(), c.Ref()).NamedOutputs()
}

// New returns an initialised AWS S3 cloud formation template
// - https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-s3-bucket.html
func New(environment, imageName string) *ContainerRepository {
	return &ContainerRepository{
		StoredName: fmt.Sprintf("%s/%s", environment, imageName),
	}
}
