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

// Helper for optional resources
func addOperationIfNotEmpty(operations store.Operations, root, filePath string, content []byte) {
	if len(content) == 0 {
		return
	}

	operations.StoreBytes(path.Join(root, filePath), content)
}

// SaveApplication applies the application to the file system
func (s *applicationStore) SaveApplication(application *client.ScaffoldedApplication) (*store.Report, error) {
	// TODO: acquire env
	appDir := path.Join(s.paths.BaseDir, application.ApplicationName)
	baseDir := config.DefaultApplicationBaseDir
	overlayDir := config.DefaultApplicationOverlayDir

	err := s.fs.MkdirAll(path.Join(appDir, baseDir), 0o744)
	if err != nil {
		return nil, fmt.Errorf("creating overlay directory: %w", err)
	}

	err = s.fs.MkdirAll(path.Join(appDir, overlayDir), 0o744)
	if err != nil {
		return nil, fmt.Errorf("creating overlay directory: %w", err)
	}

	operations := store.NewFileSystem(appDir, s.fs)

	// TODO: clean up. Do we need patches? things? meaning of life?

	addOperationIfNotEmpty(operations, baseDir, "deployment.yaml", application.Deployment)
	addOperationIfNotEmpty(operations, baseDir, "argocd-application.yaml", application.ArgoCDResource)
	addOperationIfNotEmpty(operations, baseDir, "volumes.yaml", application.Volume)
	addOperationIfNotEmpty(operations, baseDir, "ingress.yaml", application.Ingress)
	addOperationIfNotEmpty(operations, baseDir, "service.yaml", application.Service)

	addOperationIfNotEmpty(operations, overlayDir, "ingress-patch.json", application.IngressPatch)
	addOperationIfNotEmpty(operations, overlayDir, "service-patch.json", application.ServicePatch)

	report, err := operations.Do()
	if err != nil {
		return nil, fmt.Errorf("storing application resources: %w", err)
	}

	return report, nil
}

// RemoveApplication removes an application from the file system
func (s *applicationStore) RemoveApplication(applicationName string) (*store.Report, error) {
	report, err := store.NewFileSystem(path.Join(s.paths.BaseDir, applicationName), s.fs).
		Remove("").
		Do()
	if err != nil {
		return nil, fmt.Errorf("removing application: %w", err)
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
