package secret_test

import (
	"testing"

	"sigs.k8s.io/yaml"

	"github.com/stretchr/testify/assert"

	"github.com/oslokommune/okctl/pkg/kube/manifests/secret"
	"github.com/sebdah/goldie/v2"
)

func TestNew(t *testing.T) {
	testCases := []struct {
		name   string
		secret *secret.Secret
	}{
		{
			name: "Should work",
			secret: secret.New("name", "default",
				secret.NewManifest(
					"name",
					"default",
					map[string]string{
						"hi": "there",
					},
					nil,
				),
			),
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			g := goldie.New(t)

			dm := tc.secret.Manifest
			db, err := yaml.Marshal(dm)
			assert.NoError(t, err)
			g.Assert(t, "secret.yaml", db)
		})
	}
}
