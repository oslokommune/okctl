package crd

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"path"

	"github.com/oslokommune/okctl/pkg/helm/charts/argocd"
	"github.com/oslokommune/okctl/pkg/lib/paths"
	"github.com/oslokommune/okctl/pkg/scaffold"
	"github.com/spf13/afero"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/clients/kubectl"
	"sigs.k8s.io/yaml"
)

// Put adds an application manifest to etcd
func (a applicationState) Put(application v1alpha1.Application) error {
	// We probably want to place the application manifest in the application's own namespace at some point to ensure a
	// folder by feature like hierarchy. However, to do this here we'd need to imperatively create the application's
	// namespace due to not being allowed to apply this manifest to a namespace that does not exist. This approach
	// requires more consideration and this is good enough for the needs right now.
	application.Metadata.Namespace = defaultAppManifestNamespace

	rawApplication, err := yaml.Marshal(application)
	if err != nil {
		return fmt.Errorf("marshalling: %w", err)
	}

	err = a.kubectl.Apply(bytes.NewReader(rawApplication))
	if err != nil {
		return fmt.Errorf("applying: %w", err)
	}

	return nil
}

// Get retrieves an application manifest from etcd
func (a applicationState) Get(name string) (v1alpha1.Application, error) {
	resource, err := a.kubectl.Get(kubectl.Resource{
		Name:      name,
		Namespace: defaultAppManifestNamespace,
	})
	if err != nil {
		return v1alpha1.Application{}, fmt.Errorf("retrieving resource: %w", err)
	}

	rawResource, err := io.ReadAll(resource)
	if err != nil {
		return v1alpha1.Application{}, fmt.Errorf("buffering: %w", err)
	}

	var appManifest v1alpha1.Application

	err = yaml.Unmarshal(rawResource, &appManifest)
	if err != nil {
		return v1alpha1.Application{}, fmt.Errorf("unmarshalling app manifest: %w", err)
	}

	return appManifest, nil
}

// Delete removes an application manifest from etcd
func (a applicationState) Delete(name string) error {
	err := a.kubectl.DeleteByResource(kubectl.Resource{
		Name:      name,
		Kind:      "application.okctl.io",
		Namespace: defaultAppManifestNamespace,
	})
	if err != nil {
		return fmt.Errorf("deleting: %w", err)
	}

	return nil
}

// List returns all application manifests in etcd
func (a applicationState) List() ([]v1alpha1.Application, error) {
	resources, err := a.kubectl.Get(kubectl.Resource{
		Namespace: defaultAppManifestNamespace,
		Kind:      "application.okctl.io",
	})
	if err != nil {
		return nil, fmt.Errorf("retrieving: %w", err)
	}

	rawResources, err := io.ReadAll(resources)
	if err != nil {
		return nil, fmt.Errorf("buffering: %w", err)
	}

	var manifests []v1alpha1.Application

	err = yaml.Unmarshal(rawResources, &manifests)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling: %w", err)
	}

	return manifests, nil
}

// Initialize knows how to apply the CustomResourceDefinition, so we can store application manifests in etcd
func (a applicationState) Initialize(clusterManifest v1alpha1.Cluster, iacRootRepositoryDirectoryPath string) error {
	relativeOkctlDir := paths.GetRelativeClusterOkctlConfigurationDirectory(clusterManifest)
	absoluteOkctlDir := path.Join(iacRootRepositoryDirectoryPath, relativeOkctlDir)

	err := ensureDirectory(a.fs, absoluteOkctlDir)
	if err != nil {
		return fmt.Errorf("ensuring okctl directory: %w", err)
	}

	err = ensureOkctlDirArgoCDTracking(a.fs, clusterManifest, iacRootRepositoryDirectoryPath)
	if err != nil {
		return fmt.Errorf("ensuring okctl directory tracking: %w", err)
	}

	applicationsCustomResourceDefinitionManifestPath := path.Join(
		absoluteOkctlDir,
		defaultApplicationCustomResourceDefinitionFilename,
	)

	err = a.fs.WriteFile(
		applicationsCustomResourceDefinitionManifestPath,
		applicationCRDTemplate,
		paths.DefaultFilePermissions,
	)
	if err != nil {
		return fmt.Errorf("writing application CRD: %w", err)
	}

	return nil
}

func ensureOkctlDirArgoCDTracking(fs *afero.Afero, clusterManifest v1alpha1.Cluster, iacRepositoryRootDir string) error {
	relativeOkctlDir := paths.GetRelativeClusterOkctlConfigurationDirectory(clusterManifest)
	absoluteApplicationsDir := path.Join(iacRepositoryRootDir, paths.GetRelativeClusterApplicationsDirectory(clusterManifest))

	manifest, err := scaffold.GenerateArgoCDApplicationManifest(scaffold.GenerateArgoCDApplicationManifestOpts{
		Name:          "okctl-config",
		Namespace:     argocd.Namespace,
		IACRepoURL:    clusterManifest.Github.URL(),
		SourceSyncDir: relativeOkctlDir,
	})
	if err != nil {
		return fmt.Errorf("generating ArgoCD manifest: %w", err)
	}

	err = fs.WriteReader(
		path.Join(absoluteApplicationsDir, defaultOkctlConfigurationDirArgoCDApplicationManifestFilename),
		manifest,
	)
	if err != nil {
		return fmt.Errorf("writing manifest to applications dir: %w", err)
	}

	return nil
}

func ensureDirectory(fs *afero.Afero, path string) error {
	err := fs.MkdirAll(path, paths.DefaultDirectoryPermissions)
	if err != nil {
		return fmt.Errorf("creating folder: %w", err)
	}

	return nil
}

// NewApplicationState returns an initialized application state instance
func NewApplicationState(fs *afero.Afero, kubectlClient kubectl.Client) client.ApplicationState {
	return &applicationState{
		fs:      fs,
		kubectl: kubectlClient,
	}
}

type applicationState struct {
	fs      *afero.Afero
	kubectl kubectl.Client
}

//go:embed application-crd-template.yaml
var applicationCRDTemplate []byte

const (
	defaultAppManifestNamespace                                   = "okctl"
	defaultApplicationCustomResourceDefinitionFilename            = "application-manifest-crd.yaml"
	defaultOkctlConfigurationDirArgoCDApplicationManifestFilename = "okctl-config.yaml"
)
