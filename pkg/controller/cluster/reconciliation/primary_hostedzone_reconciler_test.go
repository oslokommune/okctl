package reconciliation

import (
	"context"
	"testing"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/client"
	clientCore "github.com/oslokommune/okctl/pkg/client/core"
	"github.com/stretchr/testify/assert"
)

//nolint:funlen
func TestPrimaryHostedZoneReconciler(t *testing.T) {
	// Componentflag is always true for primary hosted zone
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
			name:                "Should delete when purge and existing",
			withPurge:           true,
			withComponentExists: true,
			expectCreations:     0,
			expectDeletions:     1,
		},
		{
			name:                "Should create when not purge and not existing",
			withPurge:           false,
			withComponentExists: false,
			expectCreations:     1,
			expectDeletions:     0,
		},
		{
			name:                "Should noop when not purge and existing",
			withComponentExists: true,
			expectCreations:     0,
			expectDeletions:     0,
		},
		{
			name:                "Should noop when purge and not existing",
			withPurge:           true,
			withComponentExists: false,
			expectCreations:     0,
			expectDeletions:     0,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			creations := 0
			deletions := 0

			meta := generateTestMeta(tc.withPurge, v1alpha1.ClusterIntegrations{})

			state := &clientCore.StateHandlers{
				Cluster: &mockClusterState{exists: tc.withDependenciesMet},
				Domain:  &mockDomainState{exists: tc.withComponentExists},
			}

			reconciler := NewZoneReconciler(&mockDomainService{
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

type mockDomainService struct {
	creationBump func()
	deletionBump func()
}

func (m mockDomainService) CreatePrimaryHostedZone(_ context.Context, _ client.CreatePrimaryHostedZoneOpts) (*client.HostedZone, error) {
	m.creationBump()

	return nil, nil
}

func (m mockDomainService) DeletePrimaryHostedZone(_ context.Context, _ client.DeletePrimaryHostedZoneOpts) error {
	m.deletionBump()

	return nil
}

func (m mockDomainService) GetPrimaryHostedZone(_ context.Context) (*client.HostedZone, error) {
	panic("implement me")
}

func (m mockDomainService) SetHostedZoneDelegation(_ context.Context, _ string, _ bool) error {
	panic("implement me")
}
