package jsonpatch

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshallPatch(t *testing.T) {
	testCases := []struct {
		name string

		withPatch Patch

		expectResult []byte
	}{
		{
			name: "Should correctly serialize a replace operation",

			withPatch: Patch{Operations: []Operation{
				{
					Type:  OperationTypeReplace,
					Path:  "/spec/type",
					Value: "LoadBalancer",
				},
			}},

			expectResult: []byte(`[{"op":"replace","path":"/spec/type","value":"LoadBalancer"}]`),
		},
		{
			name: "Should correctly serialize an add operation on an object",

			withPatch: Patch{Operations: []Operation{
				{
					Type: OperationTypeAdd,
					Path: "/metadata/annotations",
					Value: map[string]string{
						"something.arn": "arn:acc:regi:something",
					},
				},
			}},

			expectResult: []byte(`[{"op":"add","path":"/metadata/annotations","value":{"something.arn":"arn:acc:regi:something"}}]`),
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			result, _ := json.Marshal(tc.withPatch)

			assert.Equal(t, tc.expectResult, result)
		})
	}
}

func TestHasOperations(t *testing.T) {
	testCases := []struct {
		name           string
		withPatch      Patch
		expectContains []Operation
	}{
		{
			name: "Should work",
			withPatch: Patch{
				Operations: []Operation{
					{
						Type:  OperationTypeAdd,
						Path:  "/hello",
						Value: "howdy",
					},
					{
						Type:  OperationTypeAdd,
						Path:  "/bye",
						Value: "bye",
					},
				},
			},
			expectContains: []Operation{
				{
					Type:  OperationTypeAdd,
					Path:  "/bye",
					Value: "bye",
				},
			},
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			assert.True(t, tc.withPatch.HasOperations(tc.expectContains))
		})
	}
}
