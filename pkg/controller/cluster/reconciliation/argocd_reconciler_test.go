package reconciliation

import (
	"context"
	"testing"

	"github.com/oslokommune/okctl/pkg/clients/kubectl"

	"github.com/oslokommune/okctl/pkg/github"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/client"
	clientCore "github.com/oslokommune/okctl/pkg/client/core"
	"github.com/stretchr/testify/assert"
)

//nolint:funlen
func TestArgocdReconciler_Reconcile(t *testing.T) {
	testCases := []struct {
		name                           string
		withComponentFlag              bool
		withAlreadyExists              bool
		withClusterExists              bool
		withPrimaryHostedZoneExists    bool
		withPrimaryHostedZoneDelegated bool
		expectCreations                int
		expectDeletions                int
		withPurge                      bool
	}{
		{
			name:                           "Should noop when requested and existing",
			withComponentFlag:              true,
			withAlreadyExists:              true,
			withClusterExists:              true,
			withPrimaryHostedZoneExists:    true,
			withPrimaryHostedZoneDelegated: true,
			expectCreations:                0,
			expectDeletions:                0,
		},
		{
			name:                           "Should noop when not requested and not existing",
			withComponentFlag:              false,
			withAlreadyExists:              false,
			withClusterExists:              true,
			withPrimaryHostedZoneExists:    true,
			withPrimaryHostedZoneDelegated: true,
			expectCreations:                0,
			expectDeletions:                0,
		},
		{
			name:                           "Should create when dependencies met and requested",
			withComponentFlag:              true,
			withAlreadyExists:              false,
			withClusterExists:              true,
			withPrimaryHostedZoneExists:    true,
			withPrimaryHostedZoneDelegated: true,
			expectCreations:                1,
			expectDeletions:                0,
		},
		{
			name:                           "Should noop when dependencies met and not requested",
			withComponentFlag:              false,
			withAlreadyExists:              false,
			withClusterExists:              true,
			withPrimaryHostedZoneExists:    true,
			withPrimaryHostedZoneDelegated: true,
			expectCreations:                0,
			expectDeletions:                0,
		},
		{
			name:                           "Should noop when dependencies not met but requested",
			withComponentFlag:              true,
			withAlreadyExists:              false,
			withClusterExists:              true,
			withPrimaryHostedZoneExists:    false,
			withPrimaryHostedZoneDelegated: true,
			expectCreations:                0,
			expectDeletions:                0,
		},
		{
			name:                           "Should delete when existing and toggled false",
			withComponentFlag:              false,
			withAlreadyExists:              true,
			withClusterExists:              true,
			withPrimaryHostedZoneExists:    true,
			withPrimaryHostedZoneDelegated: true,
			expectCreations:                0,
			expectDeletions:                1,
		},
		{
			name:                           "Should delete when existing and purge is true",
			withPurge:                      true,
			withComponentFlag:              true,
			withAlreadyExists:              true,
			withClusterExists:              true,
			withPrimaryHostedZoneExists:    true,
			withPrimaryHostedZoneDelegated: true,
			expectCreations:                0,
			expectDeletions:                1,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			creations := 0
			deletions := 0

			meta := generateTestMeta(tc.withPurge, v1alpha1.ClusterIntegrations{ArgoCD: tc.withComponentFlag})

			state := &clientCore.StateHandlers{
				ArgoCD:          &mockArgoCDState{exists: tc.withAlreadyExists},
				Cluster:         &mockClusterState{exists: tc.withClusterExists},
				IdentityManager: &mockArgoIdentityManagerState{},
				Domain: &mockDomainState{
					exists:      tc.withPrimaryHostedZoneExists,
					isDelegated: tc.withPrimaryHostedZoneDelegated,
				},
			}

			reconciler := NewArgocdReconciler(
				&mockArgocdService{
					creationBump: func() { creations++ },
					deletionBump: func() { deletions++ },
				},
				&mockGithubService{},
			)

			_, err := reconciler.Reconcile(context.Background(), meta, state)
			assert.NoError(t, err)

			assert.Equal(t, tc.expectCreations, creations)
			assert.Equal(t, tc.expectDeletions, deletions)
		})
	}
}

