package storage

import (
	"os"
	"path"

	"github.com/spf13/afero"
)

type Storer interface {
	Create(dir, name string, perms os.FileMode) (afero.File, error)
	Exists(name string) (bool, error)
}

type Storage struct {
	Fs afero.Fs
}

func (s *Storage) Create(baseDir, name string, perms os.FileMode) (afero.File, error) {
	err := s.Fs.MkdirAll(baseDir, 0755)
	if err != nil {
		return nil, err
	}

	return s.Fs.OpenFile(path.Join(baseDir, name), os.O_RDWR|os.O_CREATE, perms)
}

func (s *Storage) Exists(name string) (bool, error) {
	return afero.Exists(s.Fs, name)
}

func NewFileSystemStorage(path string) *Storage {
	return &Storage{
		Fs: afero.NewBasePathFs(afero.NewOsFs(), path),
	}
}

func NewEphemeralStorage() *Storage {
	return &Storage{
		Fs: afero.NewMemMapFs(),
	}
}
