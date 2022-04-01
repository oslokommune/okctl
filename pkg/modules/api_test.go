package modules

import (
	"fmt"
	"testing"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/sebdah/goldie/v2"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestDownload(t *testing.T) {
	testCases := []struct {
		name                string
		withModuleName      string
		withTargetDirectory string
	}{
		{
			name:           "Should work",
			withModuleName: "sqs",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			moduleFs := memfs.New()
			fs := &afero.Afero{Fs: afero.NewMemMapFs()}

			err := fs.MkdirAll(tc.withTargetDirectory, 0o700)
			assert.NoError(t, err)

			err = acquireModules(moduleFs)
			assert.NoError(t, err)
		})
	}
}

func TestCopyToFs(t *testing.T) {
	testCases := []struct {
		name               string
		withModule         string
		withDestinationDir string

		expectFiles []string
	}{
		{
			name:               "Should download and install a module to a filesystem",
			withModule:         "sqs",
			withDestinationDir: "/",
			expectFiles: []string{
				"/sqs.tf",
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			moduleFs := memfs.New()
			fs := &afero.Afero{Fs: afero.NewMemMapFs()}

			err := acquireModules(moduleFs)
			assert.NoError(t, err)

			err = copyModuleToFs(moduleFs, fs, tc.withModule, "/")
			assert.NoError(t, err)

			g := goldie.New(t)

			for index, filepath := range tc.expectFiles {
				exists, err := fs.Exists(filepath)
				assert.NoError(t, err)

				assert.True(t, exists)

				content, err := fs.ReadFile(filepath)
				assert.NoError(t, err)

				g.Assert(t, fmt.Sprintf("%s%d", tc.name, index), content)
			}
		})
	}
}
