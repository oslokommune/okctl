package binaries

import (
	"io"

	"github.com/spf13/afero"
)

type Storage interface {
	Create(name string) (io.ReadWriteSeeker, error)
	Clean() error
}

type ephemeralStorage struct {
	fs afero.Fs
}

func (h *ephemeralStorage) Create(name string) (io.ReadWriteSeeker, error) {
	return h.fs.Create(name)
}

func (h *ephemeralStorage) Clean() error {
	h.fs = afero.NewMemMapFs()
	return nil
}

func NewEphemeralStorage() Storage {
	return &ephemeralStorage{
		fs: afero.NewMemMapFs(),
	}
}
