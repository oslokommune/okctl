package core_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"path"
	"path/filepath"
	"testing"

	"github.com/oslokommune/okctl/pkg/clients/kubectl"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"

	"github.com/oslokommune/okctl/pkg/commands"

	"github.com/oslokommune/okctl/pkg/config/constant"

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
	absoluteApplicationsDir := path.Join(absoluteOutputDir, constant.DefaultApplicationsOutputDir)

	appManifestService := core.NewApplicationManifestService(fs, absoluteOutputDir)

	service := core.NewApplicationService(
		fs,
		&mockKubectlClient{},
		appManifestService,
		absoluteRepoDir,
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
	appBaseDir := filepath.Join(appDir, constant.DefaultApplicationBaseDir)
	appOverlayDir := filepath.Join(appDir, constant.DefaultApplicationOverlayDir, cluster.Metadata.Name)
	clusterArgoCDConfigDir := filepath.Join(
		absoluteOutputDir,
		cluster.Metadata.Name,
		constant.DefaultArgoCDClusterConfigDir,
	)
	clusterApplicationsDir := filepath.Join(clusterArgoCDConfigDir, constant.DefaultArgoCDClusterConfigApplicationsDir)
	clusterNamespacesDir := filepath.Join(clusterArgoCDConfigDir, constant.DefaultArgoCDClusterConfigNamespacesDir)

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

func TestDeleteApplication(t *testing.T) {
	ctx := context.Background()
	fs := &afero.Afero{Fs: afero.NewMemMapFs()}
	stdin := bytes.NewBufferString(defaultTemplate)

	clusterManifest := generateMockClusterManifest()
	applicationManifest, err := commands.InferApplicationFromStdinOrFile(clusterManifest, stdin, fs, "-")
	assert.NilError(t, err)

	absoluteRepoDir := "/"
	absoluteOutputDir := path.Join(absoluteRepoDir, clusterManifest.Github.OutputPath)
	absoluteApplicationsDir := path.Join(absoluteOutputDir, constant.DefaultApplicationsOutputDir)

	manifestService := core.NewApplicationManifestService(fs, absoluteOutputDir)
	appService := core.NewApplicationService(fs, mockKubectlClient{}, manifestService, absoluteRepoDir)

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
	appBaseDir := filepath.Join(appDir, constant.DefaultApplicationBaseDir)
	appOverlayDir := filepath.Join(appDir, constant.DefaultApplicationOverlayDir, clusterManifest.Metadata.Name)
	clusterArgoCDConfigDir := filepath.Join(
		absoluteOutputDir,
		clusterManifest.Metadata.Name,
		constant.DefaultArgoCDClusterConfigDir,
	)
	clusterApplicationsDir := filepath.Join(clusterArgoCDConfigDir, constant.DefaultArgoCDClusterConfigApplicationsDir)
	clusterNamespacesDir := filepath.Join(clusterArgoCDConfigDir, constant.DefaultArgoCDClusterConfigNamespacesDir)

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
	assert.Equal(t, false, fileExists(t, fs, filepath.Join(
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
)
