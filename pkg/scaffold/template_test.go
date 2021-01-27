package scaffold_test

import (
	"io/ioutil"
	"testing"

	"github.com/spf13/afero"

	"github.com/oslokommune/okctl/pkg/scaffold"
	"github.com/stretchr/testify/assert"
)

func TestInterpolateTemplate(t *testing.T) {
	testCases := []struct {
		name string

		withTemplate []byte
		withOpts     *scaffold.InterpolationOpts

		expect []byte
	}{
		{
			name: "Should interpolate the specified field when specified",

			withTemplate: []byte("name: a name\nurl: my-domain.io\nsomethingelse: yes\n"),
			withOpts:     &scaffold.InterpolationOpts{Domain: "works.com"},

			expect: []byte("name: a name\nurl: <app-name>.works.com\nsomethingelse: yes\n"),
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			result, err := scaffold.InterpolateTemplate(tc.withTemplate, tc.withOpts)
			if err != nil {
				t.Error(err)
			}

			assert.Equal(t, tc.expect, result)
		})
	}
}

func TestSaveTemplate(t *testing.T) {
	testCases := []struct {
		name string

		withPath     string
		withTemplate []byte
	}{
		{
			name: "Should find the file and the expected contents",

			withPath:     "application.yaml",
			withTemplate: []byte("name: test-app\nurl: https://example.com\n"),
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			fs := &afero.Afero{Fs: afero.NewMemMapFs()}

			err := scaffold.SaveTemplate(fs, tc.withPath, tc.withTemplate)
			if err != nil {
				t.Fatal(err)
			}

			_, err = fs.Stat(tc.withPath)
			if err != nil {
				// error happens if stat is attempted on a missing file
				t.Fatal(err)
			}

			f, err := fs.Open(tc.withPath)
			if err != nil {
				t.Fatal(err)
			}

			data, err := ioutil.ReadAll(f)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tc.withTemplate, data)
		})
	}
}
