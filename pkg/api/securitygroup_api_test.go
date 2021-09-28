package api

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createValidRule() Rule {
	return Rule{
		Description: "Required to be cool",
		FromPort:    1337,
		ToPort:      1337,
		CidrIP:      "0.0.0.0/0",
		Protocol:    RuleProtocolTCP,
	}
}

func TestRuleValidation(t *testing.T) {
	testCases := []struct {
		name        string
		withRule    Rule
		expectError string
	}{
		{
			name:     "Should accept a valid rule IPv4 rule",
			withRule: createValidRule(),
		},
		{
			name: "Should err when missing both cidr and sourceSecurityGroupID",
			withRule: func() Rule {
				r := createValidRule()

				r.SourceSecurityGroupID = ""
				r.CidrIP = ""

				return r
			}(),
			expectError: "CidrIp: required when SourceSecurityGroupID is empty; SourceSecurityGroupId: required when CidrIP is empty.",
		},
		{
			name: "Should err when both cidr and sourceSecurityGroupID are specified",
			withRule: func() Rule {
				r := createValidRule()

				r.SourceSecurityGroupID = "sg-fe321fea421"
				r.CidrIP = "10.0.0.0/24"

				return r
			}(),
			expectError: "CidrIp: must be blank if SourceSecurityGroupID is specified; SourceSecurityGroupId: must be blank if CidrIP is specified.",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			err := tc.withRule.Validate()

			if tc.expectError != "" {
				assert.NotNil(t, err, fmt.Sprintf("expected error: %s", tc.expectError))

				if err != nil {
					assert.Equal(t, tc.expectError, err.Error())
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
