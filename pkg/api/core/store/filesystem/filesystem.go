package filesystem

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
