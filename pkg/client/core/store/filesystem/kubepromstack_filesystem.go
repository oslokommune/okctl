package filesystem

import (
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/client/store"
	"github.com/spf13/afero"
)

// KubePromStack contains the state written to the outputs
type KubePromStack struct {
	ID       api.ID
	Hostname string
}

type kubePromStackStore struct {
	paths     Paths
	helmPaths Paths
	fs        *afero.Afero
}

func (s *kubePromStackStore) SaveKubePromStack(stack *client.KubePromStack) (*store.Report, error) {
	argo := &KubePromStack{
		ID:       stack.ID,
		Hostname: stack.Hostname,
	}

	chart := &Helm{
		ID: stack.ID,
	}

	report, err := store.NewFileSystem(s.paths.BaseDir, s.fs).
		// Outputs
		StoreStruct(s.paths.OutputFile, argo, store.ToJSON()).
		// Chart
		AlterStore(store.SetBaseDir(s.helmPaths.BaseDir)).
		StoreStruct(s.helmPaths.OutputFile, chart, store.ToJSON()).
		StoreStruct(s.helmPaths.ReleaseFile, stack.Chart.Release, store.ToJSON()).
		StoreStruct(s.helmPaths.ChartFile, stack.Chart.Chart, store.ToJSON()).
		Do()
	if err != nil {
		return nil, err
	}

	return report, nil
}

// NewKubePromStackStore returns an initialised store
func NewKubePromStackStore(
	helmPaths, paths Paths,
	fs *afero.Afero,
) client.KubePromStackStore {
	return &kubePromStackStore{
		paths:     paths,
		helmPaths: helmPaths,
		fs:        fs,
	}
}
