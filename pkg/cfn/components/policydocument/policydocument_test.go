package policydocument_test

import (
	"encoding/base64"
	"testing"

	"github.com/oslokommune/okctl/pkg/cfn/components/policydocument"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
)

func TestPolicyDocument(t *testing.T) {
	testCases := []struct {
		name     string
		golden   string
		document *policydocument.PolicyDocument
	}{
		{
			name:   "Validate action and resource",
			golden: "action-resource.json",
			document: &policydocument.PolicyDocument{
				Version: policydocument.Version,
				ID:      "4CFB8EB0-B7A8-4A6D-870A-A2C23FCD01CC",
				Statement: []policydocument.StatementEntry{
					{
						Sid:    "1",
						Effect: policydocument.EffectTypeAllow,
						Action: []string{
							"*",
						},
						Resource: []string{
							"*",
						},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.document.JSON()
			assert.NoError(t, err)

			g := goldie.New(t)
			g.Assert(t, tc.golden, got)
		})
	}
}

func TestRef(t *testing.T) {
	testCases := []struct {
		name   string
		ref    string
		expect string
	}{
		{
			name:   "AwsRegionRef",
			ref:    policydocument.AwsRegionRef(),
			expect: "{ \"Ref\": \"AWS::Region\" }",
		},
		{
			name:   "AwsAccountIDRef",
			ref:    policydocument.AwsAccountIDRef(),
			expect: "{ \"Ref\": \"AWS::AccountId\" }",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			got, err := base64.StdEncoding.DecodeString(tc.ref)
			assert.NoError(t, err)
			assert.Equal(t, tc.expect, string(got))
		})
	}
}
