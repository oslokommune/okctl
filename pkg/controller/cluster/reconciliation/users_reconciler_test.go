package reconciliation

import (
	"context"
	"testing"

	clientCore "github.com/oslokommune/okctl/pkg/client/core"
	"github.com/oslokommune/okctl/pkg/controller/common/reconciliation"

	"github.com/stretchr/testify/assert"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/client"
)

//nolint:funlen
func TestUsersReconciler(t *testing.T) {
	testCases := []struct {
		name string

		withPurge                bool
		withDisabledIdentityPool bool
		withDeclaredUsers        []v1alpha1.ClusterUser
		withExistingUsers        []client.IdentityPoolUser

		expectCreations int
		expectDeletions int
	}{
		{
			name: "Should do nothing when nothing is defined and nothing is existing",

			withDeclaredUsers: []v1alpha1.ClusterUser{},
			withExistingUsers: []client.IdentityPoolUser{},
			expectCreations:   0,
			expectDeletions:   0,
		},
		{
			name: "Should create when one user is defined and nothing is existing",

			withDeclaredUsers: []v1alpha1.ClusterUser{
				{
					Email: "dummy@email.com",
				},
			},
			withExistingUsers: []client.IdentityPoolUser{},
			expectCreations:   1,
			expectDeletions:   0,
		},
		{
			name: "Should delete when one user is existing but nothing is defined",

			withDeclaredUsers: []v1alpha1.ClusterUser{},
			withExistingUsers: []client.IdentityPoolUser{
				{
					Email: "dummy@email.com",
				},
			},
			expectCreations: 0,
			expectDeletions: 1,
		},
		{
			name: "Should create some, delete some, and ignore some when the situation arise",

			withDeclaredUsers: []v1alpha1.ClusterUser{
				{Email: "dummy@email.com"},
				{Email: "create@me.com"},
			},
			withExistingUsers: []client.IdentityPoolUser{
				{Email: "dummy@email.com"},
				{Email: "delete@me.com"},
				{Email: "delete@me.too"},
			},
			expectCreations: 1,
			expectDeletions: 2,
		},
		{
			name: "Should delete users even if Identity Pool doesn't exist (cleans up cfn stacks)",

			withPurge:                true,
			withDisabledIdentityPool: true,
			withDeclaredUsers: []v1alpha1.ClusterUser{
				{Email: "dummy@email.com"},
				{Email: "another@dummy.com"},
			},
			withExistingUsers: []client.IdentityPoolUser{
				{Email: "dummy@email.com"},
				{Email: "another@dummy.com"},
			},
			expectCreations: 0,
			expectDeletions: 2,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			creations := 0
			deletions := 0

			reconciler := NewUsersReconciler(&mockIdentityManagerService{
				createUserBump: func() { creations++ },
				deleteUserBump: func() { deletions++ },
			})

			meta := reconciliation.Metadata{
				Purge: tc.withPurge,
				ClusterDeclaration: &v1alpha1.Cluster{
					Users: tc.withDeclaredUsers,
					Integrations: &v1alpha1.ClusterIntegrations{
						Cognito: !tc.withDisabledIdentityPool,
					},
				},
			}

			state := &clientCore.StateHandlers{
				IdentityManager: mockIdentityManagerState{
					existingIdentityPool: !tc.withDisabledIdentityPool,
					existingUsers:        tc.withExistingUsers,
				},
			}

			_, err := reconciler.Reconcile(context.Background(), meta, state)
			assert.NoError(t, err)

			assert.Equal(t, tc.expectCreations, creations)
			assert.Equal(t, tc.expectDeletions, deletions)
		})
	}
}
