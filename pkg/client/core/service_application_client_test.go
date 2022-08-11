package core_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/oslokommune/okctl/pkg/paths"

	"github.com/oslokommune/okctl/pkg/clients/kubectl"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"

	"github.com/oslokommune/okctl/pkg/commands"

	"github.com/sebdah/goldie/v2"
	"github.com/spf13/afero"

	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/client/core"
	"gotest.tools/assert"
)

// nolint:funlen,lll
func TestNewApplicationService(t *testing.T) {
	testInputBuffer := bytes.NewBufferString(defaultTemplate)

	fs := &afero.Afero{Fs: afero.NewMemMapFs()}

	cluster := v1alpha1.Cluster{
		Metadata: v1alpha1.ClusterMeta{
			Name:      "test",
			Region:    "eu-west-1",
			AccountID: "012345678912",
		},
		Github: v1alpha1.ClusterGithub{
			Organisation: "test",
			Repository:   "repo.git",
			OutputPath:   "infrastructure",
		},
		ClusterRootDomain: "kjoremiljo.oslo.systems",
	}

	absoluteRepoDir := "/"
	absoluteOutputDir := path.Join(absoluteRepoDir, cluster.Github.OutputPath)
	absoluteApplicationsDir := path.Join(absoluteOutputDir, paths.DefaultApplicationsOutputDir)

	appManifestService := core.NewApplicationManifestService(fs, absoluteOutputDir)

	service := core.NewApplicationService(
		fs,
		&mockKubectlClient{},
		appManifestService,
		absoluteRepoDir,
		func(_ string, _ string, _ string) error { return nil },
	)

	application, err := commands.InferApplicationFromStdinOrFile(cluster, testInputBuffer, fs, "-")
	assert.NilError(t, err)

	err = service.ScaffoldApplication(context.Background(), &client.ScaffoldApplicationOpts{
		Cluster:        cluster,
		Application:    application,
		CertificateARN: defaultMockARN,
	})
	assert.NilError(t, err)

	err = service.CreateArgoCDApplicationManifest(client.CreateArgoCDApplicationManifestOpts{
		Cluster:     cluster,
		Application: application,
	})
	assert.NilError(t, err)

	g := goldie.New(t)

	appDir := filepath.Join(absoluteApplicationsDir, application.Metadata.Name)
	appBaseDir := filepath.Join(appDir, paths.DefaultApplicationBaseDir)
	appOverlayDir := filepath.Join(appDir, paths.DefaultApplicationOverlayDir, cluster.Metadata.Name)
	clusterArgoCDConfigDir := filepath.Join(
		absoluteOutputDir,
		cluster.Metadata.Name,
		paths.DefaultArgoCDClusterConfigDir,
	)
	clusterApplicationsDir := filepath.Join(clusterArgoCDConfigDir, paths.DefaultArgoCDClusterConfigApplicationsDir)
	clusterNamespacesDir := filepath.Join(clusterArgoCDConfigDir, paths.DefaultArgoCDClusterConfigNamespacesDir)

	g.Assert(t, "kustomization-base.yaml", readFile(t, fs, filepath.Join(appBaseDir, "kustomization.yaml")))
	g.Assert(t, "deployment.yaml", readFile(t, fs, filepath.Join(appBaseDir, "deployment.yaml")))
	g.Assert(t, "volumes.yaml", readFile(t, fs, filepath.Join(appBaseDir, "volumes.yaml")))
	g.Assert(t, "ingress.yaml", readFile(t, fs, filepath.Join(appBaseDir, "ingress.yaml")))
	g.Assert(t, "service.yaml", readFile(t, fs, filepath.Join(appBaseDir, "service.yaml")))
	g.Assert(t, "service-monitor.yaml", readFile(t, fs, filepath.Join(appBaseDir, "service-monitor.yaml")))

	g.Assert(t, "kustomization-overlay.yaml", readFile(t, fs, filepath.Join(appOverlayDir, "kustomization.yaml")))
	g.Assert(t, "deployment-patch.yaml", readFile(t, fs, filepath.Join(appOverlayDir, "deployment-patch.json")))
	g.Assert(t, "ingress-patch.yaml", readFile(t, fs, filepath.Join(appOverlayDir, "ingress-patch.json")))
	g.Assert(t, "argocd-application.yaml", readFile(t, fs, filepath.Join(
		clusterApplicationsDir,
		fmt.Sprintf("%s.yaml", application.Metadata.Name),
	)))
	g.Assert(t, "namespace.yaml", readFile(t, fs, filepath.Join(
		clusterNamespacesDir, fmt.Sprintf("%s.yaml", application.Metadata.Namespace),
	)))
}

