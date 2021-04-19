package commands

import (
	"bytes"
	"testing"

	"github.com/spf13/afero"

	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
)

func TestApplicationDeclarationScaffold(t *testing.T) {
	testCases := []struct {
		name string

		withOpts     ScaffoldApplicationOpts
		expectGolden string
	}{
		{
			name: "Should scaffold application declaration based on opts",
			withOpts: ScaffoldApplicationOpts{
				PrimaryHostedZone: "okctl.io",
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer

			err := ScaffoldApplicationDeclaration(&buf, tc.withOpts)
			assert.NoError(t, err)

			g := goldie.New(t)
			g.Assert(t, tc.name, buf.Bytes())
		})
	}
}

func TestEnsureValidDefaultApplicationTemplate(t *testing.T) {
	var buf bytes.Buffer

	err := ScaffoldApplicationDeclaration(&buf, ScaffoldApplicationOpts{})
	assert.NoError(t, err)

	application, err := InferApplicationFromStdinOrFile(&buf, &afero.Afero{}, "-")
	assert.NoError(t, err)

	err = application.Validate()
	assert.NoError(t, err)
}
