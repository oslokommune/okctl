package core

import (
	"testing"

	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/api/core"
	"github.com/stretchr/testify/assert"
)

func TestAcquireStateLock(t *testing.T) {
	testCases := []struct {
		name             string
		withExistingLock bool
		expectError      string
	}{
		{
			name:             "Should acquire lock without problems when there is no existing lock",
			withExistingLock: false,
			expectError:      "",
		},
		{
			name:             "Should get an error when trying to acquire a lock when a lock already exist",
			withExistingLock: true,
			expectError:      "state is locked",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			kvService := core.NewKeyValueStoreService(&mockDynamoProvider{})

			remoteStateServiceInstance := NewRemoteStateService(kvService, nil)

			if tc.withExistingLock {
				err := remoteStateServiceInstance.AcquireStateLock(mockClusterID())
				assert.NoError(t, err)
			}

			err := remoteStateServiceInstance.AcquireStateLock(mockClusterID())

			if tc.expectError != "" {
				assert.Equal(t, tc.expectError, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

type mockDynamoProvider struct {
	hasLock bool
}

func (m *mockDynamoProvider) ListStores() ([]string, error)           { panic("implement me") }
func (m *mockDynamoProvider) CreateStore(_ api.CreateStoreOpts) error { return nil }
func (m *mockDynamoProvider) DeleteStore(_ api.DeleteStoreOpts) error { return nil }
func (m *mockDynamoProvider) GetString(_ api.GetStringOpts) (string, error) {
	if m.hasLock {
		return "lock", nil
	}

	return "", errors.E(errors.New("not found"), errors.NotExist)
}

func (m *mockDynamoProvider) InsertItem(_ api.InsertItemOpts) error {
	m.hasLock = true

	return nil
}

func (m *mockDynamoProvider) RemoveItem(_ api.DeleteItemOpts) error {
	m.hasLock = false

	return nil
}

func mockClusterID() api.ID {
	return api.ID{
		ClusterName:  "test-cluster",
		Region:       "eu-test-1",
		AWSAccountID: "012345678912",
	}
}
