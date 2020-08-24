package filesystem

import (
	"encoding/json"
	"fmt"
	"path"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/spf13/afero"
)

type managedPolicy struct {
	externalSecrets      Paths
	albIngressController Paths
	fs                   *afero.Afero
}

// ManagedPolicy contains the state that is stored to
// and retrieved from the filesystem
type ManagedPolicy struct {
	StackName   string
	Repository  string
	Environment string
	PolicyARN   string
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
		StackName:   policy.StackName,
		Repository:  policy.Repository,
		Environment: policy.Environment,
		PolicyARN:   policy.PolicyARN,
	}

	data, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("failed to marshal policy: %w", err)
	}

	err = m.fs.MkdirAll(paths.BaseDir, 0744)
	if err != nil {
		return fmt.Errorf("failed to create policy directory: %w", err)
	}

	err = m.fs.WriteFile(path.Join(paths.BaseDir, paths.OutputFile), data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write policy: %w", err)
	}

	err = m.fs.WriteFile(path.Join(paths.BaseDir, paths.CloudFormationFile), policy.CloudFormationTemplate, 0644)
	if err != nil {
		return fmt.Errorf("failed to write cloud formation template: %w", err)
	}

	return nil
}

func (m *managedPolicy) getPolicy(paths Paths) (*api.ManagedPolicy, error) {
	data, err := m.fs.ReadFile(path.Join(paths.BaseDir, paths.OutputFile))
	if err != nil {
		return nil, fmt.Errorf("failed to read policy file: %w", err)
	}

	p := &ManagedPolicy{}

	err = json.Unmarshal(data, p)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal policy: %w", err)
	}

	policy := &api.ManagedPolicy{
		StackName:   p.StackName,
		Repository:  p.Repository,
		Environment: p.Environment,
		PolicyARN:   p.PolicyARN,
	}

	template, err := m.fs.ReadFile(path.Join(paths.BaseDir, paths.CloudFormationFile))
	if err != nil {
		return nil, fmt.Errorf("failed to read cloud formation template: %w", err)
	}

	policy.CloudFormationTemplate = template

	return policy, nil
}

// NewManagedPolicyStore returns an initialised managed policy store
func NewManagedPolicyStore(externalSecrets Paths, albIngressController Paths, fs *afero.Afero) api.ManagedPolicyStore {
	return &managedPolicy{
		externalSecrets:      externalSecrets,
		albIngressController: albIngressController,
		fs:                   fs,
	}
}
