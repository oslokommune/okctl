package binaries

import (
	"io"

	"github.com/spf13/afero"
)

// Storage knows how to store stuff
type Storage interface {
	Create(name string) (io.ReadWriteSeeker, error)
	Clean() error
}

type ephemeralStorage struct {
	fs afero.Fs
}

// Create a file in ephemeral storage
func (h *ephemeralStorage) Create(name string) (io.ReadWriteSeeker, error) {
	return h.fs.Create(name)
}

// Clean by removing everything
func (h *ephemeralStorage) Clean() error {
	h.fs = afero.NewMemMapFs()
	return nil
}

// NewEphemeralStorage returns an in memory storage
func NewEphemeralStorage() Storage {
	return &ephemeralStorage{
		fs: afero.NewMemMapFs(),
	}
}
