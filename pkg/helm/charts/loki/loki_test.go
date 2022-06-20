package loki_test

import (
	"testing"
	"time"

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
			name: "Default values should generate valid yaml",
			values: func() *loki.Values {
				v := loki.NewDefaultValues("mock-bucket-name", "mock-prefix_", "eu-mock-1")

				t, _ := time.Parse("2006-01-02", "2006-01-02")
				v.FromDate = t.Format("2006-01-02")

				return v
			}(),
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
