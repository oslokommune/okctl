package rest

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mishudark/errors"
	"github.com/stretchr/testify/assert"
)

type testMarshalledErrorServer struct {
	kind errors.Kind
}

func (t testMarshalledErrorServer) ServeHTTP(writer http.ResponseWriter, _ *http.Request) {
	err := errors.E(errors.New("something bad happened"), t.kind)

	raw, _ := json.Marshal(err)

	// Must be >= 400 to trigger error
	writer.WriteHeader(400)

	_, _ = writer.Write(raw)
}

func TestErrorSerialization(t *testing.T) {
	testCases := []struct {
		name string

		withKind                            errors.Kind
		expectDifferentAfterDeserialization bool
	}{
		{
			name: "Should deserialize into Timeout error",

			withKind: errors.Timeout,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(testMarshalledErrorServer{tc.withKind})

			defer server.Close()

			client := HTTPClient{
				BaseURL: server.URL,
				Client:  &http.Client{},
			}

			err := client.Do(http.MethodPost, "/dummy", "", "")

			assert.True(t, errors.IsKind(err, tc.withKind))
		})
	}
}

func TestKindSanity(t *testing.T) {
	server := httptest.NewServer(testMarshalledErrorServer{errors.Timeout})

	client := HTTPClient{
		BaseURL: server.URL,
		Client:  &http.Client{},
	}

	err := client.Do(http.MethodPost, "/dummy", "", "")

	assert.False(t, errors.IsKind(err, errors.Decrypt))
}
