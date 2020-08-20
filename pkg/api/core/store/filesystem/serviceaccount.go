package filesystem

import (
	"encoding/json"
	"fmt"
	"path"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/spf13/afero"
	"sigs.k8s.io/yaml"
)

type serviceAccount struct {
	externalSecretsServiceAccountOutputFile string
	externalSecretsServiceAccountConfigFile string
	externalSecretsBaseDir                  string
	fs                                      *afero.Afero
}

// ServiceAccount contains the data that should
// be serialised to the output file
type ServiceAccount struct {
	ClusterName  string
	Environment  string
	Region       string
	AWSAccountID string
	PolicyArn    string
}

func (m *serviceAccount) SaveExternalSecretsServiceAccount(account *api.ServiceAccount) error {
	p := &ServiceAccount{
		ClusterName:  account.ClusterName,
		Environment:  account.Environment,
		Region:       account.Region,
		AWSAccountID: account.AWSAccountID,
		PolicyArn:    account.PolicyArn,
	}

	data, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("failed to marshal service account: %w", err)
	}

	err = m.fs.MkdirAll(m.externalSecretsBaseDir, 0744)
	if err != nil {
		return fmt.Errorf("failed to create service account directory: %w", err)
	}

	err = m.fs.WriteFile(path.Join(m.externalSecretsBaseDir, m.externalSecretsServiceAccountOutputFile), data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write service account: %w", err)
	}

	d, err := account.Config.YAML()
	if err != nil {
		return fmt.Errorf("failed to marshal service account: %w", err)
	}

	err = m.fs.WriteFile(path.Join(m.externalSecretsBaseDir, m.externalSecretsServiceAccountConfigFile), d, 0644)
	if err != nil {
		return fmt.Errorf("failed to write template: %w", err)
	}

	return nil
}

func (m *serviceAccount) GetExternalSecretsServiceAccount() (*api.ServiceAccount, error) {
	data, err := m.fs.ReadFile(path.Join(m.externalSecretsBaseDir, m.externalSecretsServiceAccountOutputFile))
	if err != nil {
		return nil, fmt.Errorf("failed to read service account file: %w", err)
	}

	a := &ServiceAccount{}

	err = json.Unmarshal(data, a)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal service account: %w", err)
	}

	account := &api.ServiceAccount{
		ClusterName:  a.ClusterName,
		Environment:  a.Environment,
		Region:       a.Region,
		AWSAccountID: a.AWSAccountID,
		PolicyArn:    a.PolicyArn,
	}

	template, err := m.fs.ReadFile(path.Join(m.externalSecretsBaseDir, m.externalSecretsServiceAccountConfigFile))
	if err != nil {
		return nil, fmt.Errorf("failed to read template: %w", err)
	}

	cfg := &api.ClusterConfig{}

	err = yaml.Unmarshal(template, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal service account config: %w", err)
	}

	account.Config = cfg

	return account, nil
}

// NewServiceAccountStore returns an initialised service account store
// nolint: lll
func NewServiceAccountStore(externalSecretsServiceAccountOutputFile, externalSecretsServiceAccountConfigFile, externalSecretsBaseDir string, fs *afero.Afero) api.ServiceAccountStore {
	return &serviceAccount{
		externalSecretsServiceAccountOutputFile: externalSecretsServiceAccountOutputFile,
		externalSecretsServiceAccountConfigFile: externalSecretsServiceAccountConfigFile,
		externalSecretsBaseDir:                  externalSecretsBaseDir,
		fs:                                      fs,
	}
}
