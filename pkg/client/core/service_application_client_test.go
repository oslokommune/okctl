package core_test

import (
	"bytes"
	"context"
	"io/ioutil"
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
)

const defaultTemplate = `
# A name that identifies your app
name: my-app
# An URI for your app Docker image
image: docker.pkg.github.com/my-org/my-repo/my-package
# The version of your app which is available as an image
version: 0.0.1

# The URL your app should be available on
# Change to something other than https to disable configuring TLS
# Comment this out to avoid setting up an ingress
url: https://my-domain.io

# The port your app listens on
# Comment this out to avoid setting up a service (required if url is specified)
port: 3000

# How many replicas of your application should we scaffold
#replicas: 3 # 1 by default

# A namespace where your app will live
#namespace: my-namespace

# A Docker repository secret for pulling your image
#imagePullSecret: my-pull-secret-name

# The environment your app requires
#environment:
#  MY_VARIABLE: my-value

# Volumes to mount
#volumes:
#  - /path/to/mount/volume: # Requests 1Gi by default
#  - /path/to/mount/volume: 24Gi

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
	fs := afero.Afero{Fs: afero.NewMemMapFs()}
	mockPaths := clientFilesystem.Paths{BaseDir: "infrastructure/base"}

	service := core.NewApplicationService(
		spin,
		mockPaths,
		mockCertService{},
		clientFilesystem.NewApplicationStore(mockPaths, &fs),
		mockAppReporter{},
	)

	err := service.ScaffoldApplication(context.Background(), &client.ScaffoldApplicationOpts{
		In:                  testInputBuffer,
		Out:                 testOutputBuffer,
		ApplicationFilePath: "-",
		RepoDir:             "infrastructure",
		ID: &api.ID{
			Region:       "eu-west-1",
			AWSAccountID: "012345678912",
			Environment:  "test",
			Repository:   "not blank",
			ClusterName:  "dummy-cluster",
		},
		HostedZoneID: "dummyID",
		IACRepoURL:   "git@dummy.com:test/repo.git",
	})
	assert.NilError(t, err)

	k8sResourcePath := filepath.Join(mockPaths.BaseDir, "my-app", "my-app.yaml")
	k8sResources, err := fs.Open(k8sResourcePath)
	assert.NilError(t, err)

	k8sResourcesAsBytes, err := ioutil.ReadAll(k8sResources)
	assert.NilError(t, err)

	argocdResourcePath := filepath.Join(mockPaths.BaseDir, "my-app", "my-app-application.yaml")
	argocdResources, err := fs.Open(argocdResourcePath)
	assert.NilError(t, err)

	argocdResourcesAsBytes, err := ioutil.ReadAll(argocdResources)
	assert.NilError(t, err)

	g := goldie.New(t)
	g.Assert(t, "my-app.yaml", k8sResourcesAsBytes)
	g.Assert(t, "my-app-application.yaml", argocdResourcesAsBytes)
}

type mockCertService struct{}

func (m mockCertService) DeleteCertificate(_ context.Context, _ api.DeleteCertificateOpts) error {
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
