package config_test

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/oslokommune/okctl/pkg/config"
	"github.com/stretchr/testify/assert"
)

// nolint
var confContent = `
binaries:
  - name: eksctl
    version: 0.18.0
    bufferSize: 100mb
    urlPattern: https://pattern
    archive:
      type: .tar.gz
      target: eksctl
    checksums:
      - os: darwin
        arch: amd64
        type: sha256
        digest: something
`

func TestLoadConfiguration(t *testing.T) {
	testCases := []struct {
		name        string
		fn          func(c *config.UserConfig) interface{}
		expectError bool
		expect      interface{}
	}{
		{
			name: "binary checksums",
			fn: func(c *config.UserConfig) interface{} {
				return c.Binaries[0].Checksums
			},
			expect: []config.Checksum{
				{
					Os:     "darwin",
					Arch:   "amd64",
					Type:   "sha256",
					Digest: "something",
				},
			},
		},
	}

	dir, err := ioutil.TempDir("", "config")
	assert.NoError(t, err)

	err = os.Mkdir(path.Join(dir, ".okctl"), 0755)
	assert.NoError(t, err)

	defer func() {
		_ = os.RemoveAll(dir)
	}()

	err = ioutil.WriteFile(path.Join(dir, ".okctl", "user.yml"), []byte(confContent), 0600)
	assert.NoError(t, err)
	err = os.Chdir(dir)
	assert.NoError(t, err)

	c, err := config.LoadUserConfiguration(dir)
	assert.NoError(t, err)

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			if tc.expectError {
				assert.Nil(t, c)
				assert.Equal(t, tc.expect, err.Error())
			} else {
				assert.NoError(t, err, tc.name)
				assert.Equal(t, tc.expect, tc.fn(c))
			}
		})
	}
}
