package filesystem

import (
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/client/store"
	"github.com/spf13/afero"
)

// ArgoCD contains the state written to the outputs
type ArgoCD struct {
	ID         api.ID
	ArgoDomain string
	ArgoURL    string
}

type argoCDStore struct {
	paths     Paths
	helmPaths Paths
	fs        *afero.Afero
}

// nolint: funlen
func (s *argoCDStore) SaveArgoCD(cd *client.ArgoCD) (*store.Report, error) {
	argo := &ArgoCD{
		ID:         cd.ID,
		ArgoDomain: cd.ArgoDomain,
		ArgoURL:    cd.ArgoURL,
	}

	chart := &Helm{
		ID: cd.ID,
	}

	report, err := store.NewFileSystem(s.paths.BaseDir, s.fs).
		// Outputs
		StoreStruct(s.paths.OutputFile, argo, store.ToJSON()).
		// Chart
		AlterStore(store.SetBaseDir(s.helmPaths.BaseDir)).
		StoreStruct(s.helmPaths.OutputFile, chart, store.ToJSON()).
		StoreStruct(s.helmPaths.ReleaseFile, cd.Chart.Release, store.ToJSON()).
		StoreStruct(s.helmPaths.ChartFile, cd.Chart.Chart, store.ToJSON()).
		Do()
	if err != nil {
		return nil, err
	}

	return report, nil
}

// NewArgoCDStore returns an initialised store
func NewArgoCDStore(helmPaths, paths Paths, fs *afero.Afero) client.ArgoCDStore { // nolint: lll
	return &argoCDStore{
		paths:     paths,
		helmPaths: helmPaths,
		fs:        fs,
	}
}
