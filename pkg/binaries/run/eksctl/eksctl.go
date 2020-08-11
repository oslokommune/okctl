// Package eksctl knows how to run the eksctl cli
package eksctl

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/binaries/run"
	"github.com/oslokommune/okctl/pkg/credentials/aws"
	"github.com/oslokommune/okctl/pkg/storage"
)

const (
	// Name sets the name of the binary/cli
	Name = "eksctl"
	// Version sets the currently used version of the binary/cli
	Version = "0.25.0"

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
	CmdPath    []string
	CmdEnv     []string
	DoDebug    bool
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

	if len(e.CmdPath) > 0 {
		envs = append(envs, fmt.Sprintf("PATH=$PATH:%s", strings.Join(e.CmdPath, ":")))
	}

	if len(e.CmdEnv) > 0 {
		envs = append(envs, e.CmdEnv...)
	}

	return run.New(e.Store.Path(), e.BinaryPath, envs, e.CmdFn), nil
}

func (e *Eksctl) writeClusterConfig(cfg *api.ClusterConfig) error {
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

func (e *Eksctl) run(args []string, cfg *api.ClusterConfig) ([]byte, error) {
	var err error

	defer func() {
		err = e.Store.Clean()
	}()

	verbosity := "--verbose=3"
	if e.DoDebug {
		verbosity = "--verbose=4"
	}

	args = append(args, verbosity)

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

// AddToPath strips the base of the binary and adds it to the PATH
func (e *Eksctl) AddToPath(binaryPaths ...string) {
	for _, binaryPath := range binaryPaths {
		e.CmdPath = append(e.CmdPath, filepath.Dir(binaryPath))
	}
}

// AddToEnv adds additional environment variables to the command
func (e *Eksctl) AddToEnv(envs ...string) {
	e.CmdEnv = append(e.CmdEnv, envs...)
}

// Debug sets whether we should increase log output from eksctl,
// the default behavior is off
func (e *Eksctl) Debug(enable bool) {
	e.DoDebug = enable
}

// DeleteCluster invokes eksctl delete cluster using the provided
// cluster configuration as input
func (e *Eksctl) DeleteCluster(cfg *api.ClusterConfig) ([]byte, error) {
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
func (e *Eksctl) CreateCluster(kubeConfigPath string, cfg *api.ClusterConfig) ([]byte, error) {
	args := []string{
		"create",
		"cluster",
		"--write-kubeconfig=true",
		fmt.Sprintf("--kubeconfig=%s", kubeConfigPath),
	}

	b, err := e.run(args, cfg)
	if err != nil {
		return nil, errors.E(err, fmt.Sprintf("failed to create: %s", string(b)), errors.IO)
	}

	return b, nil
}
