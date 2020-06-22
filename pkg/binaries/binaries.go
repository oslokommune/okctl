package binaries

import (
	"fmt"
	"io"
	"path"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/oslokommune/okctl/pkg/binaries/digest"
	"github.com/oslokommune/okctl/pkg/config/application"
	"github.com/oslokommune/okctl/pkg/storage"
	"github.com/pkg/errors"
)

type Provider interface {
	Fetch(name, version string) (string, error)
}

func NewStager(dest io.WriteCloser, s Storage, f Fetcher, v Verifier, d Decompressor) *Stager {
	return &Stager{
		Destination:  dest,
		Storage:      s,
		Fetcher:      f,
		Verifier:     v,
		Decompressor: d,
	}
}

type Stager struct {
	BinaryPath string

	Destination  io.WriteCloser
	Storage      Storage
	Fetcher      Fetcher
	Verifier     Verifier
	Decompressor Decompressor
}

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
		return err
	}

	_, err = s.Fetcher.Fetch(raw)
	if err != nil {
		return err
	}

	if _, err = raw.Seek(0, 0); err != nil {
		return err
	}

	err = s.Verifier.Verify(raw)
	if err != nil {
		return err
	}

	if _, err = raw.Seek(0, 0); err != nil {
		return err
	}

	err = s.Decompressor.Decompress(raw, s.Destination)
	if err != nil {
		return err
	}

	return err
}

type ErrorProvider struct{}

func (s *ErrorProvider) Fetch(name, version string) (string, error) {
	return "", fmt.Errorf("this is an error operator, couldn't retrieve binary for: %s, version: %s", name, version)
}

func NewErrorProvider() *ErrorProvider {
	return &ErrorProvider{}
}

type DefaultProvider struct {
	Host     application.Host
	Store    storage.Storer
	Binaries map[string]*Stager
}

func New(host application.Host, store storage.Storer) *DefaultProvider {
	return &DefaultProvider{
		Host:     host,
		Store:    store,
		Binaries: map[string]*Stager{},
	}
}

func (s *DefaultProvider) Stager(baseDir string, bufferSize int64, binary application.Binary) (*Stager, error) {
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

func (s *DefaultProvider) FromConfig(preload bool, binaries []application.Binary) (*DefaultProvider, error) {
	for _, binary := range binaries {
		binaryBaseDir := path.Join("binaries", binary.Name, binary.Version, s.Host.Os, s.Host.Arch)
		binaryPath := path.Join(binaryBaseDir, binary.Name)

		exists, err := s.Store.Exists(binaryPath)
		if err != nil {
			return nil, err
		}

		if exists {
			s.Binaries[binaryIndex(binary.Name, binary.Version)] = &Stager{
				BinaryPath: s.Store.Abs(binaryPath),
			}

			continue
		}

		bufferSize, err := humanize.ParseBytes(binary.BufferSize)
		if err != nil {
			return nil, err
		}

		stager, err := s.Stager(binaryBaseDir, int64(bufferSize), binary)
		if err != nil {
			return nil, err
		}

		if preload {
			err = errors.Wrap(stager.Fetch(), "failed to preload binaries")
			if err != nil {
				return nil, err
			}
		}

		s.Binaries[binaryIndex(binary.Name, binary.Version)] = stager
	}

	return s, nil
}

func (s *DefaultProvider) Fetch(name, version string) (string, error) {
	binary, hasKey := s.Binaries[binaryIndex(name, version)]
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
