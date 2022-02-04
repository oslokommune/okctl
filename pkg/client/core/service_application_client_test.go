package core_test

import (
	"bytes"
	"context"
	"io"
	"path"
	"path/filepath"
	"testing"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"

	"github.com/oslokommune/okctl/pkg/commands"

	"github.com/oslokommune/okctl/pkg/config/constant"

	"github.com/sebdah/goldie/v2"
	"github.com/spf13/afero"

	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/client/core"
	"gotest.tools/assert"
)

const defaultTemplate = `apiVersion: okctl.io/v1alpha1
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

	appManifestService := core.NewApplicationManifestService(fs, absoluteApplicationsDir)

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
		CertificateARN: "arn:which:isnt:an:arn",
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

	g.Assert(t, "kustomization-base.yaml", readFile(t, fs, filepath.Join(appBaseDir, "kustomization.yaml")))
	g.Assert(t, "namespace.yaml", readFile(t, fs, filepath.Join(appBaseDir, "namespace.yaml")))
	g.Assert(t, "deployment.yaml", readFile(t, fs, filepath.Join(appBaseDir, "deployment.yaml")))
	g.Assert(t, "volumes.yaml", readFile(t, fs, filepath.Join(appBaseDir, "volumes.yaml")))
	g.Assert(t, "ingress.yaml", readFile(t, fs, filepath.Join(appBaseDir, "ingress.yaml")))
	g.Assert(t, "service.yaml", readFile(t, fs, filepath.Join(appBaseDir, "service.yaml")))
	g.Assert(t, "service-monitor.yaml", readFile(t, fs, filepath.Join(appBaseDir, "service-monitor.yaml")))

	g.Assert(t, "kustomization-overlay.yaml", readFile(t, fs, filepath.Join(appOverlayDir, "kustomization.yaml")))
	g.Assert(t, "deployment-patch.yaml", readFile(t, fs, filepath.Join(appOverlayDir, "deployment-patch.json")))
	g.Assert(t, "ingress-patch.yaml", readFile(t, fs, filepath.Join(appOverlayDir, "ingress-patch.json")))
	g.Assert(t, "argocd-application.yaml", readFile(t, fs, filepath.Join(appOverlayDir, "argocd-application.yaml")))
}

func readFile(t *testing.T, fs *afero.Afero, path string) []byte {
	result, err := fs.ReadFile(path)
	assert.NilError(t, err)

	return result
}

type mockKubectlClient struct{}

func (m mockKubectlClient) Apply(_ io.Reader) error  { panic("implement me") }
func (m mockKubectlClient) Delete(_ io.Reader) error { panic("implement me") }
