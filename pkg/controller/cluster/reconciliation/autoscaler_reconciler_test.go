package reconciliation

import (
	"context"
	"testing"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	clientCore "github.com/oslokommune/okctl/pkg/client/core"
	"github.com/stretchr/testify/assert"
)

//nolint:funlen
func TestAutoscalerReconciler(t *testing.T) {
	testCases := []struct {
		name                string
		withComponentFlag   bool
		withComponentExists bool
		withDependenciesMet bool
		withPurge           bool
		expectCreations     int
		expectDeletions     int
	}{
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
			name:                "Should create when dependencies met, not existing and requested",
			withComponentFlag:   true,
			withComponentExists: false,
			withDependenciesMet: true,
			expectCreations:     1,
			expectDeletions:     0,
		},
		{
			name:                "Should delete when dependencies met, existing and de-requested",
			withComponentFlag:   false,
			withComponentExists: true,
			withDependenciesMet: true,
			expectCreations:     0,
			expectDeletions:     1,
		},
		{
			name:                "Should delete when dependencies met, existing, requested but with purge",
			withPurge:           true,
			withComponentFlag:   true,
			withComponentExists: true,
			withDependenciesMet: true,
			expectCreations:     0,
			expectDeletions:     1,
		},
		{
			name:                "Should noop when missing dependencies, not existing but requested",
			withPurge:           true,
			withComponentFlag:   true,
			withComponentExists: false,
			withDependenciesMet: false,
			expectCreations:     0,
			expectDeletions:     0,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			creations := 0
			deletions := 0

			meta := generateTestMeta(tc.withPurge, v1alpha1.ClusterIntegrations{Autoscaler: tc.withComponentFlag})

			state := &clientCore.StateHandlers{
				Autoscaler: mockAutoscalerState{exists: tc.withComponentExists},
				Cluster:    mockClusterState{exists: tc.withDependenciesMet},
			}

			reconciler := NewAutoscalerReconciler(
				&mockAutoscalerService{
					creationBump: func() { creations++ },
					deletionBump: func() { deletions++ },
				},
			)

			_, err := reconciler.Reconcile(context.Background(), meta, state)
			assert.NoError(t, err)

			assert.Equal(t, tc.expectCreations, creations)
			assert.Equal(t, tc.expectDeletions, deletions)
		})
	}
}

type mockAutoscalerService struct {
	creationBump func()
	deletionBump func()
}

func (m mockAutoscalerService) CreateAutoscaler(_ context.Context, _ client.CreateAutoscalerOpts) (*client.Autoscaler, error) {
	m.creationBump()

	return nil, nil
}

func (m mockAutoscalerService) DeleteAutoscaler(_ context.Context, _ api.ID) error {
	m.deletionBump()

	return nil
}

type mockAutoscalerState struct {
	exists bool
}

func (m mockAutoscalerState) HasAutoscaler() (bool, error) {
	return m.exists, nil
}