//nolint:funlen
func TestDeleteApplication(t *testing.T) {
	ctx := context.Background()
	fs := &afero.Afero{Fs: afero.NewMemMapFs()}
	stdin := bytes.NewBufferString(defaultTemplate)

	clusterManifest := generateMockClusterManifest()
	applicationManifest, err := commands.InferApplicationFromStdinOrFile(clusterManifest, stdin, fs, "-")
	assert.NilError(t, err)

	absoluteRepoDir := "/"
	absoluteOutputDir := path.Join(absoluteRepoDir, clusterManifest.Github.OutputPath)
	absoluteApplicationsDir := path.Join(absoluteOutputDir, paths.DefaultApplicationsOutputDir)

	manifestService := core.NewApplicationManifestService(fs, absoluteOutputDir)
	appService := core.NewApplicationService(
		fs,
		mockKubectlClient{},
		manifestService,
		absoluteRepoDir,
		func(_ string, _ string, _ string) error { return nil },
	)

	err = appService.ScaffoldApplication(context.Background(), &client.ScaffoldApplicationOpts{
		Cluster:        clusterManifest,
		Application:    applicationManifest,
		CertificateARN: defaultMockARN,
	})
	assert.NilError(t, err)

	err = appService.DeleteApplicationManifests(ctx, client.DeleteApplicationManifestsOpts{
		Cluster:     clusterManifest,
		Application: applicationManifest,
	})
	assert.NilError(t, err)

	appDir := filepath.Join(absoluteApplicationsDir, applicationManifest.Metadata.Name)
	appBaseDir := filepath.Join(appDir, paths.DefaultApplicationBaseDir)
	appOverlayDir := filepath.Join(appDir, paths.DefaultApplicationOverlayDir, clusterManifest.Metadata.Name)
	clusterArgoCDConfigDir := filepath.Join(
		absoluteOutputDir,
		clusterManifest.Metadata.Name,
		paths.DefaultArgoCDClusterConfigDir,
	)
	clusterApplicationsDir := filepath.Join(clusterArgoCDConfigDir, paths.DefaultArgoCDClusterConfigApplicationsDir)
	clusterNamespacesDir := filepath.Join(clusterArgoCDConfigDir, paths.DefaultArgoCDClusterConfigNamespacesDir)

	assert.Equal(t, false, fileExists(t, fs, filepath.Join(appBaseDir, "kustomization.yaml")))
	assert.Equal(t, false, fileExists(t, fs, filepath.Join(appBaseDir, "deployment.yaml")))
	assert.Equal(t, false, fileExists(t, fs, filepath.Join(appBaseDir, "volumes.yaml")))
	assert.Equal(t, false, fileExists(t, fs, filepath.Join(appBaseDir, "ingress.yaml")))
	assert.Equal(t, false, fileExists(t, fs, filepath.Join(appBaseDir, "service.yaml")))
	assert.Equal(t, false, fileExists(t, fs, filepath.Join(appBaseDir, "service-monitor.yaml")))

	assert.Equal(t, false, fileExists(t, fs, filepath.Join(appOverlayDir, "kustomization.yaml")))
	assert.Equal(t, false, fileExists(t, fs, filepath.Join(appOverlayDir, "deployment-patch.json")))
	assert.Equal(t, false, fileExists(t, fs, filepath.Join(appOverlayDir, "ingress-patch.json")))
	assert.Equal(t, false, fileExists(t, fs, filepath.Join(
		clusterApplicationsDir,
		fmt.Sprintf("%s.yaml", applicationManifest.Metadata.Name),
	)))

	// Namespace should still exist. User might have resources not tracked by okctl here, and the namespaces will get
	// cleaned up by `delete cluster` later anyhow.
	assert.Equal(t, true, fileExists(t, fs, filepath.Join(
		clusterNamespacesDir, fmt.Sprintf("%s.yaml", applicationManifest.Metadata.Namespace),
	)))
}

func fileExists(t *testing.T, fs *afero.Afero, path string) bool {
	result, err := fs.Exists(path)
	assert.NilError(t, err)

	return result
}

