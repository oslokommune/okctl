package helm_test

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"testing"
	"time"

	"github.com/oslokommune/okctl/pkg/config"
	"github.com/oslokommune/okctl/pkg/helm"
	"github.com/oslokommune/okctl/pkg/integration"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestEstablishEnv(t *testing.T) {
	testCases := []struct {
		name   string
		envs   map[string]string
		expect []string
	}{
		{
			name: "Should work",
			envs: map[string]string{
				"SOMETHING": "cool",
			},
			expect: []string{
				"SOMETHING=cool",
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			before := os.Environ()

			restoreFn, err := helm.EstablishEnv(tc.envs)
			assert.NoError(t, err)
			assert.Equal(t, tc.expect, os.Environ())

			err = restoreFn()
			assert.NoError(t, err)
			assert.Equal(t, before, os.Environ())
		})
	}
}

// nolint: funlen
func TestHelm(t *testing.T) {
	dir, err := ioutil.TempDir("", "testHelm")
	assert.NoError(t, err)

	defer func() {
		_ = os.RemoveAll(dir)
	}()

	cluster := integration.NewKubernetesCluster()

	err = cluster.Create(2 * time.Minute)

	defer func() {
		err := cluster.Cleanup()
		assert.NoError(t, err)
	}()

	assert.NoError(t, err)

	kubeConfPath, err := cluster.KubeConfig()
	assert.NoError(t, err)

	testCases := []struct {
		name      string
		helm      *helm.Helm
		expect    interface{}
		expectErr bool
	}{
		{
			name: "Should work",
			helm: helm.New(&helm.Config{
				Namespace:            "test-helm",
				KubeConfig:           kubeConfPath,
				HomeDir:              dir,
				HelmPluginsDirectory: path.Join(dir, config.DefaultHelmBaseDir, config.DefaultHelmPluginsDirectory),
				HelmRegistryConfig:   path.Join(dir, config.DefaultHelmBaseDir, config.DefaultHelmRegistryConfig),
				HelmRepositoryConfig: path.Join(dir, config.DefaultHelmBaseDir, config.DefaultHelmRepositoryConfig),
				HelmRepositoryCache:  path.Join(dir, config.DefaultHelmBaseDir, config.DefaultHelmRepositoryCache),
				HelmBaseDir:          path.Join(dir, config.DefaultHelmBaseDir),
				Debug:                true,
				DebugOutput:          os.Stderr,
			}, &afero.Afero{
				Fs: afero.NewOsFs(),
			}),
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			err := tc.helm.RepoAdd("stable", "https://kubernetes-charts.storage.googleapis.com")
			assert.NoError(t, err)

			err = tc.helm.RepoUpdate()
			assert.NoError(t, err)

			mysql := helm.Mysql(&helm.MysqlValues{
				MysqlRootPassword: "admin@123",
				Persistence: helm.MysqlPersistence{
					Enabled: false,
				},
				ImagePullPolicy: "Always",
			})

			cfg, err := mysql.InstallConfig()
			assert.NoError(t, err)

			release, err := tc.helm.Install(cfg)
			assert.NoError(t, err)
			log.Printf("Released: %s, to namespace: %s", release.Name, release.Namespace)
		})
	}
}
