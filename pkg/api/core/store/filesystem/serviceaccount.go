package filesystem

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/client/store"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/spf13/afero"
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

	_, err := store.NewFileSystem(paths.BaseDir, s.fs).
		StoreStruct(paths.OutputFile, p, store.ToJSON()).
		StoreStruct(paths.ConfigFile, account.Config, store.ToJSON()).
		Do()
	if err != nil {
		return fmt.Errorf("failed to store service account: %w", err)
	}

	return nil
}

func (s *serviceAccount) getServiceAccount(paths Paths) (*api.ServiceAccount, error) {
	cfg := &api.ClusterConfig{}

	a := &ServiceAccount{}

	_, err := store.NewFileSystem(paths.BaseDir, s.fs).
		GetStruct(paths.OutputFile, a, store.FromJSON()).
		GetStruct(paths.ConfigFile, cfg, store.FromYAML()).
		Do()
	if err != nil {
		return nil, fmt.Errorf("failed to get service account: %w", err)
	}

	return &api.ServiceAccount{
		ID:        a.ID,
		PolicyArn: a.PolicyArn,
		Config:    cfg,
	}, nil
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
