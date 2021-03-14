// Package s3bucket knows how to create an AWS S3 bucket
// cloud formation resource
package s3bucket

import (
	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/awslabs/goformation/v4/cloudformation/s3"
	"github.com/oslokommune/okctl/pkg/cfn"
)

// S3Bucket contains the state required for building
// the cloud formation template
type S3Bucket struct {
	StoredName string
	BucketName string
}

// Resource returns the cloud formation template
func (s *S3Bucket) Resource() cloudformation.Resource {
	return &s3.Bucket{
		AccessControl: "BucketOwnerFullControl",
		BucketName:    s.BucketName,
	}
}

// Name returns the name of the resource
func (s *S3Bucket) Name() string {
	return s.StoredName
}

// Ref returns an AWS intrinsic ref to the resource
func (s *S3Bucket) Ref() string {
	return cloudformation.Ref(s.Name())
}

// NamedOutputs returns the named outputs
func (s *S3Bucket) NamedOutputs() map[string]cloudformation.Output {
	return cfn.NewValue(s.Name(), s.Ref()).NamedOutputs()
}

// New returns an initialised AWS S3 cloud formation template
// - https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-s3-bucket.html
func New(resourceName, bucketName string) *S3Bucket {
	return &S3Bucket{
		StoredName: resourceName,
		BucketName: bucketName,
	}
}
