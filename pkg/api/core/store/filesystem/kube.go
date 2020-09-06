package filesystem

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/client/store"

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

	manifests := map[string][]byte{}

	callback := func(name string, data []byte) error {
		manifests[name] = data
		return nil
	}

	process := func(data interface{}, operations store.Operations) error {
		if d, ok := data.(*Kube); ok {
			for _, name := range d.Manifests {
				operations.GetBytes(name, callback)
			}

			return nil
		}

		return fmt.Errorf("failed to cast to *Kube")
	}

	_, err := store.NewFileSystem(paths.BaseDir, k.fs).
		GetStruct(paths.OutputFile, &o, store.FromJSON()).
		ProcessGetStruct(paths.OutputFile, process).
		Do()
	if err != nil {
		return nil, fmt.Errorf("faileed to get kube: %w", err)
	}

	return &api.Kube{
		ID:           o.ID,
		HostedZoneID: o.HostedZoneID,
		DomainFilter: o.DomainFilter,
		Manifests:    manifests,
	}, nil
}

func (k *kubeStore) saveKubeDeployment(paths Paths, kube *api.Kube) error {
	o := Kube{
		ID:           kube.ID,
		HostedZoneID: kube.HostedZoneID,
		DomainFilter: kube.DomainFilter,
	}

	manifests := make([]store.AddStoreBytes, len(kube.Manifests))

	for name, data := range kube.Manifests {
		manifests = append(manifests, store.AddStoreBytes{
			Name: name,
			Data: data,
		})

		o.Manifests = append(o.Manifests, name)
	}

	_, err := store.NewFileSystem(paths.BaseDir, k.fs).
		StoreStruct(paths.OutputFile, &o, store.ToJSON()).
		AddStoreBytes(manifests...).
		Do()
	if err != nil {
		return fmt.Errorf("failed to write kube: %w", err)
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
