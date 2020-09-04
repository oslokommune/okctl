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
	externalSecrets      Paths
	albIngressController Paths
	externalDNS          Paths
	fs                   *afero.Afero
}

// ServiceAccount contains the data that should
// be serialised to the output file
type ServiceAccount struct {
	ID        api.ID
	PolicyArn string
}

func (s *serviceAccount) SaveExternalDNSServiceAccount(account *api.ServiceAccount) error {
	return s.saveServiceAccount(s.externalDNS, account)
}

func (s *serviceAccount) GetExternalDNSServiceAccount() (*api.ServiceAccount, error) {
	return s.getServiceAccount(s.externalDNS)
}

func (s *serviceAccount) SaveAlbIngressControllerServiceAccount(account *api.ServiceAccount) error {
	return s.saveServiceAccount(s.albIngressController, account)
}

func (s *serviceAccount) GetAlbIngressControllerServiceAccount() (*api.ServiceAccount, error) {
	return s.getServiceAccount(s.albIngressController)
}

func (s *serviceAccount) SaveExternalSecretsServiceAccount(account *api.ServiceAccount) error {
	return s.saveServiceAccount(s.externalSecrets, account)
}

func (s *serviceAccount) GetExternalSecretsServiceAccount() (*api.ServiceAccount, error) {
	return s.getServiceAccount(s.externalSecrets)
}

func (s *serviceAccount) saveServiceAccount(paths Paths, account *api.ServiceAccount) error {
	p := &ServiceAccount{
		ID:        account.ID,
		PolicyArn: account.PolicyArn,
	}

	data, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("failed to marshal service account: %w", err)
	}

	err = s.fs.MkdirAll(paths.BaseDir, 0o744)
	if err != nil {
		return fmt.Errorf("failed to create service account directory: %w", err)
	}

	err = s.fs.WriteFile(path.Join(paths.BaseDir, paths.OutputFile), data, 0o644)
	if err != nil {
		return fmt.Errorf("failed to write service account: %w", err)
	}

	d, err := account.Config.YAML()
	if err != nil {
		return fmt.Errorf("failed to marshal config file: %w", err)
	}

	err = s.fs.WriteFile(path.Join(paths.BaseDir, paths.ConfigFile), d, 0o644)
	if err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

func (s *serviceAccount) getServiceAccount(paths Paths) (*api.ServiceAccount, error) {
	data, err := s.fs.ReadFile(path.Join(paths.BaseDir, paths.OutputFile))
	if err != nil {
		return nil, fmt.Errorf("failed to read service account file: %w", err)
	}

	a := &ServiceAccount{}

	err = json.Unmarshal(data, a)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal service account: %w", err)
	}

	account := &api.ServiceAccount{
		ID:        a.ID,
		PolicyArn: a.PolicyArn,
	}

	template, err := s.fs.ReadFile(path.Join(paths.BaseDir, paths.ConfigFile))
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
func NewServiceAccountStore(externalSecrets, albIngressController, externalDNS Paths, fs *afero.Afero) api.ServiceAccountStore {
	return &serviceAccount{
		externalSecrets:      externalSecrets,
		externalDNS:          externalDNS,
		albIngressController: albIngressController,
		fs:                   fs,
	}
}
