package upgrade

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/oslokommune/okctl/pkg/github"
)

// ChecksumDownloader knows how to download the checksums for a github release asset
type ChecksumDownloader struct{}

func (c ChecksumDownloader) download(checksumAsset *github.ReleaseAsset) ([]byte, error) {
	response, err := http.Get(*checksumAsset.BrowserDownloadURL)
	if err != nil {
		return nil, fmt.Errorf("http get URL: %s. %w", *checksumAsset.BrowserDownloadURL, err)
	}

	if response.StatusCode != http.StatusOK {
		return c.statusNotOkError(checksumAsset, response)
	}

	checksumsTxt, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	err = response.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("closing response body: %w", err)
	}

	return checksumsTxt, nil
}

func (c ChecksumDownloader) statusNotOkError(checksumAsset *github.ReleaseAsset, response *http.Response) ([]byte, error) {
	var err error

	var buf bytes.Buffer

	var body string

	_, err = buf.ReadFrom(response.Body)
	if err == nil {
		body = buf.String()
	} else {
		body = "(could not read body)"
	}

	return nil, fmt.Errorf("http call did not return status OK. URL: %s. Status: %s. Body: %s",
		*checksumAsset.BrowserDownloadURL,
		response.Status,
		body,
	)
}

// NewChecksumDownloader returns a new ChecksumDownloader
func NewChecksumDownloader() ChecksumDownloader {
	return ChecksumDownloader{}
}
