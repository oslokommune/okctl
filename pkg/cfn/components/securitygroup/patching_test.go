package securitygroup

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"sigs.k8s.io/yaml"
)

func TestStringOrFn(t *testing.T) {
	testCases := []struct {
		name string

		withData     []byte
		expectString string
	}{
		{
			name: "Should correctly parse string",

			withData: []byte(`
SourceSecurityGroupId: sg-0512Fasdgfdabafnafeawfea
`),
			expectString: "sg-0512Fasdgfdabafnafeawfea",
		},
		{
			name: "Should correctly parse function",

			withData: []byte(`
SourceSecurityGroupId:
    Fn::GetAtt:
    - RDSPGOutgoingSG
    - GroupId
`),
			expectString: "Fn::GetAtt: RDSPGOutgoingSG,GroupId",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			var result struct {
				Item stringOrFn `json:"SourceSecurityGroupId"`
			}

			err := yaml.Unmarshal(tc.withData, &result)
			assert.NoError(t, err)

			assert.Equal(t, tc.expectString, result.Item.String())
		})
	}
}
