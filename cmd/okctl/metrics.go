package main

import "github.com/oslokommune/okctl/pkg/metrics"

const (
	labelKeyPhase   = "phase"
	labelValueStart = "start"
	labelValueEnd   = "end"
)

func generateStartEvent(action metrics.Action) metrics.Event {
	return metrics.Event{
		Category: metrics.CategoryCommandExecution,
		Action:   action,
		Labels: map[string]string{
			labelKeyPhase: labelValueStart,
		},
	}
}

func generateEndEvent(action metrics.Action) metrics.Event {
	return metrics.Event{
		Category: metrics.CategoryCommandExecution,
		Action:   action,
		Labels: map[string]string{
			labelKeyPhase: labelValueEnd,
		},
	}
}
