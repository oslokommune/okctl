package reconciliation

import (
	"context"
	"testing"

	"github.com/oslokommune/okctl/pkg/config/constant"

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
				deletionBump: func(string) { deletions++ },
			})

			meta := reconciliation.Metadata{
				ClusterDeclaration: &v1alpha1.Cluster{
					Databases: &v1alpha1.ClusterDatabases{Postgres: tc.withDeclaredDBs},
				},
			}

			state := &clientCore.StateHandlers{
				Vpc:         &mockVPCState{exists: tc.withVPCExists},
				Component:   &mockComponentState{databases: tc.withExistingDBs},
				Application: &mockApplicationState{existingApplications: 0},
			}

			_, err := reconciler.Reconcile(context.Background(), meta, state)
			assert.NoError(t, err)

			assert.Equal(t, tc.expectCreations, creations)
			assert.Equal(t, tc.expectDeletions, deletions)
		})
	}
}

//nolint:funlen
func TestDeletePostgresWithExistingApp(t *testing.T) {
	testCases := []struct {
		name string

		withApplications []v1alpha1.Application
		withExistingDBs  []v1alpha1.ClusterDatabasesPostgres

		expectDeletions int
	}{
		{
			name: "Should do nothing with one application depending on one database",

			withApplications: []v1alpha1.Application{
				{
					Postgres: "mockdb",
				},
			},
			withExistingDBs: []v1alpha1.ClusterDatabasesPostgres{
				{
					Name:      "mockdb",
					User:      "postgres",
					Namespace: "dbdetails",
				},
			},
			expectDeletions: 0,
		},
		{
			name: "Should delete one database when two are configured and only one dependecy",

			withApplications: []v1alpha1.Application{
				{
					Postgres: "mockdb",
				},
				{
					Postgres: "",
				},
			},
			withExistingDBs: []v1alpha1.ClusterDatabasesPostgres{
				{
					Name:      "mockdb",
					User:      "postgres",
					Namespace: "dbdetails",
				},
				{
					Name:      "secondmockdb",
					User:      "postgres",
					Namespace: "dbdetails",
				},
			},
			expectDeletions: 1,
		},
		{
			name: "Should delete everything if no apps are dependent",

			withApplications: []v1alpha1.Application{
				{
					Postgres: "",
				},
				{
					Postgres: "",
				},
			},
			withExistingDBs: []v1alpha1.ClusterDatabasesPostgres{
				{
					Name:      "mockdb",
					User:      "postgres",
					Namespace: "dbdetails",
				},
				{
					Name:      "mockdb2",
					User:      "postgres",
					Namespace: "dbdetails",
				},
				{
					Name:      "mockdb3",
					User:      "postgres",
					Namespace: "dbdetails",
				},
			},
			expectDeletions: 3,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			deletions := 0

			meta := reconciliation.Metadata{
				Purge: true,
				ClusterDeclaration: &v1alpha1.Cluster{
					Databases: &v1alpha1.ClusterDatabases{Postgres: tc.withExistingDBs},
				},
			}

			existingDbs := make([]*client.PostgresDatabase, len(tc.withExistingDBs))

			for index, db := range tc.withExistingDBs {
				existingDbs[index] = &client.PostgresDatabase{
					ApplicationName: db.Name,
					UserName:        db.User,
					Namespace:       db.Namespace,
				}
			}

			state := &clientCore.StateHandlers{
				Vpc:         &mockVPCState{exists: true},
				Component:   &mockComponentState{databases: existingDbs},
				Application: &mockPgAppState{applications: tc.withApplications},
			}

			reconciler := NewPostgresReconciler(&mockPostgresService{
				deletionBump: func(name string) {
					newExistingDbs := make([]*client.PostgresDatabase, 0)

					for _, db := range existingDbs {
						if db.ApplicationName != name {
							newExistingDbs = append(newExistingDbs, db)
						} else {
							deletions++
						}
					}

					existingDbs = newExistingDbs
				},
			})

			for i := 0; i < constant.DefaultMaxReconciliationRequeues; i++ {
				result, err := reconciler.Reconcile(context.Background(), meta, state)
				assert.NoError(t, err)

				if result.Requeue == false {
					break
				}
			}

			assert.Equal(t, tc.expectDeletions, deletions)
		})
	}
}

type mockPgAppState struct {
	applications []v1alpha1.Application
}

func (m mockPgAppState) Initialize(v1alpha1.Cluster, string) error { panic("implement me") }
func (m mockPgAppState) Put(v1alpha1.Application) error            { panic("implement me") }
func (m mockPgAppState) Get(string) (v1alpha1.Application, error)  { panic("implement me") }
func (m mockPgAppState) Delete(string) error                       { panic("implement me") }
func (m mockPgAppState) List() ([]v1alpha1.Application, error) {
	return m.applications, nil
}

type mockPostgresService struct {
	creationBump func()
	deletionBump func(string)
}

func (m mockPostgresService) CreatePostgresDatabase(_ context.Context, _ client.CreatePostgresDatabaseOpts) (*client.PostgresDatabase, error) {
	m.creationBump()

	return nil, nil
}

func (m mockPostgresService) DeletePostgresDatabase(_ context.Context, opts client.DeletePostgresDatabaseOpts) error {
	m.deletionBump(opts.ApplicationName)

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
