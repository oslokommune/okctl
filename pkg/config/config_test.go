package config_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/oslokommune/okctl/pkg/config"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func stableAppCfg(id string) *config.AppConfig {
	c := config.NewDefaultAppCfg()
	c.User.ID = id

	return c
}

func contentFromStruct(t *testing.T, content interface{}) string {
	c, err := yaml.Marshal(content)
	assert.NoError(t, err)

	return string(c)
}

func createAppTestConfig(t *testing.T, content, fileName string) string {
	dir, err := ioutil.TempDir("", "config")
	assert.NoError(t, err)

	err = os.MkdirAll(path.Join(dir, config.DefaultAppDir), 0744)
	assert.NoError(t, err)

	err = ioutil.WriteFile(path.Join(dir, config.DefaultAppDir, fileName), []byte(content), 0600)
	assert.NoError(t, err)

	err = os.Chdir(dir)
	assert.NoError(t, err)

	return dir
}

// nolint
func TestLoadApp(t *testing.T) {
	testCases := []struct {
		name        string
		fileName    string
		content     string
		preFn       func()
		appCfgFn    config.AppCfgFn
		notFoundFn  config.AppNotFoundFn
		expectError bool
		expect      interface{}
	}{
		{
			name:     "Full config",
			fileName: config.DefaultAppConfig,
			appCfgFn: config.NewDefaultAppCfgFn(),
			content:  contentFromStruct(t, stableAppCfg("1")),
			expect:   stableAppCfg("1"),
		},
		{
			name:     "Empty config",
			fileName: config.DefaultAppConfig,
			appCfgFn: func() *config.AppConfig {
				return stableAppCfg("1")
			},
			expect: stableAppCfg("1"),
		},
		{
			name:     "Envvar override",
			fileName: config.DefaultAppConfig,
			preFn: func() {
				err := os.Setenv(fmt.Sprintf("%s_USER_ID", config.DefaultAppEnvPrefix), "2")
				assert.NoError(t, err)
			},
			appCfgFn: func() *config.AppConfig {
				return stableAppCfg("1")
			},
			expect: stableAppCfg("2"),
		},
		{
			name:     "No configuration file",
			fileName: "nope.yml",
			appCfgFn: func() *config.AppConfig {
				return stableAppCfg("1")
			},
			notFoundFn:  config.NewDefaultNotFoundFn(),
			expectError: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			os.Clearenv()

			if tc.preFn != nil {
				tc.preFn()
			}

			got, err := config.LoadApp(createAppTestConfig(t, tc.content, tc.fileName), tc.appCfgFn, tc.notFoundFn)
			if tc.expectError {
				if tc.expect == nil {
					assert.NotNil(t, err)
				} else {
					assert.Equal(t, tc.expect, err.Error())
				}
			} else {
				if got != nil {
					got.BaseDir = "" // Dont like this
				}
				assert.NoError(t, err)
				assert.Equal(t, tc.expect, got)
			}
		})
	}
}
