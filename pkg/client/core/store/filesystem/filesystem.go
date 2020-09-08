package filesystem

import "github.com/oslokommune/okctl/pkg/api"

// Paths contains the paths where the output should
// be stored
type Paths struct {
	OutputFile         string
	ConfigFile         string
	CloudFormationFile string
	ReleaseFile        string
	ChartFile          string
	BaseDir            string
}

// ManagedPolicy contains the state that is stored to
// and retrieved from the filesystem
type ManagedPolicy struct {
	ID        api.ID
	StackName string
	PolicyARN string
}

// ServiceAccount contains the data that should
// be serialised to the output file
type ServiceAccount struct {
	ID        api.ID
	PolicyArn string
}

// Helm contains the outputs we will store
type Helm struct {
	ID api.ID
}

// Kube contains the stored state for a kube deployment
// nolint: godox
// FIXME: We need to add a metadata map[string]interface{} type to the kube thingy
type Kube struct {
	ID           api.ID
	HostedZoneID string
	DomainFilter string
	Manifests    []string
}
