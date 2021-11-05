package core

import (
	"io"

	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/api"
)

type objectStorageService struct {
	provider api.ObjectStorageCloudProvider
}

func (o objectStorageService) CreateBucket(opts api.CreateBucketOpts) (string, error) {
	err := opts.Validate()
	if err != nil {
		return "", errors.E(err, validatingObjectStorageOptsErrFormat)
	}

	return o.provider.CreateBucket(opts)
}

func (o objectStorageService) DeleteBucket(opts api.DeleteBucketOpts) error {
	err := opts.Validate()
	if err != nil {
		return errors.E(err, validatingObjectStorageOptsErrFormat)
	}

	return o.provider.DeleteBucket(opts)
}

func (o objectStorageService) EmptyBucket(opts api.EmptyBucketOpts) error {
	err := opts.Validate()
	if err != nil {
		return errors.E(err, validatingObjectStorageOptsErrFormat)
	}

	return o.provider.EmptyBucket(opts)
}

func (o objectStorageService) PutObject(opts api.PutObjectOpts) error {
	err := opts.Validate()
	if err != nil {
		return errors.E(err, validatingObjectStorageOptsErrFormat)
	}

	return o.provider.PutObject(opts)
}

func (o objectStorageService) GetObject(opts api.GetObjectOpts) (io.Reader, error) {
	err := opts.Validate()
	if err != nil {
		return nil, errors.E(err, validatingObjectStorageOptsErrFormat)
	}

	return o.provider.GetObject(opts)
}

func (o objectStorageService) DeleteObject(opts api.DeleteObjectOpts) error {
	err := opts.Validate()
	if err != nil {
		return errors.E(err, validatingObjectStorageOptsErrFormat)
	}

	return o.provider.DeleteObject(opts)
}

// NewObjectStorageService returns an initialized ObjectStorageService
func NewObjectStorageService(provider api.ObjectStorageCloudProvider) api.ObjectStorageService {
	return &objectStorageService{provider: provider}
}

const validatingObjectStorageOptsErrFormat = "validating opts"
