// Package fetch knows how get, verify and stage binaries
package fetch

import (
	"fmt"
	"io"
	"path"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/oslokommune/okctl/pkg/config/state"

	"github.com/dustin/go-humanize"
	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/binaries/digest"
	"github.com/oslokommune/okctl/pkg/storage"
)

// Provider defines the interface for fetching a binary
// and returning a path to its location
type Provider interface {
	Fetch(name, version string) (string, error)
}

// Processor is a provider that knows how to fetch binaries via https
type Processor struct {
	Host           state.Host
	Store          storage.Storer
	Preload        bool
	Binaries       []state.Binary
	LoadedBinaries map[string]*Stager
	Logger         *logrus.Logger
	Progress       io.Writer
}

// New returns a provider that knows how to fetch binaries via https
func New(progress io.Writer, logger *logrus.Logger, preload bool, host state.Host, binaries []state.Binary, store storage.Storer) (Provider, error) {
	p := &Processor{
		Host:           host,
		Store:          store,
		Preload:        preload,
		Binaries:       binaries,
		LoadedBinaries: map[string]*Stager{},
		Logger:         logger,
		Progress:       progress,
	}

	return p.prepareAndLoad()
}

// Stager returns a configured stager
func (s *Processor) Stager(baseDir string, bufferSize int64, binary state.Binary) (*Stager, error) {
	var d Decompressor

	switch binary.Archive.Type {
	case ".tar.gz":
		d = NewGzipTarDecompressor(binary.Archive.Target, bufferSize)
	case ".zip":
		d = NewZipDecompressor(binary.Archive.Target, bufferSize)
	default:
		d = NewNoopDecompressor()
	}

	binaryWriter, err := s.Store.Create(baseDir, binary.Name, 0o755)
	if err != nil {
		return nil, err
	}

	ver := NewVerifier(
		checksumsFor(s.Host, binary.Checksums),
	)

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
		ver,
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
			s.Logger.Debugf("binary already exists: %s (%s)", binary.Name, binary.Version)

			s.LoadedBinaries[binaryKey(binary.Name, binary.Version)] = &Stager{
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

		if s.Preload {
			msg := fmt.Sprintf("preloading missing binary: %s (%s)", binary.Name, binary.Version)

			_, err := fmt.Fprintln(s.Progress, msg)
			if err != nil {
				return nil, err
			}

			s.Logger.Debugln(msg)

			err = errors.E(stager.Fetch(), "failed to preload binaries", errors.IO)
			if err != nil {
				return nil, err
			}
		}

		stager.BinaryPath = s.Store.Abs(binaryPath)

		s.LoadedBinaries[binaryKey(binary.Name, binary.Version)] = stager
	}

	return s, nil
}

// Fetch attempts to download and verify the binary
func (s *Processor) Fetch(name, version string) (string, error) {
	binary, hasKey := s.LoadedBinaries[binaryKey(name, version)]
	if !hasKey {
		return "", fmt.Errorf("could not find configuration for binary: %s, with version: %s", name, version)
	}

	err := binary.Fetch()
	if err != nil {
		return "", err
	}

	return binary.BinaryPath, nil
}

func checksumsFor(host state.Host, checksums []state.Checksum) map[digest.Type]string {
	out := map[digest.Type]string{}

	for _, checksum := range checksums {
		if strings.EqualFold(checksum.Arch, host.Arch) && strings.EqualFold(checksum.Os, host.Os) {
			out[digest.Type(checksum.Type)] = checksum.Digest
		}
	}

	return out
}

func replaceVars(content string, vars map[string]string) string {
	for v, r := range vars {
		content = strings.ReplaceAll(content, v, r)
	}

	return content
}

func binaryKey(name, version string) string {
	return fmt.Sprintf("%s-%s", name, version)
}
