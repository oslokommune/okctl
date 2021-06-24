package upgrade

import (
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
		return nil, fmt.Errorf("http call did not return status OK. URL: %s. Status: %s. Body: %s",
			*checksumAsset.BrowserDownloadURL,
			response.Status,
			response.Body,
		)
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

// NewChecksumDownloader returns a new ChecksumDownloader
func NewChecksumDownloader() ChecksumDownloader {
	return ChecksumDownloader{}
}
