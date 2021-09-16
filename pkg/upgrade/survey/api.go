// Package survey knows how to get input from the user
package survey

import (
	"fmt"
	"io"

	"github.com/AlecAivazis/survey/v2"
)

// Surveyor knows how to get input from the user
type Surveyor interface {
	PromptUser(message string) (bool, error)
}

// TerminalSurveyor knows how to get input from the user via a terminal
type TerminalSurveyor struct {
	out         io.Writer
	autoConfirm bool
}

// PromptUser prompts the user if they want to continue
func (s TerminalSurveyor) PromptUser(message string) (bool, error) {
	if s.autoConfirm {
		return true, nil
	}

	answer := false
	prompt := &survey.Confirm{
		Message: message,
	}

	err := survey.AskOne(prompt, &answer)
	if err != nil {
		return false, err
	}

	_, _ = fmt.Fprintln(s.out, "")

	return answer, nil
}

// NewTerminalSurveyor creates a new TerminalSurveyor
func NewTerminalSurveyor(out io.Writer, autoConfirm bool) TerminalSurveyor {
	return TerminalSurveyor{
		out,
		autoConfirm,
	}
}
