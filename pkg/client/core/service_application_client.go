package core

import (
	"bytes"
	"context"
	"encoding/json"
	stderrors "errors"
	"fmt"
	"os"
	"path"

	"github.com/go-git/go-git/v5/plumbing/format/index"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/oslokommune/okctl/pkg/clients/kubectl"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"

	"github.com/mishudark/errors"

	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/config/constant"
	"github.com/oslokommune/okctl/pkg/jsonpatch"
	"github.com/oslokommune/okctl/pkg/scaffold"
	"github.com/spf13/afero"
)

const finalizerCascadingDelete = "resources-finalizer.argocd.argoproj.io"

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
	absoluteArgoCDApplicationManifestPath := path.Join(
		s.absoluteRepositoryDir,
		getRelativeArgoCDManifestPath(opts.Cluster, opts.Application),
	)

	err := s.fs.MkdirAll(path.Dir(absoluteArgoCDApplicationManifestPath), 0o700)
	if err != nil {
		return errors.E(err, "ensuring cluster applications directory: %w", err)
	}

	manifest, err := scaffold.GenerateArgoCDApplicationManifest(scaffold.GenerateArgoCDApplicationManifestOpts{
		Name:          opts.Application.Metadata.Name,
		Namespace:     opts.Application.Metadata.Namespace,
		IACRepoURL:    opts.Cluster.Github.URL(),
		SourceSyncDir: getRelativeOverlayDirectory(opts.Cluster, opts.Application),
	})
	if err != nil {
		return errors.E(err, "generating ArgoCD Application manifest")
	}

	err = s.fs.WriteReader(absoluteArgoCDApplicationManifestPath, manifest)
	if err != nil {
		return fmt.Errorf("writing manifest: %w", err)
	}

	return nil
}

// DeleteArgoCDApplicationManifest removes necessary files related to the ArgoCD integration
// This function adds a finalizer to the ArgoCD application manifest before deletion, ensuring a cascading delete.
// https://argo-cd.readthedocs.io/en/stable/user-guide/app_deletion/#about-the-deletion-finalizer
func (s *applicationService) DeleteArgoCDApplicationManifest(opts client.DeleteArgoCDApplicationManifestOpts) error {
	relativeArgoCDApplicationManifestPath := getRelativeArgoCDManifestPath(opts.Cluster, opts.Application)
	absoluteArgoCDApplicationManifestPath := path.Join(s.absoluteRepositoryDir, relativeArgoCDApplicationManifestPath)

	err := s.deleteFileFromGitRepository(opts.Cluster, opts.Application, relativeArgoCDApplicationManifestPath)
	if err != nil {
		return fmt.Errorf("deleting application manifest in repository: %w", err)
	}

	err = s.fs.Remove(absoluteArgoCDApplicationManifestPath)
	if err != nil {
		if !stderrors.Is(err, os.ErrNotExist) {
			return errors.E(err, "removing ArgoCD application manifest")
		}
	}

	err = s.patchArgoCDApplicationManifest(opts.Application.Metadata.Name)
	if err != nil {
		if !stderrors.Is(err, kubectl.ErrNotFound) {
			return fmt.Errorf("adding finalizer to application manifest: %w", err)
		}
	}

	manifest, err := scaffold.GenerateArgoCDApplicationManifest(scaffold.GenerateArgoCDApplicationManifestOpts{
		Name:          opts.Application.Metadata.Name,
		Namespace:     opts.Application.Metadata.Namespace,
		IACRepoURL:    opts.Cluster.Github.URL(),
		SourceSyncDir: getRelativeOverlayDirectory(opts.Cluster, opts.Application),
	})
	if err != nil {
		return fmt.Errorf("generating ArgoCD application manifest: %w", err)
	}

	err = s.kubectl.Delete(manifest)
	if err != nil {
		if !stderrors.Is(err, kubectl.ErrNotFound) {
			return fmt.Errorf("deleting ArgoCD application manifest from cluster: %w", err)
		}
	}

	return nil
}

// HasArgoCDIntegration checks if an application has been set up with ArgoCD
func (s *applicationService) HasArgoCDIntegration(_ context.Context, opts client.HasArgoCDIntegrationOpts) (bool, error) {
	absoluteArgoCDApplicationManifestPath := path.Join(
		s.absoluteRepositoryDir,
		getRelativeArgoCDManifestPath(opts.Cluster, opts.Application),
	)

	_, err := s.fs.Stat(absoluteArgoCDApplicationManifestPath)
	if err != nil {
		if stderrors.Is(err, os.ErrNotExist) {
			return false, nil
		}

		return false, errors.E(err, "checking existence of ArgoCD application manifest: %w", err)
	}

	return true, nil
}

// deleteFileFromGitRepository creates a commit on the main branch of a repository that removes a file path
func (s *applicationService) deleteFileFromGitRepository(cluster v1alpha1.Cluster, app v1alpha1.Application, path string) error {
	repo, err := git.Clone(memory.NewStorage(), memfs.New(), &git.CloneOptions{
		URL:   cluster.Github.URL(),
		Depth: 1,
	})
	if err != nil {
		return fmt.Errorf("cloning repository: %w", err)
	}

	tree, _ := repo.Worktree()

	_, err = tree.Remove(path)
	if err != nil {
		if stderrors.Is(err, index.ErrEntryNotFound) {
			return nil
		}

		return fmt.Errorf("removing file: %w", err)
	}

	_, err = tree.Commit(
		fmt.Sprintf("❌ Remove ArgoCD application manifest for %s", app.Metadata.Name),
		&git.CommitOptions{},
	)
	if err != nil {
		return fmt.Errorf("committing changes: %w", err)
	}

	err = repo.Push(&git.PushOptions{})
	if err != nil {
		return fmt.Errorf("pushing to repository: %w", err)
	}

	return nil
}

// patchArgoCDApplicationManifest adds a cascading deletion finalizer to an ArgoCD application manifest with a patch
// operation
func (s *applicationService) patchArgoCDApplicationManifest(applicationName string) error {
	patch := jsonpatch.New()
	patch.Add(jsonpatch.Operation{
		Type:  jsonpatch.OperationTypeAdd,
		Path:  "/metadata/finalizers",
		Value: []string{finalizerCascadingDelete},
	})

	rawPatch, err := json.Marshal(patch)
	if err != nil {
		return fmt.Errorf("marshalling patch: %w", err)
	}

	err = s.kubectl.Patch(kubectl.PatchOpts{
		Resource: kubectl.Resource{
			Namespace: constant.DefaultArgoCDNamespace,
			Kind:      "application",
			Name:      applicationName,
		},
		Patch: bytes.NewReader(rawPatch),
	})
	if err != nil {
		return fmt.Errorf("patching ArgoCD application manifest: %w", err)
	}

	return nil
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

func getRelativeArgoCDManifestPath(cluster v1alpha1.Cluster, app v1alpha1.Application) string {
	return path.Join(
		cluster.Github.OutputPath,
		cluster.Metadata.Name,
		constant.DefaultArgoCDClusterConfigDir,
		constant.DefaultArgoCDClusterConfigApplicationsDir,
		fmt.Sprintf("%s.yaml", app.Metadata.Name),
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
