package fetch_test

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/sirupsen/logrus"

	"github.com/oslokommune/okctl/pkg/config/state"

	"github.com/jarcoal/httpmock"
	"github.com/oslokommune/okctl/pkg/binaries/fetch"
	"github.com/oslokommune/okctl/pkg/storage"
	"github.com/stretchr/testify/assert"
)

type nopCloser struct {
	io.Writer
}

func (nopCloser) Close() error { return nil }

// NopCloser returns a WriteCloser with a no-op Close method wrapping
// the provided Writer r.
func NopCloser(w io.Writer) io.WriteCloser {
	return nopCloser{w}
}

func TestStagerFetch(t *testing.T) {
	testCases := []struct {
		name        string
		stager      *fetch.Stager
		expect      interface{}
		expectError bool
	}{
		{
			name: "Nop pipeline",
			stager: &fetch.Stager{
				Storage:      fetch.NewEphemeralStorage(),
				Fetcher:      fetch.NewStaticFetcher([]byte("hi there")),
				Verifier:     fetch.NewNoopVerifier(),
				Decompressor: fetch.NewNoopDecompressor(),
			},
			expect: "hi there",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := bytes.NewBuffer(nil)

			tc.stager.Destination = NopCloser(got)

			err := tc.stager.Fetch()

			if tc.expectError {
				assert.Error(t, err)
				assert.Equal(t, tc.expect, err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expect, got.String())
			}
		})
	}
}

func readBytesFromFile(t *testing.T, file string) []byte {
	//nolint: gosec
	b, err := ioutil.ReadFile(file)
	assert.NoError(t, err)

	return b
}

//nolint: funlen
func TestProcessor(t *testing.T) {
	logger := logrus.StandardLogger()
	logger.Out = ioutil.Discard

	host := state.Host{
		Os:   "darwin",
		Arch: "amd64",
	}

	binaries := []state.Binary{
		{
			Name:       "myBinary",
			Version:    "v0.1.0",
			BufferSize: "10mb",
			URLPattern: "https://localhost/myFile.tar.gz",
			Archive: state.Archive{
				Type:   ".tar.gz",
				Target: "myFile",
			},
			Checksums: []state.Checksum{
				{
					Os:     "darwin",
					Arch:   "amd64",
					Type:   "sha56",
					Digest: "21a47cdf40727a37cde83642e7abab3a1e72a954b905542644f6c543c416a189",
				},
			},
		},
	}

	testCases := []struct {
		preload     bool
		expectError bool
		name        string
		processor   fetch.Provider
		binary      string
		version     string
		preFn       func()
		expect      interface{}
	}{
		{
			name: "Should work",
			processor: func() fetch.Provider {
				p, err := fetch.New(ioutil.Discard, logger, false, host, binaries, storage.NewEphemeralStorage())
				assert.NoError(t, err)
				return p
			}(),
			binary:  "myBinary",
			version: "v0.1.0",
			expect:  "/binaries/myBinary/v0.1.0/darwin/amd64/myBinary",
			preFn: func() {
				responder := httpmock.NewBytesResponder(200, readBytesFromFile(t, "testdata/myFile.tar.gz"))
				httpmock.RegisterResponder(http.MethodGet, "https://localhost/myFile.tar.gz", responder)
			},
		},
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.processor.Fetch(tc.binary, tc.version)
			if tc.expectError {
				assert.Error(t, err)
				assert.Equal(t, tc.expect, err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expect, got)
			}
		})
	}
}
