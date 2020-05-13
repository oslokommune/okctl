package storage

import (
	"io"
	"os"
	"path"

	"github.com/spf13/afero"
)

type Storage interface {
	Create(dir, name string) (io.ReadWriteSeeker, error)
	Exists(name string) (bool, error)
}

type storage struct {
	fs afero.Fs
}

func (s *storage) Create(baseDir, name string) (io.ReadWriteSeeker, error) {
	err := s.fs.MkdirAll(baseDir, 0755)
	if err != nil {
		return nil, err
	}

	return s.fs.OpenFile(path.Join(baseDir, name), os.O_RDWR|os.O_CREATE, 0755)
}

func (s *storage) Exists(name string) (bool, error) {
	return afero.Exists(s.fs, name)
}

func NewFileSystemStorage(path string) Storage {
	return &storage{
		fs: afero.NewBasePathFs(afero.NewOsFs(), path),
	}
}

func NewEphemeralStorage() Storage {
	return &storage{
		fs: afero.NewMemMapFs(),
	}
}
