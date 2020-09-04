package filesystem

import (
	"encoding/json"
	"fmt"
	"path"
	"strings"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/spf13/afero"
)

type parameter struct {
	paths Paths
	fs    *afero.Afero
}

// SecretParameter contains the data we store to the outputs
type SecretParameter struct {
	ID      api.ID
	Name    string
	Path    string
	Version int64
}

func (p *parameter) SaveSecret(parameter *api.SecretParameter) error {
	param := SecretParameter{
		ID:      parameter.ID,
		Name:    parameter.Name,
		Path:    parameter.Path,
		Version: parameter.Version,
	}

	data, err := json.Marshal(param)
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}

	basePath := path.Join(p.paths.BaseDir, strings.ReplaceAll(parameter.Path, "/", "-"))

	err = p.fs.MkdirAll(basePath, 0o744)
	if err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	err = p.fs.WriteFile(path.Join(basePath, p.paths.OutputFile), data, 0o644)
	if err != nil {
		return fmt.Errorf("failed to write outputs: %w", err)
	}

	return nil
}

// NewParameterStore returns an initialised parameter store
func NewParameterStore(paths Paths, fs *afero.Afero) api.ParameterStore {
	return &parameter{
		paths: paths,
		fs:    fs,
	}
}
