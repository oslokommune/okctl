package awsalbingresscontroller_test

import (
	"testing"

	"github.com/oslokommune/okctl/pkg/helm/charts/awsalbingresscontroller"
	"github.com/sanathkr/go-yaml"
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
			values: awsalbingresscontroller.NewDefaultValues("my-cluster", "my-service-account"),
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
