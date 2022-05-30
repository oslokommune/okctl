package aws

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/cfn"
	"github.com/oslokommune/okctl/pkg/cfn/components"
	"github.com/oslokommune/okctl/pkg/s3api"
)

type objectStorageProvider struct {
	provider v1alpha1.CloudProvider
}

// CreateBucket produces and deploys necessary CFN template(s) for S3 bucket creation
func (o objectStorageProvider) CreateBucket(opts api.CreateBucketOpts) (bucketID string, err error) {
	composer := components.NewS3BucketComposer(opts.BucketName, "S3Bucket", opts.Encrypted)
	composer.BlockAllPublicAccess = opts.Private

	b := cfn.New(composer)

	stackName := cfn.NewStackNamer().S3Bucket(opts.BucketName, opts.ClusterID.ClusterName)

	template, err := b.Build()
	if err != nil {
		return "", fmt.Errorf("building CloudFormation template: %w", err)
	}

	r := cfn.NewRunner(o.provider)

	err = r.CreateIfNotExists(opts.ClusterID.ClusterName, stackName, template, nil, defaultTimeOut)
	if err != nil {
		return "", fmt.Errorf("creating CloudFormation template: %w", err)
	}

	err = r.Outputs(stackName, map[string]cfn.ProcessOutputFn{
		"BucketARN": cfn.String(&bucketID),
	})
	if err != nil {
		return "", fmt.Errorf("acquiring bucket outputs: %w", err)
	}

	return bucketID, nil
}

// DeleteBucket knows how to delete an existing S3 bucket via Cloudformation
func (o objectStorageProvider) DeleteBucket(opts api.DeleteBucketOpts) error {
	stackName := cfn.NewStackNamer().S3Bucket(opts.BucketName, opts.ClusterID.ClusterName)

	r := cfn.NewRunner(o.provider)

	err := r.Delete(stackName)
	if err != nil {
		return fmt.Errorf("deleting bucket CloudFormation template: %w", err)
	}

	return nil
}

// PutObject knows how to add content to a certain path in a bucket
func (o objectStorageProvider) PutObject(opts api.PutObjectOpts) error {
	s3 := s3api.New(o.provider)

	raw, err := io.ReadAll(opts.Content)
	if err != nil {
		return fmt.Errorf("buffering content: %w", err)
	}

	err = s3.PutObject(opts.BucketName, opts.Path, bytes.NewReader(raw))
	if err != nil {
		return fmt.Errorf("putting object into bucket: %w", err)
	}

	return nil
}

// GetObject knows how to retrieve content from a certain path in a bucket
func (o objectStorageProvider) GetObject(opts api.GetObjectOpts) (io.Reader, error) {
	s3 := s3api.New(o.provider)

	object, err := s3.GetObject(opts.BucketName, opts.Path)
	if err != nil {
		var aerr awserr.RequestFailure

		if errors.As(err, &aerr) {
			if aerr.StatusCode() == http.StatusNotFound {
				return nil, fmt.Errorf("bucket missing: %w", err)
			}
		}

		return nil, fmt.Errorf("retrieving object: %w", err)
	}

	return object, nil
}

// DeleteObject knows how to delete an object from a bucket
func (o objectStorageProvider) DeleteObject(opts api.DeleteObjectOpts) error {
	s3 := s3api.New(o.provider)

	err := s3.DeleteObject(opts.BucketName, opts.Path)
	if err != nil {
		return fmt.Errorf("deleting object: %w", err)
	}

	return nil
}

// EmptyBucket knows how to purge a bucket for content
func (o objectStorageProvider) EmptyBucket(opts api.EmptyBucketOpts) error {
	s3 := s3api.New(o.provider)

	err := s3.EmptyBucket(opts.BucketName)
	if err != nil {
		if errors.Is(err, s3api.ErrBucketDoesNotExist) {
			return api.ErrObjectStorageBucketNotExist
		}

		return fmt.Errorf("emptying bucket: %w", err)
	}

	return nil
}

// NewObjectStorageCloudProvider initializes an Object Storage Provider
func NewObjectStorageCloudProvider(provider v1alpha1.CloudProvider) api.ObjectStorageCloudProvider {
	return &objectStorageProvider{provider: provider}
}
