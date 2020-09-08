package noop

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/api"
)

type serviceAccountStore struct{}

func (s *serviceAccountStore) SaveExternalSecretsServiceAccount(account *api.ServiceAccount) error {
	return nil
}

func (s *serviceAccountStore) GetExternalSecretsServiceAccount() (*api.ServiceAccount, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *serviceAccountStore) SaveAlbIngressControllerServiceAccount(account *api.ServiceAccount) error {
	return nil
}

func (s *serviceAccountStore) GetAlbIngressControllerServiceAccount() (*api.ServiceAccount, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *serviceAccountStore) SaveExternalDNSServiceAccount(account *api.ServiceAccount) error {
	return nil
}

func (s *serviceAccountStore) GetExternalDNSServiceAccount() (*api.ServiceAccount, error) {
	return nil, fmt.Errorf("not implemented")
}

// NewServiceAccountStore returns a no operation store
func NewServiceAccountStore() api.ServiceAccountStore {
	return &serviceAccountStore{}
}
