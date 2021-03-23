package filesystem

import (
	"path"

	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/client/store"
	"github.com/spf13/afero"
)

type containerRepositoryStore struct {
	paths Paths
	fs    *afero.Afero
}

func (c *containerRepositoryStore) SaveContainerRepository(repository *client.ContainerRepository) (*store.Report, error) {
	report, err := store.NewFileSystem(path.Join(c.paths.BaseDir, repository.ImageName), c.fs).
		StoreStruct(c.paths.OutputFile, repository, store.ToJSON()).
		StoreBytes(c.paths.CloudFormationFile, []byte(repository.CloudFormationTemplate)).
		Do()
	if err != nil {
		return nil, err
	}

	return report, nil
}

func (c *containerRepositoryStore) RemoveContainerRepository(imageName string) (*store.Report, error) {
	return store.NewFileSystem(path.Join(c.paths.BaseDir, imageName), c.fs).
		Remove(c.paths.OutputFile).
		Remove(c.paths.CloudFormationFile).
		AlterStore(store.SetBaseDir(c.paths.BaseDir)).
		Remove(imageName).
		Do()
}

// NewContainerRepositoryStore returns an initialised component store
func NewContainerRepositoryStore(paths Paths, fs *afero.Afero) client.ContainerRepositoryStore {
	return &containerRepositoryStore{
		paths: paths,
		fs:    fs,
	}
}
