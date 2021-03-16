package main

import (
	"testing"

	"gotest.tools/assert"
)

func createValidOptsBuilder() *DeleteClusterOptsBuilder {
	return &DeleteClusterOptsBuilder{
		opts: DeleteClusterOpts{
			Region:       "eu-west-1",
			AWSAccountID: "123456789012",
			Environment:  "test",
			Repository:   "kjoremiljo",
			ClusterName:  "kjoremiljo-test",
		},
	}
}

type DeleteClusterOptsBuilder struct {
	opts DeleteClusterOpts
}

func (b *DeleteClusterOptsBuilder) ClusterName(newName string) *DeleteClusterOptsBuilder {
	b.opts.ClusterName = newName

	return b
}

func TestKM170(t *testing.T) {
	// Ensures that the expected format of opts.ClusterName is correct. This bug was due to using the
	// o.RepoStateWithEnv.GetCluster().Name instead of o.RepoStateWithEnv.GetClusterName()
	testCases := []struct {
		name string

		withOpts  *DeleteClusterOptsBuilder
		expectErr string
	}{
		{
			name: "Should pass with all required and valid opts",

			withOpts: createValidOptsBuilder(),
		},
		{
			name: "Should break if clustername doesnt equal name-env",

			withOpts:  createValidOptsBuilder().ClusterName("thiswontwork"),
			expectErr: `ClusterName: internal error: cluster name needs to be in the format repository-environment, but was "thiswontwork".`,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			err := tc.withOpts.opts.Validate()

			if tc.expectErr != "" {
				assert.Error(t, err, tc.expectErr)
			} else {
				assert.NilError(t, err)
			}
		})
	}
}
