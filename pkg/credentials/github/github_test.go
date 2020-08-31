package github_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
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
	validBody, err := json.Marshal(&github.TokenVerification{
		Scopes: github.RequiredScopes(),
	})
	assert.NoError(t, err)

	invalidBody, err := json.Marshal(&github.TokenVerification{
		Scopes: []string{},
	})
	assert.NoError(t, err)

	testCases := []struct {
		name         string
		credentials  *github.Credentials
		expect       interface{}
		expectErr    bool
		requestBody  string
		requestURL   string
		responseBody string
		responseCode int
		responseErr  error
	}{
		{
			name: "Should work",
			credentials: &github.Credentials{
				AccessToken: "tokenabc12345",
				ClientID:    "clientxyz12345",
			},
			requestURL:   "https://api.github.com/applications/clientxyz12345/token",
			requestBody:  "access_token=tokenabc12345",
			responseCode: http.StatusOK,
			responseBody: string(validBody),
		},
		{
			name: "Invalid token response code",
			credentials: &github.Credentials{
				AccessToken: "tokenabc12345",
				ClientID:    "clientxyz12345",
			},
			expect:       "HTTP error 404 (Not Found) when requesting token validation",
			expectErr:    true,
			requestBody:  "access_token=tokenabc12345",
			requestURL:   "https://api.github.com/applications/clientxyz12345/token",
			responseBody: "",
			responseCode: http.StatusNotFound,
		},
		{
			name: "Invalid response body",
			credentials: &github.Credentials{
				AccessToken: "tokenabc12345",
				ClientID:    "clientxyz12345",
			},
			expect:       "failed to parse token response: unexpected EOF",
			expectErr:    true,
			requestBody:  "access_token=tokenabc12345",
			requestURL:   "https://api.github.com/applications/clientxyz12345/token",
			responseBody: "{",
			responseCode: http.StatusOK,
		},
		{
			name: "Invalid scopes",
			credentials: &github.Credentials{
				AccessToken: "tokenabc12345",
				ClientID:    "clientxyz12345",
			},
			expect:       "token does not contain required scopes: repo, read:org",
			expectErr:    true,
			requestBody:  "access_token=tokenabc12345",
			requestURL:   "https://api.github.com/applications/clientxyz12345/token",
			responseBody: string(invalidBody),
			responseCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			client := &ClientMock{
				Fn: func(r *http.Request) {
					assert.Equal(t, tc.requestURL, r.URL.String())
					data, err := ioutil.ReadAll(r.Body)
					assert.NoError(t, err)
					assert.Equal(t, tc.requestBody, string(data))
				},
				Response: &http.Response{
					StatusCode:    tc.responseCode,
					Body:          ioutil.NopCloser(strings.NewReader(tc.responseBody)),
					ContentLength: int64(len(tc.responseBody)),
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
	validBody, err := json.Marshal(&github.TokenVerification{
		Scopes: github.RequiredScopes(),
	})
	assert.NoError(t, err)

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
						Body:          ioutil.NopCloser(strings.NewReader(string(validBody))),
						ContentLength: int64(len(validBody)),
					}
				}(),
			}, github.NewAuthStatic(&github.Credentials{
				AccessToken: "token",
				ClientID:    "client_id",
			})),
			expect: &github.Credentials{
				AccessToken: "token",
				ClientID:    "client_id",
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
			},
			keyring: func() keyring.Keyringer {
				k, err := keyring.New(keyring.NewInMemoryKeyring())
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
