package pgbouncer_test

import (
	"testing"

	"github.com/oslokommune/okctl/pkg/kube/manifests/pgbouncer"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
	"sigs.k8s.io/yaml"
)

func TestResources(t *testing.T) {
	testCases := []struct {
		name     string
		resource interface{}
	}{
		{
			name: "default-pod",
			resource: pgbouncer.Pod(
				"myBouncer",
				"",
				"test",
				"someSecret",
				"paramsConfigMap",
				"paramsSecret",
				map[string]string{"label": "value"},
				5432,
			),
		},
		{
			name: "pod-with-database",
			resource: pgbouncer.Pod(
				"myBouncer",
				"myDatabase",
				"test",
				"someSecret",
				"paramsConfigMap",
				"paramsSecret",
				map[string]string{"label": "value"},
				5432,
			),
		},
		{
			name: "secret",
			resource: pgbouncer.Secret(
				"mySecret",
				"test",
				"administrator",
				"secret",
			),
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			g := goldie.New(t)

			got, err := yaml.Marshal(tc.resource)
			assert.NoError(t, err)
			g.Assert(t, tc.name, got)
		})
	}
}
