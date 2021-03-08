package filesystem

import (
	"path"

	"github.com/gosimple/slug"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/client/store"
	"github.com/spf13/afero"
)

type managedPolicyStore struct {
	paths Paths
	fs    *afero.Afero
}

func (m *managedPolicyStore) SaveCreatePolicy(p *api.ManagedPolicy) (*store.Report, error) {
	return store.NewFileSystem(path.Join(m.paths.BaseDir, slug.Make(p.StackName)), m.fs).
		StoreStruct(
			m.paths.OutputFile,
			&ManagedPolicy{
				ID:        p.ID,
				StackName: p.StackName,
				PolicyARN: p.PolicyARN,
			},
			store.ToJSON(),
		).
		StoreBytes(m.paths.CloudFormationFile, p.CloudFormationTemplate).
		Do()
}

func (m *managedPolicyStore) RemoveDeletePolicy(stackName string) (*store.Report, error) {
	return store.NewFileSystem(path.Join(m.paths.BaseDir, slug.Make(stackName)), m.fs).
		Remove(m.paths.OutputFile).
		Remove(m.paths.CloudFormationFile).
		AlterStore(store.SetBaseDir(m.paths.BaseDir)).
		RemoveDir(slug.Make(stackName)).
		Do()
}

// NewManagedPolicyStore
func NewManagedPolicyStore(paths Paths, fs *afero.Afero) client.ManagedPolicyStore {
	return &managedPolicyStore{
		paths: paths,
		fs:    fs,
	}
}
