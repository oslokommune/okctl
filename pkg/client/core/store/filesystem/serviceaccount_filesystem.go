package filesystem

import (
	"path"

	"github.com/gosimple/slug"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/client/store"
	"github.com/spf13/afero"
)

type serviceAccountStore struct {
	paths Paths
	fs    *afero.Afero
}

func (m *serviceAccountStore) SaveCreateServiceAccount(p *api.ServiceAccount) (*store.Report, error) {
	return store.NewFileSystem(path.Join(m.paths.BaseDir, slug.Make(p.Name)), m.fs).
		StoreStruct(
			m.paths.OutputFile,
			&ServiceAccount{
				ID:        p.ID,
				Name:      p.Name,
				PolicyArn: p.PolicyArn,
			},
			store.ToJSON(),
		).
		StoreStruct(m.paths.ConfigFile, p.Config, store.ToJSON()).
		Do()
}

func (m *serviceAccountStore) RemoveDeleteServiceAccount(name string) (*store.Report, error) {
	return store.NewFileSystem(path.Join(m.paths.BaseDir, slug.Make(name)), m.fs).
		Remove(m.paths.OutputFile).
		Remove(m.paths.ConfigFile).
		AlterStore(store.SetBaseDir(m.paths.BaseDir)).
		RemoveDir(slug.Make(name)).
		Do()
}

// NewServiceAccountStore
func NewServiceAccountStore(paths Paths, fs *afero.Afero) client.ServiceAccountStore {
	return &serviceAccountStore{
		paths: paths,
		fs:    fs,
	}
}
