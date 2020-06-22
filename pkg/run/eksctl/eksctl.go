// Package eksctl knows how to run the eksctl cli
package eksctl

import (
	"io"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/binaries"
	"github.com/oslokommune/okctl/pkg/credentials"
	"github.com/oslokommune/okctl/pkg/run"
	"github.com/oslokommune/okctl/pkg/storage"
	"github.com/sirupsen/logrus"
)

const (
	// Name sets the name of the binary/cli
	Name = "eksctl"
	// Version sets the currently used version of the binary/cli
	Version = "0.21.0"
)

// Eksctl stores state for working with the eksctl cli
type Eksctl struct {
	BinaryPath  string
	WorkingDir  string
	Credentials credentials.Provider
	Runner      run.Runner
	Store       storage.StoreCleaner
	Logger      *logrus.Logger
}

// New returns a new wrapper around the eksctl cli
func New(logger *logrus.Logger, credentials credentials.Provider, binaries binaries.Provider) (*Eksctl, error) {
	binaryPath, err := binaries.Fetch(Name, Version)
	if err != nil {
		return nil, err
	}

	envs, err := credentials.AsEnv()
	if err != nil {
		return nil, err
	}

	store, err := storage.NewTemporaryStorage()
	if err != nil {
		return nil, err
	}

	return &Eksctl{
		BinaryPath:  binaryPath,
		Credentials: credentials,
		Runner:      run.New(logger, store.Path, binaryPath, envs),
		Store:       store,
		Logger:      logger,
	}, nil
}

func (e *Eksctl) writeTemporaryClusterConfig(cfg *v1alpha1.ClusterConfig) error {
	data, err := cfg.YAML()
	if err != nil {
		return err
	}

	file, err := e.Store.Create("", "cluster-config.yml", 0644)
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

// DeleteCluster deletes an EKS cluster using eksctl
func (e *Eksctl) DeleteCluster(progress io.Writer, cfg *v1alpha1.ClusterConfig) error {
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
		e.Store.Abs("cluster-config.yml"),
	}

	_, err = e.Runner.Run(progress, args)
	if err != nil {
		return err
	}

	return nil
}

// CreateCluster creates an EKS cluster using eksctl
func (e *Eksctl) CreateCluster(progress io.Writer, cfg *v1alpha1.ClusterConfig) error {
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
		e.Store.Abs("cluster-config.yml"),
	}

	_, err = e.Runner.Run(progress, args)
	if err != nil {
		return err
	}

	return nil
}
