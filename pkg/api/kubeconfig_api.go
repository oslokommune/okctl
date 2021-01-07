package api

import "github.com/oslokommune/okctl/pkg/kubeconfig"

// KubeConfig represents a kubeconfig
type KubeConfig struct {
	Path    string
	Content string
}

// KubeConfigStore defines the storage operations on a kubeconfig
type KubeConfigStore interface {
	SaveKubeConfig(cfg *kubeconfig.Config) error
	GetKubeConfig() (*KubeConfig, error)
	DeleteKubeConfig() error
}
