package stager

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"strings"
)

// Decompressor provides an interface for decompressing into a single file.
type Decompressor interface {
	Decompress(io.Reader, io.Writer) error
}

type noopDecompressor struct{}

func (d *noopDecompressor) Decompress(reader io.Reader, writer io.Writer) error {
	_, err := io.Copy(writer, reader)
	return err
}

func NewNoopDecompressor() Decompressor {
	return &noopDecompressor{}
}

type zipDecompressor struct {
	file       string
	bufferSize int64
}

func (z *zipDecompressor) Decompress(r io.Reader, w io.Writer) error {
	// Zip requires a ReaderAt interface, we provide this
	// by reading into a buffer, getting a string to that buffer
	// and creating a reader from the string.
	buf := new(bytes.Buffer)

	_, err := io.Copy(buf, io.LimitReader(r, z.bufferSize))
	if err != nil {
		return err
	}

	s := buf.String()

	zipReader, err := zip.NewReader(strings.NewReader(s), int64(len(s)))
	if err != nil {
		return err
	}

	for _, f := range zipReader.File {
		if f.Name == z.file {
			h, err := f.Open()
			if err != nil {
				return err
			}

			_, err = io.Copy(w, io.LimitReader(h, z.bufferSize))
			if err != nil {
				return err
			}

			return nil
		}
	}

	return fmt.Errorf("couldn't find: %s, in archive", z.file)
}

func NewZipDecompressor(file string, bufferSize int64) Decompressor {
	return &zipDecompressor{
		file:       file,
		bufferSize: bufferSize,
	}
}

type gzipTarDecompressor struct {
	file       string
	bufferSize int64
}

func (g *gzipTarDecompressor) Decompress(reader io.Reader, writer io.Writer) error {
	zs, err := gzip.NewReader(reader)
	if err != nil {
		return err
	}

	tr := tar.NewReader(zs)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		if header.Typeflag == tar.TypeReg && header.Name == g.file {
			_, err := io.Copy(writer, io.LimitReader(tr, g.bufferSize))
			if err != nil {
				return err
			}

			return nil
		}
	}

	return fmt.Errorf("couldn't find: %s, in archive", g.file)
}

func NewGzipTarDecompressor(file string, bufferSize int64) Decompressor {
	return &gzipTarDecompressor{
		file:       file,
		bufferSize: bufferSize,
	}
}
