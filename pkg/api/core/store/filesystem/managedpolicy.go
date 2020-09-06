package filesystem

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client/store"
	"github.com/spf13/afero"
)

type managedPolicy struct {
	externalSecrets      Paths
	albIngressController Paths
	externalDNS          Paths
	fs                   *afero.Afero
}

// ManagedPolicy contains the state that is stored to
// and retrieved from the filesystem
type ManagedPolicy struct {
	ID        api.ID
	StackName string
	PolicyARN string
}

func (m *managedPolicy) SaveExternalDNSPolicy(policy *api.ManagedPolicy) error {
	return m.savePolicy(m.externalDNS, policy)
}

func (m *managedPolicy) GetExternalDNSPolicy() (*api.ManagedPolicy, error) {
	return m.getPolicy(m.externalDNS)
}

func (m *managedPolicy) SaveAlbIngressControllerPolicy(policy *api.ManagedPolicy) error {
	return m.savePolicy(m.albIngressController, policy)
}

func (m *managedPolicy) GetAlbIngressControllerPolicy() (*api.ManagedPolicy, error) {
	return m.getPolicy(m.albIngressController)
}

func (m *managedPolicy) SaveExternalSecretsPolicy(policy *api.ManagedPolicy) error {
	return m.savePolicy(m.externalSecrets, policy)
}

func (m *managedPolicy) GetExternalSecretsPolicy() (*api.ManagedPolicy, error) {
	return m.getPolicy(m.externalSecrets)
}

func (m *managedPolicy) savePolicy(paths Paths, policy *api.ManagedPolicy) error {
	p := &ManagedPolicy{
		ID:        policy.ID,
		StackName: policy.StackName,
		PolicyARN: policy.PolicyARN,
	}

	_, err := store.NewFileSystem(paths.BaseDir, m.fs).
		StoreStruct(paths.OutputFile, &p, store.ToJSON()).
		StoreBytes(paths.CloudFormationFile, policy.CloudFormationTemplate).
		Do()
	if err != nil {
		return fmt.Errorf("failed to store policy: %w", err)
	}

	return nil
}

func (m *managedPolicy) getPolicy(paths Paths) (*api.ManagedPolicy, error) {
	var template []byte

	callback := func(_ string, data []byte) error {
		template = data
		return nil
	}

	p := &ManagedPolicy{}

	_, err := store.NewFileSystem(paths.BaseDir, m.fs).
		GetStruct(paths.OutputFile, p, store.FromJSON()).
		GetBytes(paths.CloudFormationFile, callback).
		Do()
	if err != nil {
		return nil, fmt.Errorf("failed to get managed policy: %w", err)
	}

	return &api.ManagedPolicy{
		ID:                     p.ID,
		StackName:              p.StackName,
		PolicyARN:              p.PolicyARN,
		CloudFormationTemplate: template,
	}, nil
}

// NewManagedPolicyStore returns an initialised managed policy store
func NewManagedPolicyStore(externalSecrets, albIngressController, externalDNS Paths, fs *afero.Afero) api.ManagedPolicyStore {
	return &managedPolicy{
		externalSecrets:      externalSecrets,
		albIngressController: albIngressController,
		externalDNS:          externalDNS,
		fs:                   fs,
	}
}
