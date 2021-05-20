package externalsecrets_test

import (
	"testing"

	"github.com/oslokommune/okctl/pkg/helm/charts/externalsecrets"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
)

func TestDefaultExternalSecretsValues(t *testing.T) {
	testCases := []struct {
		name   string
		values *externalsecrets.Values
		golden string
	}{
		{
			name:   "External secrets value are valid",
			values: externalsecrets.NewDefaultValues("eu-west-1"),
			golden: "external-secrets-values.yml",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.values.RawYAML()
			assert.NoError(t, err)

			g := goldie.New(t)
			g.Assert(t, tc.golden, got)
		})
	}
}
