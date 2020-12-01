package shellgetter

type shellCmdGetter interface {
	// Get returns a shell command based on the user's environment
	Get() (string, error)
}
