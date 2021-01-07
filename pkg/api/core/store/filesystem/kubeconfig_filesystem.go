package filesystem

import (
	"fmt"
	"path"

	"github.com/oslokommune/okctl/pkg/kubeconfig"

	"github.com/oslokommune/okctl/pkg/client/store"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/spf13/afero"
)

type kubeConfig struct {
	kubeConfigFileName string
	baseDir            string
	fs                 *afero.Afero
}

// This is not good, we need to rewrite this, together with
// much of the API
func (k *kubeConfig) SaveKubeConfig(config *kubeconfig.Config) error {
	cfg, err := config.Bytes()
	if err != nil {
		return fmt.Errorf("creating kubeconfig: %w", err)
	}

	_, err = store.NewFileSystem(k.baseDir, k.fs).
		StoreBytes(k.kubeConfigFileName, cfg).
		Do()
	if err != nil {
		return fmt.Errorf("failed to create kubeconfig: %w", err)
	}

	return nil
}

func (k *kubeConfig) GetKubeConfig() (*api.KubeConfig, error) {
	kube := &api.KubeConfig{
		Path: path.Join(k.baseDir, k.kubeConfigFileName),
	}

	_, err := store.NewFileSystem(k.baseDir, k.fs).
		GetBytes(k.kubeConfigFileName, func(_ string, data []byte) {
			kube.Content = string(data)
		}).
		Do()
	if err != nil {
		return nil, fmt.Errorf("failed to get kubeconfig: %w", err)
	}

	return kube, nil
}

func (k *kubeConfig) DeleteKubeConfig() error {
	_, err := store.NewFileSystem(k.baseDir, k.fs).
		Remove(k.kubeConfigFileName).
		Do()
	if err != nil {
		return fmt.Errorf("failed to remove kubeconfig: %w", err)
	}

	return nil
}

// NewKubeConfigStore returns an initialised kubeconfig store
func NewKubeConfigStore(kubeConfigFileName, baseDir string, fs *afero.Afero) api.KubeConfigStore {
	return &kubeConfig{
		kubeConfigFileName: kubeConfigFileName,
		baseDir:            baseDir,
		fs:                 fs,
	}
}
