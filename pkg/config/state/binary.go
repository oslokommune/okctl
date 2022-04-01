package state

// KnownBinaries returns a list of known binaries
func KnownBinaries() (binaries []Binary) {
	binaries = append(binaries, EksctlKnownBinaries()...)
	binaries = append(binaries, AWSIamAuthenticatorKnownBinaries()...)
	binaries = append(binaries, KubectlKnownBinaries()...)
	binaries = append(binaries, KubensKnownBinaries()...)
	binaries = append(binaries, TerraformKnownBinaries()...)

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

// TerraformKnownBinaries returns the known binaries
func TerraformKnownBinaries() []Binary {
	return []Binary{
		{
			Name:       "terraform",
			Version:    "1.1.7",
			BufferSize: "100mb",
			URLPattern: "https://releases.hashicorp.com/terraform/#{ver}/terraform_#{ver}_#{os}_#{arch}.zip",
			Archive: Archive{
				Type:   ".zip",
				Target: "terraform",
			},
			Checksums: []Checksum{
				{
					Os:     "darwin",
					Arch:   "amd64",
					Type:   "sha256",
					Digest: "5e7e939e084ae29af7fd86b00a618433d905477c52add2d4ea8770692acbceac",
				},
				{
					Os:     "linux",
					Arch:   "amd64",
					Type:   "sha256",
					Digest: "e4add092a54ff6febd3325d1e0c109c9e590dc6c38f8bb7f9632e4e6bcca99d4",
				},
			},
		},
	}
}
