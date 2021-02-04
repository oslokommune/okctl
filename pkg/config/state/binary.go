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
			Version:    "0.37.0",
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
					Digest: "e63349d961a5f9c7f13690b1915b8c6bd345d552673dae18b360d2a4b63cf5a8",
				},
				{
					Os:     "linux",
					Arch:   "amd64",
					Type:   "sha256",
					Digest: "4ef5775f237fd9c5ea3f3984e9a1bd9e7fe1b8106634665b27153eb58c245437",
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
