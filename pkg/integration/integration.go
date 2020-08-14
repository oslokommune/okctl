package integration

import (
	"context"
	"fmt"
	"math/rand"
	"path"
	"strings"
	"time"

	"github.com/rancher/k3d/v3/cmd/util"
	k3dCluster "github.com/rancher/k3d/v3/pkg/cluster"
	"github.com/rancher/k3d/v3/pkg/runtimes"
	k3d "github.com/rancher/k3d/v3/pkg/types"
	"github.com/spf13/afero"
)

// KubernetesCluster contains all state for creating a kubernetes
// cluster
type KubernetesCluster struct {
	cluster        *k3d.Cluster
	kubeConfigDir  string
	kubeConfigPath string

	ctx            context.Context
	fs *afero.Afero
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
		Name: RandStringBytes(10),
		Network: k3d.ClusterNetwork{
			Name:     RandStringBytes(10),
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
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}

	return string(b)
}
