package github_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/oslokommune/okctl/pkg/keyring"

	"github.com/stretchr/testify/assert"

	"github.com/oslokommune/okctl/pkg/credentials/github"
)

// ClientMock provides a mock of HTTPClient
type ClientMock struct {
	Fn       func(r *http.Request)
	Response *http.Response
}

// Do performs the work
func (c *ClientMock) Do(req *http.Request) (*http.Response, error) {
	if c.Fn != nil {
		c.Fn(req)
	}

	return c.Response, nil
}

// nolint: funlen
func TestAreValid(t *testing.T) {
	testCases := []struct {
		name         string
		credentials  *github.Credentials
		expect       interface{}
		expectErr    bool
		requestURL   string
		responseCode int
	}{
		{
			name: "Should work",
			credentials: &github.Credentials{
				AccessToken: "tokenabc12345",
			},
			requestURL:   "https://api.github.com",
			responseCode: http.StatusOK,
		},
		{
			name: "Invalid token response code",
			credentials: &github.Credentials{
				AccessToken: "tokenabc12345",
			},
			expect:       "HTTP error 404 (Not Found) when requesting token validation",
			expectErr:    true,
			requestURL:   "https://api.github.com",
			responseCode: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			client := &ClientMock{
				Fn: func(r *http.Request) {
					assert.Equal(t, tc.requestURL, r.URL.String())
				},
				Response: &http.Response{
					StatusCode:    tc.responseCode,
					Body:          ioutil.NopCloser(bytes.NewReader(nil)),
					ContentLength: 0,
				},
			}

			err := github.AreValid(tc.credentials, client)
			if tc.expectErr {
				assert.Error(t, err)
				assert.Equal(t, tc.expect, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAuthRaw(t *testing.T) {
	testCases := []struct {
		name      string
		auth      github.Authenticator
		expect    interface{}
		expectErr bool
	}{
		{
			name: "Should work",
			auth: github.New(github.NewInMemoryPersister(), &ClientMock{
				// nolint: bodyclose
				Response: func() *http.Response {
					return &http.Response{
						StatusCode:    http.StatusOK,
						Body:          ioutil.NopCloser(bytes.NewReader(nil)),
						ContentLength: 0,
					}
				}(),
			}, github.NewAuthStatic(&github.Credentials{
				AccessToken: "token",
				ClientID:    "client_id",
				Type:        github.CredentialsTypeDeviceFlow,
			})),
			expect: &github.Credentials{
				AccessToken: "token",
				ClientID:    "client_id",
				Type:        github.CredentialsTypeDeviceFlow,
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.auth.Raw()
			if tc.expectErr {
				assert.Error(t, err)
				assert.Equal(t, tc.expect, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expect, got)
			}
		})
	}
}

func TestNewKeyringPersister(t *testing.T) {
	testCases := []struct {
		name        string
		keyring     keyring.Keyringer
		credentials *github.Credentials
	}{
		{
			name: "Should work",
			credentials: &github.Credentials{
				AccessToken: "token",
				ClientID:    "client_id",
				Type:        github.CredentialsTypeDeviceFlow,
			},
			keyring: func() keyring.Keyringer {
				k, err := keyring.New(keyring.NewInMemoryKeyring(), false)
				assert.NoError(t, err)

				return k
			}(),
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			p := github.NewKeyringPersister(tc.keyring)

			// Get before save should fail
			got, err := p.Get()
			assert.Error(t, err)
			assert.Nil(t, got)

			// Save should succeed
			err = p.Save(tc.credentials)
			assert.NoError(t, err)

			// Get should work now
			got, err = p.Get()
			assert.NoError(t, err)
			assert.Equal(t, tc.credentials, got)
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
				"GITHUB_TOKEN": "dummytoken",
			},

			expectValid: true,
		},
		{
			name: "Should be invalid when missing token",

			withEnv: map[string]string{
				"GTHB_TKN": "misspelled-token-key",
			},

			expectValid: false,
		},
		{
			name: "Should be invalid when token is blank",

			withEnv: map[string]string{
				"GITHUB_TOKEN": "",
			},

			expectValid: false,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			getter := func(key string) string {
				return tc.withEnv[key]
			}

			auth := github.NewAuthEnvironment(getter)

			assert.Equal(t, tc.expectValid, auth.Valid())
		})
	}
}
