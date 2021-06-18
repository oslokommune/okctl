package upgrade

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/oslokommune/okctl/pkg/config/state"
	"regexp"
	"strings"
)

func parseChecksums(checksumBytes []byte) ([]state.Checksum, error) {
	reader := bytes.NewReader(checksumBytes)
	scanner := bufio.NewScanner(reader)
	var checksums []state.Checksum

	for scanner.Scan() {
		line := scanner.Text()

		parts := strings.Fields(line)
		if len(parts) != 2 {
			return nil, fmt.Errorf(
				"expected 2 substrings when splitting digest line on whitespace ( ), got %d in string '%s'",
				len(parts), line,
			)
		}

		digest := parts[0]   // Example: 1eaad82bd6e082936cfb4c108b9e5e46bba98ef19f33492ca2041de04803b86b
		filename := parts[1] // Example: okctl-upgrade_0.0.63_Darwin_amd64.tar.gz

		err := validateDigest(digest)
		if err != nil {
			return nil, fmt.Errorf("invalid digest '%s': %w", digest, err)
		}

		ugradeFile, err := parseOkctlUpgradeFilename(filename)
		if err != nil {
			return nil, fmt.Errorf("parsing upgrade filename: %w", err)
		}

		// TODO REMOVE
		//if ugradeFile.version != expectedVersion {
		//	return nil, fmt.Errorf("expected version '%s' but got '%s' in checksum file '%s'",
		//		expectedVersion, ugradeFile.version, filename)
		//}

		checksum := state.Checksum{
			Os:     ugradeFile.os,
			Arch:   ugradeFile.arch,
			Type:   ugradeFile.extension,
			Digest: digest,
		}

		checksums = append(checksums, checksum)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanning: %w", err)
	}

	return checksums, nil
}

func validateDigest(digest string) error {
	re, err := regexp.Compile(`^[0-9a-z]+$`)
	if err != nil {
		return fmt.Errorf("compiling regex: %w", err)
	}

	for range re.FindAllString(digest, -1) {
		return nil
	}

	return errors.New("invalid digest")
}
