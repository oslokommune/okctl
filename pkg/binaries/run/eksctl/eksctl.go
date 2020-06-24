// Package eksctl knows how to run the eksctl cli
package eksctl

import (
	"io"

	"github.com/oslokommune/okctl/pkg/api/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/binaries/run"
	"github.com/oslokommune/okctl/pkg/storage"
)

const (
	// Name sets the name of the binary/cli
	Name = "eksctl"
	// Version sets the currently used version of the binary/cli
	Version = "0.21.0"

	defaultClusterConfig = "cluster-config.yml"
)

// Eksctl stores state for working with the eksctl cli
type Eksctl struct {
	Progress   io.Writer
	BinaryPath string
	WorkingDir string
	Runner     run.Runner
	Store      storage.StoreCleaner
}

// New returns a new wrapper around the eksctl cli
func New(progress io.Writer, binaryPath string, envs []string) (*Eksctl, error) {
	store, err := storage.NewTemporaryStorage()
	if err != nil {
		return nil, err
	}

	return &Eksctl{
		Progress:   progress,
		BinaryPath: binaryPath,
		Runner:     run.New(store.Path, binaryPath, envs),
		Store:      store,
	}, nil
}

func (e *Eksctl) writeTemporaryClusterConfig(cfg *v1alpha1.ClusterConfig) error {
	data, err := cfg.YAML()
	if err != nil {
		return err
	}

	file, err := e.Store.Create("", defaultClusterConfig, 0644)
	if err != nil {
		return err
	}

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	err = file.Close()
	if err != nil {
		return err
	}

	return nil
}

// DeleteCluster invokes eksctl delete cluster using the provided
// cluster configuration as input
func (e *Eksctl) DeleteCluster(cfg *v1alpha1.ClusterConfig) error {
	var err error

	defer func() {
		err = e.Store.Clean()
	}()

	err = e.writeTemporaryClusterConfig(cfg)
	if err != nil {
		return err
	}

	args := []string{
		"delete",
		"cluster",
		"--config-file",
		e.Store.Abs(defaultClusterConfig),
	}

	_, err = e.Runner.Run(e.Progress, args)
	if err != nil {
		return err
	}

	return nil
}

// CreateCluster invokes eksctl create cluster using the provided
// cluster configuration as input
func (e *Eksctl) CreateCluster(cfg *v1alpha1.ClusterConfig) error {
	var err error

	defer func() {
		err = e.Store.Clean()
	}()

	err = e.writeTemporaryClusterConfig(cfg)
	if err != nil {
		return err
	}

	args := []string{
		"create",
		"cluster",
		"--write-kubeconfig=false",
		"--config-file",
		e.Store.Abs(defaultClusterConfig),
	}

	_, err = e.Runner.Run(e.Progress, args)
	if err != nil {
		return err
	}

	return nil
}
