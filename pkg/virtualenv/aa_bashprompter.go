package virtualenv

type BashPrompter struct {

}

func (b BashPrompter) BuildPrompt(opts VirtualEnvironmentOpts) ([]byte, error) {
	panic("implement me")
}

func (b BashPrompter) SetPrompt(bytes []byte) error {
	panic("implement me")
}

func (b BashPrompter) Cleanup() error {
	panic("implement me")
}

func newBashPrompter() *BashPrompter {
	return &BashPrompter{}
}

