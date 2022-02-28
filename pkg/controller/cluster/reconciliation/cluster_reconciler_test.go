package reconciliation

import (
	"context"
	"testing"

	"github.com/oslokommune/okctl/pkg/api"

	"github.com/oslokommune/okctl/pkg/config/constant"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/servicequotas"
	"github.com/oslokommune/okctl/pkg/mock"

	"github.com/oslokommune/okctl/pkg/client"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	clientCore "github.com/oslokommune/okctl/pkg/client/core"
	"github.com/stretchr/testify/assert"
)

//nolint:funlen
func TestClusterReconciler(t *testing.T) {
	testCases := []struct {
		name string

		withPurge                   bool
		withCreationDependenciesMet bool
		withDeletionDependenciesMet bool
		withComponentExists         bool
		withApplications            int

		expectCreations int
		expectDeletions int
	}{
		{
			name: "Should create when not purge and not exists",

			withCreationDependenciesMet: true,
			withComponentExists:         false,
			expectCreations:             1,
			expectDeletions:             0,
		},
		{
			name: "Should delete when purge and exists",

			withPurge:                   true,
			withDeletionDependenciesMet: true,
			withComponentExists:         true,
			expectCreations:             0,
			expectDeletions:             1,
		},
		{
			name: "Should noop when not existing and waiting for VPC",

			withComponentExists:         false,
			withCreationDependenciesMet: false,
			expectCreations:             0,
			expectDeletions:             0,
		},
		{
			name: "Should noop when existing with purge and waiting for cluster to shut down",

			withPurge:                   true,
			withComponentExists:         true,
			withDeletionDependenciesMet: false,
			expectCreations:             0,
			expectDeletions:             0,
		},
		{
			name: "Should wait when purge and applications exist",

			withApplications:            1,
			withPurge:                   true,
			withComponentExists:         true,
			withDeletionDependenciesMet: true,
			expectCreations:             0,
			expectDeletions:             0,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			creations := 0
			deletions := 0

			meta := generateTestMeta(tc.withPurge, v1alpha1.ClusterIntegrations{})
			meta.ClusterDeclaration.VPC = &v1alpha1.ClusterVPC{CIDR: constant.DefaultClusterCIDR}

			state := &clientCore.StateHandlers{
				Cluster: &mockClusterState{exists: tc.withComponentExists},

				Vpc: &mockVPCState{exists: tc.withCreationDependenciesMet},

				ExternalDNS:               &mockExternalDNSState{exists: !tc.withDeletionDependenciesMet},
				Monitoring:                &mockMonitoringState{exists: !tc.withDeletionDependenciesMet},
				Autoscaler:                &mockAutoscalerState{exists: !tc.withDeletionDependenciesMet},
				AWSLoadBalancerController: &mockAWSLoadBalancerState{exists: !tc.withDeletionDependenciesMet},
				Blockstorage:              &mockBlockstorageState{exists: !tc.withDeletionDependenciesMet},
				ExternalSecrets:           &mockExternalSecretsState{exists: !tc.withDeletionDependenciesMet},
				ArgoCD:                    &mockArgoCDState{exists: !tc.withDeletionDependenciesMet},
				Loki:                      &mockLokiState{exists: !tc.withDeletionDependenciesMet},
				Promtail:                  &mockPromtailState{exists: !tc.withDeletionDependenciesMet},
				Tempo:                     &mockTempoState{exists: !tc.withDeletionDependenciesMet},
				Application:               &mockApplicationState{existingApplications: tc.withApplications},
			}

			reconciler := NewClusterReconciler(
				&mockClusterService{
					creationBump: func() { creations++ },
					deletionBump: func() { deletions++ },
				},
				createCloudProvider(true),
			)

			_, err := reconciler.Reconcile(context.Background(), meta, state)
			assert.NoError(t, err)

			assert.Equal(t, tc.expectCreations, creations, "number of creations")
			assert.Equal(t, tc.expectDeletions, deletions, "number of deletions")
		})
	}
}

type mockClusterService struct {
	creationBump func()
	deletionBump func()
}

func (m mockClusterService) CreateCluster(_ context.Context, _ client.ClusterCreateOpts) (*client.Cluster, error) {
	m.creationBump()

	return nil, nil
}

func (m mockClusterService) DeleteCluster(_ context.Context, _ client.ClusterDeleteOpts) error {
	m.deletionBump()

	return nil
}

func (m mockClusterService) GetClusterSecurityGroupID(_ context.Context, _ client.GetClusterSecurityGroupIDOpts) (*api.ClusterSecurityGroupID, error) {
	panic("implement me")
}

func createCloudProvider(serviceQuotaOK bool) v1alpha1.CloudProvider {
	cloudProvider := mock.NewGoodCloudProvider()

	var available *float64
	if serviceQuotaOK {
		available = aws.Float64(100)
	} else {
		available = aws.Float64(0)
	}

	cloudProvider.SQAPI = &mock.SQAPI{
		GetServiceQuotaFn: func(*servicequotas.GetServiceQuotaInput) (*servicequotas.GetServiceQuotaOutput, error) {
			return &servicequotas.GetServiceQuotaOutput{
				Quota: &servicequotas.ServiceQuota{
					Value: available,
				},
			}, nil
		},
	}

	return cloudProvider
}
