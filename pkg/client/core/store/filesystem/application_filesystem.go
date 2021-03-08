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
	absoluteApplicationDir := path.Join(s.paths.BaseDir, application.ApplicationName)
	relativeApplicationBaseDir := config.DefaultApplicationBaseDir
	relativeApplicationOverlayDir := path.Join(config.DefaultApplicationOverlayDir, application.Environment)

	err := s.fs.MkdirAll(path.Join(absoluteApplicationDir, relativeApplicationBaseDir), 0o744)
	if err != nil {
		return nil, fmt.Errorf("creating overlay directory: %w", err)
	}

	err = s.fs.MkdirAll(path.Join(absoluteApplicationDir, relativeApplicationOverlayDir), 0o744)
	if err != nil {
		return nil, fmt.Errorf("creating overlay directory: %w", err)
	}

	operations := store.NewFileSystem(absoluteApplicationDir, s.fs)

	// TODO: clean up. Do we need patches? things? meaning of life?
	addOperationIfNotEmpty(operations, "", "argocd-application.yaml", application.ArgoCDResource)

	addOperationIfNotEmpty(operations, relativeApplicationBaseDir, "deployment.yaml", application.Deployment)
	addOperationIfNotEmpty(operations, relativeApplicationBaseDir, "volumes.yaml", application.Volume)
	addOperationIfNotEmpty(operations, relativeApplicationBaseDir, "ingress.yaml", application.Ingress)
	addOperationIfNotEmpty(operations, relativeApplicationBaseDir, "service.yaml", application.Service)
	addOperationIfNotEmpty(operations, relativeApplicationBaseDir, "kustomization.yaml", application.BaseKustomization)

	// TODO: figure out path problem
	addOperationIfNotEmpty(operations, relativeApplicationOverlayDir, "kustomization.yaml", application.OverlayKustomization)
	addOperationIfNotEmpty(operations, relativeApplicationOverlayDir, config.DefaultIngressPatchFilename, application.IngressPatch)
	addOperationIfNotEmpty(operations, relativeApplicationOverlayDir, "service-patch.json", application.ServicePatch)

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
