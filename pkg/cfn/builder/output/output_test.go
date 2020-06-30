package output_test

import (
	"encoding/base64"
	"testing"

	"github.com/oslokommune/okctl/pkg/cfn/builder/output"
	"github.com/stretchr/testify/assert"
)

func TestJoined_Outputs(t *testing.T) {
	testCases := []struct {
		name     string
		outputer output.Outputer
		expect   map[string]map[string]interface{}
	}{
		{
			name:     "Joined",
			outputer: output.NewJoined("JoinedTest").Add("value"),
			expect: map[string]map[string]interface{}{
				"JoinedTest": {
					"Value": base64.StdEncoding.EncodeToString([]byte("{ \"Fn::Join\": [ \",\", [ \"value\" ] ] }")),
				},
			},
		},
		{
			name:     "Value",
			outputer: output.NewValue("ValueTest", "value"),
			expect: map[string]map[string]interface{}{
				"ValueTest": {
					"Value": "value",
				},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expect, tc.outputer.NamedOutputs())
		})
	}
}
