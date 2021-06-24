package upgrade

import (
	"fmt"
	"strings"

	ghPkg "github.com/google/go-github/v32/github"
	"github.com/oslokommune/okctl/pkg/binaries/run/okctlupgrade"
	"github.com/oslokommune/okctl/pkg/config/state"
	"github.com/oslokommune/okctl/pkg/github"
)

// GithubReleaseParser parses github releases
type GithubReleaseParser struct {
	checksumDownloader ChecksumDownloader
}

func (g GithubReleaseParser) toUpgradeBinaries(releases []*github.RepositoryRelease) ([]okctlUpgradeBinary, error) {
	upgrades := make([]okctlUpgradeBinary, len(releases))

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

func (g GithubReleaseParser) parseRelease(release *github.RepositoryRelease) (okctlUpgradeBinary, error) {
	if release.Name == nil || len(*release.Name) == 0 {
		return okctlUpgradeBinary{}, fmt.Errorf("release ID '%d' name must be non-empty", *release.ID)
	}

	if release.TagName == nil {
		return okctlUpgradeBinary{}, fmt.Errorf("release '%s' tag name must be non-empty", *release.Name)
	}

	if len(release.Assets) < expectedMinimumAssetsForAGithubRelease {
		return okctlUpgradeBinary{}, fmt.Errorf(
			"release '%s' must have at least %d assets (binary and checksum) ",
			*release.Name,
			expectedMinimumAssetsForAGithubRelease)
	}

	releaseUpgradeVersion := *release.TagName
	binaryName := fmt.Sprintf(okctlupgrade.BinaryNameFormat, releaseUpgradeVersion)

	var binaryChecksums []state.Checksum

	var err error

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

	return okctlUpgradeBinary{
		name:          binaryName,
		fileExtension: ".tar.gz",
		version:       releaseUpgradeVersion,
		checksums:     binaryChecksums,
	}, nil
}

func (g GithubReleaseParser) fetchChecksums(asset *ghPkg.ReleaseAsset) ([]state.Checksum, error) {
	checksumsAsBytes, err := g.checksumDownloader.download(asset)
	if err != nil {
		return nil, fmt.Errorf("downloading checksum file: %w", err)
	}

	checksums, err := parseChecksums(checksumsAsBytes)
	if err != nil {
		return nil, fmt.Errorf("converting to checksum: %w", err)
	}

	return checksums, nil
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
