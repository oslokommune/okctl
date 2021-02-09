package aws_test

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
	"github.com/google/go-cmp/cmp"
	"github.com/oslokommune/okctl/pkg/api/mock"
	"github.com/oslokommune/okctl/pkg/credentials/aws"
	awsmock "github.com/oslokommune/okctl/pkg/mock"
	"github.com/spf13/afero"
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
			expect: awsmock.DefaultCredentials(),
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
	c := awsmock.DefaultCredentials()
	c.Expires = time.Now().Add(60 * time.Minute)

	testCases := []struct {
		name        string
		auth        aws.Authenticator
		expect      interface{}
		expectError bool
	}{
		{
			name:        "Should work",
			auth:        aws.New(aws.NewInMemoryStorage(), aws.NewAuthStatic(c)),
			expect:      c,
			expectError: false,
		},
		{
			name:        "Should fail, because the creds have expired",
			auth:        aws.New(aws.NewInMemoryStorage(), aws.NewAuthStatic(awsmock.DefaultCredentials())),
			expect:      "no valid credentials: authenticator[0]: expired credentials",
			expectError: true,
		},
		{
			name: "Should work, because one set of creds are valid",
			auth: aws.New(
				aws.NewInMemoryStorage(),
				aws.NewAuthStatic(awsmock.DefaultCredentials()),
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

func TestPersister(t *testing.T) {
	testCases := []struct {
		name      string
		persister aws.Persister
		creds     *aws.Credentials
	}{
		{
			name: "Ini persister",
			persister: aws.NewIniPersister(aws.NewFileSystemIniStorer(
				"conf",
				"creds",
				"/",
				&afero.Afero{Fs: afero.NewMemMapFs()},
			)),
			creds: awsmock.DefaultCredentials(),
		},
		{
			name:      "In memory persister",
			persister: aws.NewInMemoryStorage(),
			creds:     awsmock.DefaultCredentials(),
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Run("Get before save should fail", func(t *testing.T) {
				got, err := tc.persister.Get()
				assert.Error(t, err)
				assert.Nil(t, got)
			})

			t.Run("Save then get should succeed", func(t *testing.T) {
				err := tc.persister.Save(tc.creds)
				assert.NoError(t, err)

				got, err := tc.persister.Get()
				assert.NoError(t, err)
				assert.Empty(t, cmp.Diff(tc.creds, got))
			})
		})
	}
}

func TestAuthEnvironment(t *testing.T) {
	testCases := []struct {
		name string

		withEnv map[string]string

		expectValid bool
	}{
		{
			name: "Should be valid when necessary values are available",

			withEnv: map[string]string{
				"AWS_ACCESS_KEY_ID":     "dummyid",
				"AWS_SECRET_ACCESS_KEY": "dummy-secret",
			},

			expectValid: true,
		},
		{
			name: "Should be invalid when missing secret",

			withEnv: map[string]string{
				"AWS_ACCESS_KEY_ID": "dummyid",
			},

			expectValid: false,
		},
		{
			name: "Should be invalid when missing id",

			withEnv: map[string]string{
				"AWS_ACCESS_KEY_ID": "dummyid",
			},

			expectValid: false,
		},
		{
			name: "Should be invalid when missing everything",

			withEnv: map[string]string{},

			expectValid: false,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			getter := func(key string) string {
				return tc.withEnv[key]
			}

			auth := aws.NewAuthEnvironment("dummy-region", getter)

			assert.Equal(t, tc.expectValid, auth.Valid())
		})
	}
}
