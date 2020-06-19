package digest

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash"
	"io"

	"github.com/pkg/errors"
)

// Digester defines the functions related to creating
// digests of some given input.
type Digester interface {
	Digest(reader io.Reader) (map[DigestType]string, error)
}

// DigestType enumerates the supported digesters.
type DigestType string

const (
	TypeSHA256 DigestType = "sha256"
	TypeSHA512 DigestType = "sha512"
)

type digest struct {
	digesters []DigestType
}

// NewDigester creates a new digester with the given digest types.
func NewDigester(digesters ...DigestType) Digester {
	uniqueDigesters := map[DigestType]struct{}{}
	for _, dig := range digesters {
		uniqueDigesters[dig] = struct{}{}
	}

	toDigest := make([]DigestType, len(uniqueDigesters))
	i := 0

	for dig := range uniqueDigesters {
		toDigest[i] = dig
		i++
	}

	return &digest{
		digesters: toDigest,
	}
}

// Digest returns the hashes of the loaded data given
// the provided digesters.
func (d *digest) Digest(reader io.Reader) (map[DigestType]string, error) {
	if reader == nil {
		return nil, fmt.Errorf("reader is nil")
	}

	type Digested struct {
		Type DigestType
		Hash string

		hasher hash.Hash
	}

	digested := make([]*Digested, len(d.digesters))
	hashers := make([]io.Writer, len(d.digesters))

	for i, d := range d.digesters {
		var h hash.Hash

		switch d {
		case TypeSHA256:
			h = sha256.New()
		case TypeSHA512:
			h = sha512.New()
		default:
			return nil, fmt.Errorf("unsupported digester: %s", d)
		}

		hashers[i] = h
		digested[i] = &Digested{
			Type:   d,
			hasher: h,
		}
	}

	writer := io.MultiWriter(hashers...)
	if _, err := io.Copy(writer, reader); err != nil {
		return nil, errors.Wrap(err, "failed to write content to digesters")
	}

	digests := map[DigestType]string{}
	for _, d := range digested {
		digests[d.Type] = hex.EncodeToString(d.hasher.Sum(nil))
	}

	return digests, nil
}
