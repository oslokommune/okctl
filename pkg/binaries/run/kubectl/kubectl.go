// Package kubectl provides a convenient way of running kubectl commands
package kubectl

import "github.com/oslokommune/okctl/pkg/config/state"

const (
	// Name sets the name of the binary/cli
	Name = "kubectl"
	// Version sets the currently used version of the binary/cli
	Version = "1.16.8"
)

// Kubectl stores state for running the cli
type Kubectl struct {
	BinaryPath string
}

// New creates a new kubectl cli wrapper
func New(binaryPath string) *Kubectl {
	return &Kubectl{
		BinaryPath: binaryPath,
	}
}

// KnownBinaries returns the known binaries
func KnownBinaries() []state.Binary {
	return []state.Binary{
		{
			Name:       "kubectl",
			Version:    "1.16.8",
			BufferSize: "100mb",
			URLPattern: "https://amazon-eks.s3.us-west-2.amazonaws.com/#{ver}/2020-04-16/bin/#{os}/#{arch}/kubectl",
			Checksums: []state.Checksum{
				{
					Os:     "darwin",
					Arch:   "amd64",
					Type:   "sha256",
					Digest: "6e8439099c5a7d8d2f8f550f2f04301f9b0bb229a5f7c56477743a2cd11de2aa",
				},
				{
					Os:     "linux",
					Arch:   "amd64",
					Type:   "sha256",
					Digest: "e29544e1334f68e81546b8c8774c2484cbf82e8e5723d2a7e654f8a8fd79a7b2",
				},
			},
		},
	}
}
