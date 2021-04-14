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

// nolint: funlen
func TestOperations(t *testing.T) {
	fs := &afero.Afero{Fs: afero.NewMemMapFs()}

	inline := &TestStruct{}

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
				Data: map[string]interface{}{},
			},
			expectContent: []string{
				"{\n  \"Name\": \"hi\"\n}",
				"Name: ho\n",
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
						Description: "task.1 Remove file 'doesNotExist' (path: test/doesNotExist)",
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
						Description: "task.3 Remove file 'file' (path: test/file)",
					},
					{
						Name:        "file",
						Path:        "test/file",
						Type:        "StoreBytes",
						Description: "task.4 StoreBytes to file 'file' (path: test/file)",
					},
				},
				Data: map[string]interface{}{},
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
						Type:        "Alter[alterer=SetBaseDir]",
						Description: "task.2 Alter[alterer=SetBaseDir]",
					},
					{
						Name:        "second",
						Path:        "new/second",
						Type:        "StoreBytes",
						Description: "task.3 StoreBytes to file 'second' (path: new/second)",
					},
				},
				Data: map[string]interface{}{},
			},
			expectContent: []string{
				"hello",
				"", // Alter
				"hi",
			},
		},
		{
			name: "Add operations should work",
			operations: store.NewFileSystem("test", fs).
				AddStoreStruct(store.AddStoreStruct{
					Name:         "test.json",
					Data:         &TestStruct{Name: "hi"},
					PreProcessor: store.ToJSON(),
				}).
				AddStoreBytes(store.AddStoreBytes{
					Name: "plain",
					Data: []byte("hello"),
				}),
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
						Name:        "plain",
						Path:        "test/plain",
						Type:        "StoreBytes",
						Description: "task.2 StoreBytes to file 'plain' (path: test/plain)",
					},
				},
				Data: map[string]interface{}{},
			},
			expectContent: []string{
				"{\n  \"Name\": \"hi\"\n}",
				"hello",
			},
		},
		{
			name: "Read should work",
			operations: store.NewFileSystem("test", fs).
				StoreBytes("first", []byte("first")).
				GetBytes("first", nil),
			expect: &store.Report{
				Type:          "FileSystem",
				Configuration: "CreateDirectories: true\nOverWriteExisting: true\n",
				Actions: []store.Action{
					{
						Name:        "first",
						Path:        "test/first",
						Type:        "StoreBytes",
						Description: "task.1 StoreBytes to file 'first' (path: test/first)",
					},
					{
						Name:        "first",
						Path:        "test/first",
						Type:        "GetBytes",
						Description: "task.2 GetBytes from file 'first' (path: test/first)",
					},
				},
				Data: map[string]interface{}{
					"first": []byte("first"),
				},
			},
			expectContent: []string{
				"first",
				"first",
			},
			expectErr: false,
		},
		{
			name: "Process should work",
			operations: store.NewFileSystem("", fs).
				StoreBytes("post-file", []byte("some content")).
				StoreStruct("something", &TestStruct{Name: "post-file"}, store.ToJSON()).
				GetStruct("something", inline, store.FromJSON()).
				ProcessGetStruct("something", func(data interface{}, operations store.Operations) error {
					_ = operations.GetBytes("post-file", nil)
					return nil
				}),
			expect: &store.Report{
				Type:          "FileSystem",
				Configuration: "CreateDirectories: true\nOverWriteExisting: true\n",
				Actions: []store.Action{
					{
						Name:        "post-file",
						Path:        "post-file",
						Type:        "StoreBytes",
						Description: "task.1 StoreBytes to file 'post-file' (path: post-file)",
					},
					{
						Name:        "something",
						Path:        "something",
						Type:        "StoreStruct[preprocessing=json]",
						Description: "task.2 StoreStruct[preprocessing=json] to file 'something' (path: something)",
					},
					{
						Name:        "something",
						Path:        "something",
						Type:        "GetStruct[postprocessor=json]",
						Description: "task.3 GetStruct[postprocessor=json] from file 'something' (path: something)",
					},
					{
						Name:        "something",
						Type:        "ProcessGetStruct",
						Description: "task.4 ProcessGetStruct on name 'something",
					},
					{
						Name:        "post-file",
						Path:        "post-file",
						Type:        "GetBytes",
						Description: "task.5 GetBytes from file 'post-file' (path: post-file)",
					},
				},
				Data: map[string]interface{}{
					"something": &TestStruct{Name: "post-file"},
					"post-file": []byte("some content"),
				},
			},
			expectContent: []string{
				"some content",
				"{\n  \"Name\": \"post-file\"\n}",
				"{\n  \"Name\": \"post-file\"\n}",
				"",
				"some content",
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
					if task.Type == "Remove" || task.Type == "Alter[SetBaseDir]" || task.Type == "ProcessGetStruct" {
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

// nolint: funlen
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
				Data: map[string]interface{}{},
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
				Data: map[string]interface{}{},
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

func TestRemoveDir(t *testing.T) {
	const perm = 0o744

	fs := &afero.Afero{Fs: afero.NewMemMapFs()}
	_ = fs.Mkdir("testing/sub1", perm)

	fileSystem := store.NewFileSystem("testing", fs)

	isEmpty, _ := fs.IsEmpty("testing/sub1")
	isDirectory, _ := fs.IsDir("testing/sub1")

	assert.True(t, isEmpty)
	assert.True(t, isDirectory)

	_, _ = fileSystem.RemoveDir("sub1").Do()

	isDirectory, _ = fs.IsDir("testing/sub1")

	assert.False(t, isDirectory)
}

func TestRemoveDirWithContent(t *testing.T) {
	const perm = 0o744

	fs := &afero.Afero{Fs: afero.NewMemMapFs()}

	_ = fs.Mkdir("testing/sub1", perm)
	_ = fs.Mkdir("testing/sub1", perm)
	_ = fs.WriteFile("testing/sub1/hellofile", []byte("hello"), perm)

	fileSystem := store.NewFileSystem("testing", fs)
	_, _ = fileSystem.RemoveDir("sub1").Do()

	isDirectory, _ := fs.IsDir("testing/sub1")
	isEmpty, _ := fs.IsEmpty("testing/sub1")

	assert.True(t, isDirectory)
	assert.False(t, isEmpty)
}
