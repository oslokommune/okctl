package main

import "github.com/oslokommune/okctl/pkg/metrics"

func generateStartEvent(action metrics.Action) metrics.Event {
	return metrics.Event{
		Category: metrics.CategoryCommandExecution,
		Action:   action,
		Labels: map[string]string{
			metrics.LabelPhaseKey: metrics.LabelPhaseStart,
		},
	}
}

func generateEndEvent(action metrics.Action) metrics.Event {
	return metrics.Event{
		Category: metrics.CategoryCommandExecution,
		Action:   action,
		Labels: map[string]string{
			metrics.LabelPhaseKey: metrics.LabelPhaseEnd,
		},
	}
}
