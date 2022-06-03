// Package s3bucket knows how to create an AWS S3 bucket
// cloud formation resource
package s3bucket

import (
	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/awslabs/goformation/v4/cloudformation/s3"
	"github.com/oslokommune/okctl/pkg/cfn"
)

const (
	seeAlgorithmAES256      = "AES256"
	bucketVersioningEnabled = "Enabled"
)

// S3Bucket contains the state required for building
// the cloud formation template
type S3Bucket struct {
	StoredName           string
	BucketName           string
	BlockAllPublicAccess bool
	Encrypt              bool
	Versioning           bool
}

// Resource returns the cloud formation template
func (s *S3Bucket) Resource() cloudformation.Resource {
	bucket := &s3.Bucket{
		AccessControl: "BucketOwnerFullControl",
		BucketName:    s.BucketName,
	}

	if s.BlockAllPublicAccess {
		bucket.PublicAccessBlockConfiguration = &s3.Bucket_PublicAccessBlockConfiguration{
			BlockPublicAcls:       s.BlockAllPublicAccess,
			BlockPublicPolicy:     s.BlockAllPublicAccess,
			IgnorePublicAcls:      s.BlockAllPublicAccess,
			RestrictPublicBuckets: s.BlockAllPublicAccess,
		}
	}

	if s.Encrypt {
		bucket.BucketEncryption = &s3.Bucket_BucketEncryption{
			ServerSideEncryptionConfiguration: []s3.Bucket_ServerSideEncryptionRule{
				{
					BucketKeyEnabled: true,
					ServerSideEncryptionByDefault: &s3.Bucket_ServerSideEncryptionByDefault{
						SSEAlgorithm: seeAlgorithmAES256,
					},
				},
			},
		}
	}

	if s.Versioning {
		bucket.VersioningConfiguration = &s3.Bucket_VersioningConfiguration{Status: bucketVersioningEnabled}
	}

	return bucket
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
	valueMap := cfn.NewValueMap()

	valueMap.Add(cfn.NewValue(s.Name(), s.Ref()))
	valueMap.Add(cfn.NewValue("BucketARN", cloudformation.GetAtt(s.Name(), "Arn")))

	return valueMap.NamedOutputs()
}

// New returns an initialised AWS S3 cloud formation template
// - https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-s3-bucket.html
func New(resourceName, bucketName string, encrypt bool, versioning bool) *S3Bucket {
	return &S3Bucket{
		StoredName: resourceName,
		BucketName: bucketName,
		Encrypt:    encrypt,
		Versioning: versioning,
	}
}
