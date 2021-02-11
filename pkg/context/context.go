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
	// DefaultAWSCredentialsType if set will ensure where okctl seeks AWS authentication credentials
	DefaultAWSCredentialsType = "OKCTL_AWS_CREDENTIALS_TYPE"
	// DefaultGithubCredentialsType if set will ensure where okctl seeks Github authentication credentials
	DefaultGithubCredentialsType = "OKCTL_GITHUB_CREDENTIALS_TYPE"
)

const (
	// AWSCredentialsTypeSAML represents using SAML for AWS authentication
	AWSCredentialsTypeSAML = "saml"
	// AWSCredentialsTypeAccessKey represents using access key environment variables for AWS authentication
	AWSCredentialsTypeAccessKey = "access-key"

	// GithubCredentialsTypeDeviceAuthentication represents using the device authentication flow to authenticate with
	// Github
	GithubCredentialsTypeDeviceAuthentication = "device-authentication"
	// GithubCredentialsTypeToken represents using a Github Token (GH Actions) or a PAT for authentication
	GithubCredentialsTypeToken = "token"
)

// Context provides access to ephemeral state
// nolint: maligned
type Context struct {
	FileSystem *afero.Afero

	Debug                 bool
	AWSCredentialsType    string
	GithubCredentialsType string

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
	awsCredentialsType := os.Getenv(DefaultAWSCredentialsType)
	githubCredentialsType := os.Getenv(DefaultGithubCredentialsType)

	logger := logrus.New()

	logger.Out = os.Stderr
	logger.Formatter = &logrus.TextFormatter{}
	logger.Level = logrus.WarnLevel

	if debug {
		logger.Level = logrus.DebugLevel
	}

	return &Context{
		FileSystem:            &afero.Afero{Fs: afero.NewOsFs()},
		Debug:                 debug,
		AWSCredentialsType:    awsCredentialsType,
		GithubCredentialsType: githubCredentialsType,
		In:                    os.Stdin,
		Out:                   os.Stdout,
		Err:                   os.Stderr,
		Ctx:                   context.Background(),
		Logger:                logger,
		LogLevel:              logger.Level,
	}
}
