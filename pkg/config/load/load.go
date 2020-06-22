package load

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/oslokommune/okctl/pkg/config"
)

// DataNotFoundFn determines what should be done if
// data is not found
type DataNotFoundFn func(*config.Config) error

// DataNotFoundErr stores error returned when data
// is not found at location
type DataNotFoundErr struct {
	err error
}

// Error returns the stored error
func (e *DataNotFoundErr) Error() string {
	return e.err.Error()
}

// PromptContinue asks the user a question and does
// not abort if a positive response is provided
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
