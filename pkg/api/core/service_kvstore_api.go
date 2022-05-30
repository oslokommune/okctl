package core

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/api"
)

type kvStoreService struct {
	cloudProvider api.KeyValueStoreCloudProvider
}

func (k *kvStoreService) CreateStore(opts api.CreateStoreOpts) error {
	err := opts.Validate()
	if err != nil {
		return fmt.Errorf(validatingOptsErrFormat, err)
	}

	return k.cloudProvider.CreateStore(opts)
}

func (k *kvStoreService) DeleteStore(opts api.DeleteStoreOpts) error {
	err := opts.Validate()
	if err != nil {
		return fmt.Errorf(validatingOptsErrFormat, err)
	}

	return k.cloudProvider.DeleteStore(opts)
}

func (k *kvStoreService) ListStores() ([]string, error) {
	return k.cloudProvider.ListStores()
}

func (k *kvStoreService) InsertItem(opts api.InsertItemOpts) error {
	err := opts.Validate()
	if err != nil {
		return fmt.Errorf(validatingOptsErrFormat, err)
	}

	return k.cloudProvider.InsertItem(opts)
}

func (k *kvStoreService) GetString(opts api.GetStringOpts) (string, error) {
	err := opts.Validate()
	if err != nil {
		return "", fmt.Errorf(validatingOptsErrFormat, err)
	}

	return k.cloudProvider.GetString(opts)
}

func (k *kvStoreService) RemoveItem(opts api.DeleteItemOpts) error {
	err := opts.Validate()
	if err != nil {
		return fmt.Errorf(validatingOptsErrFormat, err)
	}

	return k.cloudProvider.RemoveItem(opts)
}

// NewKeyValueStoreService returns an initialized KeyValueStoreService
func NewKeyValueStoreService(provider api.KeyValueStoreCloudProvider) api.KeyValueStoreService {
	return &kvStoreService{cloudProvider: provider}
}

const validatingOptsErrFormat = "validating opts: %w"
