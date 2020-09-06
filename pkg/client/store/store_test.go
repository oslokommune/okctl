package store_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/spf13/afero"

	"github.com/oslokommune/okctl/pkg/client/store"
)

type TestStruct struct {
	Name string
}

func TestOperations(t *testing.T) {
	fs := &afero.Afero{Fs: afero.NewMemMapFs()}

	testCases := []struct {
		name          string
		operations    store.Operations
		expect        interface{}
		expectContent []string
		expectErr     bool
	}{
		{
			name: "Store with defaults",
			operations: store.NewFileSystem("test", fs).
				StoreStruct("test.json", &TestStruct{Name: "hi"}, store.ToJSON()).
				StoreStruct("test.yml", &TestStruct{Name: "ho"}, store.ToYAML()).
				StoreBytes("plain", []byte("hello")),
			expect: &store.Report{
				Type:          "FileSystem",
				Configuration: "CreateDirectories: true\nOverWriteExisting: true\n",
				Actions: []store.Action{
					{
						Name:        "test.json",
						Path:        "test/test.json",
						Type:        "StoreStruct[preprocessing=json]",
						Description: "task.1 StoreStruct[preprocessing=json] to file 'test.json' (path: test/test.json)",
					},
					{
						Name:        "test.yml",
						Path:        "test/test.yml",
						Type:        "StoreStruct[preprocessing=yaml]",
						Description: "task.2 StoreStruct[preprocessing=yaml] to file 'test.yml' (path: test/test.yml)",
					},
					{
						Name:        "plain",
						Path:        "test/plain",
						Type:        "StoreBytes",
						Description: "task.3 StoreBytes to file 'plain' (path: test/plain)",
					},
				},
			},
			expectContent: []string{
				"{\n  \"Name\": \"hi\"\n}",
				"name: ho\n",
				"hello",
			},
		},
		{
			name: "Do not create directories",
			operations: store.NewFileSystem("fail/this", fs, store.FileSystemCreateDirectories(false)).
				StoreBytes("fail", []byte("should fail")),
			expect:    "failed to process task StoreBytes(fail): directory does not exist 'fail/this' and create directories disabled",
			expectErr: true,
		},
		{
			name: "Do not overwrite existing",
			operations: store.NewFileSystem("", fs, store.FileSystemOverwriteExisting(false)).
				StoreBytes("myfile", []byte("content")).
				StoreBytes("myfile", []byte("new content")),
			expect:    "failed to process task StoreBytes(myfile): file 'myfile' exists and overwrite is disabled",
			expectErr: true,
		},
		{
			name: "Remove should work",
			operations: store.NewFileSystem("test", fs, store.FileSystemOverwriteExisting(false)).
				Remove("doesNotExist").
				StoreBytes("file", []byte("content")).
				Remove("file").
				StoreBytes("file", []byte("new content")),
			expect: &store.Report{
				Type:          "FileSystem",
				Configuration: "CreateDirectories: true\nOverWriteExisting: false\n",
				Actions: []store.Action{
					{
						Name:        "doesNotExist",
						Path:        "test/doesNotExist",
						Type:        "Remove",
						Description: "task.1 Remove to file 'doesNotExist' (path: test/doesNotExist)",
					},
					{
						Name:        "file",
						Path:        "test/file",
						Type:        "StoreBytes",
						Description: "task.2 StoreBytes to file 'file' (path: test/file)",
					},
					{
						Name:        "file",
						Path:        "test/file",
						Type:        "Remove",
						Description: "task.3 Remove to file 'file' (path: test/file)",
					},
					{
						Name:        "file",
						Path:        "test/file",
						Type:        "StoreBytes",
						Description: "task.4 StoreBytes to file 'file' (path: test/file)",
					},
				},
			},
			expectContent: []string{
				"",            // Remove
				"new content", // This is hacky, only reads last state, obviously
				"",            // Remove
				"new content",
			},
		},
		{
			name: "Alter should work",
			operations: store.NewFileSystem("test", fs).
				StoreBytes("plain", []byte("hello")).
				AlterStore(store.SetBaseDir("new")).
				StoreBytes("second", []byte("hi")),
			expect: &store.Report{
				Type:          "FileSystem",
				Configuration: "CreateDirectories: true\nOverWriteExisting: true\n",
				Actions: []store.Action{
					{
						Name:        "plain",
						Path:        "test/plain",
						Type:        "StoreBytes",
						Description: "task.1 StoreBytes to file 'plain' (path: test/plain)",
					},
					{
						Name:        "n/a",
						Path:        "n/a",
						Type:        "Alter[SetBaseDir]",
						Description: "task.2 Alter[SetBaseDir] to file 'n/a' (path: n/a)",
					},
					{
						Name:        "second",
						Path:        "new/second",
						Type:        "StoreBytes",
						Description: "task.3 StoreBytes to file 'second' (path: new/second)",
					},
				},
			},
			expectContent: []string{
				"hello",
				"", // Alter
				"hi",
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.operations.Do()
			if tc.expectErr {
				assert.Error(t, err)
				assert.Nil(t, got)
				assert.Equal(t, tc.expect, err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expect, got)
				for i, task := range got.Actions {
					if task.Type == "Remove" || task.Type == "Alter[SetBaseDir]" {
						continue
					}
					content, err := fs.ReadFile(task.Path)
					assert.NoError(t, err)
					assert.Equal(t, tc.expectContent[i], string(content))
				}
			}

			// Reset file system between tests
			fs.Fs = afero.NewMemMapFs()
		})
	}
}

func TestWithFilePermissionsMode(t *testing.T) {
	fs := &afero.Afero{Fs: afero.NewMemMapFs()}

	testCases := []struct {
		name              string
		operations        store.Operations
		expect            interface{}
		expectPermissions os.FileMode
		expectErr         bool
	}{
		{
			name: "Default permissions",
			operations: store.NewFileSystem("test", fs).
				StoreBytes("plain", []byte("hello")),
			expect: &store.Report{
				Type:          "FileSystem",
				Configuration: "CreateDirectories: true\nOverWriteExisting: true\n",
				Actions: []store.Action{
					{
						Name:        "plain",
						Path:        "test/plain",
						Type:        "StoreBytes",
						Description: "task.1 StoreBytes to file 'plain' (path: test/plain)",
					},
				},
			},
			expectPermissions: 0o644,
		},
		{
			name: "Override permissions",
			operations: store.NewFileSystem("test", fs).
				StoreBytes("plain", []byte("hello"), store.WithFilePermissionsMode(0o400)),
			expect: &store.Report{
				Type:          "FileSystem",
				Configuration: "CreateDirectories: true\nOverWriteExisting: true\n",
				Actions: []store.Action{
					{
						Name:        "plain",
						Path:        "test/plain",
						Type:        "StoreBytes",
						Description: "task.1 StoreBytes to file 'plain' (path: test/plain)",
					},
				},
			},
			expectPermissions: 0o400,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.operations.Do()
			if tc.expectErr {
				assert.Error(t, err)
				assert.Nil(t, got)
				assert.Equal(t, tc.expect, err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expect, got)
				assert.Len(t, got.Actions, 1)
				s, err := fs.Stat(got.Actions[0].Path)
				assert.NoError(t, err)
				assert.Equal(t, tc.expectPermissions, s.Mode())
			}

			// Reset file system between tests
			fs.Fs = afero.NewMemMapFs()
		})
	}
}