func readFile(t *testing.T, fs *afero.Afero, path string) []byte {
	result, err := fs.ReadFile(path)
	assert.NilError(t, err)

	return result
}

func generateMockClusterManifest() v1alpha1.Cluster {
	return v1alpha1.Cluster{
		Metadata: v1alpha1.ClusterMeta{
			Name:      "test",
			Region:    "eu-west-1",
			AccountID: "012345678912",
		},
		Github: v1alpha1.ClusterGithub{
			Organisation: "test",
			Repository:   "repo.git",
			OutputPath:   "infrastructure",
		},
		ClusterRootDomain: "kjoremiljo.oslo.systems",
	}
}

type mockKubectlClient struct{}

func (m mockKubectlClient) Apply(_ io.Reader) error                 { panic("implement me") }
func (m mockKubectlClient) Delete(_ io.Reader) error                { panic("implement me") }
func (m mockKubectlClient) Patch(_ kubectl.PatchOpts) error         { panic("implement me") }
func (m mockKubectlClient) Exists(_ kubectl.Resource) (bool, error) { panic("implement me") }

const (
	defaultMockARN  = "arn:which:isnt:an:arn"
	defaultTemplate = `apiVersion: okctl.io/v1alpha1
kind: Application

metadata:
  # A name that identifies your app
  name: my-app
  # The Kubernetes namespace where your app will live
  namespace: okctl

image:
  uri: docker.pkg.github.com/my-org/my-repo/my-package:0.0.1

# The URL your app should be available on
# Change to something other than https to disable configuring TLS
# Comment this out to avoid setting up an ingress
subDomain: okctl

# The port your app listens on
# Comment this out to avoid setting up a service (required if url is specified)
port: 3000

# Enable prometheus scraping of metrics
prometheus:
  path: /metrics

# How many replicas of your application should we scaffold
#replicas: 3 # 1 by default

# A Docker repository secret for pulling your image
#imagePullSecret: my-pull-secret-name

# The environment your app requires
#environment:
#  MY_VARIABLE: my-value

# Volumes to mount
volumes:
- /path/to/mount/volume: 24Gi
#  - /path/to/mount/volume: # Requests 1Gi by default

# Annotations for your ingress
#ingress:
#  annotations:
#    nginx.ingress.kubernetes.io/cors-allow-origin: http://localhost:8080
#    cert-manager.io/cluster-issuer: letsencrypt-production
`
	defaultArgoCDApplicationTemplate = `apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: mock-app
  namespace: argocd
spec:
  destination:
    namespace: mock-namespace
    server: https://kubernetes.default.svc
  project: default
  source:
    path: infrastructure/applications/mock-app/overlays/mock-cluster
    repoURL: git@github.com:mock-org/mock-iac-repo
    targetRevision: HEAD
  syncPolicy:
    automated:
      prune: false
      selfHeal: false

---
`
)

//nolint:unparam
func addApplication(t *testing.T, fs *afero.Afero, name string) {
	err := fs.MkdirAll(path.Join("/", "infrastructure", "applications", name, "overlays"), 0o700)
	assert.NilError(t, err)
}

//nolint:unparam
func addClusterToAppOverlays(t *testing.T, fs *afero.Afero, clusterName, appName string) {
	absOverlaysDir := path.Join("/", "infrastructure", "applications", appName, "overlays", clusterName)

	err := fs.MkdirAll(absOverlaysDir, 0o700)
	assert.NilError(t, err)

	err = fs.WriteReader(path.Join(absOverlaysDir, "kustomization.yaml"), strings.NewReader(""))
	assert.NilError(t, err)
}

//nolint:unparam
func addAppToClusterApplications(t *testing.T, fs *afero.Afero, clusterName string, appName string) {
	absClusterApplicationsDir := path.Join("/", "infrastructure", clusterName, "argocd", "applications")

	err := fs.MkdirAll(absClusterApplicationsDir, 0o700)
	assert.NilError(t, err)

	err = fs.WriteReader(path.Join(absClusterApplicationsDir, appName), strings.NewReader(defaultArgoCDApplicationTemplate))
	assert.NilError(t, err)
}

const defaultMockAppName = "mock-app"

