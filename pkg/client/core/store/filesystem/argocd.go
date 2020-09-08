package filesystem

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/config/repository"

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
	repoState *repository.Data
	repoPaths Paths
	cert      client.CertificateStore
	manifest  client.ManifestStore
	param     client.ParameterStore
	fs        *afero.Afero
}

// nolint: funlen
func (s *argoCDStore) SaveArgoCD(cd *client.ArgoCD) ([]*store.Report, error) {
	argo := &ArgoCD{
		ID:         cd.ID,
		ArgoDomain: cd.ArgoDomain,
		ArgoURL:    cd.ArgoURL,
	}

	chart := &Helm{
		ID: cd.ID,
	}

	cluster, ok := s.repoState.Clusters[cd.ID.Environment]
	if !ok {
		return nil, fmt.Errorf("failed to find cluster: %s", cd.ID.ClusterName)
	}

	cluster.ArgoCD = &repository.ArgoCD{
		SiteURL: cd.ArgoURL,
		Domain:  cd.ArgoDomain,
		SecretKey: &repository.SecretKeySecret{
			Name:    cd.SecretKey.Name,
			Path:    cd.SecretKey.Path,
			Version: cd.SecretKey.Version,
		},
	}

	r1, err := store.NewFileSystem(s.paths.BaseDir, s.fs).
		// Outputs
		StoreStruct(s.paths.OutputFile, argo, store.ToJSON()).
		// Chart
		AlterStore(store.SetBaseDir(s.helmPaths.BaseDir)).
		StoreStruct(s.helmPaths.OutputFile, chart, store.ToJSON()).
		StoreStruct(s.helmPaths.ReleaseFile, cd.Chart.Release, store.ToJSON()).
		StoreStruct(s.helmPaths.ChartFile, cd.Chart.Chart, store.ToJSON()).
		// State
		AlterStore(store.SetBaseDir(s.repoPaths.BaseDir)).
		StoreStruct(s.repoPaths.ConfigFile, s.repoState, store.ToYAML()).
		Do()
	if err != nil {
		return nil, err
	}

	manifest := &client.ExternalSecret{
		ID:        cd.ExternalSecret.ID,
		Manifests: cd.ExternalSecret.Manifests,
	}

	r2, err := s.manifest.SaveExternalSecret(manifest)
	if err != nil {
		return nil, err
	}

	r3, err := s.cert.SaveCertificate(cd.Certificate)
	if err != nil {
		return nil, err
	}

	r4, err := s.param.SaveSecret(cd.SecretKey)
	if err != nil {
		return nil, err
	}

	return []*store.Report{
		r1, r2, r3, r4,
	}, nil
}

// NewArgoCDStore returns an initialised store
func NewArgoCDStore(helmPaths, paths Paths, repoState *repository.Data, repoPaths Paths, param client.ParameterStore, cert client.CertificateStore, manifest client.ManifestStore, fs *afero.Afero) client.ArgoCDStore { // nolint: lll
	return &argoCDStore{
		paths:     paths,
		helmPaths: helmPaths,
		repoState: repoState,
		repoPaths: repoPaths,
		cert:      cert,
		manifest:  manifest,
		param:     param,
		fs:        fs,
	}
}
