package filesystem

import (
	"path"

	"github.com/oslokommune/okctl/pkg/client"

	"github.com/oslokommune/okctl/pkg/client/store"

	"github.com/spf13/afero"
)

type manifestStore struct {
	paths Paths
	fs    *afero.Afero
}

func (s *manifestStore) SaveKubernetesManifests(manifest *client.KubernetesManifest) error {
	_, err := store.NewFileSystem(path.Join(s.paths.BaseDir, manifest.Type.String(), manifest.Name), s.fs).
		StoreBytes(manifest.Name, manifest.Content).
		Do()
	if err != nil {
		return err
	}

	return nil
}

func (s *manifestStore) RemoveKubernetesManifests(name string, typ client.ManifestType) error {
	_, err := store.NewFileSystem(path.Join(s.paths.BaseDir, typ.String(), name), s.fs).
		RemoveDir("").
		Do()
	if err != nil {
		return err
	}

	return nil
}

// NewManifestStore returns an initialised store
func NewManifestStore(paths Paths, fs *afero.Afero) client.ManifestStore {
	return &manifestStore{
		paths: paths,
		fs:    fs,
	}
}
