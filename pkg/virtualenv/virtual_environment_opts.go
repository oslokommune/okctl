package virtualenv

import (
	"github.com/oslokommune/okctl/pkg/storage"
)

// VirtualEnvironmentOpts contains the required inputs
type VirtualEnvironmentOpts struct {
	OsEnvVars       map[string]string
	EtcStorage      storage.Storer
	UserDirStorage  storage.Storer
	TmpStorage      storage.Storer
	Environment     string
	CurrentUsername string
}
