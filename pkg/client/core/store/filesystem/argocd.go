package filesystem

import (
	"path"

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
	paths    Paths
	cert     client.CertificateStore
	manifest client.ManifestStore
	fs       *afero.Afero
}

func (s *argoCDStore) SaveArgoCD(cd *client.ArgoCD) ([]*store.Report, error) {
	argo := &ArgoCD{
		ID:         cd.ID,
		ArgoDomain: cd.ArgoDomain,
		ArgoURL:    cd.ArgoURL,
	}

	chart := &Helm{
		ID: cd.ID,
	}

	r1, err := store.NewFileSystem(s.paths.BaseDir, s.fs).
		// Outputs
		StoreStruct(s.paths.OutputFile, argo, store.ToJSON()).
		// Chart
		AlterStore(store.SetBaseDir(path.Join(s.paths.BaseDir, "helm"))).
		StoreStruct(s.paths.OutputFile, chart, store.ToJSON()).
		StoreStruct(s.paths.ReleaseFile, cd.Chart.Release, store.ToJSON()).
		StoreStruct(s.paths.ChartFile, cd.Chart.Chart, store.ToJSON()).
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

	return []*store.Report{
		r1, r2, r3,
	}, nil
}

// NewArgoCDStore returns an initialised store
func NewArgoCDStore(paths Paths, cert client.CertificateStore, manifest client.ManifestStore, fs *afero.Afero) client.ArgoCDStore {
	return &argoCDStore{
		paths:    paths,
		cert:     cert,
		manifest: manifest,
		fs:       fs,
	}
}
