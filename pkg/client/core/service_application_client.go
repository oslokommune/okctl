package core

import (
	"bytes"
	"context"
	stderrors "errors"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/oslokommune/okctl/pkg/clients/kubectl"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"

	"github.com/mishudark/errors"

	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/config/constant"
	"github.com/oslokommune/okctl/pkg/jsonpatch"
	"github.com/oslokommune/okctl/pkg/scaffold"
	"github.com/spf13/afero"
)

const (
	defaultArgoCDApplicationManifestPermissions = 0o644 // u+rw g+r o+r
	defaultArgoCDApplicationManifestFilename    = "argocd-application.yaml"
)

type applicationService struct {
	appManifestService    client.ApplicationManifestService
	fs                    *afero.Afero
	absoluteRepositoryDir string
	kubectl               kubectl.Client
}

// ScaffoldApplication turns a file path into Kubernetes resources
// nolint: funlen
func (s *applicationService) ScaffoldApplication(ctx context.Context, opts *client.ScaffoldApplicationOpts) error {
	err := opts.Validate()
	if err != nil {
		return err
	}

	manifestSaver := generateManifestSaver(ctx, s.appManifestService, opts.Application.Metadata.Name)
	patchSaver := generatePatchSaver(ctx, s.appManifestService, opts.Cluster.Metadata.Name, opts.Application.Metadata.Name)

	err = scaffold.GenerateApplicationBase(scaffold.GenerateApplicationBaseOpts{
		SaveManifest: manifestSaver,
		Application:  opts.Application,
	})
	if err != nil {
		return fmt.Errorf("generating application base resources: %w", err)
	}

	err = scaffold.GenerateApplicationOverlay(scaffold.GenerateApplicationOverlayOpts{
		SavePatch:      patchSaver,
		Application:    opts.Application,
		Domain:         opts.Cluster.ClusterRootDomain,
		CertificateARN: opts.CertificateARN,
	})
	if err != nil {
		return fmt.Errorf("generating application overlay: %w", err)
	}

	return nil
}

// DeleteApplicationManifests removes manifests related to an application
func (s *applicationService) DeleteApplicationManifests(_ context.Context, opts client.DeleteApplicationManifestsOpts) error {
	absoluteApplicationDir := path.Join(
		s.absoluteRepositoryDir,
		getRelativeApplicationDirectory(opts.Cluster, opts.Application),
	)

	err := s.fs.RemoveAll(absoluteApplicationDir)
	if err != nil {
		return errors.E(err, "removing application directory: %w", err)
	}

	return nil
}

// CreateArgoCDApplicationManifest creates necessary files for the ArgoCD integration
func (s *applicationService) CreateArgoCDApplicationManifest(opts client.CreateArgoCDApplicationManifestOpts) error {
	relativeOverlayDir := getRelativeOverlayDirectory(opts.Cluster, opts.Application)
	absoluteOverlayDir := path.Join(s.absoluteRepositoryDir, relativeOverlayDir)
	absoluteArgoCDApplicationManifestPath := path.Join(absoluteOverlayDir, defaultArgoCDApplicationManifestFilename)

	err := s.fs.MkdirAll(absoluteOverlayDir, 0o700)
	if err != nil {
		return errors.E(err, "ensuring overlay directory: %w", err)
	}

	err = scaffold.GenerateArgoCDApplicationManifest(scaffold.GenerateArgoCDApplicationManifestOpts{
		Saver: func(content []byte) error {
			return s.fs.WriteFile(absoluteArgoCDApplicationManifestPath, content, defaultArgoCDApplicationManifestPermissions)
		},
		Application:                   opts.Application,
		IACRepoURL:                    opts.Cluster.Github.URL(),
		RelativeApplicationOverlayDir: relativeOverlayDir,
	})
	if err != nil {
		return errors.E(err, "generating ArgoCD Application manifest")
	}

	return nil
}

