package core

import (
	"bytes"
	"context"
	"encoding/json"
	stderrors "errors"
	"fmt"
	fsPkg "io/fs"
	"os"
	"path"
	"strings"

	"github.com/oslokommune/okctl/pkg/argocd"

	"github.com/oslokommune/okctl/pkg/paths"
	"github.com/oslokommune/okctl/pkg/scaffold/resources"

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
	gitDeleteRemoteFileFn client.GitDeleteRemoteFileFn
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

	namespace := resources.CreateNamespace(opts.Application)

	rawNamespace, err := scaffold.ResourceAsBytes(namespace)
	if err != nil {
		return fmt.Errorf("serializing namespace: %w", err)
	}

	err = s.appManifestService.SaveNamespace(ctx, client.SaveNamespaceOpts{
		Filename:    fmt.Sprintf("%s.yaml", namespace.Name),
		ClusterName: opts.Cluster.Metadata.Name,
		Payload:     bytes.NewReader(rawNamespace),
	})
	if err != nil {
		return fmt.Errorf("storing namespace: %w", err)
	}

	return nil
}

func getApplicationOverlayClusters(fs *afero.Afero, absoluteApplicationDir string) ([]string, error) {
	absoluteApplicationOverlaysDirectory := path.Join(absoluteApplicationDir, paths.DefaultApplicationOverlayDir)

	clusters := make([]string, 0)

	err := fs.Walk(absoluteApplicationOverlaysDirectory, func(currentPath string, info fsPkg.FileInfo, err error) error {
		if err != nil {
			if stderrors.Is(err, os.ErrNotExist) {
				return nil
			}

			return err
		}

		if contains(recognizedKustomizationFileNames(), path.Base(currentPath)) {
			dir, _ := path.Split(currentPath)

			clusters = append(clusters, path.Base(dir))
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walking overlays directory: %w", err)
	}

	return clusters, nil
}

func getClusterApplications(fs *afero.Afero, absoluteRepositoryRootDir string, cluster v1alpha1.Cluster) ([]string, error) {
	absoluteApplicationsDir := path.Join(absoluteRepositoryRootDir, paths.GetRelativeArgoCDApplicationsDir(cluster))

	applications := make([]string, 0)

	items, err := fs.ReadDir(absoluteApplicationsDir)
	if err != nil {
		if stderrors.Is(err, os.ErrNotExist) {
			return applications, nil
		}

		return nil, fmt.Errorf("listing applications directory: %w", err)
	}

	for _, item := range items {
		valid, err := argocd.IsArgoCDApplication(fs, path.Join(absoluteApplicationsDir, item.Name()))
		if err != nil {
			return nil, fmt.Errorf("checking if ArgoCD application: %w", err)
		}

		if valid {
			applications = append(applications, strings.Replace(item.Name(), path.Ext(item.Name()), "", 1))
		}
	}

	return applications, nil
}

// AmountAssociatedClusters knows how to count the number of associated clusters a specific application has. A cluster is
// classified as associated iff the cluster has an ArgoCD application referencing the app and the app has an overlay folder
// referencing the cluster.
func AmountAssociatedClusters(fs *afero.Afero, absoluteRepositoryRootDir string, clusterContext v1alpha1.Cluster, app v1alpha1.Application) (int, error) {
	absoluteApplicationDir := path.Join(absoluteRepositoryRootDir, paths.GetRelativeApplicationDir(clusterContext, app))

	clusters, err := getApplicationOverlayClusters(fs, absoluteApplicationDir)
	if err != nil {
		return 0, fmt.Errorf("retrieving clusters in overlay directory: %w", err)
	}

	associatedClusters := 0

	for _, cluster := range clusters {
		clusterApplications, err := getClusterApplications(fs, absoluteRepositoryRootDir, v1alpha1.Cluster{
			Metadata: v1alpha1.ClusterMeta{Name: cluster},
			Github:   v1alpha1.ClusterGithub{OutputPath: clusterContext.Github.OutputPath},
		})
		if err != nil {
			return 0, fmt.Errorf("listing cluster applications: %w", err)
		}

		if contains(clusterApplications, app.Metadata.Name) {
			associatedClusters++
		}
	}

	return associatedClusters, nil
}

// DeleteApplicationManifests removes manifests related to an application
func (s *applicationService) DeleteApplicationManifests(_ context.Context, opts client.DeleteApplicationManifestsOpts) error {
	absoluteApplicationOverlaysDirectory := path.Join(s.absoluteRepositoryDir, getRelativeOverlayDirectory(opts.Cluster, opts.Application))
	absoluteApplicationRootDirectory := path.Join(
		s.absoluteRepositoryDir,
		getRelativeApplicationDirectory(opts.Cluster, opts.Application),
	)

	associatedClusters, err := AmountAssociatedClusters(s.fs, s.absoluteRepositoryDir, opts.Cluster, opts.Application)
	if err != nil {
		return fmt.Errorf("counting associated clusters: %w", err)
	}

	var targetDir string

	if associatedClusters == 0 {
		targetDir = absoluteApplicationRootDirectory
	} else {
		targetDir = absoluteApplicationOverlaysDirectory
	}

	err = s.fs.RemoveAll(targetDir)
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

	err := s.patchArgoCDApplicationManifest(opts.Application.Metadata.Name)
	if err != nil {
		if !stderrors.Is(err, kubectl.ErrNotFound) {
			return fmt.Errorf("adding finalizer to application manifest: %w", err)
		}
	}

	err = s.gitDeleteRemoteFileFn(
		opts.Cluster.Github.URL(),
		relativeArgoCDApplicationManifestPath,
		fmt.Sprintf("❌ Remove ArgoCD application manifest for %s", opts.Application.Metadata.Name),
	)
	if err != nil {
		return fmt.Errorf("deleting application manifest in repository: %w", err)
	}

	err = s.fs.Remove(absoluteArgoCDApplicationManifestPath)
	if err != nil {
		if !stderrors.Is(err, os.ErrNotExist) {
			return errors.E(err, "removing ArgoCD application manifest")
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
		paths.DefaultApplicationsOutputDir,
		app.Metadata.Name,
	)
}

func getRelativeOverlayDirectory(cluster v1alpha1.Cluster, app v1alpha1.Application) string {
	return path.Join(
		getRelativeApplicationDirectory(cluster, app),
		paths.DefaultApplicationOverlayDir, cluster.Metadata.Name,
	)
}

func getRelativeArgoCDManifestPath(cluster v1alpha1.Cluster, app v1alpha1.Application) string {
	return path.Join(
		cluster.Github.OutputPath,
		cluster.Metadata.Name,
		paths.DefaultArgoCDClusterConfigDir,
		paths.DefaultArgoCDClusterConfigApplicationsDir,
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
	gitDeleteRemoteFileFn client.GitDeleteRemoteFileFn,
) client.ApplicationService {
	return &applicationService{
		kubectl:               kubectlClient,
		appManifestService:    appManifestService,
		fs:                    fs,
		absoluteRepositoryDir: absoluteRepositoryDir,
		gitDeleteRemoteFileFn: gitDeleteRemoteFileFn,
	}
}

func recognizedKustomizationFileNames() []string {
	return []string{
		"kustomization.yaml",
		"kustomization.yml",
		"Kustomization",
	}
}
