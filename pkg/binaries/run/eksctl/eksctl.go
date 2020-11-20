// Package eksctl knows how to run the eksctl cli
package eksctl

import (
	"fmt"
	"io"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/oslokommune/okctl/pkg/api/okctl.io/v1alpha1"

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
	Logger     *logrus.Logger
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
func New(logger *logrus.Logger, store storage.StoreCleaner, progress io.Writer, binaryPath string, auth aws.Authenticator, fn run.CmdFn) *Eksctl {
	return &Eksctl{
		Progress:   progress,
		Logger:     logger,
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

	return run.New(e.Logger, e.Store.Path(), e.BinaryPath, envs, e.CmdFn), nil
}

func (e *Eksctl) writeClusterConfig(cfg *v1alpha1.ClusterConfig) error {
	data, err := cfg.YAML()
	if err != nil {
		return err
	}

	file, err := e.Store.Create("", defaultClusterConfig, 0o644)
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

func (e *Eksctl) run(args []string) ([]byte, error) {
	var err error

	verbosity := "--verbose=3"
	if e.DoDebug {
		verbosity = "--verbose=4"
	}

	args = append(args, verbosity)

	runner, err := e.runner()
	if err != nil {
		return nil, err
	}

	return runner.Run(e.Progress, args)
}

func (e *Eksctl) runWithConfig(args []string, cfg *v1alpha1.ClusterConfig) ([]byte, error) {
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

// CreateServiceAccount invokes eksctl create iamserviceaccount using the provided
// configuration file
func (e *Eksctl) CreateServiceAccount(cfg *v1alpha1.ClusterConfig) ([]byte, error) {
	args := []string{
		"create",
		"iamserviceaccount",
		"--override-existing-serviceaccounts",
		"--approve",
	}

	b, err := e.runWithConfig(args, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create service account: %s, because: %w", string(b), err)
	}

	return b, nil
}

// DeleteServiceAccount invokes eksctl delete iamserviceaccount using the provided
// configuration file
func (e *Eksctl) DeleteServiceAccount(cfg *v1alpha1.ClusterConfig) ([]byte, error) {
	args := []string{
		"delete",
		"iamserviceaccount",
		"--approve",
	}

	b, err := e.runWithConfig(args, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to delete service account: %s, because: %w", string(b), err)
	}

	return b, nil
}

// DeleteCluster invokes eksctl delete cluster using the provided
// cluster configuration as input
func (e *Eksctl) DeleteCluster(clusterName string) ([]byte, error) {
	args := []string{
		"delete",
		"cluster",
		fmt.Sprintf("--name=%s", clusterName),
		"--wait",
	}

	b, err := e.run(args)
	if err != nil {
		return nil, fmt.Errorf("failed to delete: %s, because: %w", string(b), err)
	}

	return b, nil
}

// DeleteFargateProfiles invokes eksctl delete fargate profile for cluster
func (e *Eksctl) DeleteFargateProfiles(clusterName, fragetProfileName string) ([]byte, error) {
	args := []string{
		"delete",
		"fargateprofile",
		fmt.Sprintf("--cluster=%s", clusterName),
		fmt.Sprintf("--name=%s", fragetProfileName),
		"--wait",
	}

	b, err := e.run(args)
	if err != nil {
		return nil, fmt.Errorf("failed to delete: %s, because: %w", string(b), err)
	}

	return b, nil
}

// CreateCluster invokes eksctl create cluster using the provided
// cluster configuration as input
func (e *Eksctl) CreateCluster(cfg *v1alpha1.ClusterConfig) ([]byte, error) {
	args := []string{
		"create",
		"cluster",
	}

	b, err := e.runWithConfig(args, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create: %s, because: %w", string(b), err)
	}

	return b, nil
}

// HasCluster invokes eksctl get cluster using the provided
// cluster config as input and returns an error if the cluster
// does not exist.
func (e *Eksctl) HasCluster(cfg *v1alpha1.ClusterConfig) (bool, error) {
	re := regexp.MustCompile(`(.*)ResourceNotFoundException: No cluster found for name(.*)`)

	args := []string{
		"get",
		"cluster",
		"--name",
		cfg.Metadata.Name,
		"--region",
		cfg.Metadata.Region,
	}

	out, err := e.run(args)
	if err != nil {
		if re.Match(out) && strings.Contains(string(out), cfg.Metadata.Name) {
			return false, nil
		}

		return false, fmt.Errorf("failed to get cluster information: %s: %w", string(out), err)
	}

	return true, nil
}
