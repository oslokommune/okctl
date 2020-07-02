package fetch_test

import (
	"bytes"
	"io"
	"io/ioutil"
	"testing"

	"github.com/oslokommune/okctl/pkg/binaries/fetch"
	"github.com/stretchr/testify/assert"
)

const (
	zipTestFile   = "testdata/myFile.zip"
	tarGzTestFile = "testdata/myFile.tar.gz"
)

func readFile(t *testing.T, file string) io.Reader {
	//nolint: gosec
	b, err := ioutil.ReadFile(file)
	assert.NoError(t, err)

	return bytes.NewReader(b)
}

//nolint: funlen
func TestDecompressor(t *testing.T) {
	testCases := []struct {
		name         string
		decompressor fetch.Decompressor
		data         io.Reader
		expect       interface{}
		expectError  bool
	}{
		{
			name:         "NoopDecompressor returns raw data",
			decompressor: fetch.NewNoopDecompressor(),
			data:         bytes.NewReader([]byte{'a'}),
			expect:       []byte{'a'},
		},
		{
			name:         "ZipDecompressor returns targeted file",
			decompressor: fetch.NewZipDecompressor("myFile", 1000),
			data:         readFile(t, zipTestFile),
			expect:       []byte("some content\n"),
		},
		{
			name:         "ZipDecompressor returns error when file not found",
			decompressor: fetch.NewZipDecompressor("does_not_exist", 1000),
			data:         readFile(t, zipTestFile),
			expect:       "couldn't find: does_not_exist, in archive",
			expectError:  true,
		},
		{
			name:         "ZipDecompressor returns error when size exceeds buffer",
			decompressor: fetch.NewZipDecompressor("myFile", 40),
			data:         readFile(t, zipTestFile),
			expect:       "zip: not a valid zip file",
			expectError:  true,
		},
		{
			name:         "GzipTarDecompressor returns targeted file",
			decompressor: fetch.NewGzipTarDecompressor("myFile", 1000),
			data:         readFile(t, tarGzTestFile),
			expect:       []byte("some content\n"),
		},
		{
			name:         "GzipTarDecompressor returns errors when file not found",
			decompressor: fetch.NewGzipTarDecompressor("does_not_exist", 1000),
			data:         readFile(t, tarGzTestFile),
			expect:       "couldn't find: does_not_exist, in archive",
			expectError:  true,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			got := new(bytes.Buffer)

			err := tc.decompressor.Decompress(tc.data, got)

			if tc.expectError {
				assert.Error(t, err)
				assert.Equal(t, tc.expect, err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expect, got.Bytes())
			}
		})
	}
}
