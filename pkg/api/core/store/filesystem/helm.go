package filesystem

import (
	"encoding/json"
	"fmt"
	"path"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/spf13/afero"
)

type helmStore struct {
	extSecBaseDir     string
	extSecReleaseFile string
	extSecChartFile   string
	extSecOutputFile  string
	fs                *afero.Afero
}

type Helm struct {
	Repository  string
	Environment string
}

func (s *helmStore) SaveExternalSecretsHelmChart(helm *api.Helm) error {
	h := &Helm{
		Repository:  helm.Repository,
		Environment: helm.Environment,
	}

	outputs, err := json.Marshal(h)
	if err != nil {
		return fmt.Errorf("failed to marshal outputs: %w", err)
	}

	err = s.fs.MkdirAll(s.extSecBaseDir, 0744)
	if err != nil {
		return fmt.Errorf("failed to directory: %w", err)
	}

	err = s.fs.WriteFile(path.Join(s.extSecBaseDir, s.extSecOutputFile), outputs, 0644)
	if err != nil {
		return fmt.Errorf("failed to write outputs: %w", err)
	}

	release, err := json.MarshalIndent(helm.Release, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal release: %w", err)
	}

	err = s.fs.WriteFile(path.Join(s.extSecBaseDir, s.extSecOutputFile), release, 0644)
	if err != nil {
		return fmt.Errorf("failed to write release: %w", err)
	}

	chart, err := json.MarshalIndent(helm.Chart, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal chart: %w", err)
	}

	err = s.fs.WriteFile(path.Join(s.extSecBaseDir, s.extSecChartFile), chart, 0644)
	if err != nil {
		return fmt.Errorf("failed to write chart: %w", err)
	}

	return nil
}

func NewHelmStore(extSecOutputFile, extSecChartFile, extSecReleaseFile, extSecBaseDir string, fs *afero.Afero) api.HelmStore {
	return &helmStore{
		extSecBaseDir:     extSecBaseDir,
		extSecReleaseFile: extSecReleaseFile,
		extSecChartFile:   extSecChartFile,
		extSecOutputFile:  extSecOutputFile,
		fs:                fs,
	}
}
