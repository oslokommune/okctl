package filesystem

import (
	"fmt"
	"path"

	"github.com/oslokommune/okctl/pkg/config/constant"

	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/client/store"
	"github.com/spf13/afero"
)

type applicationStore struct {
	paths Paths
	fs    *afero.Afero
}

// Helper for optional resources
func addOperationIfNotEmpty(operations store.Operations, filePath string, content []byte) {
	if len(content) == 0 {
		return
	}

	operations.StoreBytes(filePath, content)
}

// SaveApplication applies the application to the file system
func (s *applicationStore) SaveApplication(application *client.ScaffoldedApplication) (*store.Report, error) {
	absoluteApplicationDir := path.Join(s.paths.BaseDir, application.ApplicationName)
	relativeApplicationBaseDir := constant.DefaultApplicationBaseDir
	relativeApplicationOverlayDir := path.Join(constant.DefaultApplicationOverlayDir, application.ClusterName)

	operations := store.NewFileSystem(absoluteApplicationDir, s.fs)
	addOperationIfNotEmpty(operations, "argocd-application.yaml", application.ArgoCDResource)

	operations.AlterStore(store.SetBaseDir(path.Join(absoluteApplicationDir, relativeApplicationBaseDir)))
	addOperationIfNotEmpty(operations, "deployment.yaml", application.Deployment)
	addOperationIfNotEmpty(operations, "volumes.yaml", application.Volume)
	addOperationIfNotEmpty(operations, "ingress.yaml", application.Ingress)
	addOperationIfNotEmpty(operations, "service.yaml", application.Service)
	addOperationIfNotEmpty(operations, "service-monitor.yaml", application.ServiceMonitor)
	addOperationIfNotEmpty(operations, "kustomization.yaml", application.BaseKustomization)

	operations.AlterStore(store.SetBaseDir(path.Join(absoluteApplicationDir, relativeApplicationOverlayDir)))
	addOperationIfNotEmpty(operations, "kustomization.yaml", application.OverlayKustomization)
	addOperationIfNotEmpty(operations, constant.DefaultIngressPatchFilename, application.IngressPatch)

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
