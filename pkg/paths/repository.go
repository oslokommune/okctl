package paths

import (
	"bytes"
	"fmt"
	"os/exec"
)

// GetAbsoluteRepositoryRootDirectory returns the absolute path of the repository root no matter what the current working
// directory of the repository the user is in.
func GetAbsoluteRepositoryRootDirectory() (string, error) {
	result, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return "", fmt.Errorf("getting repository root directory: %w", err)
	}

	pathAsString := string(bytes.Trim(result, "\n"))

	return pathAsString, nil
}
