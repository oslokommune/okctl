package aws_test

import (
	"testing"
	"time"

	aws2 "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
	"github.com/oslokommune/okctl/pkg/api/mock"
	"github.com/oslokommune/okctl/pkg/credentials/aws"
	awsmock "github.com/oslokommune/okctl/pkg/mock"
	"github.com/stretchr/testify/assert"
)

func TestNewAuthSAML(t *testing.T) {
	testCases := []struct {
		name        string
		retriever   aws.Retriever
		provider    aws.StsProviderFn
		expect      interface{}
		expectError bool
	}{
		{
			name: "SAML retriever should work",
			retriever: aws.NewAuthSAML(
				mock.DefaultAWSAccountID,
				mock.DefaultRegion,
				awsmock.NewGoodScraper(),
				func(session *session.Session) stsiface.STSAPI {
					return awsmock.NewGoodSTSAPI()
				},
				aws.Static("byr999999", "the", "123456"),
			),
			expect: awsmock.DefaultStsCredentials(),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.retriever.Retrieve()
			if tc.expectError {
				assert.Error(t, err)
				assert.Equal(t, tc.expect, err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expect, got)
			}
		})
	}
}

func TestAuthRaw(t *testing.T) {
	c := awsmock.DefaultStsCredentials()
	c.Expiration = aws2.Time(time.Now().Add(60 * time.Minute))

	testCases := []struct {
		name        string
		auth        aws.Authenticator
		expect      interface{}
		expectError bool
	}{
		{
			name:        "Should work",
			auth:        aws.New(aws.NewAuthStatic(c)),
			expect:      c,
			expectError: false,
		},
		{
			name:        "Should fail, because the creds have expired",
			auth:        aws.New(aws.NewAuthStatic(awsmock.DefaultStsCredentials())),
			expect:      "no valid credentials: authenticator[0]: expired credentials",
			expectError: true,
		},
		{
			name: "Should work, because one set of creds are valid",
			auth: aws.New(
				aws.NewAuthStatic(awsmock.DefaultStsCredentials()),
				aws.NewAuthStatic(c),
			),
			expect: c,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.auth.Raw()
			if tc.expectError {
				assert.Error(t, err)
				assert.Equal(t, tc.expect, err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expect, got)
			}
		})
	}
}
