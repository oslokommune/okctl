package externalsecrets_test

import (
	"testing"

	"github.com/oslokommune/okctl/pkg/helm/charts/externalsecrets"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestDefaultExternalSecretsValues(t *testing.T) {
	testCases := []struct {
		name   string
		values *externalsecrets.Values
		golden string
	}{
		{
			name:   "External secrets value are valid",
			values: externalsecrets.DefaultExternalSecretsValues(),
			golden: "external-secrets-values.yml",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			b, err := yaml.Marshal(tc.values)
			assert.NoError(t, err)

			g := goldie.New(t)
			g.Assert(t, tc.golden, b)
		})
	}
}
