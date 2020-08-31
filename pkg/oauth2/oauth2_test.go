package oauth2_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/oslokommune/okctl/pkg/oauth2"
	oauth2pkg "golang.org/x/oauth2"

	"github.com/stretchr/testify/assert"
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
func TestRequestDeviceCode(t *testing.T) {
	validCode := &oauth2.DeviceCode{
		DeviceCode:      "device_code",
		UserCode:        "user_code",
		VerificationURI: "https://verify",
		Interval:        5,
	}

	validBody, err := json.Marshal(validCode)
	assert.NoError(t, err)

	testCases := []struct {
		name         string
		config       *oauth2.Config
		requestURL   string
		requestBody  string
		responseCode int
		responseBody string
		expectErr    bool
		expect       interface{}
	}{
		{
			name: "Should work",
			config: &oauth2.Config{
				Config: &oauth2pkg.Config{
					ClientID: "client_id",
					Scopes:   []string{"something"},
				},
				DeviceEndpoint: oauth2.DeviceEndpoint{
					CodeURL: "https://something",
				},
			},
			requestURL:   "https://something",
			requestBody:  "client_id=client_id&scope=something",
			responseCode: http.StatusOK,
			responseBody: string(validBody),
			expect:       validCode,
		},
		{
			name: "Invalid response code",
			config: &oauth2.Config{
				Config: &oauth2pkg.Config{
					ClientID: "client_id",
					Scopes:   []string{"something"},
				},
				DeviceEndpoint: oauth2.DeviceEndpoint{
					CodeURL: "https://something",
				},
			},
			requestURL:   "https://something",
			requestBody:  "client_id=client_id&scope=something",
			responseCode: http.StatusBadRequest,
			expectErr:    true,
			expect:       "request for device code authorisation returned status 400 (Bad Request)",
		},
		{
			name: "Invalid response body",
			config: &oauth2.Config{
				Config: &oauth2pkg.Config{
					ClientID: "client_id",
					Scopes:   []string{"something"},
				},
				DeviceEndpoint: oauth2.DeviceEndpoint{
					CodeURL: "https://something",
				},
			},
			requestURL:   "https://something",
			requestBody:  "client_id=client_id&scope=something",
			responseCode: http.StatusOK,
			responseBody: "{",
			expectErr:    true,
			expect:       "failed to decode response: {, because: unexpected EOF",
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

			got, err := oauth2.RequestDeviceCode(client, tc.config)
			if tc.expectErr {
				assert.Error(t, err)
				assert.Equal(t, tc.expect, err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expect, got)
			}
		})
	}
}

// nolint: funlen
func TestWaitForDeviceAuthorization(t *testing.T) {
	validToken := &oauth2pkg.Token{
		AccessToken: "access_token",
	}

	validBody, err := json.Marshal(validToken)
	assert.NoError(t, err)

	testCases := []struct {
		name         string
		config       *oauth2.Config
		code         *oauth2.DeviceCode
		requestURL   string
		requestBody  string
		responseCode int
		responseBody string
		expectErr    bool
		expect       interface{}
	}{
		{
			name: "Should work",
			config: &oauth2.Config{
				Config: &oauth2pkg.Config{
					ClientID: "client_id",
					Endpoint: oauth2pkg.Endpoint{
						TokenURL: "https://token",
					},
				},
			},
			code: &oauth2.DeviceCode{
				DeviceCode: "device_code",
				Interval:   5,
			},
			requestURL:   "https://token",
			requestBody:  "client_id=client_id&device_code=device_code&grant_type=urn%3Aietf%3Aparams%3Aoauth%3Agrant-type%3Adevice_code",
			responseCode: http.StatusOK,
			responseBody: string(validBody),
			expect:       validToken,
		},
		{
			name: "Invalid status code",
			config: &oauth2.Config{
				Config: &oauth2pkg.Config{
					ClientID: "client_id",
					Endpoint: oauth2pkg.Endpoint{
						TokenURL: "https://token",
					},
				},
			},
			code: &oauth2.DeviceCode{
				DeviceCode: "device_code",
				Interval:   5,
			},
			requestURL:   "https://token",
			requestBody:  "client_id=client_id&device_code=device_code&grant_type=urn%3Aietf%3Aparams%3Aoauth%3Agrant-type%3Adevice_code",
			responseCode: http.StatusBadRequest,
			expectErr:    true,
			expect:       "HTTP error 400 (Bad Request) when polling for OAuth token",
		},
		{
			name: "Invalid response body",
			config: &oauth2.Config{
				Config: &oauth2pkg.Config{
					ClientID: "client_id",
					Endpoint: oauth2pkg.Endpoint{
						TokenURL: "https://token",
					},
				},
			},
			code: &oauth2.DeviceCode{
				DeviceCode: "device_code",
				Interval:   5,
			},
			requestURL:   "https://token",
			requestBody:  "client_id=client_id&device_code=device_code&grant_type=urn%3Aietf%3Aparams%3Aoauth%3Agrant-type%3Adevice_code",
			responseCode: http.StatusOK,
			responseBody: "{",
			expectErr:    true,
			expect:       "failed to decode polling response: unexpected EOF",
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

			got, err := oauth2.WaitForDeviceAuthorization(client, tc.config, tc.code)
			if tc.expectErr {
				assert.Error(t, err)
				assert.Equal(t, tc.expect, err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expect, got)
			}
		})
	}
}
