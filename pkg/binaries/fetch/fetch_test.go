package fetch_test

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/oslokommune/okctl/pkg/binaries/fetch"
	"github.com/stretchr/testify/assert"
)

func TestHttpFetcherFetch(t *testing.T) {
	testCases := []struct {
		name      string
		url       string
		expectErr bool
		expect    interface{}
		response  httpmock.Responder
	}{
		{
			name:     "Valid request",
			url:      "https://valid",
			expect:   "hi there",
			response: httpmock.NewBytesResponder(200, []byte("hi there")),
		},
		{
			name:      "Invalid URL scheme",
			url:       "http://invalid",
			expectErr: true,
			expect:    "a valid pkgURL must begin https://, got: http://invalid",
		},
		{
			name:      "Connection error",
			url:       "https://connection",
			expectErr: true,
			expect:    "Get \"https://connection\": no responder found",
			response:  httpmock.ConnectionFailure,
		},
		{
			name:      "Internal error",
			url:       "https://internal",
			expectErr: true,
			expect:    "bad status: 500",
			response:  httpmock.NewBytesResponder(http.StatusInternalServerError, []byte("oops")),
		},
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			httpmock.RegisterResponder(http.MethodGet, tc.url, tc.response)

			buf := &bytes.Buffer{}
			_, err := fetch.NewHTTPFetcher(tc.url).Fetch(buf)

			if tc.expectErr {
				assert.Equal(t, tc.expect, err.Error())
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tc.expect, buf.String())
			}
		})
	}
}
