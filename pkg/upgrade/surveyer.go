package upgrade

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"io"
)

// Surveyor knows how to get input from the user
type Surveyor interface {
	AskUserIfReady() (bool, error)
}

// TerminalSurveyor knows how to get input from the user via a terminal
type TerminalSurveyor struct {
	out         io.Writer
	autoConfirm bool
}

// AskUserIfReady prompts the user if they want to continue
func (s TerminalSurveyor) AskUserIfReady() (bool, error) {
	if s.autoConfirm {
		return true, nil
	}

	doContinue := false
	prompt := &survey.Confirm{
		Message: "This will upgrade your okctl cluster, are you sure you want to continue?",
	}

	err := survey.AskOne(prompt, &doContinue)
	if err != nil {
		return false, err
	}

	_, _ = fmt.Fprintln(s.out, "")

	return doContinue, nil
}

func NewTerminalSurveyor(out io.Writer, autoConfirm bool) TerminalSurveyor {
	return TerminalSurveyor{
		out,
		autoConfirm,
	}
}
