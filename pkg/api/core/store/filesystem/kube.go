package filesystem

import (
	"encoding/json"
	"fmt"
	"path"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/spf13/afero"
)

type kubeStore struct {
	externalDNS     Paths
	externalSecrets Paths
	fs              *afero.Afero
}

// Kube contains the stored state for a kube deployment
type Kube struct {
	ID           api.ID
	HostedZoneID string
	DomainFilter string
	Manifests    []string
}

// This is not good, we need to refactor
func (k *kubeStore) SaveExternalSecrets(kube *api.Kube) error {
	return k.saveKubeDeployment(k.externalSecrets, kube)
}

func (k *kubeStore) SaveExternalDNSKubeDeployment(kube *api.Kube) error {
	return k.saveKubeDeployment(k.externalDNS, kube)
}

func (k *kubeStore) GetExternalDNSKubeDeployment() (*api.Kube, error) {
	return k.getKubeDeployment(k.externalDNS)
}

func (k *kubeStore) getKubeDeployment(paths Paths) (*api.Kube, error) {
	o := Kube{}

	outputs, err := k.fs.ReadFile(path.Join(paths.BaseDir, paths.OutputFile))
	if err != nil {
		return nil, fmt.Errorf("failed to read outputs: %w", err)
	}

	err = json.Unmarshal(outputs, &o)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal outputs: %w", err)
	}

	manifests := map[string][]byte{}

	for _, file := range o.Manifests {
		d, err := k.fs.ReadFile(path.Join(paths.BaseDir, file))
		if err != nil {
			return nil, fmt.Errorf("failed to read manifest: %w", err)
		}

		manifests[file] = d
	}

	return &api.Kube{
		ID:           o.ID,
		HostedZoneID: o.HostedZoneID,
		DomainFilter: o.DomainFilter,
		Manifests:    manifests,
	}, nil
}

func (k *kubeStore) saveKubeDeployment(paths Paths, kube *api.Kube) error {
	err := k.fs.MkdirAll(paths.BaseDir, 0o744)
	if err != nil {
		return fmt.Errorf("failed to create base directory: %w", err)
	}

	o := Kube{
		ID:           kube.ID,
		HostedZoneID: kube.HostedZoneID,
		DomainFilter: kube.DomainFilter,
	}

	for file, data := range kube.Manifests {
		err = k.fs.WriteFile(path.Join(paths.BaseDir, file), data, 0o644)
		if err != nil {
			return fmt.Errorf("failed to write manifest: %w", err)
		}

		o.Manifests = append(o.Manifests, file)
	}

	outputs, err := json.Marshal(&o)
	if err != nil {
		return fmt.Errorf("failed to serialise outputs: %w", err)
	}

	err = k.fs.WriteFile(path.Join(paths.BaseDir, paths.OutputFile), outputs, 0o644)
	if err != nil {
		return fmt.Errorf("failed to write outputs: %w", err)
	}

	return nil
}

// NewKubeStore returns an initialised kube store
func NewKubeStore(externalDNS, externalSecrets Paths, fs *afero.Afero) api.KubeStore {
	return &kubeStore{
		externalSecrets: externalSecrets,
		externalDNS:     externalDNS,
		fs:              fs,
	}
}
