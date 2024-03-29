package blockstorage_test

import (
	"testing"

	"github.com/oslokommune/okctl/pkg/helm/charts/blockstorage"

	"github.com/oslokommune/okctl/pkg/api/mock"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
)

func TestNewDefaultValues(t *testing.T) {
	testCases := []struct {
		name   string
		values *blockstorage.Values
		golden string
	}{
		{
			name:   "Default values should generate valid yaml",
			values: blockstorage.NewDefaultValues(mock.DefaultRegion, mock.DefaultClusterName, "blockstorage"),
			golden: "blockstorage.yml",
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
