// Package awsiamauthenticator knows how to invoke the aws-iam-authenticator CLI
package awsiamauthenticator

import "github.com/oslokommune/okctl/pkg/config/state"

const (
	// Name sets the name of the binary/cli
	Name = "aws-iam-authenticator"
	// Version sets the currently used version of the binary/cli
	Version = "0.5.1"
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

// KnownBinaries returns the known binaries
func KnownBinaries() []state.Binary {
	return []state.Binary{
		{
			Name:       "aws-iam-authenticator",
			Version:    "0.5.1",
			BufferSize: "100mb",
			URLPattern: "https://github.com/kubernetes-sigs/aws-iam-authenticator/releases/download/v#{ver}/aws-iam-authenticator_#{ver}_#{os}_#{arch}",
			Checksums: []state.Checksum{
				{
					Os:     "darwin",
					Arch:   "amd64",
					Type:   "sha256",
					Digest: "e6050faee00732d1da88e4ba9910bcb03f0fc40eaf192e39dd55dfdf6cf6f681",
				},
				{
					Os:     "linux",
					Arch:   "amd64",
					Type:   "sha256",
					Digest: "afb16f35071c977554f1097cbb84ca4f38f9ce42142c8a0612716ae66bb9fdb9",
				},
			},
		},
	}
}
