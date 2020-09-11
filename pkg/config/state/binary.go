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
			Version:    "0.5.1",
			BufferSize: "100mb",
			URLPattern: "https://github.com/kubernetes-sigs/aws-iam-authenticator/releases/download/v#{ver}/aws-iam-authenticator_#{ver}_#{os}_#{arch}",
			Checksums: []Checksum{
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

// EksctlKnownBinaries returns the known binaries
func EksctlKnownBinaries() []Binary {
	return []Binary{
		{
			Name:       "eksctl",
			Version:    "0.25.0",
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
					Digest: "e232f48e4995f711620ea34c09f582b097e5b006f45fbe82a11fc8955636c9c4",
				},
				{
					Os:     "linux",
					Arch:   "amd64",
					Type:   "sha256",
					Digest: "e94e4ec335c036d8f511ea214d5a55dfd097e2053747d7d04d6db49fff107531",
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
			Version:    "1.16.8",
			BufferSize: "100mb",
			URLPattern: "https://amazon-eks.s3.us-west-2.amazonaws.com/#{ver}/2020-04-16/bin/#{os}/#{arch}/kubectl",
			Checksums: []Checksum{
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
