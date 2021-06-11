package reconciliation

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/servicequotas"
	"github.com/oslokommune/okctl/pkg/config/constant"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/client"
	clientCore "github.com/oslokommune/okctl/pkg/client/core"
	"github.com/oslokommune/okctl/pkg/mock"
	"github.com/stretchr/testify/assert"
)

//nolint:funlen
func TestVPCReconciler(t *testing.T) {
	// Componentflag is ignored and always true for the VPC
	testCases := []struct {
		name                      string
		withPurge                 bool
		withComponentExists       bool
		expectCreations           int
		expectDeletions           int
		withDeleteDependenciesMet bool
	}{
		{
			name:                "Should noop when existing",
			withComponentExists: true,
			expectCreations:     0,
			expectDeletions:     0,
		},
		{
			name:                "Should create when not purge and not existing",
			withPurge:           false,
			withComponentExists: false,
			expectCreations:     1,
			expectDeletions:     0,
		},
		{
			name:                      "Should delete when purge and existing",
			withPurge:                 true,
			withComponentExists:       true,
			withDeleteDependenciesMet: true,
			expectCreations:           0,
			expectDeletions:           1,
		},
		{
			name:                      "Should noop when purge, existing and cluster exists",
			withPurge:                 true,
			withComponentExists:       true,
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

			meta := generateTestMeta(tc.withPurge, v1alpha1.ClusterIntegrations{})
			meta.ClusterDeclaration.VPC = &v1alpha1.ClusterVPC{
				CIDR: constant.DefaultClusterCIDR,
			}

			state := &clientCore.StateHandlers{
				Cluster: &mockClusterState{exists: !tc.withDeleteDependenciesMet},
				Vpc:     &mockVPCState{exists: tc.withComponentExists},
			}

			reconciler := NewVPCReconciler(
				&mockVPCService{
					creationBump: func() { creations++ },
					deletionBump: func() { deletions++ },
				},
				createHappyServiceQuotaCloudProvider(),
			)

			_, err := reconciler.Reconcile(context.Background(), meta, state)
			assert.NoError(t, err)

			assert.Equal(t, tc.expectCreations, creations, "creations")
			assert.Equal(t, tc.expectDeletions, deletions, "deletions")
		})
	}
}

type mockVPCService struct {
	creationBump func()
	deletionBump func()
}

func (m mockVPCService) CreateVpc(_ context.Context, _ client.CreateVpcOpts) (*client.Vpc, error) {
	m.creationBump()

	return nil, nil
}

func (m mockVPCService) DeleteVpc(_ context.Context, _ client.DeleteVpcOpts) error {
	m.deletionBump()

	return nil
}

func (m mockVPCService) GetVPC(_ context.Context, _ api.ID) (*client.Vpc, error) {
	panic("implement me")
}

func createHappyServiceQuotaCloudProvider() v1alpha1.CloudProvider {
	cloudProvider := mock.NewGoodCloudProvider()
	cloudProvider.SQAPI = &mock.SQAPI{
		GetServiceQuotaFn: func(*servicequotas.GetServiceQuotaInput) (*servicequotas.GetServiceQuotaOutput, error) {
			return &servicequotas.GetServiceQuotaOutput{
				Quota: &servicequotas.ServiceQuota{
					Value: aws.Float64(4),
				},
			}, nil
		},
	}

	return cloudProvider
}
