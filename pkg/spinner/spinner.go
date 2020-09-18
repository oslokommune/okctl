package spinner

import (
	"fmt"
	"io"
	"time"

	"github.com/hako/durafmt"
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
	spinner *yacspin.Spinner
	parent  *spinner
	exit    chan struct{}
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
		return nil, fmt.Errorf("creating spinner: %w", err)
	}

	return &spinner{
		spinner: s,
	}, nil
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
	if s.isChild() {
		return nil
	}

	s.timer(component)

	return s.spinner.Start()
}

func (s *spinner) Stop() error {
	if s.isChild() {
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

// Timer creates and associates a timer with the spinner
func (s *spinner) timer(component string) {
	exit := make(chan struct{})

	go func(ch chan struct{}, start time.Time) {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ch:
				return
			case <-ticker.C:
				s.spinner.Message(component + " (elapsed: " + durafmt.Parse(time.Since(start)).LimitFirstN(2).String() + ")") // nolint: gomnd
			}
		}
	}(exit, time.Now())

	s.exit = exit
}
