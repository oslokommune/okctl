package upgrade

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/gosuri/uitable/util/strutil"
	"github.com/oslokommune/okctl/pkg/config/state"
	"github.com/oslokommune/okctl/pkg/github"
	"io"
	"net/http"
	"regexp"
	"strings"
)

type checksumFetcher struct {
}

func (c checksumFetcher) getFor(release *github.RepositoryRelease) ([]state.Checksum, error) {
	checksumAsset, err := c.getChecksumAsset(release)
	if err != nil {
		return nil, fmt.Errorf("getting checksum asset: %w", err)
	}

	checksumsAsBytes, err := c.downloadChecksumsFile(checksumAsset)
	if err != nil {
		return nil, fmt.Errorf("downloading checksum file: %w", err)
	}

	checksums, err := c.parseChecksums(checksumsAsBytes, *release.TagName)
	if err != nil {
		return nil, fmt.Errorf("converting to checksum: %w", err)
	}

	return checksums, nil
}

func (c checksumFetcher) getChecksumAsset(release *github.RepositoryRelease) (*github.ReleaseAsset, error) {
	var checksumAsset *github.ReleaseAsset

	for _, asset := range release.Assets {
		if *asset.Name == "okctl_upgrade_checksums.txt" {
			checksumAsset = asset
		}
	}

	if checksumAsset == nil {
		var assetNameArr []string

		for _, asset := range release.Assets {
			assetNameArr = append(assetNameArr, *asset.Name)
		}

		assetNames := strutil.Join(assetNameArr, ",")

		return nil, fmt.Errorf(
			"could not find checksum asset for release %s (assets: %s)", *release.Name, assetNames)
	}

	return checksumAsset, nil
}

func (c checksumFetcher) downloadChecksumsFile(checksumAsset *github.ReleaseAsset) ([]byte, error) {
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

func (c checksumFetcher) parseChecksums(checksumBytes []byte, expectedVersion string) ([]state.Checksum, error) {
	reader := bytes.NewReader(checksumBytes)
	scanner := bufio.NewScanner(reader)
	var checksums []state.Checksum

	for scanner.Scan() {
		line := scanner.Text()

		parts := strings.Split(line, " ")
		if len(parts) != 2 {
			return nil, fmt.Errorf(
				"expected 2 substrings when splitting digest line on whitespace ( ), got %d in string '%s'",
				len(parts), line,
			)
		}

		digest := parts[0]   // Example: 1eaad82bd6e082936cfb4c108b9e5e46bba98ef19f33492ca2041de04803b86b
		filename := parts[1] // Example: okctl-upgrade_0.0.63_Darwin_amd64.tar.gz

		err := c.validateDigest(digest)
		if err != nil {
			return nil, fmt.Errorf("invalid digest '%s': %w", digest, err)
		}

		ugradeFile, err := parseOkctlUpgradeFilename(filename)
		if err != nil {
			return nil, fmt.Errorf("parsing okctl upgrade filename: %w", err)
		}

		if ugradeFile.version != expectedVersion {
			return nil, fmt.Errorf("expected version '%s' but got '%s' in checksum file '%s'",
				expectedVersion, ugradeFile.version, filename)
		}

		checksum := state.Checksum{
			Os:     ugradeFile.os,
			Arch:   ugradeFile.arch,
			Type:   ugradeFile.extension,
			Digest: digest,
		}

		checksums = append(checksums, checksum)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanning: %w", err)
	}

	return checksums, nil
}

func (c checksumFetcher) validateDigest(digest string) error {
	re, err := regexp.Compile(`^[0-9a-z]+$`)
	if err != nil {
		return fmt.Errorf("compiling regex: %w", err)
	}

	for range re.FindAllString(digest, -1) {
		return nil
	}

	return errors.New("invalid digest")
}
