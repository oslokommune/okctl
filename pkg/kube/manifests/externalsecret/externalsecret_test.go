package externalsecret_test

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"

	"github.com/oslokommune/okctl/pkg/kube/manifests/externalsecret"

	"github.com/oslokommune/okctl/pkg/integration"
	"github.com/oslokommune/okctl/pkg/kube"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
	"sigs.k8s.io/yaml"
)

func TestNew(t *testing.T) {
	testCases := []struct {
		name string
		ext  *externalsecret.ExternalSecret
	}{
		{
			name: "Should work",
			ext: externalsecret.New("name", "default", nil, nil, map[string]string{
				"from": "/path",
			}),
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			g := goldie.New(t)

			dm := tc.ext.SecretManifest()
			db, err := yaml.Marshal(dm)
			assert.NoError(t, err)
			g.Assert(t, "secret.yaml", db)
		})
	}
}

// nolint: funlen
func TestExternalDNS(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping external dns kube integration tests in short mode")
	}

	dir, err := ioutil.TempDir("", "externalDNS")
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
		name string
		ext  *externalsecret.ExternalSecret
	}{
		// This will not actually work unless we also add the CRD for
		// external secrets:
		// - https://github.com/external-secrets/kubernetes-external-secrets/blob/master/charts/kubernetes-external-secrets/crds/kubernetes-client.io_externalsecrets_crd.yaml
		{
			name: "Should work",
			ext: externalsecret.New("name", "default", nil, nil, map[string]string{
				"from": "/path",
			}),
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			// Skipping for now, since this won't work
			t.Skip()

			k, err := kube.New(kube.NewFromKubeConfig(kubeConfPath))
			assert.NoError(t, err)

			// Need to add k.Watch to the output her, but that means bringing up localstack, I think?
			_, err = k.Apply(kube.Applier{Fn: tc.ext.CreateSecret})
			assert.NoError(t, err)

			outputs, err := cluster.Debug("kube-system")
			assert.NoError(t, err)

			for title, data := range outputs {
				log.Println(title)
				log.Println(data)
			}
		})
	}
}
