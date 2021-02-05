package awslbc_test

import (
	"testing"

	"github.com/oslokommune/okctl/pkg/helm/charts/awslbc"

	"gopkg.in/yaml.v3"

	"github.com/oslokommune/okctl/pkg/api/mock"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
)

func TestNewDefaultValues(t *testing.T) {
	testCases := []struct {
		name   string
		values *awslbc.Values
		golden string
	}{
		{
			name:   "Default values should generate valid yaml",
			values: awslbc.NewDefaultValues("my-cluster", mock.DefaultVpcID, mock.DefaultRegion),
			golden: "aws-load-balancer-controller.yml",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			got, err := yaml.Marshal(tc.values)
			assert.NoError(t, err)

			g := goldie.New(t)
			g.Assert(t, tc.golden, got)
		})
	}
}
