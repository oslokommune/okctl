package shellgetter

type envShellCmdGetter struct {
	shellCmd string
}

func (e *envShellCmdGetter) Get() (string, error) {
	return e.shellCmd, nil
}