//nolint:funlen
func TestAmountAssociatedClusters(t *testing.T) {
	testCases := []struct {
		name                          string
		withFs                        *afero.Afero
		withAppToTest                 string
		expectAssociatedClusterAmount int
	}{
		{
			name: "Should return 0 when theres no cluster overlays and no argocd applications",
			withFs: func() *afero.Afero {
				fs := &afero.Afero{Fs: afero.NewMemMapFs()}

				addApplication(t, fs, defaultMockAppName)

				return fs
			}(),
			withAppToTest:                 defaultMockAppName,
			expectAssociatedClusterAmount: 0,
		},
		{
			name: "Should return 0 when theres no cluster overlays and but a relevant argocd application",
			withFs: func() *afero.Afero {
				fs := &afero.Afero{Fs: afero.NewMemMapFs()}

				appName := defaultMockAppName

				addApplication(t, fs, appName)
				addAppToClusterApplications(t, fs, "mock-cluster", appName)

				return fs
			}(),
			withAppToTest:                 defaultMockAppName,
			expectAssociatedClusterAmount: 0,
		},
		{
			name: "Should return 0 when theres a cluster overlay and but no relevant argocd applications",
			withFs: func() *afero.Afero {
				fs := &afero.Afero{Fs: afero.NewMemMapFs()}

				appName := defaultMockAppName

				addApplication(t, fs, appName)
				addAppToClusterApplications(t, fs, "mock-cluster", appName)

				return fs
			}(),
			withAppToTest:                 defaultMockAppName,
			expectAssociatedClusterAmount: 0,
		},
		{
			name: "Should return 1 when theres a cluster overlay and a relevant argocd application",
			withFs: func() *afero.Afero {
				fs := &afero.Afero{Fs: afero.NewMemMapFs()}

				appName := defaultMockAppName
				clusterName := "mock-cluster"

				addApplication(t, fs, appName)
				addClusterToAppOverlays(t, fs, clusterName, appName)
				addAppToClusterApplications(t, fs, clusterName, appName)

				return fs
			}(),
			withAppToTest:                 defaultMockAppName,
			expectAssociatedClusterAmount: 1,
		},
		{
			name: "Should return 1 when theres three cluster overlays but only one relevant argocd application",
			withFs: func() *afero.Afero {
				fs := &afero.Afero{Fs: afero.NewMemMapFs()}

				appName := defaultMockAppName
				clusterName := "valid-cluster"

				addApplication(t, fs, appName)
				addClusterToAppOverlays(t, fs, clusterName, appName)
				addAppToClusterApplications(t, fs, clusterName, appName)

				addClusterToAppOverlays(t, fs, "not-relevant-cluster", appName)
				addClusterToAppOverlays(t, fs, "random-cluster", appName)

				return fs
			}(),
			withAppToTest:                 defaultMockAppName,
			expectAssociatedClusterAmount: 1,
		},
		{
			name: "Should return 3 when theres three cluster overlays and three relevant argocd applications",
			withFs: func() *afero.Afero {
				fs := &afero.Afero{Fs: afero.NewMemMapFs()}

				appName := defaultMockAppName

				clusterName := "valid-cluster"
				addApplication(t, fs, appName)
				addClusterToAppOverlays(t, fs, clusterName, appName)
				addAppToClusterApplications(t, fs, clusterName, appName)

				clusterName = "valid-second-cluster"
				addApplication(t, fs, appName)
				addClusterToAppOverlays(t, fs, clusterName, appName)
				addAppToClusterApplications(t, fs, clusterName, appName)

				clusterName = "valid-third-cluster"
				addApplication(t, fs, appName)
				addClusterToAppOverlays(t, fs, clusterName, appName)
				addAppToClusterApplications(t, fs, clusterName, appName)

				return fs
			}(),
			withAppToTest:                 defaultMockAppName,
			expectAssociatedClusterAmount: 3,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			result, err := core.AmountAssociatedClusters(
				tc.withFs,
				"/",
				v1alpha1.Cluster{
					Metadata: v1alpha1.ClusterMeta{Name: "mock-cluster"},
					Github:   v1alpha1.ClusterGithub{OutputPath: "/infrastructure"},
				},
				v1alpha1.Application{Metadata: v1alpha1.ApplicationMeta{Name: tc.withAppToTest}},
			)
			assert.NilError(t, err)

			assert.Equal(t, tc.expectAssociatedClusterAmount, result)
		})
	}
}
