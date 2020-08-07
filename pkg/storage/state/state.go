// Package state implements common storage operations in a consistent way
package state

import (
	"bytes"
	"fmt"
	"io"
	"path"

	"github.com/oslokommune/okctl/pkg/config/application"
	"github.com/oslokommune/okctl/pkg/config/repository"
	"github.com/oslokommune/okctl/pkg/storage"
)

// PersisterProvider handles interactions with state for the app and repo
type PersisterProvider interface {
	Application() ApplicationPersister
	Repository() RepositoryPersister
}

// Persister provides a set of common operations for all persisters
type Persister interface {
	WriteConfig(b []byte) error
	ReadFromDefault(def string) ([]byte, error)
	WriteToDefault(def string, b []byte) error
	DeleteDefault(def string) error
	GetDefaultPath(def string) (string, error)
}

// RepositoryPersister defines an interface for working on the repo state
type RepositoryPersister interface {
	State() *repository.Data
	SaveState() error
	Persister
}

// ApplicationPersister defines an interface for working on the app state
type ApplicationPersister interface {
	State() *application.Data
	SaveState() error
	Persister
}

type persistors struct {
	app *appStore
	rep *repoStore
}

// Application returns the application persister
func (p *persistors) Application() ApplicationPersister {
	return p.app
}

// Repository returns the repository persister
func (p *persistors) Repository() RepositoryPersister {
	return p.rep
}

// Opts defines the common state for a persister
type Opts struct {
	BaseDir    string
	ConfigFile string
	Defaults   map[string]string
}

// New returns a new persister provider
func New(repOpts RepoStoreOpts, appOpts AppStoreOpts) PersisterProvider {
	return &persistors{
		app: &appStore{
			store: newStore(appOpts.BaseDir, appOpts.ConfigFile, appOpts.Defaults),
			state: appOpts.State,
		},
		rep: &repoStore{
			store: newStore(repOpts.BaseDir, repOpts.ConfigFile, repOpts.Defaults),
			state: repOpts.State,
		},
	}
}

// RepoStoreOpts defines the inputs for creating a repository persister
type RepoStoreOpts struct {
	Opts
	State *repository.Data
}

type repoStore struct {
	*store
	state *repository.Data
}

// GetDefaultPath returns the path to the default
func (r *repoStore) GetDefaultPath(def string) (string, error) {
	loc, ok := r.defaults[def]
	if !ok {
		return "", fmt.Errorf("no default for: %s, when trying to get", def)
	}

	return path.Join(r.baseDir, loc), nil
}

// SaveState will write the current repository state
func (r *repoStore) SaveState() error {
	data, err := r.state.YAML()
	if err != nil {
		return err
	}

	return r.WriteConfig(data)
}

// State returns the current repository state
func (r *repoStore) State() *repository.Data {
	return r.state
}

// AppStoreOpts defines the inputs for creating an application persister
type AppStoreOpts struct {
	Opts
	State *application.Data
}

type appStore struct {
	*store
	state *application.Data
}

// GetDefaultPath returns the path to the default
func (a *appStore) GetDefaultPath(def string) (string, error) {
	loc, ok := a.defaults[def]
	if !ok {
		return "", fmt.Errorf("no default for: %s, when trying to get", def)
	}

	return path.Join(a.baseDir, loc), nil
}

// SaveState will write the current application state
func (a *appStore) SaveState() error {
	data, err := a.state.YAML()
	if err != nil {
		return err
	}

	return a.WriteConfig(data)
}

// State returns the current application state
func (a *appStore) State() *application.Data {
	return a.state
}

type store struct {
	store      storage.Storer
	baseDir    string
	configFile string
	defaults   map[string]string
}

func newStore(baseDir, configFile string, defaults map[string]string) *store {
	return &store{
		store:      storage.NewFileSystemStorage(baseDir),
		baseDir:    baseDir,
		configFile: configFile,
		defaults:   defaults,
	}
}

// DeleteDefault knows how to remove default locations
func (s *store) DeleteDefault(def string) error {
	loc, ok := s.defaults[def]
	if !ok {
		return fmt.Errorf("no default for: %s, when trying to delete", def)
	}

	return s.store.RemoveAll(loc)
}

// ReadFromDefault knows how to load a file from a predefined default
func (s *store) ReadFromDefault(def string) ([]byte, error) {
	loc, ok := s.defaults[def]
	if !ok {
		return nil, fmt.Errorf("no default for: %s, when trying to read", def)
	}

	return s.store.ReadAll(loc)
}

// WriteToDefault knows to write to predefined default file
func (s *store) WriteToDefault(def string, b []byte) error {
	loc, ok := s.defaults[def]
	if !ok {
		return fmt.Errorf("no default for: %s, when trying to write", def)
	}

	dir, file := path.Split(loc)

	return s.recreate(dir, file, b)
}

// WriteConfig knows how to write the main config
func (s *store) WriteConfig(b []byte) error {
	return s.recreate("", s.configFile, b)
}

func (s *store) recreate(dir, file string, b []byte) error {
	writer, err := s.store.Recreate(dir, file, 0644)
	if err != nil {
		return err
	}

	defer func() {
		err = writer.Close()
	}()

	_, err = io.Copy(writer, bytes.NewReader(b))
	if err != nil {
		return err
	}

	return nil
}
