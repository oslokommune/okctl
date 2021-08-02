package core_test

import (
	"bytes"
	"context"
	"path"
	"path/filepath"
	"testing"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"

	"github.com/oslokommune/okctl/pkg/commands"

	"github.com/oslokommune/okctl/pkg/config/constant"

	"github.com/sebdah/goldie/v2"
	"github.com/spf13/afero"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/client/core"
	clientFilesystem "github.com/oslokommune/okctl/pkg/client/core/store/filesystem"
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

// nolint: lll
func TestNewApplicationService(t *testing.T) {
	testInputBuffer := bytes.NewBufferString(defaultTemplate)

	aferoFs := afero.Afero{
		Fs: afero.NewMemMapFs(),
	}

	mockPaths := clientFilesystem.Paths{
		BaseDir: path.Join("infrastructure", "applications"),
	}

	service := core.NewApplicationService(
		&mockCertService{},
		clientFilesystem.NewApplicationStore(mockPaths, &aferoFs),
	)

	cluster := v1alpha1.Cluster{
		ClusterRootDomain: "kjoremiljo.oslo.systems",
	}

	application, err := commands.InferApplicationFromStdinOrFile(cluster, testInputBuffer, &aferoFs, "-")
	assert.NilError(t, err)

	clusterName := "test"
	err = service.ScaffoldApplication(context.Background(), &client.ScaffoldApplicationOpts{
		ID: &api.ID{
			Region:       "eu-west-1",
			AWSAccountID: "012345678912",
			ClusterName:  clusterName,
		},
		HostedZoneID:     "dummyID",
		HostedZoneDomain: "kjoremiljo.oslo.systems",
		IACRepoURL:       "git@dummy.com:test/repo.git",
		Application:      application,
		OutputDir:        "infrastructure",
	})
	assert.NilError(t, err)

	g := goldie.New(t)

	appDir := filepath.Join(mockPaths.BaseDir, application.Metadata.Name)
	appBaseDir := filepath.Join(appDir, constant.DefaultApplicationBaseDir)
	appOverlayDir := filepath.Join(appDir, constant.DefaultApplicationOverlayDir, clusterName)

	g.Assert(t, "kustomization-base.yaml", readFile(t, &aferoFs, filepath.Join(appBaseDir, "kustomization.yaml")))
	g.Assert(t, "namespace.yaml", readFile(t, &aferoFs, filepath.Join(appBaseDir, "namespace.yaml")))
	g.Assert(t, "deployment.yaml", readFile(t, &aferoFs, filepath.Join(appBaseDir, "deployment.yaml")))
	g.Assert(t, "volumes.yaml", readFile(t, &aferoFs, filepath.Join(appBaseDir, "volumes.yaml")))
	g.Assert(t, "ingress.yaml", readFile(t, &aferoFs, filepath.Join(appBaseDir, "ingress.yaml")))
	g.Assert(t, "service.yaml", readFile(t, &aferoFs, filepath.Join(appBaseDir, "service.yaml")))
	g.Assert(t, "service-monitor.yaml", readFile(t, &aferoFs, filepath.Join(appBaseDir, "service-monitor.yaml")))

	g.Assert(t, "kustomization-overlay.yaml", readFile(t, &aferoFs, filepath.Join(appOverlayDir, "kustomization.yaml")))
	g.Assert(t, "deployment-patch.yaml", readFile(t, &aferoFs, filepath.Join(appOverlayDir, "deployment-patch.json")))
	g.Assert(t, "ingress-patch.yaml", readFile(t, &aferoFs, filepath.Join(appOverlayDir, "ingress-patch.json")))
	g.Assert(t, "argocd-application.yaml", readFile(t, &aferoFs, filepath.Join(appOverlayDir, "argocd-application.yaml")))
}

// nolint: funlen
func TestCertificateCreation(t *testing.T) {
	testCases := []struct {
		name string

		withApplication   func() v1alpha1.Application
		expectCreateCount int
	}{
		{
			name: "Should request certificate upon subDomain specification",

			withApplication: func() v1alpha1.Application {
				app := createValidApplication()

				app.SubDomain = "dummyapp"

				return app
			},

			expectCreateCount: 1,
		},
		{
			name: "Should not request certificate when subDomain is not specified",

			withApplication: func() v1alpha1.Application {
				app := createValidApplication()

				app.SubDomain = ""

				return app
			},

			expectCreateCount: 0,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			certService := &mockCertService{CreateCounter: 0}

			service := core.NewApplicationService(
				certService,
				clientFilesystem.NewApplicationStore(clientFilesystem.Paths{}, &afero.Afero{Fs: afero.NewMemMapFs()}),
			)

			err := service.ScaffoldApplication(context.Background(), &client.ScaffoldApplicationOpts{
				OutputDir: "infrastructure",
				ID: &api.ID{
					Region:       "azeroth",
					AWSAccountID: "012345678912",
					ClusterName:  "dummy-dev",
				},
				HostedZoneID:     "somedummyid",
				HostedZoneDomain: "okctl.io",
				IACRepoURL:       "dummyurl",
				Application:      tc.withApplication(),
			})

			assert.NilError(t, err)
			assert.Equal(t, tc.expectCreateCount, certService.CreateCounter)
		})
	}
}

func createValidApplication() v1alpha1.Application {
	app := v1alpha1.NewApplication(v1alpha1.NewCluster())

	app.Metadata.Name = "dummy-app"
	app.Metadata.Namespace = "dummyns"

	return app
}

func readFile(t *testing.T, fs *afero.Afero, path string) []byte {
	result, err := fs.ReadFile(path)
	assert.NilError(t, err)

	return result
}

type mockCertService struct {
	CreateCounter int
}

func (m *mockCertService) DeleteCertificate(_ context.Context, _ client.DeleteCertificateOpts) error {
	return nil
}

func (m *mockCertService) DeleteCognitoCertificate(_ context.Context, _ client.DeleteCognitoCertificateOpts) error {
	return nil
}

func (m *mockCertService) CreateCertificate(_ context.Context, _ client.CreateCertificateOpts) (*client.Certificate, error) {
	m.CreateCounter++

	return &client.Certificate{
		ARN: "arn:which:isnt:an:arn",
	}, nil
}
