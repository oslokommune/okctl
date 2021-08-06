package reconciliation

import (
	"context"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/controller/common/reconciliation"
)

func generateTestMeta(purge bool, integrations v1alpha1.ClusterIntegrations) reconciliation.Metadata {
	return reconciliation.Metadata{
		ClusterDeclaration: &v1alpha1.Cluster{
			Integrations: &integrations,
		},
		Purge: purge,
	}
}

type generalizedTestCase struct {
	name                string
	withPurge           bool
	withComponentFlag   bool
	withComponentExists bool
	withDependenciesMet bool
	expectCreations     int
	expectDeletions     int
}

// # Mocks

// ## State
type mockVPCState struct {
	exists bool
}

func (m mockVPCState) HasVPC(_ string) (bool, error) {
	return m.exists, nil
}

func (m mockVPCState) GetVpc(_ string) (*client.Vpc, error) {
	return &client.Vpc{
		VpcID: "dummy-id",
	}, nil
}

func (m mockVPCState) SaveVpc(_ *client.Vpc) error { panic("implement me") }
func (m mockVPCState) RemoveVpc(_ string) error    { panic("implement me") }

type mockClusterState struct {
	exists bool
}

func (m mockClusterState) HasCluster(_ string) (bool, error) {
	return m.exists, nil
}

func (m mockClusterState) SaveCluster(_ *client.Cluster) error          { panic("implement me") }
func (m mockClusterState) GetCluster(_ string) (*client.Cluster, error) { panic("implement me") }
func (m mockClusterState) RemoveCluster(_ string) error                 { panic("implement me") }

type mockDomainState struct {
	exists      bool
	isDelegated bool
}

func (m mockDomainState) GetPrimaryHostedZone() (*client.HostedZone, error) {
	if m.exists {
		return &client.HostedZone{
			IsDelegated: m.isDelegated,
		}, nil
	}

	return nil, nil
}

func (m mockDomainState) HasPrimaryHostedZone() (bool, error) {
	return m.exists, nil
}

func (m mockDomainState) SaveHostedZone(_ *client.HostedZone) error          { panic("implement me") }
func (m mockDomainState) UpdateHostedZone(_ *client.HostedZone) error        { panic("implement me") }
func (m mockDomainState) RemoveHostedZone(_ string) error                    { panic("implement me") }
func (m mockDomainState) GetHostedZone(_ string) (*client.HostedZone, error) { panic("implement me") }
func (m mockDomainState) GetHostedZones() ([]*client.HostedZone, error)      { panic("implement me") }

type mockIdentityManagerState struct {
	existingIdentityPool bool
	existingUsers        []client.IdentityPoolUser
}

func (m mockIdentityManagerState) GetIdentityPoolUsers() ([]client.IdentityPoolUser, error) {
	return m.existingUsers, nil
}

func (m mockIdentityManagerState) HasIdentityPool() (bool, error) {
	return m.existingIdentityPool, nil
}

func (m mockIdentityManagerState) GetIdentityPool(_ string) (*client.IdentityPool, error) {
	return &client.IdentityPool{
		UserPoolID: "dummyID",
	}, nil
}

func (m mockIdentityManagerState) SaveIdentityPool(_ *client.IdentityPool) error {
	panic("implement me")
}

func (m mockIdentityManagerState) RemoveIdentityPool(_ string) error { panic("implement me") }

func (m mockIdentityManagerState) SaveIdentityPoolClient(_ *client.IdentityPoolClient) error {
	panic("implement me")
}

func (m mockIdentityManagerState) GetIdentityPoolClient(_ string) (*client.IdentityPoolClient, error) {
	panic("implement me")
}

func (m mockIdentityManagerState) RemoveIdentityPoolClient(_ string) error { panic("implement me") }

func (m mockIdentityManagerState) SaveIdentityPoolUser(_ *client.IdentityPoolUser) error {
	panic("implement me")
}

func (m mockIdentityManagerState) GetIdentityPoolUser(_ string) (*client.IdentityPoolUser, error) {
	panic("implement me")
}

func (m mockIdentityManagerState) RemoveIdentityPoolUser(_ string) error { panic("implement me") }

// ## Services
type mockIdentityManagerService struct {
	createIdentityPoolBump func()
	deleteIdentityPoolBump func()
	createUserBump         func()
	deleteUserBump         func()
}

func (m *mockIdentityManagerService) CreateIdentityPoolUser(_ context.Context, _ client.CreateIdentityPoolUserOpts) (*client.IdentityPoolUser, error) {
	m.createUserBump()

	return nil, nil
}

func (m *mockIdentityManagerService) DeleteIdentityPoolUser(_ context.Context, _ client.DeleteIdentityPoolUserOpts) error {
	m.deleteUserBump()

	return nil
}

func (m *mockIdentityManagerService) CreateIdentityPool(_ context.Context, _ client.CreateIdentityPoolOpts) (*client.IdentityPool, error) {
	m.createIdentityPoolBump()

	return nil, nil
}

func (m *mockIdentityManagerService) DeleteIdentityPool(_ context.Context, _ api.ID) error {
	m.deleteIdentityPoolBump()

	return nil
}

func (m *mockIdentityManagerService) CreateIdentityPoolClient(_ context.Context, _ client.CreateIdentityPoolClientOpts) (*client.IdentityPoolClient, error) {
	panic("implement me")
}

func (m *mockIdentityManagerService) DeleteIdentityPoolClient(_ context.Context, _ client.DeleteIdentityPoolClientOpts) error {
	panic("implement me")
}

type mockMonitoringService struct {
	createKubePromStackBump func()
	deleteKubePromStackBump func()
	createLokiBump          func()
	deleteLokiBump          func()
	createPromtailBump      func()
	deletePromtailBump      func()
	createTempoBump         func()
	deleteTempoBump         func()
}

func (m mockMonitoringService) CreateKubePromStack(_ context.Context, _ client.CreateKubePromStackOpts) (*client.KubePromStack, error) {
	m.createKubePromStackBump()

	return nil, nil
}

func (m mockMonitoringService) DeleteKubePromStack(_ context.Context, _ client.DeleteKubePromStackOpts) error {
	m.deleteKubePromStackBump()

	return nil
}

func (m mockMonitoringService) CreateLoki(_ context.Context, _ api.ID) (*client.Helm, error) {
	m.createLokiBump()

	return nil, nil
}

func (m mockMonitoringService) DeleteLoki(_ context.Context, _ api.ID) error {
	m.deleteLokiBump()

	return nil
}

func (m mockMonitoringService) CreatePromtail(_ context.Context, _ api.ID) (*client.Helm, error) {
	m.createPromtailBump()

	return nil, nil
}

func (m mockMonitoringService) DeletePromtail(_ context.Context, _ api.ID) error {
	m.deletePromtailBump()

	return nil
}

func (m mockMonitoringService) CreateTempo(_ context.Context, _ api.ID) (*client.Helm, error) {
	m.createTempoBump()

	return nil, nil
}

func (m mockMonitoringService) DeleteTempo(_ context.Context, _ api.ID) error {
	m.deleteTempoBump()

	return nil
}
