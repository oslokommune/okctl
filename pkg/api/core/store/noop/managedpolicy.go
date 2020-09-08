package noop

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/api"
)

type managedPolicyStore struct{}

func (m *managedPolicyStore) SaveExternalSecretsPolicy(policy *api.ManagedPolicy) error {
	return nil
}

func (m *managedPolicyStore) GetExternalSecretsPolicy() (*api.ManagedPolicy, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *managedPolicyStore) SaveAlbIngressControllerPolicy(policy *api.ManagedPolicy) error {
	return nil
}

func (m *managedPolicyStore) GetAlbIngressControllerPolicy() (*api.ManagedPolicy, error) {
	return nil, fmt.Errorf("not implemented")
}

func (m *managedPolicyStore) SaveExternalDNSPolicy(policy *api.ManagedPolicy) error {
	return nil
}

func (m *managedPolicyStore) GetExternalDNSPolicy() (*api.ManagedPolicy, error) {
	return nil, fmt.Errorf("not implemented")
}

// NewManagedPolicyStore returns a no operation store
func NewManagedPolicyStore() api.ManagedPolicyStore {
	return &managedPolicyStore{}
}
