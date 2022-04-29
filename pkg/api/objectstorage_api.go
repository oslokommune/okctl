package api

import (
	"errors"
	"io"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Bucketer defines functionality related to buckets
type Bucketer interface {
	// CreateBucket knows how to create a bucket
	CreateBucket(CreateBucketOpts) (bucketID string, err error)
	// DeleteBucket knows how to delete a bucket
	DeleteBucket(DeleteBucketOpts) error
	// EmptyBucket knows how to remove all content from a bucket
	EmptyBucket(EmptyBucketOpts) error
}

// Objecter defines functionality related to objects
type Objecter interface {
	// PutObject knows how to insert an object into a bucket
	PutObject(PutObjectOpts) error
	// GetObject knows how to retrieve an object from a bucket
	GetObject(GetObjectOpts) (io.Reader, error)
	// DeleteObject knows how to remove an object from a bucket
	DeleteObject(DeleteObjectOpts) error
}

// CreateBucketOpts defines necessary data for bucket creation
type CreateBucketOpts struct {
	// ClusterID identifies which cluster to do the operation in context of
	ClusterID ID
	// BucketName defines the name of the bucket to create
	BucketName string
	// Determines if the bucket should be protected from public access. N.B.: leaving this as false does not mean public
	// access is allowed nor enabled
	Private bool
}

// Validate ensures correct and required data
func (c CreateBucketOpts) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.BucketName, validation.Required),
		validation.Field(&c.ClusterID, validation.Required),
	)
}

// DeleteBucketOpts defines necessary data for bucket deletion
type DeleteBucketOpts struct {
	// ClusterID identifies which cluster to do the operation in context of
	ClusterID ID
	// BucketName defines the name of the bucket to delete
	BucketName string
}

// Validate ensures correct and required data
func (c DeleteBucketOpts) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.BucketName, validation.Required),
		validation.Field(&c.ClusterID, validation.Required),
	)
}

// EmptyBucketOpts defines necessary data for
type EmptyBucketOpts struct {
	// BucketName defines the name of the bucket to empty
	BucketName string
}

// Validate ensures correct and required data
func (c EmptyBucketOpts) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.BucketName, validation.Required),
	)
}

// PutObjectOpts defines necessary data to store an object
type PutObjectOpts struct {
	// ClusterID identifies which cluster to do the operation in context of
	ClusterID ID
	// BucketName defines the name of the bucket to do the operation in
	BucketName string
	// Path defines the path to the object
	Path string
	// Content defines the content of the object
	Content io.Reader
}

// Validate ensures correct and required data
func (c PutObjectOpts) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.BucketName, validation.Required),
		validation.Field(&c.Path, validation.Required),
	)
}

// GetObjectOpts defines necessary data to retrieve an object
type GetObjectOpts struct {
	// ClusterID identifies which cluster to do the operation in context of
	ClusterID ID
	// BucketName defines the name of the bucket to do the operation in
	BucketName string
	// Path defines the path of the object to retrieve
	Path string
}

// Validate ensures correct and required data
func (c GetObjectOpts) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.BucketName, validation.Required),
		validation.Field(&c.Path, validation.Required),
	)
}

// DeleteObjectOpts defines necessary data to delete an object
type DeleteObjectOpts struct {
	// BucketName defines the name of the bucket to do the operation in
	BucketName string
	// Path defines the path of the objcet to delete
	Path string
}

// Validate ensures correct and required data
func (c DeleteObjectOpts) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.BucketName, validation.Required),
		validation.Field(&c.Path, validation.Required),
	)
}

// ObjectStorageService defines necessary operations for object storage
type ObjectStorageService interface {
	Bucketer
	Objecter
}

// ObjectStorageCloudProvider provides the cloud provider layer
type ObjectStorageCloudProvider interface {
	Bucketer
	Objecter
}

// ErrObjectStorageBucketNotExist indicates that a specified bucket does not exist
var ErrObjectStorageBucketNotExist = errors.New("bucket does not exist")
