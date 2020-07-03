package kubectl

const (
	// Name sets the name of the binary/cli
	Name = "kubectl"
	// Version sets the currently used version of the binary/cli
	Version = "1.16.8"
)

type Kubectl struct {
	BinaryPath string
}

func New(binaryPath string) *Kubectl {
	return &Kubectl{
		BinaryPath: binaryPath,
	}
}
