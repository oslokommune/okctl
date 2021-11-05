package rest

import (
	"bytes"
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

type dummyBody struct {
	Cat string `json:"cat"`
	Dog string `json:"dog"`
}

type dummyCensoredRequest struct {
	Animal dummyBody `json:"animal"`
}

func (d dummyCensoredRequest) AnonymizeRequest(request interface{}) interface{} {
	r := request.(dummyCensoredRequest)

	r.Animal.Dog = "censored"

	return r
}

func TestDebugPrintRequest(t *testing.T) {
	testCases := []struct {
		name string

		withDebugFlag   bool
		withRequestBody interface{}

		expect string
	}{
		{
			name:          "Should produce nothing when debug is false",
			withDebugFlag: false,
			withRequestBody: dummyBody{
				Cat: "Mary",
				Dog: "HadADog",
			},
			expect: "",
		},
		{
			name:          "Should print out a statement when debug is true",
			withDebugFlag: true,
			withRequestBody: dummyBody{
				Cat: "Ben",
				Dog: "Affleck",
			},
			expect: "client (method: , endpoint: ) starting request: rest.dummyBody{\n  Cat: \"Ben\",\n  Dog: \"Affleck\",\n}",
		},
		{
			name:          "Should print out a censored statement when debug is true and anonymized request implemented",
			withDebugFlag: true,
			withRequestBody: dummyCensoredRequest{Animal: dummyBody{
				Cat: "So visible",
				Dog: "Not visible, I hope",
			}},
			expect: "client (method: , endpoint: ) starting request: rest.dummyCensoredRequest{\n  Animal: rest.dummyBody{\n    Cat: \"So visible\",\n    Dog: \"censored\",\n  },\n}", //nolint:lll
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			buf := bytes.Buffer{}

			meta := debugPrintMeta{
				Debug:  tc.withDebugFlag,
				Writer: &buf,
			}

			err := debugPrintRequest(meta, tc.withRequestBody)
			assert.NoError(t, err)

			assert.Equal(t, tc.expect, buf.String())
		})
	}
}

type dummyCensoredResponse struct {
	Animal dummyBody `json:"animal"`
}

func (d dummyCensoredResponse) AnonymizeResponse(response interface{}) interface{} {
	r := response.(dummyCensoredResponse)

	r.Animal.Dog = "so censored"

	return r
}

func TestDebugPrintResponse(t *testing.T) {
	testCases := []struct {
		name string

		withDebugFlag    bool
		withResponseBody interface{}

		expect string
	}{
		{
			name:          "Should produce nothing when debug is false",
			withDebugFlag: false,
			withResponseBody: dummyBody{
				Cat: "Mary",
				Dog: "HadADog",
			},
			expect: "",
		},
		{
			name:          "Should print out a statement when debug is true",
			withDebugFlag: true,
			withResponseBody: dummyBody{
				Cat: "Ben",
				Dog: "Affleck",
			},
			expect: "{\"cat\":\"Ben\",\"dog\":\"Affleck\"}",
		},
		{
			name:          "Should print out a censored statement when debug is true and anonymized response implemented",
			withDebugFlag: true,
			withResponseBody: dummyCensoredResponse{Animal: dummyBody{
				Cat: "So visible",
				Dog: "Not visible, I hope",
			}},
			expect: "{\"animal\":{\"cat\":\"So visible\",\"dog\":\"so censored\"}}",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			buf := bytes.Buffer{}

			meta := debugPrintMeta{
				Debug:  tc.withDebugFlag,
				Writer: &buf,
			}

			out, err := json.Marshal(tc.withResponseBody)
			assert.NoError(t, err)

			err = debugPrintResponse(meta, out, tc.withResponseBody)
			assert.NoError(t, err)

			assert.Equal(t, tc.expect, buf.String())
		})
	}
}
