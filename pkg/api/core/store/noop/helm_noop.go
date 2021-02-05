package noop

import "github.com/oslokommune/okctl/pkg/api"

type helmStore struct{}

func (h *helmStore) SaveExternalSecretsHelmChart(helm *api.Helm) error {
	return nil
}

func (h *helmStore) SaveAlbIngressControllerHelmChar(helm *api.Helm) error {
	return nil
}

func (h *helmStore) SaveAWSLoadBalancerControllerHelmChar(helm *api.Helm) error {
	return nil
}

func (h *helmStore) SaveArgoCD(helm *api.Helm) error {
	return nil
}

// NewHelmStore returns a no operation store
func NewHelmStore() api.HelmStore {
	return &helmStore{}
}
