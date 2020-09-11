package console

import (
	"fmt"
	"io"
	"time"

	"github.com/logrusorgru/aurora/v3"
	"github.com/oslokommune/okctl/pkg/client/store"

	"github.com/theckman/yacspin"
)

// Console stores the state for writing and closing the progress
type Console struct {
	out     io.Writer
	exit    chan struct{}
	spinner *yacspin.Spinner
}

// New returns an initialised console writer
func New(out io.Writer, exit chan struct{}, spinner *yacspin.Spinner) *Console {
	return &Console{
		out:     out,
		exit:    exit,
		spinner: spinner,
	}
}

// Report writes the content to the provided io.Writer
// nolint: funlen gocognit
func (c *Console) Report(actions []store.Action, component, description string) error {
	close(c.exit)

	if c.exit == nil {
		err := c.spinner.Pause()
		if err != nil {
			return err
		}

		defer func() {
			err = c.spinner.Unpause()
		}()
	} else {
		err := c.spinner.Stop()
		if err != nil {
			return err
		}
	}

	time.Sleep(100 * time.Millisecond) // nolint: gomnd

	_, err := fmt.Fprintf(c.out, "created %s: %s\n", component, description)
	if err != nil {
		return err
	}

	if len(actions) > 0 {
		_, err = fmt.Fprintf(c.out, "storage operations:\n\n")
		if err != nil {
			return err
		}

		for _, a := range actions {
			_, err := fmt.Fprintf(c.out, "\t%s", aurora.Gray(12, a.Type)) // nolint: gomnd
			if err != nil {
				return err
			}

			if len(a.Name) > 0 {
				_, err := fmt.Fprintf(c.out, ": %s", aurora.Blue(a.Name))
				if err != nil {
					return err
				}
			}

			if len(a.Path) > 0 {
				actionPath := a.Path
				if len(a.Path) > 33 { // nolint: gomnd
					actionPath = fmt.Sprintf("...%s", a.Path[len(a.Path)-30:]) // nolint: gomnd
				}

				_, err := fmt.Fprintf(c.out, " (%s)", actionPath)
				if err != nil {
					return err
				}
			}

			_, err = fmt.Fprint(c.out, "\n")
			if err != nil {
				return err
			}
		}

		_, err := fmt.Fprintf(c.out, "\n")
		if err != nil {
			return err
		}
	}

	return nil
}
