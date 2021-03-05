package core_test

import (
	"bytes"
	"context"
	"fmt"
	"github.com/oslokommune/okctl/pkg/config"
	"k8s.io/cli-runtime/pkg/kustomize"
	"path/filepath"
	"testing"

	"github.com/sebdah/goldie/v2"
	"github.com/spf13/afero"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/client/core"
	clientFilesystem "github.com/oslokommune/okctl/pkg/client/core/store/filesystem"
	"github.com/oslokommune/okctl/pkg/client/store"
	"github.com/oslokommune/okctl/pkg/spinner"
	"gotest.tools/assert"
	//fsKust "sigs.k8s.io/kustomize/pkg/fs"
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
	testOutputBuffer := bytes.NewBufferString("")
	testInputBuffer := bytes.NewBufferString(defaultTemplate)

	spin, _ := spinner.New("", testOutputBuffer)
	//aeroFs := afero.Afero{Fs: afero.NewMemMapFs()}
	osFs := afero.Afero{Fs: afero.NewOsFs()}
	aeroFs := afero.Afero{Fs: afero.NewBasePathFs(osFs, "/home/yngvar/yk/git/oslokommune/julius_iac")}

	mockPaths := clientFilesystem.Paths{BaseDir: "infrastructure/applications"}

	service := core.NewApplicationService(
		&aeroFs,
		spin,
		mockPaths,
		mockCertService{},
		clientFilesystem.NewApplicationStore(mockPaths, &aeroFs),
		mockAppReporter{},
	)

	env := "test"
	err := service.ScaffoldApplication(context.Background(), &client.ScaffoldApplicationOpts{
		In:                  testInputBuffer,
		Out:                 testOutputBuffer,
		ApplicationFilePath: "-",
		RepoDir:             "infrastructure",
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
	})
	assert.NilError(t, err)

	g := goldie.New(t)
	g.Assert(t, "kustomization-base.yaml", readFile(t, &aeroFs, filepath.Join(mockPaths.BaseDir, "my-app", config.DefaultApplicationBaseDir, "kustomization.yaml")))

	g.Assert(t, "deployment.yaml", readFile(t, &aeroFs, filepath.Join(mockPaths.BaseDir, "my-app", config.DefaultApplicationBaseDir, "deployment.yaml")))
	g.Assert(t, "argocd-application.yaml", readFile(t, &aeroFs, filepath.Join(mockPaths.BaseDir, "my-app", config.DefaultApplicationBaseDir, "argocd-application.yaml")))
	g.Assert(t, "volumes.yaml", readFile(t, &aeroFs, filepath.Join(mockPaths.BaseDir, "my-app", config.DefaultApplicationBaseDir, "volumes.yaml")))
	g.Assert(t, "ingress.yaml", readFile(t, &aeroFs, filepath.Join(mockPaths.BaseDir, "my-app", config.DefaultApplicationBaseDir, "ingress.yaml")))
	g.Assert(t, "service.yaml", readFile(t, &aeroFs, filepath.Join(mockPaths.BaseDir, "my-app", config.DefaultApplicationBaseDir, "service.yaml")))

	g.Assert(t, "kustomization-overlay.yaml", readFile(t, &aeroFs, filepath.Join(mockPaths.BaseDir, "my-app", config.DefaultApplicationOverlayDir, env, "kustomization.yaml")))
	g.Assert(t, "ingress-patch.yaml", readFile(t, &aeroFs, filepath.Join(mockPaths.BaseDir, "my-app", config.DefaultApplicationOverlayDir, env, "ingress-patch.json")))

	var buf bytes.Buffer

	kustomizeFs := KustomizeFs{
		afero: aeroFs,
	}

	// TODO: This fails because the kustomize dependency is 2.0.3, a really old version of kustomize.
	// See https://github.com/kubernetes/kubernetes/pull/98946.
	// We should attempt to use the kustomize lib directly.
	// Useful links:
	// https://github.com/kubernetes-sigs/kustomize/issues/142#issuecomment-426054466
	// https://github.com/kubernetes-sigs/kustomize/blob/master/kustomize/commands/build/build.go
	// https://pkg.go.dev/search?q=kustomize
	err = kustomize.RunKustomizeBuild(&buf, kustomizeFs, filepath.Join(mockPaths.BaseDir, "my-app", config.DefaultApplicationOverlayDir, env))
	assert.NilError(t, err)
	fmt.Println(err)
}

func readFile(t *testing.T, fs *afero.Afero, path string) []byte {
	result, err := fs.ReadFile(path)
	assert.NilError(t, err)

	return result
}

type mockCertService struct{}

func (m mockCertService) DeleteCertificate(_ context.Context, _ api.DeleteCertificateOpts) error {
	return nil
}

func (m mockCertService) DeleteCognitoCertificate(_ context.Context, _ api.DeleteCognitoCertificateOpts) error {
	return nil
}

func (m mockCertService) CreateCertificate(_ context.Context, _ api.CreateCertificateOpts) (*api.Certificate, error) {
	return &api.Certificate{
		CertificateARN: "arn:which:isnt:an:arn",
	}, nil
}

type mockAppReporter struct{}

func (m mockAppReporter) ReportCreateApplication(_ *client.ScaffoldedApplication, _ []*store.Report) error {
	return nil
}

func (m mockAppReporter) ReportDeleteApplication(_ []*store.Report) error {
	return nil
}

// KustomizeFs

type KustomizeFs struct {
	afero afero.Afero
}

func (c KustomizeFs) Create(name string) (fsKust.File, error) {
	return c.afero.Create(name)
}

func (c KustomizeFs) Mkdir(name string) error {
	return c.afero.Mkdir(name, 0x744)
}

func (c KustomizeFs) MkdirAll(name string) error {
	return c.afero.MkdirAll(name, 0x744)
}

func (c KustomizeFs) RemoveAll(name string) error {
	return c.afero.RemoveAll(name)
}

func (c KustomizeFs) Open(name string) (fsKust.File, error) {
	return c.afero.Open(name)
}

func (c KustomizeFs) IsDir(name string) bool {
	isDir, _ := c.afero.IsDir(name)
	return isDir
}

func (c KustomizeFs) CleanedAbs(path string) (fsKust.ConfirmedDir, string, error) {
	// Copyed from fakefs
	if c.IsDir(path) {
		return fsKust.ConfirmedDir(path), "", nil
	}
	d := filepath.Dir(path)
	if d == path {
		return fsKust.ConfirmedDir(d), "", nil
	}
	return fsKust.ConfirmedDir(d), filepath.Base(path), nil
}

func (c KustomizeFs) Exists(name string) bool {
	exists, _ := c.afero.Exists(name)
	return exists
}

func (c KustomizeFs) Glob(pattern string) ([]string, error) {
	return afero.Glob(c.afero, pattern)
}

func (c KustomizeFs) ReadFile(name string) ([]byte, error) {
	return c.afero.ReadFile(name)
}

func (c KustomizeFs) WriteFile(name string, data []byte) error {
	return c.afero.WriteFile(name, data, 0x644)
}
