package securitygrouppolicy_test

import (
	"testing"

	"github.com/oslokommune/okctl/pkg/kube/manifests/securitygrouppolicy"
	"github.com/oslokommune/okctl/pkg/kube/securitygrouppolicy/api/types/v1beta1"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
	"sigs.k8s.io/yaml"
)

func TestNew(t *testing.T) {
	testCases := []struct {
		name     string
		manifest *v1beta1.SecurityGroupPolicy
	}{
		{
			name: "Should work",
			manifest: securitygrouppolicy.Manifest(
				"name",
				"default",
				map[string]string{
					"app": "psqlclient",
				},
				[]string{
					"sg-a90eug3FAKE",
				},
			),
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			g := goldie.New(t)

			db, err := yaml.Marshal(tc.manifest)
			assert.NoError(t, err)
			g.Assert(t, "securitygrouppolicy.yaml", db)
		})
	}
}
