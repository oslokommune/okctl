package state

import (
	"github.com/oslokommune/okctl/pkg/binaries/run/eksctl"
	"github.com/oslokommune/okctl/pkg/binaries/run/kubectl"
)

// KnownBinaries returns a list of known binaries
func KnownBinaries() (binaries []Binary) {
	binaries = append(binaries, EksctlKnownBinaries()...)
	binaries = append(binaries, AWSIamAuthenticatorKnownBinaries()...)
	binaries = append(binaries, KubectlKnownBinaries()...)
	binaries = append(binaries, KubensKnownBinaries()...)

	return binaries
}

// AWSIamAuthenticatorKnownBinaries returns the known binaries
func AWSIamAuthenticatorKnownBinaries() []Binary {
	return []Binary{
		{
			Name:       "aws-iam-authenticator",
			Version:    "0.5.3",
			BufferSize: "100mb",
			URLPattern: "https://github.com/kubernetes-sigs/aws-iam-authenticator/releases/download/v#{ver}/aws-iam-authenticator_#{ver}_#{os}_#{arch}",
			Checksums: []Checksum{
				{
					Os:     "darwin",
					Arch:   "amd64",
					Type:   "sha256",
					Digest: "effd376c6ad00e90e45384000decac89f8495c76a3f52dee9eec389cfda236b7",
				},
				{
					Os:     "linux",
					Arch:   "amd64",
					Type:   "sha256",
					Digest: "20f4d8ece0f867c38b917ebe37dff934a31aabe385e26986183b14d72c70c137",
				},
			},
		},
	}
}

// EksctlKnownBinaries returns the known binaries
func EksctlKnownBinaries() []Binary {
	return []Binary{
		{
			Name:       "eksctl",
			Version:    eksctl.Version,
			BufferSize: "200mb", // Note the BufferSize must be higher than file size (this is a bug).
			URLPattern: "https://github.com/weaveworks/eksctl/releases/download/v#{ver}/eksctl_#{os}_#{arch}.tar.gz",
			Archive: Archive{
				Type:   ".tar.gz",
				Target: "eksctl",
			},
			// Example: How to find digest:
			// URL="https://github.com/weaveworks/eksctl/releases/download/v0.104.0/eksctl_darwin_amd64.tar.gz"; curl --location $URL | sha256sum
			// URL="https://github.com/weaveworks/eksctl/releases/download/v0.104.0/eksctl_linux_amd64.tar.gz"; curl --location $URL | sha256sum
			Checksums: []Checksum{
				{
					Os:     "darwin",
					Arch:   "amd64",
					Type:   "sha256",
					Digest: "a1ea933fc998ebc00502c993f6489c3a73535e846a79b64ac70dbe290189712c",
				},
				{
					Os:     "linux",
					Arch:   "amd64",
					Type:   "sha256",
					Digest: "f00539c6a00f7acfc614023961f4eec2fbcdaa836390b703aba63585b360afc3",
				},
			},
		},
	}
}

// KubectlKnownBinaries returns the known binaries
//
// For versions, see
// - https://kubernetes.io/releases/
// - https://kubernetes.io/releases/patch-releases/
func KubectlKnownBinaries() []Binary {
	return []Binary{
		{
			Name:       "kubectl",
			Version:    kubectl.Version,
			BufferSize: "100mb",
			URLPattern: "https://dl.k8s.io/release/v#{ver}/bin/#{os}/#{arch}/kubectl",
			// Example: How to find digest:
			// URL="https://dl.k8s.io/release/v1.21.14/bin/darwin/amd64/kubectl.sha256"; curl --location $URL
			// URL="https://dl.k8s.io/release/v1.21.14/bin/linux/amd64/kubectl.sha256"; curl --location $URL
			Checksums: []Checksum{
				{
					Os:     "darwin",
					Arch:   "amd64",
					Type:   "sha256",
					Digest: "30c529fe2891eb93dda99597b5c84cb10d2318bb92ae89e1e6189b3ae5fb6296",
				},
				{
					Os:     "linux",
					Arch:   "amd64",
					Type:   "sha256",
					Digest: "0c1682493c2abd7bc5fe4ddcdb0b6e5d417aa7e067994ffeca964163a988c6ee",
				},
			},
		},
	}
}

// KubensKnownBinaries returns the known binaries
func KubensKnownBinaries() []Binary {
	// Known limitation of Kubens: release-builds does not follow the same pattern as the above binaries,
	// resulting in different architecture releases. Hardcoded to `_x86_64` for now
	return []Binary{
		{
			Name:       "kubens",
			Version:    "0.9.4",
			BufferSize: "100mb",
			URLPattern: "https://github.com/ahmetb/kubectx/releases/download/v#{ver}/kubens_v#{ver}_#{os}_x86_64.tar.gz",
			Archive: Archive{
				Type:   ".tar.gz",
				Target: "kubens",
			},
			Checksums: []Checksum{
				{
					Os:     "darwin",
					Arch:   "amd64",
					Type:   "sha256",
					Digest: "ef43ab1217e09ac1b929d4b9dd2c22cbb10540ef277a3a9b484c020820c988b1",
				},
				{
					Os:     "linux",
					Arch:   "amd64",
					Type:   "sha256",
					Digest: "8b3672961fb15f8b87d5793af8bd3c1cca52c016596fbf57c46ab4ef39265fcd",
				},
			},
		},
	}
}
