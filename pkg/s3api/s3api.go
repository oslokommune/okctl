// Package s3api provides some convenience functions for interacting
// with the AWS S3 API
package s3api

import (
	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
)

// S3API contains the state required for interacting
// with the AWS S3 API
type S3API struct {
	provider v1alpha1.CloudProvider
}

// PutObject stores the provided in the given bucket/key
func (a *S3API) PutObject(bucket, key string, body io.ReadSeeker) error {
	_, err := a.provider.S3().PutObject(&s3.PutObjectInput{
		Body:   body,
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return err
	}

	return nil
}

// DeleteObject removes the object at bucket/key
func (a *S3API) DeleteObject(bucket, key string) error {
	_, err := a.provider.S3().DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return err
	}

	return nil
}

// New returns an initialised AWS S3 API client
func New(provider v1alpha1.CloudProvider) *S3API {
	return &S3API{
		provider: provider,
	}
}
