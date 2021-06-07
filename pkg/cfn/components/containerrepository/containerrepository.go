// Package containerrepository knows how to create an AWS ECRepository
// cloud formation resource
package containerrepository

import (
	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/awslabs/goformation/v4/cloudformation/ecr"
	"github.com/oslokommune/okctl/pkg/cfn"
)

type (
	// TagMutabilityConfiguration configures image tag mutability. If mutability is set to Immutable, we'll see an error
	// trying to push an already existing tag
	TagMutabilityConfiguration string
	// ImageScanConfiguration configures whether to scan pushed images for security vulnerabilities
	ImageScanConfiguration bool
)

//goland:noinspection GoUnusedConst
const (
	// TagMutabilityConfigurationMutable enables overwriting existing tags
	TagMutabilityConfigurationMutable = "MUTABLE"
	// TagMutabilityConfigurationImmutable disables overwriting existing tags
	TagMutabilityConfigurationImmutable = "IMMUTABLE"

	// ImageScanConfigurationOn enables image vulnerability scanning
	ImageScanConfigurationOn = true
	// ImageScanConfigurationOff disables image vulnerability scanning
	ImageScanConfigurationOff = false
)

// ContainerRepository contains state for building a cloud formation resource
type ContainerRepository struct {
	StoredName string
	ImageName  string
}

// Resource returns the cloud formation resource
func (c *ContainerRepository) Resource() cloudformation.Resource {
	return &ecr.Repository{
		ImageScanningConfiguration: &ecr.Repository_ImageScanningConfiguration{
			ScanOnPush: ImageScanConfigurationOn,
		},
		ImageTagMutability: TagMutabilityConfigurationImmutable,
		RepositoryName:     c.ImageName,
	}
}

// Name returns the logical identifier
func (c *ContainerRepository) Name() string {
	return c.StoredName
}

// Ref returns an aws intrinsic ref to this resource
func (c *ContainerRepository) Ref() string {
	return cloudformation.Ref(c.Name())
}

// NamedOutputs returns the named outputs
func (c *ContainerRepository) NamedOutputs() map[string]cloudformation.Output {
	return cfn.NewValue(c.Name(), c.Ref()).NamedOutputs()
}

// New returns an initialised AWS S3 cloud formation template
// - https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-s3-bucket.html
func New(imageName string) *ContainerRepository {
	return &ContainerRepository{
		StoredName: "ECRContainerRepository",
		ImageName:  imageName,
	}
}
