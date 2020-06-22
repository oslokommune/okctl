package storage

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/spf13/afero"
)

type StoreCleaner interface {
	Storer
	Cleaner
}

type Cleaner interface {
	Clean() error
}

type Storer interface {
	Create(dir, name string, perms os.FileMode) (afero.File, error)
	MkdirAll(dir string) error
	Exists(name string) (bool, error)
	Abs(name string) string
}

type Storage struct {
	Path string
	Fs   afero.Fs
}

func (s *Storage) MkdirAll(dir string) error {
	return s.Fs.MkdirAll(dir, 0755)
}

func (s *Storage) Create(dir, file string, perms os.FileMode) (afero.File, error) {
	err := s.Fs.MkdirAll(dir, 0755)
	if err != nil {
		return nil, err
	}

	return s.Fs.OpenFile(path.Join(dir, file), os.O_RDWR|os.O_CREATE, perms)
}

func (s *Storage) Recreate(dir, file string, perms os.FileMode) (afero.File, error) {
	err := s.Fs.RemoveAll(path.Join(dir, file))
	if err != nil {
		return nil, err
	}

	return s.Create(dir, file, perms)
}

func (s *Storage) Exists(name string) (bool, error) {
	return afero.Exists(s.Fs, name)
}

func (s *Storage) Abs(name string) string {
	switch fs := s.Fs.(type) {
	case *afero.BasePathFs:
		return afero.FullBaseFsPath(fs, name)
	default:
		return path.Join(s.Path, name)
	}
}

func NewFileSystemStorage(path string) *Storage {
	return &Storage{
		Path: path,
		Fs:   afero.NewBasePathFs(afero.NewOsFs(), path),
	}
}

func NewEphemeralStorage() *Storage {
	return &Storage{
		Fs: afero.NewMemMapFs(),
	}
}

type TemporaryStorage struct {
	*Storage
}

func (s *TemporaryStorage) Clean() error {
	return os.RemoveAll(s.Path)
}

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
