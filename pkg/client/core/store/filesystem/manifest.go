package filesystem

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/client"

	"github.com/oslokommune/okctl/pkg/client/store"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/spf13/afero"
)

type manifestStore struct {
	externalSecret Paths
	fs             *afero.Afero
}

// ExternalSecret is stored in the outputs file
type ExternalSecret struct {
	ID        api.ID
	Manifests []string
}

func (s *manifestStore) SaveExternalSecret(e *client.ExternalSecret) (*store.Report, error) {
	o := &ExternalSecret{
		ID:        e.ID,
		Manifests: make([]string, len(e.Manifests)),
	}

	manifests := make([]store.AddStoreBytes, len(e.Manifests))

	i := 0

	for name, data := range e.Manifests {
		manifests[i] = store.AddStoreBytes{
			Name: name,
			Data: data,
		}

		o.Manifests[i] = name

		i++
	}

	report, err := store.NewFileSystem(s.externalSecret.BaseDir, s.fs).
		StoreStruct(s.externalSecret.OutputFile, &o, store.ToJSON()).
		AddStoreBytes(manifests...).
		Do()
	if err != nil {
		return nil, fmt.Errorf("failed to write external secret: %w", err)
	}

	return report, nil
}

// NewManifestStore returns an initialised store
func NewManifestStore(externalSecret Paths, fs *afero.Afero) client.ManifestStore {
	return &manifestStore{
		externalSecret: externalSecret,
		fs:             fs,
	}
}
