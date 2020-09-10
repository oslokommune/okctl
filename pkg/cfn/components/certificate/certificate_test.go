package certificate_test

import (
	"testing"

	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/oslokommune/okctl/pkg/cfn/components/certificate"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	testCases := []struct {
		name   string
		golden string
		cert   *certificate.Certificate
	}{
		{
			name:   "Validate output",
			golden: "certificate",
			cert:   certificate.New("test.oslo.systems.", "AZ12345678"),
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			template := cloudformation.NewTemplate()

			template.Resources = map[string]cloudformation.Resource{
				tc.cert.Name(): tc.cert.Resource(),
			}

			for k, v := range tc.cert.NamedOutputs() {
				template.Outputs[k] = v
			}

			got, err := template.YAML()
			assert.NoError(t, err)

			g := goldie.New(t)
			g.Assert(t, tc.golden, got)
		})
	}
}
