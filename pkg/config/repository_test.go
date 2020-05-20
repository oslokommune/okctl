package config_test

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/oslokommune/okctl/pkg/config"
	"github.com/stretchr/testify/assert"
)

func createRepositoryTestConfig(t *testing.T, content, fileName string) string {
	dir, err := ioutil.TempDir("", "config")
	assert.NoError(t, err)

	err = ioutil.WriteFile(path.Join(dir, fileName), []byte(content), 0600)
	assert.NoError(t, err)

	err = os.Chdir(dir)
	assert.NoError(t, err)

	return dir
}

// nolint
func TestLoadRepository(t *testing.T) {
	testCases := []struct {
		name      string
		fileName  string
		content   string
		expect    interface{}
		expectErr bool
	}{
		{
			name:     "Valid configuration",
			fileName: config.DefaultRepositoryConfig,
			expect: &config.RepoConfig{
				Name:    "okctl",
				Region:  "eu-west-1",
				BaseDir: "deployment",
				Clusters: []config.Cluster{
					{
						Name: "dev",
						AWS: config.AWS{
							Account: 3456789,
							Profile: "abc",
						},
					},
				},
			},
			content: `name: okctl
region: eu-west-1
baseDir: deployment
clusters:
  - name: dev
    aws:
      account: 3456789
      profile: abc
`,
		},
		{
			name:     "Empty configuration",
			fileName: config.DefaultRepositoryConfig,
			content:  "",
			expect:   &config.RepoConfig{},
		},
		{
			name:      "No configuration",
			fileName:  "wrong.yml",
			content:   "",
			expect:    "Config File \".okctl\" Not Found in",
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got, err := config.LoadRepo(createRepositoryTestConfig(t, tc.content, tc.fileName))
			if tc.expectErr {
				assert.Contains(t, err.Error(), tc.expect)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expect, got)
			}
		})
	}
}
