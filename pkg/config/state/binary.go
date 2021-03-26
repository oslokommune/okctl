package state

// KnownBinaries returns a list of known binaries
func KnownBinaries() (binaries []Binary) {
	binaries = append(binaries, EksctlKnownBinaries()...)
	binaries = append(binaries, AWSIamAuthenticatorKnownBinaries()...)
	binaries = append(binaries, KubectlKnownBinaries()...)

	return binaries
}

// AWSIamAuthenticatorKnownBinaries returns the known binaries
func AWSIamAuthenticatorKnownBinaries() []Binary {
	return []Binary{
		{
			Name:       "aws-iam-authenticator",
			Version:    "0.5.2",
			BufferSize: "100mb",
			URLPattern: "https://github.com/kubernetes-sigs/aws-iam-authenticator/releases/download/v#{ver}/aws-iam-authenticator_#{ver}_#{os}_#{arch}",
			Checksums: []Checksum{
				{
					Os:     "darwin",
					Arch:   "amd64",
					Type:   "sha256",
					Digest: "f418c52d90947e56c9d9b01d3f32bbe52a0ba5ec02b65fc1ca9b85bff1652c2b",
				},
				{
					Os:     "linux",
					Arch:   "amd64",
					Type:   "sha256",
					Digest: "5bbe44ad7f6dd87a02e0b463a2aed9611836eb2f40d7fbe8c517460a4385621b",
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
			Version:    "0.41.0",
			BufferSize: "100mb",
			URLPattern: "https://github.com/weaveworks/eksctl/releases/download/#{ver}/eksctl_#{os}_#{arch}.tar.gz",
			Archive: Archive{
				Type:   ".tar.gz",
				Target: "eksctl",
			},
			Checksums: []Checksum{
				{
					Os:     "darwin",
					Arch:   "amd64",
					Type:   "sha256",
					Digest: "f703322bd778a0d59b133ebfb93c09170fb732d87504b5cd4cb6dded7f538556",
				},
				{
					Os:     "linux",
					Arch:   "amd64",
					Type:   "sha256",
					Digest: "0769b8a784bf75ccd91c39d5b96a12a55a6d5995dbd8fc97d58d7a14929d9d6c",
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
			Version:    "1.18.9",
			BufferSize: "100mb",
			URLPattern: "https://amazon-eks.s3.us-west-2.amazonaws.com/#{ver}/2020-11-02/bin/#{os}/#{arch}/kubectl",
			Checksums: []Checksum{
				{
					Os:     "darwin",
					Arch:   "amd64",
					Type:   "sha256",
					Digest: "f3f3919bf94d7b7f2014e2e9b318f049f4de378aed62833d609d211cf416935b",
				},
				{
					Os:     "linux",
					Arch:   "amd64",
					Type:   "sha256",
					Digest: "3dbe69e6deb35fbd6fec95b13d20ac1527544867ae56e3dae17e8c4d638b25b9",
				},
			},
		},
	}
}
