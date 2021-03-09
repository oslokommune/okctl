package helm_test

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/oslokommune/okctl/pkg/config/constant"

	"github.com/oslokommune/okctl/pkg/helm/charts/argocd"

	"github.com/oslokommune/okctl/pkg/credentials/aws"
	"github.com/oslokommune/okctl/pkg/helm"
	"github.com/oslokommune/okctl/pkg/integration"
	"github.com/oslokommune/okctl/pkg/mock"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestEstablishEnv(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skip("Skipping testing in CI environment, as this test is flaky here.")
	}

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
			after := os.Environ()
			assert.NoError(t, err)
			sort.Strings(before)
			sort.Strings(after)
			assert.Equal(t, before, after)
		})
	}
}

// nolint: funlen
func TestHelm(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skip("Skipping testing in CI environment, as this test is flaky here.")
	}

	if testing.Short() {
		t.Skip("skipping helm integration tests in short mode")
	}

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
		chart     *helm.Chart
		expect    interface{}
		expectErr bool
	}{
		{
			name: "Mysql should work",
			helm: helm.New(
				&helm.Config{
					HomeDir:              dir,
					HelmPluginsDirectory: path.Join(dir, constant.DefaultHelmBaseDir, constant.DefaultHelmPluginsDirectory),
					HelmRegistryConfig:   path.Join(dir, constant.DefaultHelmBaseDir, constant.DefaultHelmRegistryConfig),
					HelmRepositoryConfig: path.Join(dir, constant.DefaultHelmBaseDir, constant.DefaultHelmRepositoryConfig),
					HelmRepositoryCache:  path.Join(dir, constant.DefaultHelmBaseDir, constant.DefaultHelmRepositoryCache),
					HelmBaseDir:          path.Join(dir, constant.DefaultHelmBaseDir),
					Debug:                true,
					DebugOutput:          os.Stderr,
				},
				aws.New(aws.NewInMemoryStorage(), aws.NewAuthStatic(mock.DefaultValidCredentials())),
				&afero.Afero{
					Fs: afero.NewOsFs(),
				},
			),
			chart: helm.Mysql(&helm.MysqlValues{
				MysqlRootPassword: "admin@123",
				Persistence: helm.MysqlPersistence{
					Enabled: false,
				},
				ImagePullPolicy: "Always",
			}),
			expect:    nil,
			expectErr: false,
		},
		{
			name: "argocd should work",
			helm: helm.New(
				&helm.Config{
					HomeDir:              dir,
					HelmPluginsDirectory: path.Join(dir, constant.DefaultHelmBaseDir, constant.DefaultHelmPluginsDirectory),
					HelmRegistryConfig:   path.Join(dir, constant.DefaultHelmBaseDir, constant.DefaultHelmRegistryConfig),
					HelmRepositoryConfig: path.Join(dir, constant.DefaultHelmBaseDir, constant.DefaultHelmRepositoryConfig),
					HelmRepositoryCache:  path.Join(dir, constant.DefaultHelmBaseDir, constant.DefaultHelmRepositoryCache),
					HelmBaseDir:          path.Join(dir, constant.DefaultHelmBaseDir),
					Debug:                true,
					DebugOutput:          os.Stderr,
				},
				aws.New(aws.NewInMemoryStorage(), aws.NewAuthStatic(mock.DefaultValidCredentials())),
				&afero.Afero{
					Fs: afero.NewOsFs(),
				},
			),
			chart: argocd.New(argocd.NewDefaultValues(argocd.ValuesOpts{
				URL:                  "https://argocd.test.oslo.systems",
				Region:               "eu-west-1",
				HostName:             "argocd.test.oslo.systems",
				CertificateARN:       "arn:aws:acm:eu-west-1:123456789012/certificate/12345abc",
				ClientID:             "client12345",
				Organisation:         "oslokommune",
				AuthDomain:           "auth.oslo.systems",
				UserPoolID:           "GYUE78FAKE",
				RepoURL:              "git@github.com:oslokommune/test.git",
				RepoName:             "test",
				PrivateKeySecretName: "argocd-test-oslokommune-private-key",
				PrivateKeySecretKey:  "ssh-private-key",
			})),
			expect:    nil,
			expectErr: false,
		},
		//  We need to bring up localstack for these tests to pass, something to be done..
		//  {
		// 	name: "aws-alb-ingress-controller should work",
		// 	helm: helm.New(&helm.Config{
		// 		HomeDir:              dir,
		// 		HelmPluginsDirectory: path.Join(dir, config.DefaultHelmBaseDir, config.DefaultHelmPluginsDirectory),
		// 		HelmRegistryConfig:   path.Join(dir, config.DefaultHelmBaseDir, config.DefaultHelmRegistryConfig),
		// 		HelmRepositoryConfig: path.Join(dir, config.DefaultHelmBaseDir, config.DefaultHelmRepositoryConfig),
		// 		HelmRepositoryCache:  path.Join(dir, config.DefaultHelmBaseDir, config.DefaultHelmRepositoryCache),
		// 		HelmBaseDir:          path.Join(dir, config.DefaultHelmBaseDir),
		// 		Debug:                true,
		// 		DebugOutput:          os.Stderr,
		// 	},
		// 		aws.New(aws.NewInMemoryStorage(), aws.NewAuthStatic(mock.DefaultValidCredentials())),
		// 		&afero.Afero{
		// 			Fs: afero.NewOsFs(),
		// 		}),
		// 	chart: awsalbingresscontroller.New(
		// 		awsalbingresscontroller.NewDefaultValues("my-cluster", "my-service-account"),
		// 	),
		// 	expect:    nil,
		// 	expectErr: false,
		// },
		// {
		// 	name: "external-secrets should work",
		// 	helm: helm.New(&helm.Config{
		// 		HomeDir:              dir,
		// 		HelmPluginsDirectory: path.Join(dir, config.DefaultHelmBaseDir, config.DefaultHelmPluginsDirectory),
		// 		HelmRegistryConfig:   path.Join(dir, config.DefaultHelmBaseDir, config.DefaultHelmRegistryConfig),
		// 		HelmRepositoryConfig: path.Join(dir, config.DefaultHelmBaseDir, config.DefaultHelmRepositoryConfig),
		// 		HelmRepositoryCache:  path.Join(dir, config.DefaultHelmBaseDir, config.DefaultHelmRepositoryCache),
		// 		HelmBaseDir:          path.Join(dir, config.DefaultHelmBaseDir),
		// 		Debug:                true,
		// 		DebugOutput:          os.Stderr,
		// 	}, &afero.Afero{
		// 		Fs: afero.NewOsFs(),
		// 	}),
		// 	chart:     helm.ExternalSecrets(helm.DefaultExternalSecretsValues()),
		// 	expect:    nil,
		// 	expectErr: false,
		// },
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			err := tc.helm.RepoAdd(tc.chart.RepositoryName, tc.chart.RepositoryURL)
			assert.NoError(t, err)

			err = tc.helm.RepoUpdate()
			assert.NoError(t, err)

			cfg, err := tc.chart.InstallConfig()
			assert.NoError(t, err)

			release, err := tc.helm.Install(kubeConfPath, cfg)
			assert.NoError(t, err)
			assert.NotNil(t, release)

			if err != nil {
				items, err := cluster.Debug(tc.chart.Namespace)
				assert.NoError(t, err)
				for title, item := range items {
					log.Printf("debug information for: %s\n", title)
					log.Println(strings.Join(item, "\n"))
				}
			}

			err = tc.helm.Delete(kubeConfPath, tc.chart.DeleteConfig())
			assert.NoError(t, err)

			if err != nil {
				items, err := cluster.Debug(tc.chart.Namespace)
				assert.NoError(t, err)
				for title, item := range items {
					log.Printf("debug information for: %s\n", title)
					log.Println(strings.Join(item, "\n"))
				}
			}
		})
	}
}
