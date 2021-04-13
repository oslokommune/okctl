package commands

import (
	"bytes"
	"testing"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"

	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
)

func TestClusterDeclarationScaffold(t *testing.T) {
	testCases := []struct {
		name string

		withOpts     ScaffoldClusterOpts
		expectGolden string
	}{
		{
			name: "Should scaffold cluster declaration based on opts",
			withOpts: ScaffoldClusterOpts{
				Name:            "SomeClusterName",
				AWSAccountID:    "123456789012",
				Environment:     "production",
				Organization:    "oslokommune",
				RepositoryName:  "my_repo",
				OutputDirectory: "my_infrastructure",
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer

			err := ScaffoldClusterDeclaration(&buf, tc.withOpts)
			assert.NoError(t, err)

			g := goldie.New(t)
			g.Assert(t, tc.name, buf.Bytes())
		})
	}
}

func TestEnsureValidDefaultTemplate(t *testing.T) {
	var buf bytes.Buffer

	err := ScaffoldClusterDeclaration(&buf, ScaffoldClusterOpts{
		Name:            "name",
		AWSAccountID:    "123456789012",
		Environment:     "test",
		Organization:    "oslokommune",
		RepositoryName:  "my_iac_repo",
		OutputDirectory: "infrastructure",
	})
	assert.NoError(t, err)

	cluster, err := InferClusterFromStdinOrFile(&buf, "-")
	assert.NoError(t, err)

	err = cluster.Validate()
	assert.NoError(t, err)
}

// nolint: funlen
func TestRetrieveValuesFromTemplate(t *testing.T) {
	testCases := []struct {
		name string

		withTemplate string

		expectTrueTest func(v1alpha1.Cluster) bool
	}{
		{
			name: "Should correctly scaffold and parse users",

			withTemplate: `
              users:
              - email: test.user@origo.oslo.kommune.no`,

			expectTrueTest: func(c v1alpha1.Cluster) bool {
				if len(c.Users) != 1 {
					return false
				}

				if c.Users[0].Email != "test.user@origo.oslo.kommune.no" {
					return false
				}

				return true
			},
		},
		{
			name: "Should correctly scaffold and parse databases",

			withTemplate: `
              databases:
                postgres:
                - name: test
                  namespace: testns
                  user: postgres`,

			expectTrueTest: func(c v1alpha1.Cluster) bool {
				if len(c.Databases.Postgres) != 1 {
					return false
				}

				db := c.Databases.Postgres[0]

				if db.Name != "test" {
					return false
				}

				if db.Namespace != "testns" {
					return false
				}

				if db.User != "postgres" {
					return false
				}

				return true
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer

			_, _ = buf.WriteString(tc.withTemplate)

			cluster, err := InferClusterFromStdinOrFile(&buf, "-")
			assert.NoError(t, err)

			assert.True(t, tc.expectTrueTest(*cluster))
		})
	}
}
