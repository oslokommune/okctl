package loki_test

import (
	"testing"

	"github.com/oslokommune/okctl/pkg/helm/charts/loki"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
)

func TestNewDefaultValues(t *testing.T) {
	testCases := []struct {
		name   string
		values *loki.Values
		golden string
	}{
		{
			name:   "Default values should generate valid yaml",
			values: loki.NewDefaultValues(),
			golden: "loki.yml",
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
