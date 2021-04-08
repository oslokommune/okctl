package core_test

import (
	"bytes"
	"context"
	"path/filepath"
	"testing"

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

const defaultTemplate = `
# A name that identifies your app
name: my-app
# An URI for your app Docker image
image: docker.pkg.github.com/my-org/my-repo/my-package
# The version of your app which is available as an image
version: 0.0.1
# A namespace where your app will live
namespace: okctl


# The URL your app should be available on
# Change to something other than https to disable configuring TLS
# Comment this out to avoid setting up an ingress
subDomain: okctl

# The port your app listens on
# Comment this out to avoid setting up a service (required if url is specified)
port: 3000

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

func TestNewApplicationService(t *testing.T) {
	testInputBuffer := bytes.NewBufferString(defaultTemplate)

	aferoFs := afero.Afero{Fs: afero.NewMemMapFs()}

	mockPaths := clientFilesystem.Paths{BaseDir: "infrastructure/applications"}

	service := core.NewApplicationService(
		&aferoFs,
		mockPaths,
		mockCertService{},
		clientFilesystem.NewApplicationStore(mockPaths, &aferoFs),
	)

	application, err := commands.InferApplicationFromStdinOrFile(testInputBuffer, &aferoFs, "-")
	assert.NilError(t, err)

	env := "test"
	err = service.ScaffoldApplication(context.Background(), &client.ScaffoldApplicationOpts{
		OutputDir: "infrastructure",
		ID: &api.ID{
			Region:       "eu-west-1",
			AWSAccountID: "012345678912",
			Environment:  env,
			Repository:   "not blank",
			ClusterName:  "dummy-cluster",
		},
		HostedZoneID:     "dummyID",
		HostedZoneDomain: "kjoremiljo.oslo.systems",
		IACRepoURL:       "git@dummy.com:test/repo.git",
		Application:      application,
	})
	assert.NilError(t, err)

	g := goldie.New(t)
	g.Assert(t, "argocd-application.yaml", readFile(t, &aferoFs, filepath.Join(mockPaths.BaseDir, "my-app", "argocd-application.yaml")))

	g.Assert(t, "kustomization-base.yaml", readFile(t, &aferoFs, filepath.Join(mockPaths.BaseDir, "my-app", constant.DefaultApplicationBaseDir, "kustomization.yaml")))

	g.Assert(t, "deployment.yaml", readFile(t, &aferoFs, filepath.Join(mockPaths.BaseDir, "my-app", constant.DefaultApplicationBaseDir, "deployment.yaml")))
	g.Assert(t, "volumes.yaml", readFile(t, &aferoFs, filepath.Join(mockPaths.BaseDir, "my-app", constant.DefaultApplicationBaseDir, "volumes.yaml")))
	g.Assert(t, "ingress.yaml", readFile(t, &aferoFs, filepath.Join(mockPaths.BaseDir, "my-app", constant.DefaultApplicationBaseDir, "ingress.yaml")))
	g.Assert(t, "service.yaml", readFile(t, &aferoFs, filepath.Join(mockPaths.BaseDir, "my-app", constant.DefaultApplicationBaseDir, "service.yaml")))

	g.Assert(t,
		"kustomization-overlay.yaml",
		readFile(t, &aferoFs, filepath.Join(mockPaths.BaseDir, "my-app", constant.DefaultApplicationOverlayDir, env, "kustomization.yaml")),
	)
	g.Assert(t, "ingress-patch.yaml", readFile(t, &aferoFs, filepath.Join(mockPaths.BaseDir, "my-app", constant.DefaultApplicationOverlayDir, env, "ingress-patch.json")))
}

func readFile(t *testing.T, fs *afero.Afero, path string) []byte {
	result, err := fs.ReadFile(path)
	assert.NilError(t, err)

	return result
}

type mockCertService struct{}

func (m mockCertService) DeleteCertificate(_ context.Context, _ client.DeleteCertificateOpts) error {
	return nil
}

func (m mockCertService) DeleteCognitoCertificate(_ context.Context, _ client.DeleteCognitoCertificateOpts) error {
	return nil
}

func (m mockCertService) CreateCertificate(_ context.Context, _ client.CreateCertificateOpts) (*client.Certificate, error) {
	return &client.Certificate{
		ARN: "arn:which:isnt:an:arn",
	}, nil
}
