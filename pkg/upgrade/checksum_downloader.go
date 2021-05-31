package upgrade

import (
	"fmt"
	"github.com/oslokommune/okctl/pkg/github"
	"io"
	"net/http"
)

type ChecksumDownloader struct {
}

func (c ChecksumDownloader) download(checksumAsset *github.ReleaseAsset) ([]byte, error) {
	response, err := http.Get(*checksumAsset.BrowserDownloadURL)
	if err != nil {
		return nil, fmt.Errorf("http get URL: %s. %w", *checksumAsset.BrowserDownloadURL, err)
	}

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("http get URL returned non-200. URL: %s. Status: %s. Body: %s",
			*checksumAsset.BrowserDownloadURL,
			response.Status,
			response.Body,
		)
	}

	checksumsTxt, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("reading reponse body: %w", err)
	}

	err = response.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("closing response body: %w", err)
	}

	return checksumsTxt, nil
}

func NewChecksumDownloader() ChecksumDownloader {
	return ChecksumDownloader{}
}
