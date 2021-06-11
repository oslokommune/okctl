package reconciliation

import (
	"context"
	"testing"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/client"
	clientCore "github.com/oslokommune/okctl/pkg/client/core"
	"github.com/stretchr/testify/assert"
)

//nolint:funlen
func TestExternalDNSReconciler(t *testing.T) {
	testCases := []struct {
		name                      string
		withPurge                 bool
		withComponentFlag         bool
		withComponentExists       bool
		withClusterExists         bool
		withCreateDependenciesMet bool
		withDeleteDependenciesMet bool
		expectCreations           int
		expectDeletions           int
	}{
		{
			name:                      "Should noop when requested and already existing",
			withComponentFlag:         true,
			withComponentExists:       true,
			withClusterExists:         true,
			withCreateDependenciesMet: true,
			expectCreations:           0,
			expectDeletions:           0,
		},
		{
			name:                      "Should noop when not requested and not existing",
			withComponentFlag:         false,
			withComponentExists:       false,
			withClusterExists:         true,
			withCreateDependenciesMet: true,
			expectCreations:           0,
			expectDeletions:           0,
		},
		{
			name:                "Should noop when indicated but missing dependencies",
			withComponentFlag:   true,
			withComponentExists: false,
			withClusterExists:   false,
			expectCreations:     0,
			expectDeletions:     0,
		},
		{
			name:                      "Should delete when indicated but purge",
			withPurge:                 true,
			withComponentFlag:         true,
			withComponentExists:       true,
			withClusterExists:         true,
			withDeleteDependenciesMet: true,
			expectCreations:           0,
			expectDeletions:           1,
		},
		{
			name:                      "Should create when indicated and not existing",
			withComponentFlag:         true,
			withComponentExists:       false,
			withClusterExists:         true,
			withCreateDependenciesMet: true,
			expectCreations:           1,
			expectDeletions:           0,
		},
		{
			name:                      "Should delete when de indicated and existing",
			withComponentFlag:         false,
			withComponentExists:       true,
			withClusterExists:         true,
			withDeleteDependenciesMet: true,
			expectCreations:           0,
			expectDeletions:           1,
		},
		{
			name:                      "Should noop when deleting and dependencies are not yet deleted",
			withPurge:                 true,
			withComponentFlag:         true,
			withComponentExists:       true,
			withClusterExists:         true,
			withDeleteDependenciesMet: false,
			expectCreations:           0,
			expectDeletions:           0,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			creations := 0
			deletions := 0

			meta := generateTestMeta(tc.withPurge, v1alpha1.ClusterIntegrations{ExternalDNS: tc.withComponentFlag})

			state := &clientCore.StateHandlers{
				ArgoCD:     &mockArgoCDState{exists: !tc.withDeleteDependenciesMet},
				Monitoring: &mockMonitoringState{exists: !tc.withDeleteDependenciesMet},

				Cluster:     &mockClusterState{exists: tc.withClusterExists},
				ExternalDNS: &mockExternalDNSState{exists: tc.withComponentExists},
				Domain: &mockDomainState{
					exists:      tc.withCreateDependenciesMet,
					isDelegated: tc.withCreateDependenciesMet,
				},
			}

			reconciler := NewExternalDNSReconciler(&mockExternalDNSService{
				creationBump: func() { creations++ },
				deletionBump: func() { deletions++ },
			})

			_, err := reconciler.Reconcile(context.Background(), meta, state)
			assert.NoError(t, err)

			assert.Equal(t, tc.expectCreations, creations)
			assert.Equal(t, tc.expectDeletions, deletions)
		})
	}
}

type mockExternalDNSService struct {
	creationBump func()
	deletionBump func()
}

func (m mockExternalDNSService) CreateExternalDNS(_ context.Context, _ client.CreateExternalDNSOpts) (*client.ExternalDNS, error) {
	m.creationBump()

	return nil, nil
}

func (m mockExternalDNSService) DeleteExternalDNS(_ context.Context, _ api.ID) error {
	m.deletionBump()

	return nil
}

type mockExternalDNSState struct {
	exists bool
}

func (m mockExternalDNSState) HasExternalDNS() (bool, error) {
	return m.exists, nil
}

func (m mockExternalDNSState) SaveExternalDNS(_ *client.ExternalDNS) error  { panic("implement me") }
func (m mockExternalDNSState) GetExternalDNS() (*client.ExternalDNS, error) { panic("implement me") }
func (m mockExternalDNSState) RemoveExternalDNS() error                     { panic("implement me") }
