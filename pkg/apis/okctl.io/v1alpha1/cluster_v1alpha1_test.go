package v1alpha1_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"sigs.k8s.io/yaml"

	"github.com/sebdah/goldie/v2"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
)

func TestMarshalCluster(t *testing.T) {
	testCases := []struct {
		name    string
		cluster v1alpha1.Cluster
		golden  string
	}{
		{
			name:    "Empty cluster",
			cluster: v1alpha1.Cluster{},
			golden:  "empty-cluster.yml",
		},
		{
			name: "Default cluster",
			cluster: newCluster(
				"okctl-stage",
				"okctl-stage.oslo.systems",
				"oslokommune",
				"okctl-iac",
				"123456789012",
			),
			golden: "default-cluster.yml",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			got, err := yaml.Marshal(tc.cluster)
			assert.NoError(t, err)

			g := goldie.New(t)
			g.Assert(t, tc.golden, got)
		})
	}
}

func TestValidations(t *testing.T) {
	testCases := []struct {
		name      string
		with      v1alpha1.Cluster
		expectErr string
	}{
		{
			name: "Should not trigger KM302",
			with: func() v1alpha1.Cluster {
				c := newCluster(
					"test",
					"okctl.io",
					"oslokommune",
					"okctl",
					"012345678912",
				)

				c.VPC = &v1alpha1.ClusterVPC{
					CIDR:             "192.168.0.0/20",
					HighAvailability: false,
				}

				return c
			}(),
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			err := tc.with.Validate()

			if tc.expectErr == "" {
				assert.NoError(t, err)
			} else {
				assert.Equal(t, tc.expectErr, err)
			}
		})
	}
}

//nolint:funlen
func TestPostgresValidateDatabaseName(t *testing.T) {
	testCases := []struct {
		name         string
		databaseName string
		expectErr    string
	}{
		{
			name:         "Valid databaes name",
			databaseName: "my-user-database",
		},
		{
			name:         "'db' is a reserved database name",
			databaseName: "db",
			expectErr:    "name: 'db' and 'database' are reserved postgres database names.",
		},
		{
			name:         "'database' is a reserved database name",
			databaseName: "database",
			expectErr:    "name: 'db' and 'database' are reserved postgres database names.",
		},
		{
			name:         "Should fail because database name is longer than 60 characters",
			databaseName: "usersusersusersusersusersusersusersusersusersusersusersusers1",
			expectErr:    "name: database name cannot be longer than 60 characters.",
		},
		{
			name:         "Should fail because database name starts with a number",
			databaseName: "9-nine",
			expectErr:    "name: database name must start with a letter.",
		},

		{
			name:         "Should fail because database ends with a hyphen",
			databaseName: "nine-",
			expectErr:    "name: database name must not end with a hyphen.",
		},

		{
			name:         "Should fail because database name have two consecutive hyphens",
			databaseName: "nine--nine",
			expectErr:    "name: database name can not have two consecutive hyphens.",
		},
		{
			name:         "Should fail because database name have uppercase letters",
			databaseName: "aInvalidDatabase",
			expectErr:    "name: database name cannot have capital letter.",
		},
		{
			name:         "Should fail because database name starts with uppercase letters",
			databaseName: "Invaliddatabase",
			expectErr:    "name: database name cannot have capital letter.",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			database := v1alpha1.ClusterDatabasesPostgres{
				Name:      tc.databaseName,
				User:      "administrator",
				Namespace: "valid-namespace",
			}
			errs := database.Validate()
			if tc.expectErr == "" {
				assert.NoError(t, errs)
			} else {
				assert.Error(t, errs)
				assert.Equal(t, tc.expectErr, errs.Error())
			}
		})
	}
}

//nolint:funlen
func TestPostgresValidateUserName(t *testing.T) {
	testCases := []struct {
		name      string
		user      string
		expectErr string
	}{
		{
			name: "'administrator' is a valid username'",
			user: "administrator",
		},
		{
			name: "'admin2' is a valid username'",
			user: "admin2",
		},
		{
			name: "'my-admin' is a valid username'",
			user: "my-admin",
		},
		{
			name:      "'admin' is a reserved postgres username",
			user:      "admin",
			expectErr: "user: 'admin' is a reserved postgres username.",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			database := v1alpha1.ClusterDatabasesPostgres{
				Name:      "valid-database-name",
				User:      tc.user,
				Namespace: "valid-namespace",
			}
			errs := database.Validate()
			if tc.expectErr == "" {
				assert.NoError(t, errs)
			} else {
				assert.Error(t, errs)
				assert.Equal(t, tc.expectErr, errs.Error())
			}
		})
	}
}

func newCluster(name, rootDomain, githubOrganization, githubRepository, awsAccountID string) v1alpha1.Cluster {
	c := v1alpha1.NewCluster()
	c.Metadata.Name = name
	c.Metadata.AccountID = awsAccountID
	c.Github.Organisation = githubOrganization
	c.Github.Repository = githubRepository
	c.ClusterRootDomain = rootDomain

	return c
}
