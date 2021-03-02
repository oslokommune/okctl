package filesystem

import (
	"fmt"
	"path"

	"github.com/oslokommune/okctl/pkg/config"

	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/client/store"
	"github.com/spf13/afero"
)

type applicationStore struct {
	paths Paths
	fs    *afero.Afero
}

// SaveApplication applies the application to the file system
func (s *applicationStore) SaveApplication(application *client.ScaffoldedApplication) (*store.Report, error) {
	baseDir := path.Join(s.paths.BaseDir, application.ApplicationName)
	overlayDir := path.Join(baseDir, config.DefaultApplicationOverlayDir)

	report, err := store.NewFileSystem(baseDir, s.fs).
		StoreBytes(fmt.Sprintf("%s.yaml", application.ApplicationName), application.KubernetesResources).
		StoreBytes(fmt.Sprintf("%s-application.yaml", application.ApplicationName), application.ArgoCDResource).
		StoreBytes(path.Join(overlayDir, fmt.Sprintf("ingress-patch.json")), application.ArgoCDResource).
		Do()
	if err != nil {
		return nil, fmt.Errorf("failed to store application resources: %w", err)
	}

	return report, nil
}

// RemoveApplication removes an application from the file system
func (s *applicationStore) RemoveApplication(applicationName string) (*store.Report, error) {
	report, err := store.NewFileSystem(path.Join(s.paths.BaseDir, applicationName), s.fs).
		Remove("").
		Do()
	if err != nil {
		return nil, fmt.Errorf("failed to remove application: %w", err)
	}

	return report, err
}

// NewApplicationStore returns an initialized application store
func NewApplicationStore(paths Paths, fs *afero.Afero) client.ApplicationStore {
	return &applicationStore{
		paths: paths,
		fs:    fs,
	}
}
