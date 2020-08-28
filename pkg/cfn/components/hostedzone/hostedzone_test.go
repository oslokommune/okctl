package hostedzone_test

import (
	"testing"

	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/oslokommune/okctl/pkg/cfn/components/hostedzone"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	testCases := []struct {
		name   string
		golden string
		hosted *hostedzone.HostedZone
	}{
		{
			name:   "Valid output",
			golden: "hosted-zone",
			hosted: hostedzone.New("test.oslo.systems.", "some comment"),
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			template := cloudformation.NewTemplate()

			template.Resources = map[string]cloudformation.Resource{
				tc.hosted.Name(): tc.hosted.Resource(),
			}

			for k, v := range tc.hosted.NamedOutputs() {
				template.Outputs[k] = v
			}

			got, err := template.YAML()
			assert.NoError(t, err)

			patched, err := hostedzone.PatchYAML(got)
			assert.NoError(t, err)

			g := goldie.New(t)
			g.Assert(t, tc.golden, patched)
		})
	}
}
