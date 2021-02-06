// Package awsiamauthenticator knows how to invoke the aws-iam-authenticator CLI
package awsiamauthenticator

const (
	// Name sets the name of the binary/cli
	Name = "aws-iam-authenticator"
	// Version sets the currently used version of the binary/cli
	Version = "0.5.2"
)

// AwsIamAuthenticator stores state for running the cli
type AwsIamAuthenticator struct {
	BinaryPath string
}

// New creates a new kubectl cli wrapper
func New(binaryPath string) *AwsIamAuthenticator {
	return &AwsIamAuthenticator{
		BinaryPath: binaryPath,
	}
}
