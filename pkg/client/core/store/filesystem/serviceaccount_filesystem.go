package filesystem

import (
	"path"

	"github.com/gosimple/slug"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/client/store"
	"github.com/spf13/afero"
)

type serviceAccountStore struct {
	paths Paths
	fs    *afero.Afero
}

func (s *serviceAccountStore) SaveServiceAccount(sa *client.ServiceAccount) error {
	_, err := store.NewFileSystem(path.Join(s.paths.BaseDir, slug.Make(sa.Name)), s.fs).
		StoreStruct(s.paths.ConfigFile, sa.Config, store.ToJSON()).
		Do()
	if err != nil {
		return err
	}

	return nil
}

func (s *serviceAccountStore) RemoveServiceAccount(name string) error {
	_, err := store.NewFileSystem(path.Join(s.paths.BaseDir, slug.Make(name)), s.fs).
		Remove(s.paths.ConfigFile).
		AlterStore(store.SetBaseDir(s.paths.BaseDir)).
		RemoveDir(slug.Make(name)).
		RemoveDir("").
		Do()
	if err != nil {
		return err
	}

	return nil
}

// NewServiceAccountStore returns an initialised store
func NewServiceAccountStore(paths Paths, fs *afero.Afero) client.ServiceAccountStore {
	return &serviceAccountStore{
		paths: paths,
		fs:    fs,
	}
}
