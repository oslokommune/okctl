// Package s3api provides some convenience functions for interacting
// with the AWS S3 API
package s3api

import (
	"io"

	"github.com/mishudark/errors"

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

// GetObject retrieves the object at bucket/key
func (a *S3API) GetObject(bucket, key string) (io.ReadCloser, error) {
	result, err := a.provider.S3().GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}

	return result.Body, nil
}

// EmptyBucket deletes all objects and versions from a bucket, making the bucket itself deletable
func (a *S3API) EmptyBucket(bucketName string) error {
	err := a.deleteAllObjects(bucketName)
	if err != nil {
		return errors.E(err, "deleting all objects")
	}

	err = a.deleteAllVersions(bucketName)
	if err != nil {
		return errors.E(err, "deleting all versions")
	}

	return nil
}

func (a *S3API) deleteAllObjects(bucketName string) error {
	var (
		result     *s3.ListObjectsOutput
		nextMarker *string
		err        error
	)

	for {
		result, err = a.provider.S3().ListObjects(&s3.ListObjectsInput{
			Bucket: aws.String(bucketName),
			Marker: nextMarker,
		})
		if err != nil {
			return err
		}

		_, err = a.provider.S3().DeleteObjects(&s3.DeleteObjectsInput{
			Bucket: aws.String(bucketName),
			Delete: &s3.Delete{Objects: objectsAsIdentifiers(result.Contents)},
		})
		if err != nil {
			return err
		}

		if !*result.IsTruncated {
			break
		}

		nextMarker = result.Contents[len(result.Contents)-1].Key
	}

	return nil
}

func (a *S3API) deleteAllVersions(bucketName string) error {
	var (
		result     *s3.ListObjectVersionsOutput
		nextMarker *string
		err        error
	)

	for {
		result, err = a.provider.S3().ListObjectVersions(&s3.ListObjectVersionsInput{
			Bucket:          aws.String(bucketName),
			VersionIdMarker: nextMarker,
		})
		if err != nil {
			return errors.E(err, "listing object versions")
		}

		_, err = a.provider.S3().DeleteObjects(&s3.DeleteObjectsInput{
			Bucket: aws.String(bucketName),
			Delete: &s3.Delete{Objects: objectVersionsAsIdentifiers(result.Versions)},
		})
		if err != nil {
			return errors.E(err, "deleting versions")
		}

		if !*result.IsTruncated {
			break
		}

		nextMarker = result.Versions[len(result.Versions)-1].VersionId
	}

	return nil
}

func objectsAsIdentifiers(objects []*s3.Object) []*s3.ObjectIdentifier {
	identifiers := make([]*s3.ObjectIdentifier, len(objects))

	for index, object := range objects {
		identifiers[index] = &s3.ObjectIdentifier{Key: object.Key}
	}

	return identifiers
}

func objectVersionsAsIdentifiers(versions []*s3.ObjectVersion) []*s3.ObjectIdentifier {
	identifiers := make([]*s3.ObjectIdentifier, len(versions))

	for index, version := range versions {
		identifiers[index] = &s3.ObjectIdentifier{
			Key:       version.Key,
			VersionId: version.VersionId,
		}
	}

	return identifiers
}

// New returns an initialised AWS S3 API client
func New(provider v1alpha1.CloudProvider) *S3API {
	return &S3API{
		provider: provider,
	}
}
