package transport_test

import (
	"context"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/oslokommune/okctl/pkg/transport"
	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	Title string
	Body  string
}

func (t *TestStruct) Text() []byte {
	return []byte(fmt.Sprintf("Title: %s\nBody: %s", t.Title, t.Body))
}

func TestEncodeTextResponse(t *testing.T) {
	testCases := []struct {
		name     string
		response interface{}
		expect   interface{}
	}{
		{
			name: "Texter should work",
			response: &TestStruct{
				Title: "Hello",
				Body:  "This is hello.",
			},
			expect: "Title: Hello\nBody: This is hello.",
		},
		{
			name:     "No formatter should work",
			response: "hi there!",
			expect:   "\"hi there!\"",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := httptest.NewRecorder()

			err := transport.EncodeTextResponse(context.Background(), got, tc.response)

			assert.NoError(t, err)
			assert.Equal(t, tc.expect, got.Body.String())
		})
	}
}

func TestEncodeYAMLResponse(t *testing.T) {
	testCases := []struct {
		name        string
		response    interface{}
		expect      interface{}
		expectError bool
	}{
		{
			name: "Yaml should work",
			response: &TestStruct{
				Title: "Hello",
				Body:  "This is hello.",
			},
			expect: "Body: This is hello.\nTitle: Hello\n",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := httptest.NewRecorder()

			err := transport.EncodeYAMLResponse(context.Background(), got, tc.response)

			if tc.expectError {
				assert.Error(t, err)
				assert.Equal(t, tc.expect, err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expect, got.Body.String())
			}
		})
	}
}
