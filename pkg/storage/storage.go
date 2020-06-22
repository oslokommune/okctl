// Package storage provides an API towards common
// storage operations
package storage

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/spf13/afero"
)

// StoreCleaner provides an interface for
// a store that allows for common read and write
// operations and wiping everything afterwards
type StoreCleaner interface {
	Storer
	Cleaner
}

// Cleaner provides an interface that may be implemented
// in order to cleanup a store
type Cleaner interface {
	Clean() error
}

// Storer provides an interface for creating and modifying files
// and directories
type Storer interface {
	Create(dir, name string, perms os.FileMode) (afero.File, error)
	MkdirAll(dir string) error
	Exists(name string) (bool, error)
	Abs(name string) string
}

// Storage stores state about the filesystem
type Storage struct {
	Path string
	Fs   afero.Fs
}

// MkdirAll creates a directory and all preceding directories
// if they do not exist
func (s *Storage) MkdirAll(dir string) error {
	return s.Fs.MkdirAll(dir, 0755)
}

// Create will create all directories leading to a file
// and then the file itself
func (s *Storage) Create(dir, file string, perms os.FileMode) (afero.File, error) {
	err := s.Fs.MkdirAll(dir, 0755)
	if err != nil {
		return nil, err
	}

	return s.Fs.OpenFile(path.Join(dir, file), os.O_RDWR|os.O_CREATE, perms)
}

// Recreate will delete a file and then recreate it
func (s *Storage) Recreate(dir, file string, perms os.FileMode) (afero.File, error) {
	err := s.Fs.Remove(path.Join(dir, file))
	if err != nil {
		return nil, err
	}

	return s.Create(dir, file, perms)
}

// Exists will determine if a file exists
func (s *Storage) Exists(name string) (bool, error) {
	return afero.Exists(s.Fs, name)
}

// Abs will return the absolute path to a file
func (s *Storage) Abs(name string) string {
	switch fs := s.Fs.(type) {
	case *afero.BasePathFs:
		return afero.FullBaseFsPath(fs, name)
	default:
		return path.Join(s.Path, name)
	}
}

// NewFileSystemStorage will return a store to a
// base path filesystem
func NewFileSystemStorage(path string) *Storage {
	return &Storage{
		Path: path,
		Fs:   afero.NewBasePathFs(afero.NewOsFs(), path),
	}
}

// NewEphemeralStorage will return an in memory file system
func NewEphemeralStorage() *Storage {
	return &Storage{
		Fs: afero.NewMemMapFs(),
	}
}

// TemporaryStorage wraps Storage and implements
// the Cleaner interface
type TemporaryStorage struct {
	*Storage
}

// Clean removes everything at the path the filesystem was
// created from
func (s *TemporaryStorage) Clean() error {
	return os.RemoveAll(s.Path)
}

// NewTemporaryStorage creates a new temporary storage
func NewTemporaryStorage() (*TemporaryStorage, error) {
	dir, err := ioutil.TempDir("", "okctl")
	if err != nil {
		return nil, err
	}

	return &TemporaryStorage{
		&Storage{
			Path: dir,
			Fs:   afero.NewBasePathFs(afero.NewOsFs(), dir),
		},
	}, nil
}
