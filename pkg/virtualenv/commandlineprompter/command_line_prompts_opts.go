package commandlineprompter

import (
	"github.com/oslokommune/okctl/pkg/storage"
	"github.com/oslokommune/okctl/pkg/virtualenv/shellgetter"
)

// CommandLinePromptOpts contains the required inputs to create a command line prompt
type CommandLinePromptOpts struct {
	Os                   shellgetter.Os
	MacOsUserShellGetter shellgetter.MacOsUserShellCmdGetter
	OsEnvVars            map[string]string
	EtcStorage           storage.Storer
	UserDirStorage       storage.Storer
	UserHomeDirStorage   storage.Storer
	TmpStorage           storage.Storer
	Environment          string
	CurrentUsername      string
}

// NewShellGetter creates a new ShellGetter
func NewShellGetter(opts CommandLinePromptOpts) *shellgetter.ShellGetter {
	return shellgetter.NewShellGetter(opts.Os, opts.MacOsUserShellGetter, opts.OsEnvVars, opts.EtcStorage, opts.CurrentUsername)
}
