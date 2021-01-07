package filesystem

import (
	"fmt"
	"path"

	"github.com/oslokommune/okctl/pkg/client"

	"github.com/gosimple/slug"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client/store"
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

func (p *parameter) SaveSecret(s *api.SecretParameter) (*store.Report, error) {
	param := SecretParameter{
		ID:      s.ID,
		Name:    s.Name,
		Path:    s.Path,
		Version: s.Version,
	}

	report, err := store.NewFileSystem(path.Join(p.paths.BaseDir, slug.Make(s.Path)), p.fs).
		StoreStruct(p.paths.OutputFile, &param, store.ToJSON()).
		Do()
	if err != nil {
		return nil, fmt.Errorf("failed to store secret parameter: %w", err)
	}

	return report, nil
}

// NewParameterStore returns an initialised parameter store
func NewParameterStore(paths Paths, fs *afero.Afero) client.ParameterStore {
	return &parameter{
		paths: paths,
		fs:    fs,
	}
}
