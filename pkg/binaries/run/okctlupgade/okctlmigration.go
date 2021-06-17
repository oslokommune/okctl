package okctlupgade

const (
	// Name sets the name of the binary/cli
	Name = "okctl-upgrade"
)

// AwsIamAuthenticator stores state for running the cli
type OkctlUpgrade struct {
	BinaryPath string
}

// New creates a new kubectl cli wrapper
func New(binaryPath string) *OkctlUpgrade {
	return &OkctlUpgrade{
		BinaryPath: binaryPath,
	}
}
