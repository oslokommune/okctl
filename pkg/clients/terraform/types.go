package terraform

// Client defines functionality for a terraform client
type Client interface {
	Initialize(directory string) error
	Apply(directory string) error
}
