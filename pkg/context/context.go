// Package context provides an interface to ephemeral resources
package context

import (
	"context"
	"io"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

const (
	// DefaultDebugEnv if set will ensure verbose debugging output
	DefaultDebugEnv = "OKCTL_DEBUG"
	// DefaultNoInputEnv if set will ensure that no interactive dialogs are started
	DefaultNoInputEnv = "OKCTL_NO_INPUT"
)

// Context provides access to ephemeral state
// nolint: maligned
type Context struct {
	FileSystem *afero.Afero

	Debug   bool
	NoInput bool

	In  io.Reader
	Out io.Writer
	Err io.Writer

	Ctx context.Context

	Logger   *logrus.Logger
	LogLevel logrus.Level
}

// New returns a context with sensible defaults
func New() *Context {
	_, debug := os.LookupEnv(DefaultDebugEnv)
	_, noInput := os.LookupEnv(DefaultNoInputEnv)

	logger := logrus.New()

	logger.Out = os.Stderr
	logger.Formatter = &logrus.TextFormatter{}
	logger.Level = logrus.InfoLevel

	if debug {
		logger.Level = logrus.DebugLevel
	}

	return &Context{
		FileSystem: &afero.Afero{Fs: afero.NewOsFs()},
		Debug:      debug,
		NoInput:    noInput,
		In:         os.Stdin,
		Out:        os.Stdout,
		Err:        os.Stderr,
		Ctx:        context.Background(),
		Logger:     logger,
		LogLevel:   logger.Level,
	}
}
