package filesystem

import (
	"fmt"
	"path"

	"github.com/oslokommune/okctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"

	"github.com/oslokommune/okctl/pkg/kubeconfig"

	"github.com/oslokommune/okctl/pkg/client/store"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/spf13/afero"
)

type kubeConfig struct {
	kubeConfigFileName    string
	kubeConfigBaseDir     string
	clusterConfigFileName string
	clusterConfigBaseDir  string

	fs       *afero.Afero
	provider v1alpha1.CloudProvider
}

// This is not good, we need to rewrite this, together with
// much of the API
func (k *kubeConfig) SaveKubeConfig(config *kubeconfig.Config) error {
	cfg, err := config.Bytes()
	if err != nil {
		return fmt.Errorf("creating kubeconfig: %w", err)
	}

	_, err = store.NewFileSystem(k.kubeConfigBaseDir, k.fs).
		StoreBytes(k.kubeConfigFileName, cfg).
		Do()
	if err != nil {
		return fmt.Errorf("failed to create kubeconfig: %w", err)
	}

	return nil
}

func (k *kubeConfig) GetKubeConfig() (*api.KubeConfig, error) {
	// We create the kubeconfig on new every time to avoid
	// situations where we get a stale/old kubeconfig.
	// This does feel a little bit awkward.
	cfg, err := k.createKubeConfig()
	if err != nil {
		return nil, err
	}

	content, err := cfg.Bytes()
	if err != nil {
		return nil, err
	}

	c := &api.KubeConfig{
		Path:    path.Join(k.kubeConfigBaseDir, k.kubeConfigFileName),
		Content: string(content),
	}

	err = k.SaveKubeConfig(cfg)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (k *kubeConfig) createKubeConfig() (*kubeconfig.Config, error) {
	clusterConfig := &v1alpha5.ClusterConfig{}

	_, err := store.NewFileSystem(k.clusterConfigBaseDir, k.fs).
		GetStruct(k.clusterConfigFileName, &clusterConfig, store.FromYAML()).
		Do()
	if err != nil {
		return nil, err
	}

	return kubeconfig.New(clusterConfig, k.provider).Get()
}

func (k *kubeConfig) DeleteKubeConfig() error {
	_, err := store.NewFileSystem(k.kubeConfigBaseDir, k.fs).
		Remove(k.kubeConfigFileName).
		Do()
	if err != nil {
		return fmt.Errorf("failed to remove kubeconfig: %w", err)
	}

	return nil
}

// NewKubeConfigStore returns an initialised kubeconfig store
func NewKubeConfigStore(
	provider v1alpha1.CloudProvider,
	kubeConfigFileName, kubeConfigBaseDir, clusterConfigFileName, clusterConfigBaseDir string,
	fs *afero.Afero,
) api.KubeConfigStore {
	return &kubeConfig{
		kubeConfigFileName:    kubeConfigFileName,
		kubeConfigBaseDir:     kubeConfigBaseDir,
		clusterConfigFileName: clusterConfigFileName,
		clusterConfigBaseDir:  clusterConfigBaseDir,
		fs:                    fs,
		provider:              provider,
	}
}
