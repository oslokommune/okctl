package handlers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"path"
	"strings"
	"testing"

	"github.com/oslokommune/okctl/pkg/clients/kubectl"
	"github.com/spf13/afero"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/client/core"

	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
)

func TestWriteDeleteApplicationReadyCheckInfo(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		withIngress       bool
		withECRRepository bool
	}{
		{
			name: "Should work with no optionals",
		},
		{
			name:        "Should add additional information with ingress",
			withIngress: true,
		},
		{
			name:              "Should add additional information with ECR",
			withECRRepository: true,
		},
		{
			name:              "Should include all information when everything is active",
			withECRRepository: true,
			withIngress:       true,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			stdout := bytes.Buffer{}

			err := writeDeleteApplicationReadyCheckInfo(&stdout, deleteApplicationPromptTemplateOpts{
				ApplicationName: "mockapp",
				HasIngress:      tc.withIngress,
				HasECR:          tc.withECRRepository,
			})
			assert.NoError(t, err)

			g := goldie.New(t)
			g.Assert(t, t.Name(), stdout.Bytes())
		})
	}
}

func TestCorrectBehaviourWithCertainFlags(t *testing.T) {
	testCases := []struct {
		name string

		withPurge bool
	}{
		{
			name:      "Should do create/update-ish kind of operations when purge is false",
			withPurge: false,
		},
		{
			name:      "Should do delete/purge-ish kind of operations when purge is true",
			withPurge: true,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			stdout := bytes.Buffer{}

			services, evaluationFn := mockServices(t)

			err := HandleApplication(&HandleApplicationOpts{
				Out:  &stdout,
				File: "-",
				ClusterManifest: v1alpha1.Cluster{
					Integrations: &v1alpha1.ClusterIntegrations{},
				},
				ApplicationManifest: v1alpha1.Application{
					Metadata: v1alpha1.ApplicationMeta{
						Name:      "mockApp",
						Namespace: "mockNamespace",
					},
				},
				Services:      services,
				State:         mockStates(),
				Purge:         tc.withPurge,
				Confirm:       true,
				DelayFunction: func() {},
			})(nil, nil)
			assert.NoError(t, err)

			assert.Equal(t, tc.withPurge, evaluationFn())
		})
	}
}

type evaluator func() bool

func mockServices(t *testing.T) (*core.Services, evaluator) {
	service := mockService{}

	return &core.Services{
			ApplicationService:         &service,
			ApplicationPostgresService: &service,
			Certificate:                &service,
			Domain:                     &service,
			ContainerRepository:        &service,
		}, func() bool {
			if service.scaffoldApplicationCalled == service.deleteApplicationManifestsCalled {
				t.Fatal("received no clear indication of create actions or delete actions")
			}

			return service.deleteApplicationManifestsCalled
		}
}

func mockStates() *core.StateHandlers {
	return &core.StateHandlers{
		Certificate:         mockState{},
		Kubernetes:          mockState{},
		Domain:              mockState{},
		Github:              mockState{},
		ContainerRepository: mockState{},
	}
}

type mockState struct{}

func (m mockState) SaveContainerRepository(*client.ContainerRepository) error {
	panic("implement me")
}
func (m mockState) RemoveContainerRepository(string) error { panic("implement me") }
func (m mockState) GetContainerRepository(string) (*client.ContainerRepository, error) {
	panic("implement me")
}

func (m mockState) GetContainerRepositoryFromApplication(string) (*client.ContainerRepository, error) {
	panic("implement me")
}
func (m mockState) ApplicationHasImage(string) (bool, error) { return false, nil }

func (m mockState) SaveGithubRepository(*client.GithubRepository) error {
	panic("implement me")
}

func (m mockState) GetGithubRepository(string) (*client.GithubRepository, error) {
	return &client.GithubRepository{}, nil
}
func (m mockState) RemoveGithubRepository(string) error { panic("implement me") }

func (m mockState) SaveHostedZone(*client.HostedZone) error           { panic("implement me") }
func (m mockState) UpdateHostedZone(*client.HostedZone) error         { panic("implement me") }
func (m mockState) RemoveHostedZone(string) error                     { panic("implement me") }
func (m mockState) GetHostedZone(string) (*client.HostedZone, error)  { panic("implement me") }
func (m mockState) GetPrimaryHostedZone() (*client.HostedZone, error) { panic("implement me") }
func (m mockState) GetHostedZones() ([]*client.HostedZone, error)     { panic("implement me") }
func (m mockState) HasPrimaryHostedZone() (bool, error)               { return true, nil }

