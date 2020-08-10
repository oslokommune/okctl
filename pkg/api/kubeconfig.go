package api

// KubeConfig represents a kubeconfig
type KubeConfig struct {
	Path    string
	Content string
}

// KubeConfigStore defines the storage opereations on a kubeconfig
type KubeConfigStore interface {
	CreateKubeConfig() (string, error)
	GetKubeConfig() (*KubeConfig, error)
	DeleteKubeConfig() error
}
