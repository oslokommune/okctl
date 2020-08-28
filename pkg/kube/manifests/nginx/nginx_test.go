package nginx_test

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"

	"github.com/oslokommune/okctl/pkg/integration"
	"github.com/oslokommune/okctl/pkg/kube"
	"github.com/oslokommune/okctl/pkg/kube/manifests/nginx"
	"github.com/stretchr/testify/assert"
)

func TestNginx(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping nginx kube integration tests in short mode")
	}

	dir, err := ioutil.TempDir("", "nginx")
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
		ext  *nginx.Nginx
	}{
		{
			name: "Should work",
			ext:  nginx.New("default"),
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			k, err := kube.New(kubeConfPath)
			assert.NoError(t, err)

			resources, err := k.Apply(tc.ext.CreateDeployment)
			assert.NoError(t, err)

			err = k.Watch(resources, 2*time.Minute)
			assert.NoError(t, err)

			outputs, err := cluster.Debug("default")
			assert.NoError(t, err)

			time.Sleep(30 * time.Second)

			for title, data := range outputs {
				log.Println(title)
				log.Println(data)
			}
		})
	}
}