func (m mockState) HasResource(string, string, string) (bool, error) {
	return false, nil
}

func (m mockState) SaveCertificate(*client.Certificate) error {
	panic("implement me")
}

func (m mockState) GetCertificate(string) (*client.Certificate, error) {
	return &client.Certificate{}, nil
}

func (m mockState) HasCertificate(string) (bool, error) {
	return false, nil
}

func (m mockState) RemoveCertificate(string) error {
	panic("implement me")
}

type mockService struct {
	scaffoldApplicationCalled        bool
	deleteApplicationManifestsCalled bool
}

func (m mockService) CreateContainerRepository(context.Context, client.CreateContainerRepositoryOpts) (*client.ContainerRepository, error) {
	panic("implement me")
}

func (m mockService) DeleteContainerRepository(context.Context, client.DeleteContainerRepositoryOpts) error {
	panic("implement me")
}

func (m mockService) EmptyContainerRepository(context.Context, client.EmptyContainerRepositoryOpts) error {
	panic("implement me")
}

func (m mockService) AddPostgresToApplication(context.Context, client.AddPostgresToApplicationOpts) error {
	panic("implement me")
}

func (m mockService) RemovePostgresFromApplication(context.Context, client.RemovePostgresFromApplicationOpts) error {
	panic("implement me")
}

func (m mockService) HasPostgresIntegration(context.Context, client.HasPostgresIntegrationOpts) (bool, error) {
	return false, nil
}

func (m *mockService) ScaffoldApplication(context.Context, *client.ScaffoldApplicationOpts) error {
	m.scaffoldApplicationCalled = true

	return nil
}

func (m *mockService) DeleteApplicationManifests(context.Context, client.DeleteApplicationManifestsOpts) error {
	m.deleteApplicationManifestsCalled = true

	return nil
}

func (m mockService) CreateArgoCDApplicationManifest(client.CreateArgoCDApplicationManifestOpts) error {
	panic("implement me")
}

func (m mockService) DeleteArgoCDApplicationManifest(client.DeleteArgoCDApplicationManifestOpts) error {
	return nil
}

func (m mockService) HasArgoCDIntegration(context.Context, client.HasArgoCDIntegrationOpts) (bool, error) {
	panic("implement me")
}

func (m mockService) CreatePrimaryHostedZone(_ context.Context, _ client.CreatePrimaryHostedZoneOpts) (*client.HostedZone, error) {
	panic("implement me")
}

func (m mockService) GetPrimaryHostedZone(_ context.Context) (*client.HostedZone, error) {
	panic("implement me")
}

func (m mockService) DeletePrimaryHostedZone(_ context.Context, _ client.DeletePrimaryHostedZoneOpts) error {
	panic("implement me")
}

func (m mockService) SetHostedZoneDelegation(_ context.Context, _ string, _ bool) error {
	panic("implement me")
}

func (m mockService) CreateCertificate(_ context.Context, _ client.CreateCertificateOpts) (*client.Certificate, error) {
	panic("implement me")
}

func (m mockService) DeleteCertificate(_ context.Context, _ client.DeleteCertificateOpts) error {
	return nil
}

func (m mockService) DeleteCognitoCertificate(_ context.Context, _ client.DeleteCognitoCertificateOpts) error {
	panic("implement me")
}

const defaultFolderPermissions = 0o700

func createCluster(t *testing.T, fs *afero.Afero, name string) {
	absArgoCDConfigurationDir := path.Join("/", "infrastructure", name, "argocd")

	err := fs.MkdirAll(path.Join(absArgoCDConfigurationDir, "applications"), defaultFolderPermissions)
	assert.NoError(t, err)

	err = fs.MkdirAll(path.Join(absArgoCDConfigurationDir, "namespaces"), defaultFolderPermissions)
	assert.NoError(t, err)
}

func createApp(t *testing.T, fs *afero.Afero, name string) {
	absAppDir := path.Join("/", "infrastructure", "applications", name)

	err := fs.MkdirAll(path.Join(absAppDir, "base"), defaultFolderPermissions)
	assert.NoError(t, err)

	err = fs.MkdirAll(path.Join(absAppDir, "overlays"), defaultFolderPermissions)
	assert.NoError(t, err)
}

