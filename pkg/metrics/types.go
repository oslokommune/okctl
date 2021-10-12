// Package metrics handles metric management
package metrics

import (
	"context"

	"github.com/oslokommune/okctl-metrics-service/pkg/endpoints/metrics"
)

// Client knows how to publish metrics to a metrics backend
type Client interface {
	Publish(context.Context, metrics.Event) error
}
