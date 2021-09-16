package fetch

import (
	"io"

	"github.com/mishudark/errors"
)

// Stager stores the state required for fetching and verifying a binary
type Stager struct {
	BinaryPath string

	Destination  io.WriteCloser
	Storage      Storage
	Fetcher      Fetcher
	Verifier     Verifier
	Decompressor Decompressor
}

// NewStager creates a new stager for fetching and verifying a binary
func NewStager(dest io.WriteCloser, s Storage, f Fetcher, v Verifier, d Decompressor) *Stager {
	return &Stager{
		Destination:  dest,
		Storage:      s,
		Fetcher:      f,
		Verifier:     v,
		Decompressor: d,
	}
}

// Fetch the binary and ensure that no errors occurred
func (s *Stager) Fetch() error {
	// We have already fetched this binary
	if len(s.BinaryPath) > 0 {
		return nil
	}

	var err error

	defer func() {
		err = s.Storage.Clean()
		if err != nil {
			return
		}

		err = s.Destination.Close()
	}()

	raw, err := s.Storage.Create("raw-content")
	if err != nil {
		return errors.E(err, "failed to create temporary storage", errors.Transient)
	}

	_, err = s.Fetcher.Fetch(raw)
	if err != nil {
		return errors.E(err, "failed to fetch binary", errors.IO)
	}

	if _, err = raw.Seek(0, 0); err != nil {
		return errors.E(err, "failed to reset buffer", errors.Internal)
	}

	err = s.Verifier.Verify(raw)
	if err != nil {
		return errors.E(err, "failed to verify binary signature", errors.Invalid)
	}

	if _, err = raw.Seek(0, 0); err != nil {
		return errors.E(err, "failed to reset buffer", errors.Internal)
	}

	err = s.Decompressor.Decompress(raw, s.Destination)
	if err != nil {
		return errors.E(err, "failed to decompress binary", errors.Invalid)
	}

	return err
}
