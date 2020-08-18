// Package integration implements functionality for more easily running integration tests
package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	"github.com/rancher/k3d/v3/cmd/util"
	k3dCluster "github.com/rancher/k3d/v3/pkg/cluster"
	"github.com/rancher/k3d/v3/pkg/runtimes"
	k3d "github.com/rancher/k3d/v3/pkg/types"
	"github.com/spf13/afero"
)

// Localstack contains all state for managing the lifecycle
// of the localstack instance
type Localstack struct {
	pool     *dockertest.Pool
	edgePort int
	resource *dockertest.Resource
}

// NewLocalstack returns initialised localstack
// state
func NewLocalstack() *Localstack {
	return &Localstack{}
}

// Create a localstack instance
func (l *Localstack) Create(timeout time.Duration) error {
	pool, err := dockertest.NewPool("")
	if err != nil {
		return fmt.Errorf("couldn't connect to docker: %w", err)
	}

	l.pool = pool

	port, err := util.GetFreePort()
	if err != nil {
		return fmt.Errorf("failed to find available port for edge: %w", err)
	}

	l.edgePort = port

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "localstack/localstack",
		Tag:        "0.11.4",
		Env: []string{
			"EDGE_PORT=4566",
			"DEFAULT_REGION=eu-west-1",
			"LAMBDA_REMOTE_DOCKER=0",
			"START_WEB=0",
			"DEBUG=1",
		},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"4566/tcp": {
				{
					HostIP:   "0.0.0.0",
					HostPort: fmt.Sprintf("%d", l.edgePort),
				},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to start localstack container: %w", err)
	}

	err = pool.Retry(l.Health)
	if err != nil {
		return fmt.Errorf("failed to wait for localstack: %w", err)
	}

	l.resource = resource

	return nil
}

// AWSSession returns an AWS session that works with the
// localstack instance
func (l *Localstack) AWSSession() *session.Session {
	return session.Must(session.NewSession(&aws.Config{
		Region:   aws.String("eu-west-1"),
		Endpoint: aws.String(fmt.Sprintf("http://localhost:%d", l.edgePort)),
		//DisableSSL:  aws.Bool(true),
		Credentials: credentials.NewStaticCredentials("fake", "fake", "fake"),
	}))
}

// Logs retrieves the logs from the localstack instance
func (l *Localstack) Logs() (string, error) {
	var b bytes.Buffer

	err := l.pool.Client.Logs(docker.LogsOptions{
		Container:         l.resource.Container.ID,
		OutputStream:      &b,
		ErrorStream:       &b,
		InactivityTimeout: 0,
		Stdout:            true,
		Stderr:            true,
	})

	return b.String(), err
}

// Cleanup removes all resources created by localstack
func (l *Localstack) Cleanup() error {
	if l.pool != nil {
		err := l.pool.Purge(l.resource)
		if err != nil {
			return fmt.Errorf("failed to cleanup resources: %w", err)
		}
	}

	return nil
}

// Health returns nil if the localstack instance is up and running
func (l *Localstack) Health() error {
	resp, err := http.Get(fmt.Sprintf("http://localhost:%d/health?reload", l.edgePort))
	if err != nil {
		return err
	}

	var services map[string]map[string]string

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read body: %w", err)
	}

	err = json.Unmarshal(body, &services)
	if err != nil {
		return fmt.Errorf("failed to unmarshal json: %w", err)
	}

	for service, status := range services["services"] {
		if status != "running" {
			return fmt.Errorf("waiting for: %s, to get to running state, currently: %s", service, status)
		}
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("got response code from localstack: %d, not 200 OK", resp.StatusCode)
	}

	return nil
}

// KubernetesCluster contains all state for creating a kubernetes
// cluster
type KubernetesCluster struct {
	cluster        *k3d.Cluster
	kubeConfigDir  string
	kubeConfigPath string

	ctx context.Context
	fs  *afero.Afero
}

