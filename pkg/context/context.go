package context

import (
	"io"
	"os"

	"github.com/spf13/afero"
)

const (
	DefaultDebugEnv   = "OKCTL_DEBUG"
	DefaultNoInputEnv = "OKCTL_NO_INPUT"
)

type Context struct {
	FileSystem *afero.Afero

	Debug   bool
	NoInput bool

	In  io.Reader
	Out io.Writer
	Err io.Writer
}

func New() *Context {
	_, debug := os.LookupEnv(DefaultDebugEnv)
	_, noInput := os.LookupEnv(DefaultNoInputEnv)

	return &Context{
		FileSystem: &afero.Afero{Fs: afero.NewOsFs()},
		Debug:      debug,
		NoInput:    noInput,
		In:         os.Stdin,
		Out:        os.Stdout,
		Err:        os.Stderr,
	}
}
