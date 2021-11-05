package hooks

import (
	"github.com/oslokommune/okctl/pkg/metrics"
	"github.com/spf13/cobra"
)

// EmitStartCommandExecutionEvent publishes an event indicating the start of a command
func EmitStartCommandExecutionEvent(action metrics.Action) RunEer {
	return func(_ *cobra.Command, _ []string) error {
		metrics.Publish(metrics.Event{
			Category: metrics.CategoryCommandExecution,
			Action:   action,
			Labels: map[string]string{
				metrics.LabelPhaseKey: metrics.LabelPhaseStart,
			},
		})

		return nil
	}
}

// EmitEndCommandExecutionEvent publishes an event indicating the end of a command
func EmitEndCommandExecutionEvent(action metrics.Action) RunEer {
	return func(_ *cobra.Command, _ []string) error {
		metrics.Publish(metrics.Event{
			Category: metrics.CategoryCommandExecution,
			Action:   action,
			Labels: map[string]string{
				metrics.LabelPhaseKey: metrics.LabelPhaseEnd,
			},
		})

		return nil
	}
}
