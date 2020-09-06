package filesystem

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client/store"
	"github.com/spf13/afero"
)

type helmStore struct {
	externalSecrets      Paths
	albIngressController Paths
	argoCD               Paths
	fs                   *afero.Afero
}

// Helm contains the outputs we will store
type Helm struct {
	ID api.ID
}

func (s *helmStore) SaveArgoCD(helm *api.Helm) error {
	return s.saveHelmChart(s.argoCD, helm)
}

func (s *helmStore) SaveAlbIngressControllerHelmChar(helm *api.Helm) error {
	return s.saveHelmChart(s.albIngressController, helm)
}

func (s *helmStore) SaveExternalSecretsHelmChart(helm *api.Helm) error {
	return s.saveHelmChart(s.externalSecrets, helm)
}

func (s *helmStore) saveHelmChart(paths Paths, helm *api.Helm) error {
	h := &Helm{
		ID: helm.ID,
	}

	_, err := store.NewFileSystem(paths.BaseDir, s.fs).
		StoreStruct(paths.OutputFile, &h, store.ToJSON()).
		StoreStruct(paths.ReleaseFile, helm.Release, store.ToJSON()).
		StoreStruct(paths.ChartFile, helm.Chart, store.ToJSON()).
		Do()
	if err != nil {
		return fmt.Errorf("failed to store helm: %w", err)
	}

	return nil
}

// NewHelmStore returns an initialised helm store
func NewHelmStore(externalSecrets, albIngressController, argocd Paths, fs *afero.Afero) api.HelmStore {
	return &helmStore{
		externalSecrets:      externalSecrets,
		albIngressController: albIngressController,
		argoCD:               argocd,
		fs:                   fs,
	}
}
