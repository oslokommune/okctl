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
func TestAWSLoadBalancerControllerReconciler(t *testing.T) {
	testCases := []struct {
		name                      string
		withPurge                 bool
		withComponentFlag         bool
		withComponentExists       bool
		withCreateDependenciesMet bool
		withDeleteDependenciesMet bool
		expectCreations           int
		expectDeletions           int
	}{
		{
			name:                      "Should noop when requested and already existing",
			withComponentFlag:         true,
			withComponentExists:       true,
			withCreateDependenciesMet: true,
			expectCreations:           0,
			expectDeletions:           0,
		},
		{
			name:                      "Should noop when not requested and not existing",
			withComponentFlag:         false,
			withComponentExists:       false,
			withDeleteDependenciesMet: true,
			expectCreations:           0,
			expectDeletions:           0,
		},
		{
			name:                      "Should noop when indicated, not pre existing but missing dependencies",
			withComponentFlag:         true,
			withComponentExists:       false,
			withCreateDependenciesMet: false,
			expectCreations:           0,
			expectDeletions:           0,
		},
		{
			name:                      "Should create when indicated and not pre existing",
			withComponentFlag:         true,
			withComponentExists:       false,
			withCreateDependenciesMet: true,
			expectCreations:           1,
			expectDeletions:           0,
		},
		{
			name:                      "Should delete when de-indicated and pre existing",
			withComponentFlag:         false,
			withComponentExists:       true,
			withCreateDependenciesMet: true,
			withDeleteDependenciesMet: true,
			expectCreations:           0,
			expectDeletions:           1,
		},
		{
			name:                      "Should delete when indicated, pre existing but purge",
			withPurge:                 true,
			withComponentFlag:         true,
			withComponentExists:       true,
			withCreateDependenciesMet: true,
			withDeleteDependenciesMet: true,
			expectCreations:           0,
			expectDeletions:           1,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			creations := 0
			deletions := 0

			meta := generateTestMeta(tc.withPurge, v1alpha1.ClusterIntegrations{AWSLoadBalancerController: tc.withComponentFlag})

			state := &clientCore.StateHandlers{
				Cluster:                   &mockClusterState{exists: tc.withCreateDependenciesMet || tc.withDeleteDependenciesMet},
				Vpc:                       &mockVPCState{},
				AWSLoadBalancerController: &mockAWSLoadBalancerState{exists: tc.withComponentExists},
				ArgoCD:                    &mockArgoCDState{exists: !tc.withDeleteDependenciesMet},
				Monitoring:                &mockMonitoringState{exists: !tc.withDeleteDependenciesMet},
				Application:               &mockApplicationState{existingApplications: 0},
			}

			reconciler := NewAWSLoadBalancerControllerReconciler(&mockAWSLoadBalancerControllerService{
				creationBump: func() { creations++ },
				deletionBump: func() { deletions++ },
			})

			_, err := reconciler.Reconcile(context.Background(), meta, state)
			assert.NoError(t, err)

			assert.Equal(t, tc.expectCreations, creations, "creations")
			assert.Equal(t, tc.expectDeletions, deletions, "deletions")
		})
	}
}

type mockAWSLoadBalancerControllerService struct {
	creationBump func()
	deletionBump func()
}

func (m mockAWSLoadBalancerControllerService) CreateAWSLoadBalancerController(
	_ context.Context,
	_ client.CreateAWSLoadBalancerControllerOpts,
) (*client.AWSLoadBalancerController, error) {
	m.creationBump()

	return nil, nil
}

func (m mockAWSLoadBalancerControllerService) DeleteAWSLoadBalancerController(_ context.Context, _ api.ID) error {
	m.deletionBump()

	return nil
}

type mockAWSLoadBalancerState struct {
	exists bool
}

func (a mockAWSLoadBalancerState) HasAWSLoadBalancerController() (bool, error) {
	return a.exists, nil
}
