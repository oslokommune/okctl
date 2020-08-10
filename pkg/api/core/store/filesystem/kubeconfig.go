package filesystem

import (
	"path"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/spf13/afero"
)

type kubeConfig struct {
	kubeConfigFileName string
	baseDir            string
	fileSystem         *afero.Afero
}

func (k *kubeConfig) CreateKubeConfig() (string, error) {
	err := k.fileSystem.MkdirAll(k.baseDir, 0744)
	if err != nil {
		return "", err
	}

	err = k.fileSystem.WriteFile(path.Join(k.baseDir, k.kubeConfigFileName), []byte{}, 0644)
	if err != nil {
		return "", err
	}

	return path.Join(k.baseDir, k.kubeConfigFileName), nil
}

func (k *kubeConfig) GetKubeConfig() (*api.KubeConfig, error) {
	data, err := k.fileSystem.ReadFile(path.Join(k.baseDir, k.kubeConfigFileName))
	if err != nil {
		return nil, err
	}

	return &api.KubeConfig{
		Path:    path.Join(k.baseDir, k.kubeConfigFileName),
		Content: string(data),
	}, nil
}

func (k *kubeConfig) DeleteKubeConfig() error {
	return k.fileSystem.Remove(path.Join(k.baseDir, k.kubeConfigFileName))
}

// NewKubeConfigStore returns an initialised kubeconfig store
func NewKubeConfigStore(kubeConfigFileName, baseDir string, fileSystem *afero.Afero) api.KubeConfigStore {
	return &kubeConfig{
		kubeConfigFileName: kubeConfigFileName,
		baseDir:            baseDir,
		fileSystem:         fileSystem,
	}
}
