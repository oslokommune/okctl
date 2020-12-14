package cognito_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/oslokommune/okctl/pkg/mock"

	"github.com/oslokommune/okctl/pkg/cognito"
)

func TestUserPoolDomainInfo(t *testing.T) {
	testCases := []struct {
		name    string
		cognito *cognito.Cognito
		expect  *cognito.UserPoolDomainInfo
	}{
		{
			name:    "Should work",
			cognito: cognito.New(mock.NewGoodCloudProvider(), nil),
			expect: &cognito.UserPoolDomainInfo{
				CloudFrontDomainName: "cloudfront-us-east-1.something.aws.com",
				UserPoolDomain:       "auth.oslo.systems",
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.cognito.UserPoolDomainInfo("auth.oslo.systems")
			assert.NoError(t, err)
			assert.Equal(t, tc.expect, got)
		})
	}
}
