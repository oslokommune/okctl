package handlers

import (
	"bytes"
	"context"
	"testing"

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
