package eksctl

import (
	"io"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/binaries"
	"github.com/oslokommune/okctl/pkg/credentials"
	"github.com/oslokommune/okctl/pkg/run"
	"github.com/oslokommune/okctl/pkg/storage"
)

const (
	Name    = "eksctl"
	Version = "0.18.0"
)

type Eksctl struct {
	BinaryPath  string
	WorkingDir  string
	Credentials credentials.Provider
	Runner      run.Runner
	Store       storage.StoreCleaner
}

func New(credentials credentials.Provider, binaries binaries.Provider) (*Eksctl, error) {
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
		Runner:      run.New(store.Path, binaryPath, envs),
		Store:       store,
	}, nil
}

func (e *Eksctl) CreateCluster(progress io.Writer, cfg *v1alpha1.ClusterConfig) error {
	var err error

	defer func() {
		err = e.Store.Clean()
	}()

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

	args := []string{
		"create",
		"cluster",
		"--config-file",
		file.Name(),
	}

	_, err = e.Runner.Run(progress, args)
	if err != nil {
		return err
	}

	return nil
}
