package upgrade_test

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/jarcoal/httpmock"
	"github.com/oslokommune/okctl/pkg/github"
)

// mockHTTPResponsesForGithubReleases mocks Github releases from the folder [baseFolder/release.TagName]
func mockHTTPResponsesForGithubReleases(t *testing.T, baseFolder string, releases []*github.RepositoryRelease) {
	// shuffle slice to better test sorting of releases
	shuffled := shuffle(releases)

	for _, release := range shuffled {
		registerHTTPResponseFromReleaseTagFolder(
			t,
			baseFolder,
			release,
		)
	}

	httpmock.Activate()
}

func shuffle(src []*github.RepositoryRelease) []*github.RepositoryRelease {
	dest := make([]*github.RepositoryRelease, len(src))

	rand.Seed(123)

	perm := rand.Perm(len(src))
	for i, v := range perm {
		dest[v] = src[i]
	}

	return dest
}

func registerHTTPResponseFromReleaseTagFolder(t *testing.T, baseFolder string, release *github.RepositoryRelease) {
	versionFolder := *release.TagName
	mockHTTPResponseForGithubRelease(t, path.Join(baseFolder, versionFolder), release)
}

// mockHTTPResponseForGithubRelease mocks the Github release from the given folder, ignoring whatever the
// release.TagName is
func mockHTTPResponseForGithubRelease(t *testing.T, upgradeBinaryFolder string, release *github.RepositoryRelease) {
	for _, asset := range release.Assets {
		assetFilename := getAssetFilename(*asset.BrowserDownloadURL)

		data, err := readBytesFromFile(fmt.Sprintf("testdata/%s/%s", upgradeBinaryFolder, assetFilename))
		require.NoError(t, err)

		responder := httpmock.NewBytesResponder(http.StatusOK, data)
		httpmock.RegisterResponder(http.MethodGet, *asset.BrowserDownloadURL, responder)
	}
}

func getAssetFilename(url string) string {
	split := strings.Split(url, "/")
	return split[len(split)-1]
}

func readBytesFromFile(file string) ([]byte, error) {
	//nolint: gosec
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	return b, nil
}
