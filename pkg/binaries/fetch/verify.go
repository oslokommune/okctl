package fetch

import (
	"fmt"
	"io"

	"github.com/oslokommune/okctl/pkg/binaries/digest"
	"github.com/pkg/errors"
)

// Verifier provides an interface for verifying a file
type Verifier interface {
	Verify(io.Reader) error
}

type noopVerifier struct{}

// Verify does nothing, simply return nil
func (v *noopVerifier) Verify(io.Reader) error {
	return nil
}

// NewNoopVerifier creates a verifier that does nothing
func NewNoopVerifier() Verifier {
	return &noopVerifier{}
}

// NewVerifier returns a verifier that understand common
// hashing algorithms
func NewVerifier(digests map[digest.Type]string) Verifier {
	digestTypes := make([]digest.Type, len(digests))
	i := 0

	for digestType := range digests {
		digestTypes[i] = digestType
		i++
	}

	return &verifier{
		digests:     digests,
		digestTypes: digestTypes,
	}
}

type verifier struct {
	digests     map[digest.Type]string
	digestTypes []digest.Type
}

// ErrNoDigests indicates that no digests were provided
var ErrNoDigests = errors.New("no digests provided")

// VerifyDigests using the provided input.
func (v *verifier) Verify(reader io.Reader) error {
	if len(v.digests) == 0 {
		return ErrNoDigests
	}

	calculatedDigests, err := digest.NewDigester(v.digestTypes...).Digest(reader)
	if err != nil {
		return errors.Wrap(err, "failed to verify digests")
	}

	for hashType, digest := range calculatedDigests {
		if v.digests[hashType] != digest {
			return fmt.Errorf("verification failed, hash mismatch, got: %s, expected: %s", digest, v.digests[hashType])
		}

		delete(calculatedDigests, hashType)
	}

	if len(calculatedDigests) != 0 {
		return fmt.Errorf("failed to verify all hashes we produced")
	}

	return nil
}
