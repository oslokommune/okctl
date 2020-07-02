package core_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/api/core"
	"github.com/oslokommune/okctl/pkg/api/mock"
	"github.com/stretchr/testify/assert"
)

func TestClusterCreateCluster(t *testing.T) {
	testCases := []struct {
		name        string
		service     api.ClusterService
		opts        api.ClusterCreateOpts
		expect      interface{}
		expectError bool
	}{
		{
			name: "Valid request",
			service: core.NewClusterService(
				mock.NewGoodClusterStore(),
				mock.NewGoodClusterCloud(),
				mock.NewGoodClusterExe(),
			),
			opts:   mock.DefaultClusterCreateOpts(),
			expect: mock.DefaultCluster(),
		},
		{
			name: "Invalid opts",
			service: core.NewClusterService(
				mock.NewGoodClusterStore(),
				mock.NewGoodClusterCloud(),
				mock.NewGoodClusterExe(),
			),
			opts: api.ClusterCreateOpts{},
			//nolint: lll
			expect:      "failed to validate create cluster input: AWSAccountID: cannot be blank; Cidr: cannot be blank; ClusterName: cannot be blank; Environment: cannot be blank; Region: cannot be blank; RepositoryName: cannot be blank.",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.service.CreateCluster(context.Background(), tc.opts)

			if tc.expectError {
				assert.Error(t, err)
				assert.Equal(t, tc.expect, err.Error())
			} else {
				assert.NoError(t, err)
				assert.Empty(t, cmp.Diff(tc.expect, got))
			}
		})
	}
}