// NewKubernetesCluster returns an initialised struct for managing kubernetes
// clusters
func NewKubernetesCluster() *KubernetesCluster {
	return &KubernetesCluster{
		ctx: context.Background(),
		fs: &afero.Afero{
			Fs: afero.NewOsFs(),
		},
	}
}

// Create a kubernetes cluster
func (k *KubernetesCluster) Create(timeout time.Duration) error {
	createClusterOpts := &k3d.ClusterCreateOpts{
		WaitForServer: true,
		Timeout:       timeout,
	}

	serverNode := &k3d.Node{
		Role: k3d.ServerRole,
		// Find tags at: https://hub.docker.com/r/rancher/k3s/tags
		Image: "rancher/k3s:v1.18.8-rc1-k3s1",
		Args:  createClusterOpts.K3sServerArgs,
		ServerOpts: k3d.ServerOpts{
			IsInit: true,
		},
	}

	port, err := util.GetFreePort()
	if err != nil {
		return fmt.Errorf("failed to find free port: %w", err)
	}

	cluster := &k3d.Cluster{
		Name: RandStringBytes(10), //nolint: gomnd
		Network: k3d.ClusterNetwork{
			Name:     RandStringBytes(10), //nolint: gomnd
			External: true,
		},
		Nodes: []*k3d.Node{
			serverNode,
		},
		InitNode:          serverNode,
		ExternalDatastore: k3d.ExternalDatastore{},
		CreateClusterOpts: createClusterOpts,
		ExposeAPI: k3d.ExposeAPI{
			Host:   k3d.DefaultAPIHost,
			HostIP: k3d.DefaultAPIHost,
			Port:   fmt.Sprintf("%d", port),
		},
		ServerLoadBalancer: &k3d.Node{
			Role: k3d.LoadBalancerRole,
		},
	}

	err = k3dCluster.ClusterCreate(k.ctx, runtimes.SelectedRuntime, cluster)
	if err != nil {
		return fmt.Errorf("failed to create cluster: %w", err)
	}

	k.cluster = cluster

	return nil
}

// Destroy removes the cluster
func (k *KubernetesCluster) Destroy() error {
	err := k3dCluster.ClusterDelete(k.ctx, runtimes.SelectedRuntime, k.cluster)
	if err != nil {
		return fmt.Errorf("failed to destroy cluster: %w", err)
	}

	return nil
}

// KubeConfig writes the kubeconfig to a temporary file and
// returns the path to that file
func (k *KubernetesCluster) KubeConfig() (string, error) {
	if len(k.kubeConfigPath) == 0 {
		dir, err := k.fs.TempDir("", "kubeconfig")
		if err != nil {
			return "", fmt.Errorf("failed to create temporary directory for kubeconfig: %w", err)
		}

		k.kubeConfigDir = dir
		k.kubeConfigPath = path.Join(dir, "kubeconfig")
	}

	kubeConfOpts := &k3dCluster.WriteKubeConfigOptions{
		OverwriteExisting:    true,
		UpdateCurrentContext: true,
	}

	_, err := k3dCluster.KubeconfigGetWrite(k.ctx, runtimes.SelectedRuntime, k.cluster, k.kubeConfigPath, kubeConfOpts)
	if err != nil {
		return "", fmt.Errorf("failed to create kubeconfig: %w", err)
	}

	return k.kubeConfigPath, nil
}

// Cleanup removes all created resources
func (k *KubernetesCluster) Cleanup() error {
	var errors []string

	err := k.fs.RemoveAll(k.kubeConfigDir)
	if err != nil {
		errors = append(errors, fmt.Errorf("failed to cleanup kubeconfig dir: %w", err).Error())
	}

	err = k.Destroy()
	if err != nil {
		errors = append(errors, fmt.Errorf("failed to cleanup cluster: %w", err).Error())
	}

	if len(errors) > 0 {
		return fmt.Errorf("%s", strings.Join(errors, ", "))
	}

	return nil
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// RandStringBytes returns a random string
// https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go/22892986#22892986
func RandStringBytes(n int) string {
	rand.Seed(time.Now().UnixNano())

	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))] //nolint: gosec
	}

	return string(b)
}
