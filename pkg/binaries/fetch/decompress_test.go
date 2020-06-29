package fetch_test

import (
	"bytes"
	"io"
	"io/ioutil"
	"testing"

	"github.com/oslokommune/okctl/pkg/binaries/fetch"
	"github.com/stretchr/testify/assert"
)

func readFile(t *testing.T, file string) io.Reader {
	b, err := ioutil.ReadFile(file)
	assert.NoError(t, err)
	return bytes.NewReader(b)
}

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
			data:         readFile(t, "testdata/myFile.zip"),
			expect:       []byte("some content\n"),
		},
		{
			name:         "ZipDecompressor returns error when file not found",
			decompressor: fetch.NewZipDecompressor("does_not_exist", 1000),
			data:         readFile(t, "testdata/myFile.zip"),
			expect:       "couldn't find: does_not_exist, in archive",
			expectError:  true,
		},
		{
			name:         "ZipDecompressor returns error when size exceeds buffer",
			decompressor: fetch.NewZipDecompressor("myFile", 40),
			data:         readFile(t, "testdata/myFile.zip"),
			expect:       "zip: not a valid zip file",
			expectError:  true,
		},
		{
			name:         "GzipTarDecompressor returns targeted file",
			decompressor: fetch.NewGzipTarDecompressor("myFile", 1000),
			data:         readFile(t, "testdata/myFile.tar.gz"),
			expect:       []byte("some content\n"),
		},
		{
			name:         "GzipTarDecompressor returns errors when file not found",
			decompressor: fetch.NewGzipTarDecompressor("does_not_exist", 1000),
			data:         readFile(t, "testdata/myFile.tar.gz"),
			expect:       "couldn't find: does_not_exist, in archive",
			expectError:  true,
		},
	}

	for _, tc := range testCases {
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