// DeleteArgoCDApplicationManifest removes necessary files related to the ArgoCD integration
func (s *applicationService) DeleteArgoCDApplicationManifest(opts client.DeleteArgoCDApplicationManifestOpts) error {
	var manifest io.Reader

	err := scaffold.GenerateArgoCDApplicationManifest(scaffold.GenerateArgoCDApplicationManifestOpts{
		Saver: func(content []byte) error {
			manifest = bytes.NewReader(content)

			return nil
		},
		Application:                   opts.Application,
		IACRepoURL:                    opts.Cluster.Github.URL(),
		RelativeApplicationOverlayDir: getRelativeOverlayDirectory(opts.Cluster, opts.Application),
	})
	if err != nil {
		return fmt.Errorf("generating ArgoCD application manifest: %w", err)
	}

	err = s.kubectl.Delete(manifest)
	if err != nil {
		return fmt.Errorf("deleting ArgoCD application manifest from cluster: %w", err)
	}

	absoluteOverlayDir := path.Join(s.absoluteRepositoryDir, getRelativeOverlayDirectory(opts.Cluster, opts.Application))

	argoCDApplicationManifestPath := path.Join(absoluteOverlayDir, defaultArgoCDApplicationManifestFilename)

	err = s.fs.Remove(argoCDApplicationManifestPath)
	if err != nil {
		if stderrors.Is(err, os.ErrNotExist) {
			return nil
		}

		return errors.E(err, "removing ArgoCD application manifest")
	}

	return nil
}

// HasArgoCDIntegration checks if an application has been set up with ArgoCD
func (s *applicationService) HasArgoCDIntegration(_ context.Context, opts client.HasArgoCDIntegrationOpts) (bool, error) {
	absoluteOverlayDir := path.Join(s.absoluteRepositoryDir, getRelativeOverlayDirectory(opts.Cluster, opts.Application))
	absoluteArgoCDApplicationManifestPath := path.Join(absoluteOverlayDir, defaultArgoCDApplicationManifestFilename)

	_, err := s.fs.Stat(absoluteArgoCDApplicationManifestPath)
	if err != nil {
		if stderrors.Is(err, os.ErrNotExist) {
			return false, nil
		}

		return false, errors.E(err, "checking existence of ArgoCD application manifest: %w", err)
	}

	return true, nil
}

func getRelativeApplicationDirectory(cluster v1alpha1.Cluster, app v1alpha1.Application) string {
	return path.Join(
		cluster.Github.OutputPath,
		constant.DefaultApplicationsOutputDir,
		app.Metadata.Name,
	)
}

func getRelativeOverlayDirectory(cluster v1alpha1.Cluster, app v1alpha1.Application) string {
	return path.Join(
		getRelativeApplicationDirectory(cluster, app),
		constant.DefaultApplicationOverlayDir, cluster.Metadata.Name,
	)
}

func generateManifestSaver(ctx context.Context, service client.ApplicationManifestService, applicationName string) scaffold.ManifestSaver {
	return func(filename string, content []byte) error {
		return service.SaveManifest(ctx, client.SaveManifestOpts{
			ApplicationName: applicationName,
			Filename:        filename,
			Content:         content,
		})
	}
}

func generatePatchSaver(ctx context.Context, service client.ApplicationManifestService, clusterName, applicationName string) scaffold.PatchSaver {
	return func(kind string, patch jsonpatch.Patch) error {
		return service.SavePatch(ctx, client.SavePatchOpts{
			ApplicationName: applicationName,
			ClusterName:     clusterName,
			Kind:            kind,
			Patch:           patch,
		})
	}
}

// NewApplicationService initializes a new Scaffold application service
func NewApplicationService(
	fs *afero.Afero,
	kubectlClient kubectl.Client,
	appManifestService client.ApplicationManifestService,
	absoluteRepositoryDir string,
) client.ApplicationService {
	return &applicationService{
		kubectl:               kubectlClient,
		appManifestService:    appManifestService,
		fs:                    fs,
		absoluteRepositoryDir: absoluteRepositoryDir,
	}
}
