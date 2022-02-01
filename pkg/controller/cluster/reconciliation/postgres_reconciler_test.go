package reconciliation

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/client"
	clientCore "github.com/oslokommune/okctl/pkg/client/core"
	"github.com/oslokommune/okctl/pkg/controller/common/reconciliation"
)

//nolint:funlen
func TestPostgresReconciler(t *testing.T) {
	testCases := []struct {
		name string

		withVPCExists   bool
		withDeclaredDBs []v1alpha1.ClusterDatabasesPostgres
		withExistingDBs []*client.PostgresDatabase

		expectCreations int
		expectDeletions int
	}{
		{
			name: "Should do nothing when nothing is defined and nothing is existing",

			withVPCExists:   true,
			withDeclaredDBs: []v1alpha1.ClusterDatabasesPostgres{},
			withExistingDBs: []*client.PostgresDatabase{},
			expectCreations: 0,
			expectDeletions: 0,
		},
		{
			name: "Should create when one db is defined and nothing is existing",

			withVPCExists: true,
			withDeclaredDBs: []v1alpha1.ClusterDatabasesPostgres{
				{
					Name: "dummydb",
				},
			},
			withExistingDBs: []*client.PostgresDatabase{},
			expectCreations: 1,
			expectDeletions: 0,
		},
		{
			name: "Should delete when one db is existing but nothing is defined",

			withVPCExists:   true,
			withDeclaredDBs: []v1alpha1.ClusterDatabasesPostgres{},
			withExistingDBs: []*client.PostgresDatabase{
				{
					ApplicationName: "dummydb",
				},
			},
			expectCreations: 0,
			expectDeletions: 1,
		},
		{
			name: "Should create some, delete some, and ignore some when the situation arise",

			withVPCExists: true,
			withDeclaredDBs: []v1alpha1.ClusterDatabasesPostgres{
				{
					Name: "noop",
				},
				{
					Name: "create-db",
				},
			},
			withExistingDBs: []*client.PostgresDatabase{
				{
					ApplicationName: "noop",
				},
				{
					ApplicationName: "deleteone",
				},
				{
					ApplicationName: "deletetwo",
				},
			},
			expectCreations: 1,
			expectDeletions: 2,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			creations := 0
			deletions := 0

			reconciler := NewPostgresReconciler(&mockPostgresService{
				creationBump: func() { creations++ },
				deletionBump: func() { deletions++ },
			})

			meta := reconciliation.Metadata{
				ClusterDeclaration: &v1alpha1.Cluster{
					Databases: &v1alpha1.ClusterDatabases{Postgres: tc.withDeclaredDBs},
				},
			}

			state := &clientCore.StateHandlers{
				Vpc:       &mockVPCState{exists: tc.withVPCExists},
				Component: &mockComponentState{databases: tc.withExistingDBs},
			}

			_, err := reconciler.Reconcile(context.Background(), meta, state)
			assert.NoError(t, err)

			assert.Equal(t, tc.expectCreations, creations)
			assert.Equal(t, tc.expectDeletions, deletions)
		})
	}
}

//nolint:funlen
func TestPostgresReconcilerInvalidDatabaseName(t *testing.T) {
	const withVPCExists bool = true

	testCases := []struct {
		name            string
		withDeclaredDBs []v1alpha1.ClusterDatabasesPostgres
		err             string
	}{
		{
			name: "Should create no databases because the declared databases names are reserved",
			withDeclaredDBs: []v1alpha1.ClusterDatabasesPostgres{
				{
					Name: "db",
				},
				{
					Name: "database",
				},
			},
			err: "determining course of action: invalid database name: 'db' and 'database' are reserved",
		},
		{
			name: "Should create no databases because one of the declared database names are reserved",
			withDeclaredDBs: []v1alpha1.ClusterDatabasesPostgres{
				{
					Name: "db",
				},
				{
					Name: "new-user-database",
				},
			},
			err: "determining course of action: invalid database name: 'db' and 'database' are reserved",
		},
		{
			name: "Should fail because database name is longer than 60 characters",
			withDeclaredDBs: []v1alpha1.ClusterDatabasesPostgres{
				{
					Name: "usersusersusersusersusersusersusersusersusersusersusersusers1",
				},
			},
			err: "determining course of action: invalid database name: cannot be longer than 60 characters",
		},
		{
			name: "Should fail because database name starts with a number",
			withDeclaredDBs: []v1alpha1.ClusterDatabasesPostgres{
				{
					Name: "9-nine",
				},
			},
			err: "determining course of action: invalid database name: cannot start with a number",
		},
		{
			name: "Should fail because database ends with a hyphen",
			withDeclaredDBs: []v1alpha1.ClusterDatabasesPostgres{
				{
					Name: "nine-",
				},
			},
			err: "determining course of action: invalid database name: cannot end with a hyphen",
		},
		{
			name: "Should fail because database name have two consecutive hyphens",
			withDeclaredDBs: []v1alpha1.ClusterDatabasesPostgres{
				{
					Name: "nine--nine",
				},
			},
			err: "determining course of action: invalid database name: cannot have two consecutive hyphens",
		},
		{
			name: "Should fail because database name have uppercase letters",
			withDeclaredDBs: []v1alpha1.ClusterDatabasesPostgres{
				{
					Name: "aInvalidDatabase",
				},
			},
			err: "determining course of action: invalid database name: cannot contain uppercase letters",
		},
		{
			name: "Should fail because database name starts with uppercase letters",
			withDeclaredDBs: []v1alpha1.ClusterDatabasesPostgres{
				{
					Name: "Invaliddatabase",
				},
			},
			err: "determining course of action: invalid database name: cannot contain uppercase letters",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			reconciler := NewPostgresReconciler(&mockPostgresService{})

			meta := reconciliation.Metadata{
				ClusterDeclaration: &v1alpha1.Cluster{
					Databases: &v1alpha1.ClusterDatabases{Postgres: tc.withDeclaredDBs},
				},
			}

			state := &clientCore.StateHandlers{
				Vpc:       &mockVPCState{exists: withVPCExists},
				Component: &mockComponentState{},
			}

			_, err := reconciler.Reconcile(context.Background(), meta, state)
			assert.Error(t, err)
			assert.Equal(t, tc.err, err.Error())
		})
	}
}

type mockPostgresService struct {
	creationBump func()
	deletionBump func()
}

func (m mockPostgresService) CreatePostgresDatabase(_ context.Context, _ client.CreatePostgresDatabaseOpts) (*client.PostgresDatabase, error) {
	m.creationBump()

	return nil, nil
}

func (m mockPostgresService) DeletePostgresDatabase(_ context.Context, _ client.DeletePostgresDatabaseOpts) error {
	m.deletionBump()

	return nil
}

func (m mockPostgresService) GetPostgresDatabase(_ context.Context, _ client.GetPostgresDatabaseOpts) (*client.PostgresDatabase, error) {
	panic("implement me")
}

type mockComponentState struct {
	databases []*client.PostgresDatabase
}

func (m mockComponentState) GetPostgresDatabases() ([]*client.PostgresDatabase, error) {
	return m.databases, nil
}

func (m mockComponentState) SavePostgresDatabase(_ *client.PostgresDatabase) error {
	panic("implement me")
}
func (m mockComponentState) RemovePostgresDatabase(_ string) error { panic("implement me") }
func (m mockComponentState) GetPostgresDatabase(_ string) (*client.PostgresDatabase, error) {
	panic("implement me")
}
