package state

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
			Version:    "0.98.0",
			BufferSize: "200mb", // Note the BufferSize must be higher than file size (this is a bug).
			URLPattern: "https://github.com/weaveworks/eksctl/releases/download/v#{ver}/eksctl_#{os}_#{arch}.tar.gz",
			Archive: Archive{
				Type:   ".tar.gz",
				Target: "eksctl",
			},
			// Find digest by using
			// curl --location URL | sha256sum
			Checksums: []Checksum{
				{
					Os:     "darwin",
					Arch:   "amd64",
					Type:   "sha256",
					Digest: "2d948e28a3b69a71fc545aee22a0bf0a791fb5ee605212b424b505cbe45a23db",
				},
				{
					Os:     "linux",
					Arch:   "amd64",
					Type:   "sha256",
					Digest: "385de15f7bf361ec8bd723e90f64398ab1de20cd25e81977a45097ef8182e771",
				},
			},
		},
	}
}

// KubectlKnownBinaries returns the known binaries
func KubectlKnownBinaries() []Binary {
	return []Binary{
		{
			Name:       "kubectl",
			Version:    "1.20.15",
			BufferSize: "100mb",
			URLPattern: "https://dl.k8s.io/release/v#{ver}/bin/#{os}/#{arch}/kubectl",
			Checksums: []Checksum{
				{
					Os:     "darwin",
					Arch:   "amd64",
					Type:   "sha256",
					Digest: "6b6cf555a34271379b45013dfa9b580329314254aafc91b543bf2d83ebd1db74",
				},
				{
					Os:     "linux",
					Arch:   "amd64",
					Type:   "sha256",
					Digest: "d283552d3ef3b0fd47c08953414e1e73897a1b3f88c8a520bb2e7de4e37e96f3",
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
