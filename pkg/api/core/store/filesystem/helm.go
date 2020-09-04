package filesystem

import (
	"encoding/json"
	"fmt"
	"path"

	"github.com/oslokommune/okctl/pkg/api"
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

	outputs, err := json.Marshal(h)
	if err != nil {
		return fmt.Errorf("failed to marshal outputs: %w", err)
	}

	err = s.fs.MkdirAll(paths.BaseDir, 0o744)
	if err != nil {
		return fmt.Errorf("failed to directory: %w", err)
	}

	err = s.fs.WriteFile(path.Join(paths.BaseDir, paths.OutputFile), outputs, 0o644)
	if err != nil {
		return fmt.Errorf("failed to write outputs: %w", err)
	}

	release, err := json.MarshalIndent(helm.Release, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal release: %w", err)
	}

	err = s.fs.WriteFile(path.Join(paths.BaseDir, paths.ReleaseFile), release, 0o644)
	if err != nil {
		return fmt.Errorf("failed to write release: %w", err)
	}

	chart, err := json.MarshalIndent(helm.Chart, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal chart: %w", err)
	}

	err = s.fs.WriteFile(path.Join(paths.BaseDir, paths.ChartFile), chart, 0o644)
	if err != nil {
		return fmt.Errorf("failed to write chart: %w", err)
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
