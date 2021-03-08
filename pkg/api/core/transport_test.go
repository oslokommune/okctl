package core

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestDecodeStruct provides a simple struct
// for testing the decoder
type TestDecodeStruct struct {
	Name string
}

func TestDecodeStructRequest(t *testing.T) {
	testCases := []struct {
		name      string
		into      interface{}
		request   *http.Request
		expectErr bool
		expect    interface{}
	}{
		{
			name: "Should work",
			into: &TestDecodeStruct{},
			request: httptest.NewRequest(
				"",
				"/something",
				strings.NewReader(`{"Name": "Bob"}`),
			),
			expect: &TestDecodeStruct{
				Name: "Bob",
			},
		},
		{
			name: "Should fail",
			into: &TestDecodeStruct{},
			request: httptest.NewRequest(
				"",
				"/something",
				strings.NewReader(""),
			),
			expectErr: true,
			expect:    "decoding request as json: EOF",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			got, err := decodeStructRequest(tc.into)(context.Background(), tc.request)
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
