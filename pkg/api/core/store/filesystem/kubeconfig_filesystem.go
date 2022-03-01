package filesystem

import (
	"fmt"
	"path"

	"github.com/oslokommune/okctl/pkg/binaries"
	"github.com/oslokommune/okctl/pkg/binaries/run/awsiamauthenticator"

	"github.com/oslokommune/okctl/pkg/config/constant"

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

	fs               *afero.Afero
	provider         v1alpha1.CloudProvider
	clusterState     client.ClusterState
	binariesProvider binaries.Provider
}

func (k *kubeConfig) SaveKubeConfig(config *kubeconfig.Config) error {
	cfg, err := config.Bytes()
	if err != nil {
		return fmt.Errorf("creating kubeconfig: %w", err)
	}

	_, err = store.NewFileSystem(k.kubeConfigBaseDir, k.fs).
		StoreBytes(k.kubeConfigFileName, cfg, store.WithFilePermissionsMode(constant.DefaultClusterKubePermission)).
		Do()
	if err != nil {
		return fmt.Errorf("storing kubeconfig: %w", err)
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

	awsIAMAuthenticatorProvider, err := k.binariesProvider.AwsIamAuthenticator(awsiamauthenticator.Version)
	if err != nil {
		return nil, fmt.Errorf("acquiring aws-iam-authenticator provider: %w", err)
	}

	cfg, err := kubeconfig.New(awsIAMAuthenticatorProvider.BinaryPath, cluster.Config, k.provider).Get()
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
		return fmt.Errorf("removing kubeconfig: %w", err)
	}

	return nil
}

// NewKubeConfigStore returns an initialised kubeconfig store
func NewKubeConfigStore(
	provider v1alpha1.CloudProvider,
	binariesProvider binaries.Provider,
	kubeConfigFileName,
	kubeConfigBaseDir string,
	clusterState client.ClusterState,
	fs *afero.Afero,
) api.KubeConfigStore {
	return &kubeConfig{
		kubeConfigFileName: kubeConfigFileName,
		kubeConfigBaseDir:  kubeConfigBaseDir,
		fs:                 fs,
		provider:           provider,
		binariesProvider:   binariesProvider,
		clusterState:       clusterState,
	}
}
