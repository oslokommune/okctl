package commandlineprompter

import (
	"github.com/oslokommune/okctl/pkg/storage"
)

// CommandLinePromptOpts contains the required inputs to create a command line prompt
type CommandLinePromptOpts struct {
	OsEnvVars          map[string]string
	EtcStorage         storage.Storer
	UserDirStorage     storage.Storer
	UserHomeDirStorage storage.Storer
	TmpStorage         storage.Storer
	Environment        string
	CurrentUsername    string
}
