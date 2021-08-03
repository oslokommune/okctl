package upgrade

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	digestPkg "github.com/oslokommune/okctl/pkg/binaries/digest"

	ghPkg "github.com/google/go-github/v32/github"
	"github.com/oslokommune/okctl/pkg/config/state"
	"github.com/oslokommune/okctl/pkg/github"
)

// GithubReleaseParser parses github releases
type GithubReleaseParser struct {
	checksumDownloader ChecksumDownloader
}

const expectedSubStringsGithubReleaseDigestFile = 2

func (g GithubReleaseParser) toUpgradeBinaries(releases []*github.RepositoryRelease) ([]okctlUpgradeBinary, error) {
	upgrades := make([]okctlUpgradeBinary, 0, len(releases))

	for _, release := range releases {
		upgrade, err := g.parseRelease(release)
		if err != nil {
			return nil, fmt.Errorf("validating release: %w", err)
		}

		upgrades = append(upgrades, upgrade)
	}

	return upgrades, nil
}

const expectedMinimumAssetsForAGithubRelease = 2

func (g GithubReleaseParser) validateRelease(r *github.RepositoryRelease) error {
	return validation.ValidateStruct(r,
		validation.Field(&r.ID, validation.Required),
		validation.Field(&r.Name, validation.Required),
		validation.Field(&r.TagName, validation.Required),
		validation.Field(&r.Assets, validation.By(func(value interface{}) error {
			assets := value.([]*github.ReleaseAsset)
			if len(assets) < expectedMinimumAssetsForAGithubRelease {
				releaseIdentifier := ""
				if r.Name == nil {
					if r.ID == nil {
						releaseIdentifier = "(unknown)"
					} else {
						releaseIdentifier = fmt.Sprintf("ID: %d", *r.ID)
					}
				} else {
					releaseIdentifier = *r.Name
				}

				return fmt.Errorf(
					"release '%s' must have at least %d assets (binary and checksum) ",
					releaseIdentifier,
					expectedMinimumAssetsForAGithubRelease)
			}

			return nil
		})),
	)
}

func (g GithubReleaseParser) parseRelease(release *github.RepositoryRelease) (okctlUpgradeBinary, error) {
	err := g.validateRelease(release)
	if err != nil {
		return okctlUpgradeBinary{}, err
	}

	releaseUpgradeVersion := *release.TagName

	var binaryChecksums []state.Checksum

	for _, asset := range release.Assets {
		if *asset.Name == "okctl-upgrade-checksums.txt" {
			binaryChecksums, err = g.fetchChecksums(asset)
			if err != nil {
				return okctlUpgradeBinary{}, fmt.Errorf("fetching checksums: %w", err)
			}
		} else {
			err = g.validateUpgradeBinaryAsset(asset, releaseUpgradeVersion)
			if err != nil {
				return okctlUpgradeBinary{}, fmt.Errorf("validating upgrade binary asset: %w", err)
			}
		}
	}

	if binaryChecksums == nil {
		var assetNameArr []string

		for _, asset := range release.Assets {
			assetNameArr = append(assetNameArr, *asset.Name)
		}

		assetNames := strings.Join(assetNameArr, ",")

		return okctlUpgradeBinary{}, fmt.Errorf(
			"could not find checksum asset for release %s (assets: %s)", *release.Name, assetNames)
	}

	binaryVersion, err := parseUpgradeBinaryVersion(releaseUpgradeVersion)
	if err != nil {
		return okctlUpgradeBinary{}, fmt.Errorf("parsing upgrade version: %w", err)
	}

	return newOkctlUpgradeBinary(binaryVersion, binaryChecksums), nil
}

func (g GithubReleaseParser) fetchChecksums(asset *ghPkg.ReleaseAsset) ([]state.Checksum, error) {
	checksumsAsBytes, err := g.checksumDownloader.download(asset)
	if err != nil {
		return nil, fmt.Errorf("downloading checksum file: %w", err)
	}

	checksums, err := g.parseChecksums(checksumsAsBytes)
	if err != nil {
		return nil, fmt.Errorf("converting to checksum: %w", err)
	}

	return checksums, nil
}

func (g GithubReleaseParser) parseChecksums(checksumBytes []byte) ([]state.Checksum, error) {
	reader := bytes.NewReader(checksumBytes)
	scanner := bufio.NewScanner(reader)

	var checksums []state.Checksum

	for scanner.Scan() {
		line := scanner.Text()

		parts := strings.Fields(line)
		if len(parts) != expectedSubStringsGithubReleaseDigestFile {
			return nil, fmt.Errorf(
				"expected %d substrings when splitting digest line on whitespace ( ), got %d in string '%s'",
				expectedSubStringsGithubReleaseDigestFile, len(parts), line,
			)
		}

		digest := parts[0]   // Example: 1eaad82bd6e082936cfb4c108b9e5e46bba98ef19f33492ca2041de04803b86b
		filename := parts[1] // Example: okctl-upgrade_0.0.63_Darwin_amd64.tar.gz

		err := validateDigest(digest)
		if err != nil {
			return nil, fmt.Errorf("invalid digest '%s': %w", digest, err)
		}

		ugradeFile, err := parseOkctlUpgradeFilename(filename)
		if err != nil {
			return nil, fmt.Errorf("parsing upgrade filename: %w", err)
		}

		checksum := state.Checksum{
			Os:     ugradeFile.os,
			Arch:   ugradeFile.arch,
			Type:   string(digestPkg.TypeSHA256),
			Digest: digest,
		}

		checksums = append(checksums, checksum)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanning: %w", err)
	}

	return checksums, nil
}

func validateDigest(digest string) error {
	re := regexp.MustCompile(`^[0-9a-z]+$`)

	for range re.FindAllString(digest, -1) {
		return nil
	}

	return errors.New("invalid digest")
}

func (g GithubReleaseParser) validateUpgradeBinaryAsset(asset *ghPkg.ReleaseAsset, releaseUpgradeVersion string) error {
	downloadURLParts := strings.Split(*asset.BrowserDownloadURL, "/")
	if len(downloadURLParts) == 0 {
		return fmt.Errorf("expected at least 1 '/' in browser download URL %s", *asset.BrowserDownloadURL)
	}

	upgradeBinaryFilename := downloadURLParts[len(downloadURLParts)-1]

	upgradeFile, err := parseOkctlUpgradeFilename(upgradeBinaryFilename)
	if err != nil {
		return fmt.Errorf("cannot parse upgrade filename '%s': %w", upgradeBinaryFilename, err)
	}

	if upgradeFile.version != releaseUpgradeVersion {
		return fmt.Errorf("expected okctl upgrade binary version '%s' to equal release upgrade version '%s'",
			upgradeFile.version, releaseUpgradeVersion)
	}

	return nil
}

// NewGithubReleaseParser returns a new GithubReleaseParser
func NewGithubReleaseParser(checksumDownloader ChecksumDownloader) GithubReleaseParser {
	return GithubReleaseParser{
		checksumDownloader: checksumDownloader,
	}
}
