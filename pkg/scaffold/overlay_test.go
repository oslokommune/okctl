package scaffold

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMarshallPatch(t *testing.T) {
	testCases := []struct {
		name string

		withPatch Patch

		expectResult []byte
	}{
		{
			name: "Should work",

			withPatch: Patch{Operations: []Operation{
				{
					Type:  OperationTypeReplace,
					Path:  "/spec/type",
					Value: "LoadBalancer",
				},
			}},

			expectResult: []byte(`[{"op": "replace", "path": "/spec/type", "value": "LoadBalancer"}]`),
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer

			_ = json.NewEncoder(&buf).Encode(tc.withPatch)

			assert.Equal(t, tc.expectResult, buf.Bytes())
		})
	}
}
