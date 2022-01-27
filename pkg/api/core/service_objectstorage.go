package core

import (
	stderrors "errors"
	"io"

	"github.com/oslokommune/okctl/pkg/s3api"

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

	err = o.provider.EmptyBucket(opts)
	if err != nil {
		kind := errors.Internal

		if stderrors.Is(err, s3api.ErrBucketDoesNotExist) {
			kind = errors.NotExist
		}

		return errors.E(err, "calling provider's empty bucket function", kind)
	}

	return nil
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
