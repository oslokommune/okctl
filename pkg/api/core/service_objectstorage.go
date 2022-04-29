package core

import (
	"fmt"
	"io"

	"github.com/oslokommune/okctl/pkg/api"
)

type objectStorageService struct {
	provider api.ObjectStorageCloudProvider
}

func (o objectStorageService) CreateBucket(opts api.CreateBucketOpts) (string, error) {
	err := opts.Validate()
	if err != nil {
		return "", fmt.Errorf(validatingOptsErrFormat, err)
	}

	return o.provider.CreateBucket(opts)
}

func (o objectStorageService) DeleteBucket(opts api.DeleteBucketOpts) error {
	err := opts.Validate()
	if err != nil {
		return fmt.Errorf(validatingObjectStorageOptsErrFormat, err)
	}

	return o.provider.DeleteBucket(opts)
}

func (o objectStorageService) EmptyBucket(opts api.EmptyBucketOpts) error {
	err := opts.Validate()
	if err != nil {
		return fmt.Errorf(validatingObjectStorageOptsErrFormat, err)
	}

	err = o.provider.EmptyBucket(opts)
	if err != nil {
		return fmt.Errorf("calling provider's empty bucket function: %w", err)
	}

	return nil
}

func (o objectStorageService) PutObject(opts api.PutObjectOpts) error {
	err := opts.Validate()
	if err != nil {
		return fmt.Errorf(validatingObjectStorageOptsErrFormat, err)
	}

	return o.provider.PutObject(opts)
}

func (o objectStorageService) GetObject(opts api.GetObjectOpts) (io.Reader, error) {
	err := opts.Validate()
	if err != nil {
		return nil, fmt.Errorf(validatingObjectStorageOptsErrFormat, err)
	}

	return o.provider.GetObject(opts)
}

func (o objectStorageService) DeleteObject(opts api.DeleteObjectOpts) error {
	err := opts.Validate()
	if err != nil {
		return fmt.Errorf(validatingObjectStorageOptsErrFormat, err)
	}

	return o.provider.DeleteObject(opts)
}

// NewObjectStorageService returns an initialized ObjectStorageService
func NewObjectStorageService(provider api.ObjectStorageCloudProvider) api.ObjectStorageService {
	return &objectStorageService{provider: provider}
}

const validatingObjectStorageOptsErrFormat = "validating opts: %w"
