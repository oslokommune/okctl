package handlers

import (
	"bytes"
	"testing"

	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
)

func TestWriteDeleteApplicationReadyCheckInfo(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		withIngress       bool
		withECRRepository bool
	}{
		{
			name: "Should work with no optionals",
		},
		{
			name:        "Should add additional information with ingress",
			withIngress: true,
		},
		{
			name:              "Should add additional information with ECR",
			withECRRepository: true,
		},
		{
			name:              "Should include all information when everything is active",
			withECRRepository: true,
			withIngress:       true,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			stdout := bytes.Buffer{}

			err := writeDeleteApplicationReadyCheckInfo(&stdout, deleteApplicationPromptTemplateOpts{
				ApplicationName: "mockapp",
				HasIngress:      tc.withIngress,
				HasECR:          tc.withECRRepository,
			})
			assert.NoError(t, err)

			g := goldie.New(t)
			g.Assert(t, t.Name(), stdout.Bytes())
		})
	}
}