func addAppToCluster(t *testing.T, fs *afero.Afero, appName string, clusterName string) {
	absAppDir := path.Join("/", "infrastructure", appName)
	absOverlaysDir := path.Join(absAppDir, clusterName)
	absArgoCDConfigurationDir := path.Join("/", "infrastructure", clusterName, "argocd")

	err := fs.MkdirAll(path.Join(absOverlaysDir), defaultFolderPermissions)
	assert.NoError(t, err)

	err = fs.WriteReader(path.Join(absOverlaysDir, "kustomization.yaml"), strings.NewReader(""))
	assert.NoError(t, err)

	err = fs.WriteReader(
		path.Join(absArgoCDConfigurationDir, "applications", fmt.Sprintf("%s.yaml", appName)),
		strings.NewReader(""),
	)
	assert.NoError(t, err)
}

func initializeEnvironment(t *testing.T, fs *afero.Afero, applications []testApplication) {
	for _, app := range applications {
		createApp(t, fs, app.Name)

		for _, cluster := range app.Clusters {
			createCluster(t, fs, cluster)

			addAppToCluster(t, fs, app.Name, cluster)
		}
	}
}

func mockCluster(name string) v1alpha1.Cluster {
	cluster := v1alpha1.NewCluster()

	cluster.Metadata.Name = name
	cluster.Github.OutputPath = "infrastructure"

	return cluster
}

func mockApp(name string) v1alpha1.Application {
	return v1alpha1.Application{
		TypeMeta: v1alpha1.ApplicationTypeMeta(),
		Metadata: v1alpha1.ApplicationMeta{
			Name:      name,
			Namespace: "mock-namespace",
		},
	}
}

type testApplication struct {
	Name     string
	Clusters []string
}

//nolint:funlen
func TestDeleteApplicationFiles(t *testing.T) {
	testCases := []struct {
		name                   string
		withCurrentCluster     v1alpha1.Cluster
		withCurrentApplication v1alpha1.Application
		withApplications       []testApplication
		expectExistentFiles    []string
		expectNonExistentFiles []string
	}{
		{
			name:                   "Should delete expected files with one app and one cluster",
			withCurrentCluster:     mockCluster("cluster-one"),
			withCurrentApplication: mockApp("app-one"),
			withApplications: []testApplication{
				{
					Name:     "app-one",
					Clusters: []string{"cluster-one"},
				},
			},
			expectNonExistentFiles: []string{
				"/infrastructure/applications/app-one",
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			fs := &afero.Afero{Fs: afero.NewMemMapFs()}

			initializeEnvironment(t, fs, tc.withApplications)

			for _, item := range tc.expectNonExistentFiles {
				exists, err := fs.Exists(item)
				assert.NoError(t, err)

				assert.Equal(t, true, exists)
			}

			mockS, _ := mockServices(t)

			err := HandleApplication(&HandleApplicationOpts{
				Out:   io.Discard,
				Err:   io.Discard,
				Ctx:   context.Background(),
				State: mockStates(),
				Services: &core.Services{
					ApplicationService: core.NewApplicationService(
						fs,
						&mockKubectlClient{},
						nil,
						"/",
						func(_ string, _ string, _ string) error { return nil },
					),
					ApplicationPostgresService: mockS.ApplicationPostgresService,
					Certificate:                mockS.Certificate,
					Domain:                     mockS.Domain,
					ContainerRepository:        mockS.ContainerRepository,
				},
				File:                "-",
				ClusterManifest:     tc.withCurrentCluster,
				ApplicationManifest: tc.withCurrentApplication,
				Purge:               true,
				Confirm:             true,
				DelayFunction:       func() {},
			})(nil, nil)
			assert.NoError(t, err)

			for _, item := range tc.expectExistentFiles {
				exists, err := fs.Exists(item)
				assert.NoError(t, err)

				assert.Equal(t, true, exists)
			}

			for _, item := range tc.expectNonExistentFiles {
				exists, err := fs.Exists(item)
				assert.NoError(t, err)

				assert.Equal(t, false, exists)
			}
		})
	}
}

type mockKubectlClient struct{}

func (m mockKubectlClient) Apply(_ io.Reader) error { panic("implement me") }

func (m mockKubectlClient) Delete(_ io.Reader) error { return nil }

func (m mockKubectlClient) Patch(_ kubectl.PatchOpts) error { return nil }

func (m mockKubectlClient) Exists(_ kubectl.Resource) (bool, error) { panic("implement me") }
