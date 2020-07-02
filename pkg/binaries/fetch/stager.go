// Package fetch knows how get, verify and stage binaries
package fetch

import (
	"fmt"
	"io"
	"path"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/binaries/digest"
	"github.com/oslokommune/okctl/pkg/config/application"
	"github.com/oslokommune/okctl/pkg/storage"
)

// Provider defines the interface for fetching a binary
// and returning a path to its location
type Provider interface {
	Fetch(name, version string) (string, error)
}

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

// Processor is a provider that knows how to fetch binaries via https
type Processor struct {
	Host           application.Host
	Store          storage.Storer
	Preload        bool
	Binaries       []application.Binary
	LoadedBinaries map[string]*Stager
}

// New returns a provider that knows how to fetch binaries via https
func New(preload bool, host application.Host, binaries []application.Binary, store storage.Storer) (*Processor, error) {
	p := &Processor{
		Host:           host,
		Store:          store,
		Preload:        preload,
		Binaries:       binaries,
		LoadedBinaries: map[string]*Stager{},
	}

	return p.prepareAndLoad()
}

// Stager returns a configured stager
func (s *Processor) Stager(baseDir string, bufferSize int64, binary application.Binary) (*Stager, error) {
	var d Decompressor

	switch binary.Archive.Type {
	case ".tar.gz":
		d = NewGzipTarDecompressor(binary.Archive.Target, bufferSize)
	case ".zip":
		d = NewZipDecompressor(binary.Archive.Target, bufferSize)
	default:
		d = NewNoopDecompressor()
	}

	binaryWriter, err := s.Store.Create(baseDir, binary.Name, 0755)
	if err != nil {
		return nil, err
	}

	stager := NewStager(
		binaryWriter,
		NewEphemeralStorage(),
		NewHTTPFetcher(
			replaceVars(binary.URLPattern, map[string]string{
				"#{os}":   s.Host.Os,
				"#{arch}": s.Host.Arch,
				"#{ver}":  binary.Version,
			}),
		),
		NewVerifier(
			checksumsFor(s.Host, binary.Checksums),
		),
		d,
	)

	return stager, nil
}

// prepareAndLoad a set of stagers
func (s *Processor) prepareAndLoad() (*Processor, error) {
	for _, binary := range s.Binaries {
		binaryBaseDir := path.Join("binaries", binary.Name, binary.Version, s.Host.Os, s.Host.Arch)
		binaryPath := path.Join(binaryBaseDir, binary.Name)

		exists, err := s.Store.Exists(binaryPath)
		if err != nil {
			return nil, errors.E(err, "failed to determine if binary exists", errors.IO)
		}

		if exists {
			s.LoadedBinaries[binaryIndex(binary.Name, binary.Version)] = &Stager{
				BinaryPath: s.Store.Abs(binaryPath),
			}

			continue
		}

		bufferSize, err := humanize.ParseBytes(binary.BufferSize)
		if err != nil {
			return nil, errors.E(err, "failed to parse buffer size", errors.Invalid)
		}

		stager, err := s.Stager(binaryBaseDir, int64(bufferSize), binary)
		if err != nil {
			return nil, errors.E(err, "failed to create stager", errors.Invalid)
		}

		stager.BinaryPath = s.Store.Abs(binaryPath)

		if s.Preload {
			err = errors.E(stager.Fetch(), "failed to preload binaries", errors.IO)
			if err != nil {
				return nil, err
			}
		}

		s.LoadedBinaries[binaryIndex(binary.Name, binary.Version)] = stager
	}

	return s, nil
}

// Fetch attempts to download and verify the binary
func (s *Processor) Fetch(name, version string) (string, error) {
	binary, hasKey := s.LoadedBinaries[binaryIndex(name, version)]
	if !hasKey {
		return "", fmt.Errorf("could not find configuration for binary: %s, with version: %s", name, version)
	}

	err := binary.Fetch()
	if err != nil {
		return "", err
	}

	return binary.BinaryPath, nil
}

func checksumsFor(h application.Host, cs []application.Checksum) map[digest.Type]string {
	out := map[digest.Type]string{}

	for _, c := range cs {
		if c.Arch == h.Arch && c.Os == h.Os {
			out[digest.Type(c.Type)] = c.Digest
		}
	}

	return out
}

func replaceVars(content string, vars map[string]string) string {
	for v, r := range vars {
		content = strings.Replace(content, v, r, -1)
	}

	return content
}

func binaryIndex(name, version string) string {
	return fmt.Sprintf("%s-%s", name, version)
}
