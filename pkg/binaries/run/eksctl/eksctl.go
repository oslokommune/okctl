// Package eksctl knows how to run the eksctl cli
package eksctl

import (
	"fmt"
	"io"

	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/api/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/binaries/run"
	"github.com/oslokommune/okctl/pkg/credentials/aws"
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
	Store      storage.StoreCleaner
	Auth       aws.Authenticator
	CmdFn      run.CmdFn
}

// New returns a new wrapper around the eksctl cli
func New(store storage.StoreCleaner, progress io.Writer, binaryPath string, auth aws.Authenticator, fn run.CmdFn) *Eksctl {
	return &Eksctl{
		Progress:   progress,
		BinaryPath: binaryPath,
		Store:      store,
		Auth:       auth,
		CmdFn:      fn,
	}
}

func (e *Eksctl) runner() (run.Runner, error) {
	envs, err := e.Auth.AsEnv()
	if err != nil {
		return nil, err
	}

	return run.New(e.Store.Path(), e.BinaryPath, envs, e.CmdFn), nil
}

func (e *Eksctl) writeClusterConfig(cfg *v1alpha1.ClusterConfig) error {
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

func (e *Eksctl) run(args []string, cfg *v1alpha1.ClusterConfig) ([]byte, error) {
	var err error

	defer func() {
		err = e.Store.Clean()
	}()

	err = e.writeClusterConfig(cfg)
	if err != nil {
		return nil, err
	}

	runner, err := e.runner()
	if err != nil {
		return nil, err
	}

	return runner.Run(e.Progress, append(args, "--config-file", e.Store.Abs(defaultClusterConfig)))
}

// DeleteCluster invokes eksctl delete cluster using the provided
// cluster configuration as input
func (e *Eksctl) DeleteCluster(cfg *v1alpha1.ClusterConfig) ([]byte, error) {
	args := []string{
		"delete",
		"cluster",
	}

	b, err := e.run(args, cfg)
	if err != nil {
		return nil, errors.E(err, fmt.Sprintf("failed to delete: %s", string(b)), errors.IO)
	}

	return b, nil
}

// CreateCluster invokes eksctl create cluster using the provided
// cluster configuration as input
func (e *Eksctl) CreateCluster(cfg *v1alpha1.ClusterConfig) ([]byte, error) {
	args := []string{
		"create",
		"cluster",
		"--write-kubeconfig=false",
	}

	b, err := e.run(args, cfg)
	if err != nil {
		return nil, errors.E(err, fmt.Sprintf("failed to create: %s", string(b)), errors.IO)
	}

	return b, nil
}
