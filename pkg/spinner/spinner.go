package spinner

import (
	"fmt"
	"io"
	"time"

	"github.com/hako/durafmt"
	"github.com/theckman/yacspin"
)

// New creates a new spinner
func New(suffix string, out io.Writer) (*yacspin.Spinner, error) {
	cfg := yacspin.Config{
		Frequency:       100 * time.Millisecond, // nolint: gomnd
		CharSet:         yacspin.CharSets[59],
		Suffix:          " " + suffix,
		SuffixAutoColon: true,
		StopCharacter:   "✓",
		StopColors:      []string{"fgGreen"},
		Writer:          out,
	}

	spinner, err := yacspin.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create spinner")
	}

	return spinner, nil
}

// Timer creates and associated a timer with the spinner
func Timer(component string, spinner *yacspin.Spinner) chan struct{} {
	exit := make(chan struct{})

	go func(ch chan struct{}, start time.Time) {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ch:
				return
			case <-ticker.C:
				spinner.Message(component + " (elapsed: " + durafmt.Parse(time.Since(start)).LimitFirstN(2).String() + ")") // nolint: gomnd
			}
		}
	}(exit, time.Now())

	return exit
}
