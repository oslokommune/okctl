package reconciliation

import (
	"testing"

	"github.com/oslokommune/okctl/pkg/client"
	clientCore "github.com/oslokommune/okctl/pkg/client/core"

	"github.com/stretchr/testify/assert"
)

//nolint:funlen
func TestAssertDependencyExistence(t *testing.T) {
	testCases := []struct {
		name string

		withTests  []DependencyTestFn
		withExpect bool

		expectResult bool
		expectErr    string
	}{
		{
			name: "Should return true if expectence is true all tests are true",
			withTests: []DependencyTestFn{
				func() (bool, error) { return true, nil },
				func() (bool, error) { return true, nil },
				func() (bool, error) { return true, nil },
			},
			withExpect:   true,
			expectResult: true,
		},
		{
			name: "Should return false if expectence is true and one of the tests are false",
			withTests: []DependencyTestFn{
				func() (bool, error) { return true, nil },
				func() (bool, error) { return false, nil },
				func() (bool, error) { return true, nil },
			},
			withExpect:   true,
			expectResult: false,
		},
		{
			name: "Should return true if expectence is false all tests are false",
			withTests: []DependencyTestFn{
				func() (bool, error) { return false, nil },
				func() (bool, error) { return false, nil },
				func() (bool, error) { return false, nil },
			},
			withExpect:   false,
			expectResult: true,
		},
		{
			name: "Should return false if expectence is false and one of the tests are true",
			withTests: []DependencyTestFn{
				func() (bool, error) { return false, nil },
				func() (bool, error) { return true, nil },
				func() (bool, error) { return false, nil },
			},
			withExpect:   false,
			expectResult: false,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			result, err := AssertDependencyExistence(tc.withExpect, tc.withTests...)

			if tc.expectErr != "" {
				assert.NotNil(t, err)
				assert.Equal(t, tc.expectErr, err.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.expectResult, result)
		})
	}
}

func TestGeneratePrimaryDomainDelegationTest(t *testing.T) {
	testCases := []struct {
		name string

		withPrimaryHostedZoneExists bool
		withDelegated               bool

		expect bool
	}{
		{
			name: "Should create a dependency test which returns true if primary hosted zone exists and is delegated",

			withPrimaryHostedZoneExists: true,
			withDelegated:               true,

			expect: true,
		},
		{
			name: "Should create a dependency test which returns false if primary hosted zone exists and is not delegated",

			withPrimaryHostedZoneExists: true,
			withDelegated:               false,

			expect: false,
		},
		{
			name: "Should create a dependency test which returns false if primary hosted zone does not exists",

			withPrimaryHostedZoneExists: false,
			withDelegated:               true,

			expect: false,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			state := &clientCore.StateHandlers{Domain: mockDomainState{
				primaryHostedZoneExists:    tc.withPrimaryHostedZoneExists,
				primaryHostedZoneDelegated: tc.withDelegated,
			}}

			test := GeneratePrimaryDomainDelegationTest(state)

			result, err := test()
			assert.NoError(t, err)

			assert.Equal(t, tc.expect, result)
		})
	}
}

type mockDomainState struct {
	primaryHostedZoneExists    bool
	primaryHostedZoneDelegated bool
}

func (m mockDomainState) HasPrimaryHostedZone() (bool, error) {
	return m.primaryHostedZoneExists, nil
}

func (m mockDomainState) GetPrimaryHostedZone() (*client.HostedZone, error) {
	return &client.HostedZone{
		IsDelegated: m.primaryHostedZoneDelegated,
	}, nil
}

func (m mockDomainState) SaveHostedZone(_ *client.HostedZone) error          { panic("implement me") }
func (m mockDomainState) UpdateHostedZone(_ *client.HostedZone) error        { panic("implement me") }
func (m mockDomainState) RemoveHostedZone(_ string) error                    { panic("implement me") }
func (m mockDomainState) GetHostedZone(_ string) (*client.HostedZone, error) { panic("implement me") }
func (m mockDomainState) GetHostedZones() ([]*client.HostedZone, error)      { panic("implement me") }

func TestGenerateClusterExistenceTest(t *testing.T) {
	testCases := []struct {
		name string

		withClusterExistence bool
		expectResult         bool
	}{
		{
			name:                 "Should create a test which returns true if the cluster exists",
			withClusterExistence: true,
			expectResult:         true,
		},
		{
			name:                 "Should create a test which returns false if the cluster does not exists",
			withClusterExistence: false,
			expectResult:         false,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			state := &clientCore.StateHandlers{Cluster: mockClusterState{
				clusterExists: tc.withClusterExistence,
			}}

			test := GenerateClusterExistenceTest(state, "dummy")

			result, err := test()
			assert.NoError(t, err)

			assert.Equal(t, tc.expectResult, result)
		})
	}
}

type mockClusterState struct {
	clusterExists bool
}

func (m mockClusterState) HasCluster(_ string) (bool, error) {
	return m.clusterExists, nil
}

func (m mockClusterState) SaveCluster(_ *client.Cluster) error          { panic("implement me") }
func (m mockClusterState) GetCluster(_ string) (*client.Cluster, error) { panic("implement me") }
func (m mockClusterState) RemoveCluster(_ string) error                 { panic("implement me") }
