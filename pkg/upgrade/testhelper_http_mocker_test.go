package upgrade_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/jarcoal/httpmock"
	"github.com/oslokommune/okctl/pkg/github"
)

func mockHTTPResponse(t *testing.T, folder string, releases []*github.RepositoryRelease) {
	err := registerHTTPResponse(
		folder,
		releases,
	)
	require.NoError(t, err)

	httpmock.Activate()
}

func registerHTTPResponse(baseFolder string, releases []*github.RepositoryRelease) error {
	for _, release := range releases {
		versionFolder := *release.TagName

		for _, asset := range release.Assets {
			assetFilename := getAssetFilename(*asset.BrowserDownloadURL)

			data, err := readBytesFromFile(fmt.Sprintf("testdata/%s/%s/%s", baseFolder, versionFolder, assetFilename))
			if err != nil {
				return err
			}

			responder := httpmock.NewBytesResponder(http.StatusOK, data)
			httpmock.RegisterResponder(http.MethodGet, *asset.BrowserDownloadURL, responder)
		}
	}

	return nil
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
