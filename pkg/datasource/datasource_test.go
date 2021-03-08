package datasource_test

import (
	"testing"

	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
	"sigs.k8s.io/yaml"

	"github.com/oslokommune/okctl/pkg/datasource"
)

func TestNewDatasource(t *testing.T) {
	testCases := []struct {
		name string
		ds   *datasource.Datasources
	}{
		{
			name: "loki",
			ds:   datasource.NewLoki(),
		},
		{
			name: "tempo",
			ds:   datasource.NewTempo(),
		},
		{
			name: "cloudwatch",
			ds:   datasource.NewCloudWatch("eu-west-1"),
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			got, err := yaml.Marshal(tc.ds)
			assert.NoError(t, err)
			g := goldie.New(t)
			g.Assert(t, tc.name, got)
		})
	}
}
