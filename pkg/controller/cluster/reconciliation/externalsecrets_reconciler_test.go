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
func TestExternalSecretsReconciler(t *testing.T) {
	testCases := []generalizedTestCase{
		{
			name:                "Should noop when requested and already existing",
			withComponentFlag:   true,
			withComponentExists: true,
			withDependenciesMet: true,
			expectCreations:     0,
			expectDeletions:     0,
		},
		{
			name:                "Should noop when not requested and not existing",
			withComponentFlag:   false,
			withComponentExists: false,
			withDependenciesMet: true,
			expectCreations:     0,
			expectDeletions:     0,
		},
		{
			name:                "Should noop when indicated but missing dependencies",
			withComponentFlag:   true,
			withComponentExists: false,
			withDependenciesMet: false,
			expectCreations:     0,
			expectDeletions:     0,
		},
		{
			name:                "Should delete when indicated but purge",
			withPurge:           true,
			withComponentFlag:   true,
			withComponentExists: true,
			withDependenciesMet: true,
			expectCreations:     0,
			expectDeletions:     1,
		},
		{
			name:                "Should create when indicated and not existing",
			withComponentFlag:   true,
			withComponentExists: false,
			withDependenciesMet: true,
			expectCreations:     1,
			expectDeletions:     0,
		},
		{
			name:                "Should delete when de indicated and existing",
			withComponentFlag:   false,
			withComponentExists: true,
			withDependenciesMet: true,
			expectCreations:     0,
			expectDeletions:     1,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			creations := 0
			deletions := 0

			meta := generateTestMeta(tc.withPurge, v1alpha1.ClusterIntegrations{ExternalSecrets: tc.withComponentFlag})

			state := &clientCore.StateHandlers{
				Cluster:         &mockClusterState{exists: tc.withDependenciesMet},
				ExternalSecrets: &mockExternalSecretsState{exists: tc.withComponentExists},
			}

			reconciler := NewExternalSecretsReconciler(&mockExternalSecretsService{
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

type mockExternalSecretsService struct {
	creationBump func()
	deletionBump func()
}

func (m mockExternalSecretsService) CreateExternalSecrets(_ context.Context, _ client.CreateExternalSecretsOpts) (*client.ExternalSecrets, error) {
	m.creationBump()

	return nil, nil
}

func (m mockExternalSecretsService) DeleteExternalSecrets(_ context.Context, _ api.ID) error {
	m.deletionBump()

	return nil
}

type mockExternalSecretsState struct {
	exists bool
}

func (m mockExternalSecretsState) HasExternalSecrets() (bool, error) {
	return m.exists, nil
}
