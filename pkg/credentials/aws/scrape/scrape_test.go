package scrape_test

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/oslokommune/okctl/pkg/credentials/aws/scrape"
	"github.com/stretchr/testify/assert"
)

func NewResponder(t *testing.T, method, file, url string, code int) {
	//nolint: gosec
	b, err := ioutil.ReadFile(file)
	assert.NoError(t, err)

	response := httpmock.NewBytesResponder(code, b)

	httpmock.RegisterResponder(method, url, response)
}

func TestScrape(t *testing.T) {
	testCases := []struct {
		name         string
		respondersFn func()
		expect       interface{}
		expectError  bool
	}{
		{
			name: "Should work",
			respondersFn: func() {
				NewResponder(t, http.MethodGet, "testdata/login.html", scrape.DefaultURL, http.StatusOK)
				NewResponder(t, http.MethodPost, "testdata/totp.html", "https://doLogin", http.StatusOK)
				NewResponder(t, http.MethodPost, "testdata/saml.html", "https://doTotp", http.StatusOK)
			},
			expect: "SomeNiceSAMLResponse",
		},
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	for _, tc := range testCases {
		tc := tc

		httpmock.Reset()

		t.Run(tc.name, func(t *testing.T) {
			tc.respondersFn()

			got, err := scrape.New().Scrape("user", "pass", "token")
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
