package filesystem

import (
	"fmt"
	"path"

	"github.com/gosimple/slug"

	"github.com/oslokommune/okctl/pkg/client"

	"github.com/oslokommune/okctl/pkg/client/store"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/spf13/afero"
)

type manifestStore struct {
	externalSecret Paths
	storageClass   Paths
	fs             *afero.Afero
}

// StorageClass is stored in the outputs file
type StorageClass struct {
	ID   api.ID
	Name string
}

func (s *manifestStore) SaveStorageClass(sc *client.StorageClass) (*store.Report, error) {
	o := &StorageClass{
		ID:   sc.ID,
		Name: sc.Name,
	}

	report, err := store.NewFileSystem(path.Join(s.storageClass.BaseDir, slug.Make(sc.Name)), s.fs).
		StoreStruct(s.storageClass.OutputFile, o, store.ToJSON()).
		StoreBytes(s.storageClass.ConfigFile, sc.Manifest).
		Do()
	if err != nil {
		return nil, fmt.Errorf("storing storage class: %w", err)
	}

	return report, nil
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
		return nil, fmt.Errorf("writing external secret: %w", err)
	}

	return report, nil
}

// NewManifestStore returns an initialised store
func NewManifestStore(storageClass, externalSecret Paths, fs *afero.Afero) client.ManifestStore {
	return &manifestStore{
		storageClass:   storageClass,
		externalSecret: externalSecret,
		fs:             fs,
	}
}
