package filesystem

import (
	"fmt"
	"path"

	"github.com/gosimple/slug"

	"github.com/oslokommune/okctl/pkg/client"

	"github.com/oslokommune/okctl/pkg/client/store"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/spf13/afero"
)

type manifestStore struct {
	configMap      Paths
	externalSecret Paths
	storageClass   Paths
	fs             *afero.Afero
}

// ConfigMap is stored in the outputs file
type ConfigMap struct {
	ID        api.ID
	Name      string
	Namespace string
}

func (s *manifestStore) SaveConfigMap(secret *client.ConfigMap) (*store.Report, error) {
	o := &ConfigMap{
		ID:        secret.ID,
		Name:      secret.Name,
		Namespace: secret.Namespace,
	}

	report, err := store.NewFileSystem(path.Join(
		s.configMap.BaseDir,
		slug.Make(fmt.Sprintf("%s-%s", secret.Name, secret.Namespace)),
	), s.fs).
		StoreStruct(s.configMap.OutputFile, o, store.ToJSON()).
		StoreBytes(s.configMap.ConfigFile, secret.Manifest).
		Do()
	if err != nil {
		return nil, err
	}

	return report, nil
}

func (s *manifestStore) RemoveConfigMap(name, namespace string) (*store.Report, error) {
	return store.NewFileSystem(path.Join(
		s.configMap.BaseDir,
		slug.Make(fmt.Sprintf("%s-%s", name, namespace)),
	), s.fs).
		Remove(s.configMap.OutputFile).
		Remove(s.configMap.ConfigFile).
		RemoveDir("").
		Do()
}

// StorageClass is stored in the outputs file
type StorageClass struct {
	ID   api.ID
	Name string
}

func (s *manifestStore) SaveStorageClass(sc *client.StorageClass) (*store.Report, error) {
	o := &StorageClass{
		ID:   sc.ID,
		Name: sc.Name,
	}

	report, err := store.NewFileSystem(path.Join(s.storageClass.BaseDir, slug.Make(sc.Name)), s.fs).
		StoreStruct(s.storageClass.OutputFile, o, store.ToJSON()).
		StoreBytes(s.storageClass.ConfigFile, sc.Manifest).
		Do()
	if err != nil {
		return nil, fmt.Errorf("storing storage class: %w", err)
	}

	return report, nil
}

// ExternalSecret is stored in the outputs file
type ExternalSecret struct {
	ID        api.ID
	Manifests []string
}

func (s *manifestStore) RemoveExternalSecret(_ map[string]string) (*store.Report, error) {
	// There is something fishy going on here, though I believe this will mostly be correct
	// since we will rarely create secrets without setting a directory, this is not
	// very pretty or intuitive..
	return store.NewFileSystem(s.externalSecret.BaseDir, s.fs).
		RemoveDir("").
		Do()
}

func (s *manifestStore) SaveExternalSecret(e *client.ExternalSecret) (*store.Report, error) {
	o := &ExternalSecret{
		ID:        e.ID,
		Manifests: make([]string, len(e.Manifests)),
	}

	manifests := make([]store.AddStoreBytes, len(e.Manifests))

	i := 0

	for name, data := range e.Manifests {
		manifests[i] = store.AddStoreBytes{
			Name: name,
			Data: data,
		}

		o.Manifests[i] = name

		i++
	}

	report, err := store.NewFileSystem(s.externalSecret.BaseDir, s.fs).
		StoreStruct(s.externalSecret.OutputFile, &o, store.ToJSON()).
		AddStoreBytes(manifests...).
		Do()
	if err != nil {
		return nil, fmt.Errorf("writing external secret: %w", err)
	}

	return report, nil
}

// NewManifestStore returns an initialised store
func NewManifestStore(configMap, storageClass, externalSecret Paths, fs *afero.Afero) client.ManifestStore {
	return &manifestStore{
		configMap:      configMap,
		storageClass:   storageClass,
		externalSecret: externalSecret,
		fs:             fs,
	}
}
