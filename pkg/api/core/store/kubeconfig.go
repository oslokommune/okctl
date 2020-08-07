package store

import (
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/storage/state"
)

type kubeConfig struct {
	provider state.PersisterProvider
}

func (k *kubeConfig) CreateKubeConfig() (string, error) {
	err := k.provider.Application().WriteToDefault("kubeconfig", []byte{})
	if err != nil {
		return "", err
	}

	path, err := k.provider.Application().GetDefaultPath("kubeconfig")
	if err != nil {
		return "", err
	}

	return path, nil
}

func (k *kubeConfig) GetKubeConfig() (*api.KubeConfig, error) {
	path, err := k.provider.Application().GetDefaultPath("kubeconfig")
	if err != nil {
		return nil, err
	}

	data, err := k.provider.Application().ReadFromDefault("kubeconfig")
	if err != nil {
		return nil, err
	}

	return &api.KubeConfig{
		Path:    path,
		Content: string(data),
	}, nil
}

func (k *kubeConfig) DeleteKubeConfig() error {
	return k.provider.Application().DeleteDefault("kubeconfig")
}

// NewKubeConfigStore returns an initialised kubeconfig store
func NewKubeConfigStore(provider state.PersisterProvider) api.KubeConfigStore {
	return &kubeConfig{
		provider: provider,
	}
}
