package stage

import (
	"io"
	"path"
	"strings"

	humanize "github.com/dustin/go-humanize"
	"github.com/oslokommune/okctl/pkg/config/application"
	"github.com/oslokommune/okctl/pkg/digest"
	"github.com/oslokommune/okctl/pkg/storage"
)

type Stager interface {
	Run() error
}

type stager struct {
	dest         io.Writer
	storage      Storage
	fetcher      Fetcher
	verifier     Verifier
	decompressor Decompressor
}

func (s *stager) Run() error {
	var err error

	defer func() {
		err = s.storage.Clean()
	}()

	raw, err := s.storage.Create("raw-content")
	if err != nil {
		return err
	}

	_, err = s.fetcher.Fetch(raw)
	if err != nil {
		return err
	}

	if _, err = raw.Seek(0, 0); err != nil {
		return err
	}

	err = s.verifier.Verify(raw)
	if err != nil {
		return err
	}

	if _, err = raw.Seek(0, 0); err != nil {
		return err
	}

	err = s.decompressor.Decompress(raw, s.dest)
	if err != nil {
		return err
	}

	return err
}

func New(dest io.Writer, s Storage, f Fetcher, v Verifier, d Decompressor) Stager {
	return &stager{
		dest:         dest,
		storage:      s,
		fetcher:      f,
		verifier:     v,
		decompressor: d,
	}
}

func FromConfig(binaries []application.Binary, host application.Host, s storage.Storer) (map[string]Stager, error) {
	stagers := map[string]Stager{}

	for _, binary := range binaries {
		binaryBaseDir := path.Join("binaries", binary.Name, binary.Version, host.Os, host.Arch)
		binaryPath := path.Join(binaryBaseDir, binary.Name)

		exists, err := s.Exists(binaryPath)
		if err != nil {
			return nil, err
		}

		if !exists {
			checksums := map[digest.DigestType]string{}

			for _, c := range binary.Checksums {
				if c.Arch == host.Arch && c.Os == host.Os {
					checksums[digest.DigestType(c.Type)] = c.Digest
				}
			}

			bufferSize, err := humanize.ParseBytes(binary.BufferSize)
			if err != nil {
				return nil, err
			}

			binaryWriter, err := s.Create(binaryBaseDir, binary.Name, 0755)
			if err != nil {
				return nil, err
			}

			replacements := map[string]string{
				"#{os}":   host.Os,
				"#{arch}": host.Arch,
				"#{ver}":  binary.Version,
			}

			var d Decompressor

			switch binary.Archive.Type {
			case ".tar.gz":
				d = NewGzipTarDecompressor(binary.Archive.Target, int64(bufferSize))
			case ".zip":
				d = NewZipDecompressor(binary.Archive.Target, int64(bufferSize))
			default:
				d = NewNoopDecompressor()
			}

			stagers[binary.Name] = New(
				binaryWriter,
				NewEphemeralStorage(),
				NewHTTPFetcher(replaceVars(binary.URLPattern, replacements)),
				NewVerifier(checksums),
				d,
			)
		}
	}

	return stagers, nil
}

func replaceVars(content string, vars map[string]string) string {
	for v, r := range vars {
		content = strings.Replace(content, v, r, -1)
	}

	return content
}
