package storage_test

import (
	"path/filepath"
	"testing"

	"github.com/oslokommune/okctl/pkg/storage"
	"github.com/stretchr/testify/assert"
)

// nolint: funlen
func TestStorer(t *testing.T) {
	testCases := []struct {
		name    string
		storer  storage.Storer
		content string
		dir     string
		file    string
	}{
		{
			name:    "Storage should work",
			storer:  storage.NewEphemeralStorage(),
			content: "hi there",
			dir:     "test",
			file:    "stuff",
		},
		{
			name: "Temporary storage should work",
			storer: func() storage.Storer {
				s, err := storage.NewTemporaryStorage()
				assert.NoError(t, err)
				return s
			}(),
			content: "hi theree",
			dir:     "test",
			file:    "stuff",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			relative := filepath.Join(tc.dir, tc.file)

			f, err := tc.storer.Create(tc.dir, tc.file, 0644)
			assert.NoError(t, err)

			_, err = f.WriteString(tc.content)
			assert.NoError(t, err)

			yes, err := tc.storer.Exists(relative)
			assert.NoError(t, err)
			assert.True(t, yes)

			abs := tc.storer.Abs(relative)
			assert.Equal(t, filepath.Join(tc.storer.Path(), tc.dir, tc.file), abs)

			err = tc.storer.RemoveAll(relative)
			assert.NoError(t, err)

			yes, err = tc.storer.Exists(relative)
			assert.NoError(t, err)
			assert.False(t, yes)

			err = tc.storer.MkdirAll(tc.dir)
			assert.NoError(t, err)

			f, err = tc.storer.Recreate(tc.dir, tc.file, 0644)
			assert.NoError(t, err)

			_, err = f.WriteString(tc.content)
			assert.NoError(t, err)

			data, err := tc.storer.ReadAll(relative)
			assert.NoError(t, err)
			assert.Equal(t, tc.content, string(data))
		})
	}
}
