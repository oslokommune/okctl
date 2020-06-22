package load

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/oslokommune/okctl/pkg/config"
)

type DataNotFoundFn func(*config.Config) error

type DataNotFoundErr struct {
	err error
}

func (e *DataNotFoundErr) Error() string {
	return e.err.Error()
}

func PromptContinue(msg, errMsg string) error {
	var doContinue bool

	prompt := &survey.Confirm{
		Message: msg,
		Default: true,
	}

	err := survey.AskOne(prompt, &doContinue)
	if err != nil {
		return err
	}

	if !doContinue {
		return fmt.Errorf(errMsg)
	}

	return nil
}
