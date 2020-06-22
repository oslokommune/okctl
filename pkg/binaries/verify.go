package binaries

import (
	"fmt"
	"io"

	"github.com/oslokommune/okctl/pkg/binaries/digest"
	"github.com/pkg/errors"
)

type Verifier interface {
	Verify(io.Reader) error
}

type noopVerifier struct{}

func (v *noopVerifier) Verify(io.Reader) error {
	return nil
}

func NewNoopVerifier() Verifier {
	return &noopVerifier{}
}

func NewVerifier(digests map[digest.Type]string) Verifier {
	digestTypes := make([]digest.Type, len(digests))
	i := 0

	for digestType := range digests {
		digestTypes[i] = digestType
		i++
	}

	return &verifier{
		digests: digests,
	}
}

type verifier struct {
	digests     map[digest.Type]string
	digestTypes []digest.Type
}

var (
	// ErrNoDigests indicates that no digests were provided
	ErrNoDigests = errors.New("no digests provided")
)

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
