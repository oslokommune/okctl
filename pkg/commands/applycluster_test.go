package commands

import (
	"bytes"
	"testing"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/stretchr/testify/assert"
)

func TestClusterDefaultsAfterInfering(t *testing.T) {
	testCases := []struct {
		name string

		withDeclaration string
		expectFn        func(cluster *v1alpha1.Cluster) bool
	}{
		{
			name: "Should have integrations set to true by default",

			withDeclaration: ``,
			expectFn: func(cluster *v1alpha1.Cluster) bool {
				return cluster.Integrations.ArgoCD != false
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			buf := bytes.NewBufferString(tc.withDeclaration)

			cluster, err := InferClusterFromStdinOrFile(buf, "-")
			assert.NoError(t, err)

			assert.True(t, tc.expectFn(cluster))
		})
	}
}
