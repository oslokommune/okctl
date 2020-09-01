package awsalbingresscontroller_test

import (
	"testing"

	"gopkg.in/yaml.v3"

	"github.com/oslokommune/okctl/pkg/api/mock"
	"github.com/oslokommune/okctl/pkg/helm/charts/awsalbingresscontroller"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
)

func TestNewDefaultValues(t *testing.T) {
	testCases := []struct {
		name   string
		values *awsalbingresscontroller.Values
		golden string
	}{
		{
			name:   "Default values should generate valid yaml",
			values: awsalbingresscontroller.NewDefaultValues("my-cluster", mock.DefaultVpcID, mock.DefaultRegion),
			golden: "aws-alb-ingress-controller-values.yaml",
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
