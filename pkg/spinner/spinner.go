package spinner

import (
	"fmt"
	"github.com/oslokommune/okctl/pkg/config/constant"
	"io"
	"time"

	"github.com/theckman/yacspin"
)

const updateFrequency = 100 * time.Millisecond

// Spinner defines the available operations
type Spinner interface {
	Start(component string) error
	Stop() error
	Pause() error
	Unpause() error
	SubSpinner() Spinner
}

// spinner contains state for the spinner
type spinner struct {
	spinner   *yacspin.Spinner
	parent    *spinner
	exit      chan struct{}
	component string
}

// New creates a new spinner
func New(suffix string, out io.Writer) (Spinner, error) {
	cfg := yacspin.Config{
		Frequency:       updateFrequency,
		CharSet:         yacspin.CharSets[59],
		Suffix:          " " + suffix,
		SuffixAutoColon: true,
		StopCharacter:   "✓",
		StopColors:      []string{"fgGreen"},
		Writer:          out,
	}

	s, err := yacspin.New(cfg)
	if err != nil {
		return nil, fmt.Errorf(constant.CreateSpinnerError, err)
	}

	return &spinner{
		spinner: s,
	}, nil
}

func (s *spinner) String() string {
	return s.component
}

func (s *spinner) isChild() bool {
	return s.parent != nil
}

// SubSpinner returns a sub spinner
func (s *spinner) SubSpinner() Spinner {
	return &spinner{
		spinner: s.spinner,
		parent:  s,
	}
}

func (s *spinner) Start(component string) error {
	s.component = component

	if s.isChild() {
		s.spinner.Message(fmt.Sprintf("%s %s", s.parent.component, s.component))

		return nil
	}

	s.spinner.Message(s.component)

	return s.spinner.Start()
}

func (s *spinner) Stop() error {
	if s.isChild() {
		s.spinner.Message(s.parent.component)

		return nil
	}

	if s.exit != nil {
		close(s.exit)
	}

	return s.spinner.Stop()
}

func (s *spinner) Pause() error {
	return s.spinner.Pause()
}

func (s *spinner) Unpause() error {
	return s.spinner.Unpause()
}
