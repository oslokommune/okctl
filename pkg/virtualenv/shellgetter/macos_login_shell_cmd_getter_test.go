package shellgetter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// nolint: funlen
func TestCreateVirtualEnvironment(t *testing.T) {
	t.Run("Should parse dscl output correctly", func(t *testing.T) {
		cmdGetter := &macOsLoginShellCmdGetter{}

		parsed := cmdGetter.parseShellCmd("UserShell: /bin/bash")

		assert.Equal(t, "/bin/bash", parsed)
	})
}
