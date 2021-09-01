package filesystem

import (
	"fmt"
	"github.com/oslokommune/okctl/pkg/config/constant"
	"path"

	"github.com/oslokommune/okctl/pkg/client"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"

	"github.com/oslokommune/okctl/pkg/kubeconfig"

	"github.com/oslokommune/okctl/pkg/client/store"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/spf13/afero"
)

type kubeConfig struct {
	kubeConfigFileName string
	kubeConfigBaseDir  string

	fs           *afero.Afero
	provider     v1alpha1.CloudProvider
	clusterState client.ClusterState
}

func (k *kubeConfig) SaveKubeConfig(config *kubeconfig.Config) error {
	cfg, err := config.Bytes()
	if err != nil {
		return fmt.Errorf(constant.CreateKubeConfigError, err)
	}

	_, err = store.NewFileSystem(k.kubeConfigBaseDir, k.fs).
		StoreBytes(k.kubeConfigFileName, cfg).
		Do()
	if err != nil {
		return fmt.Errorf(constant.StoreKubecofigError, err)
	}

	return nil
}

func (k *kubeConfig) GetKubeConfig(clusterName string) (*api.KubeConfig, error) {
	// We create the kubeconfig on new every time to avoid
	// situations where we get a stale/old kubeconfig.
	// This does feel a little bit awkward.
	cluster, err := k.clusterState.GetCluster(clusterName)
	if err != nil {
		return nil, err
	}

	cfg, err := kubeconfig.New(cluster.Config, k.provider).Get()
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

func (k *kubeConfig) DeleteKubeConfig() error {
	_, err := store.NewFileSystem(k.kubeConfigBaseDir, k.fs).
		Remove(k.kubeConfigFileName).
		Do()
	if err != nil {
		return fmt.Errorf(constant.RemoveKubeconfigError, err)
	}

	return nil
}

// NewKubeConfigStore returns an initialised kubeconfig store
func NewKubeConfigStore(
	provider v1alpha1.CloudProvider,
	kubeConfigFileName, kubeConfigBaseDir string,
	clusterState client.ClusterState,
	fs *afero.Afero,
) api.KubeConfigStore {
	return &kubeConfig{
		kubeConfigFileName: kubeConfigFileName,
		kubeConfigBaseDir:  kubeConfigBaseDir,
		fs:                 fs,
		provider:           provider,
		clusterState:       clusterState,
	}
}
