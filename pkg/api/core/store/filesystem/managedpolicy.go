package filesystem

import (
	"encoding/json"
	"fmt"
	"path"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/spf13/afero"
)

type managedPolicy struct {
	externalSecretsFile               string
	externalSecretsCloudFormationFile string
	externalSecretsBaseDir            string
	fs                                *afero.Afero
}

// ManagedPolicy contains the state that is stored to
// and retrieved from the filesystem
type ManagedPolicy struct {
	StackName   string
	Repository  string
	Environment string
	PolicyARN   string
}

func (m *managedPolicy) SaveExternalSecretsPolicy(policy *api.ManagedPolicy) error {
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

	err = m.fs.MkdirAll(m.externalSecretsBaseDir, 0744)
	if err != nil {
		return fmt.Errorf("failed to create policy directory: %w", err)
	}

	err = m.fs.WriteFile(path.Join(m.externalSecretsBaseDir, m.externalSecretsFile), data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write policy: %w", err)
	}

	err = m.fs.WriteFile(path.Join(m.externalSecretsBaseDir, m.externalSecretsCloudFormationFile), policy.CloudFormationTemplate, 0644)
	if err != nil {
		return fmt.Errorf("failed to write cloud formation template: %w", err)
	}

	return nil
}

func (m *managedPolicy) GetExternalSecretsPolicy() (*api.ManagedPolicy, error) {
	data, err := m.fs.ReadFile(path.Join(m.externalSecretsBaseDir, m.externalSecretsFile))
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

	template, err := m.fs.ReadFile(path.Join(m.externalSecretsBaseDir, m.externalSecretsCloudFormationFile))
	if err != nil {
		return nil, fmt.Errorf("failed to read cloud formation template: %w", err)
	}

	policy.CloudFormationTemplate = template

	return policy, nil
}

// NewManagedPolicyStore returns an initialised managed policy store
func NewManagedPolicyStore(externalSecretsFile, externalSecretsCloudFormationFile, externalSecretsBaseDir string, fs *afero.Afero) api.ManagedPolicyStore {
	return &managedPolicy{
		externalSecretsFile:               externalSecretsFile,
		externalSecretsCloudFormationFile: externalSecretsCloudFormationFile,
		externalSecretsBaseDir:            externalSecretsBaseDir,
		fs:                                fs,
	}
}