type mockArgocdService struct {
	creationBump func()
	deletionBump func()
}

func (m mockArgocdService) SetupNamespacesSync(_ context.Context, _ kubectl.Client, _ v1alpha1.Cluster) error {
	return nil
}

func (m mockArgocdService) SetupApplicationsSync(_ context.Context, _ kubectl.Client, _ v1alpha1.Cluster) error {
	return nil
}

func (m mockArgocdService) CreateArgoCD(_ context.Context, _ client.CreateArgoCDOpts) (*client.ArgoCD, error) {
	m.creationBump()

	return nil, nil
}

func (m mockArgocdService) DeleteArgoCD(_ context.Context, _ client.DeleteArgoCDOpts) error {
	m.deletionBump()

	return nil
}

type mockGithubService struct{}

func (m mockGithubService) CreateGithubRepository(_ context.Context, _ client.CreateGithubRepositoryOpts) (*client.GithubRepository, error) {
	return nil, nil
}

func (m mockGithubService) DeleteGithubRepository(_ context.Context, _ client.DeleteGithubRepositoryOpts) error {
	return nil
}

func (m mockGithubService) ListReleases(_, _ string) ([]*github.RepositoryRelease, error) {
	return nil, nil
}

func (m mockGithubService) CreateRepositoryDeployKey(_ client.CreateGithubDeployKeyOpts) (*client.GithubDeployKey, error) {
	return nil, nil
}

func (m mockGithubService) DeleteRepositoryDeployKey(_ client.DeleteGithubDeployKeyOpts) error {
	return nil
}

type mockArgoCDState struct {
	exists bool
}

func (m mockArgoCDState) HasArgoCD() (bool, error) {
	return m.exists, nil
}

func (m mockArgoCDState) SaveArgoCD(_ *client.ArgoCD) error  { panic("implement me") }
func (m mockArgoCDState) GetArgoCD() (*client.ArgoCD, error) { panic("implement me") }
func (m mockArgoCDState) RemoveArgoCD() error                { panic("implement me") }

type mockArgoIdentityManagerState struct{}

func (m mockArgoIdentityManagerState) HasIdentityPool() (bool, error) { return true, nil }

func (m mockArgoIdentityManagerState) GetIdentityPool(_ string) (*client.IdentityPool, error) {
	return &client.IdentityPool{
		UserPoolID: "dummy.domain.io",
		AuthDomain: "dummy-id",
	}, nil
}

func (m mockArgoIdentityManagerState) SaveIdentityPool(_ *client.IdentityPool) error {
	panic("implement me")
}
func (m mockArgoIdentityManagerState) RemoveIdentityPool(_ string) error { panic("implement me") }

func (m mockArgoIdentityManagerState) SaveIdentityPoolClient(_ *client.IdentityPoolClient) error {
	panic("implement me")
}

func (m mockArgoIdentityManagerState) GetIdentityPoolClient(_ string) (*client.IdentityPoolClient, error) {
	panic("implement me")
}
func (m mockArgoIdentityManagerState) RemoveIdentityPoolClient(_ string) error { panic("implement me") }
func (m mockArgoIdentityManagerState) SaveIdentityPoolUser(_ *client.IdentityPoolUser) error {
	panic("implement me")
}

func (m mockArgoIdentityManagerState) GetIdentityPoolUser(_ string) (*client.IdentityPoolUser, error) {
	panic("implement me")
}
func (m mockArgoIdentityManagerState) RemoveIdentityPoolUser(_ string) error { panic("implement me") }
func (m mockArgoIdentityManagerState) GetIdentityPoolUsers() ([]client.IdentityPoolUser, error) {
	panic("implement me")
}
