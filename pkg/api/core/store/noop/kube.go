package noop

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/api"
)

type kubeStore struct{}

func (k *kubeStore) SaveExternalDNSKubeDeployment(kube *api.Kube) error {
	return nil
}

func (k *kubeStore) GetExternalDNSKubeDeployment() (*api.Kube, error) {
	return nil, fmt.Errorf("not implemented")
}

func (k *kubeStore) SaveExternalSecrets(kube *api.Kube) error {
	return nil
}

// NewKubeStore returns a no operation store
func NewKubeStore() api.KubeStore {
	return &kubeStore{}
}
