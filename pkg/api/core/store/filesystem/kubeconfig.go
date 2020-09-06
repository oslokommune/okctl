package filesystem

import (
	"fmt"
	"path"

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
func (k *kubeConfig) CreateKubeConfig() (string, error) {
	exists, err := k.fs.Exists(path.Join(k.baseDir, k.kubeConfigFileName))
	if err != nil {
		return "", fmt.Errorf("failed to determine if kubeconfig exists: %w", err)
	}

	if exists {
		return path.Join(k.baseDir, k.kubeConfigFileName), nil
	}

	_, err = store.NewFileSystem(k.baseDir, k.fs).
		StoreBytes(k.kubeConfigFileName, []byte{}).
		Do()
	if err != nil {
		return "", fmt.Errorf("failed to create kubeconfig: %w", err)
	}

	return path.Join(k.baseDir, k.kubeConfigFileName), nil
}

func (k *kubeConfig) GetKubeConfig() (*api.KubeConfig, error) {
	kube := &api.KubeConfig{
		Path: path.Join(k.baseDir, k.kubeConfigFileName),
	}

	callback := func(_ string, data []byte) error {
		kube.Content = string(data)
		return nil
	}

	_, err := store.NewFileSystem(k.baseDir, k.fs).
		GetBytes(k.kubeConfigFileName, callback).
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
