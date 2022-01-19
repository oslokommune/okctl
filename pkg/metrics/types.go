package metrics

import (
	"io"
	"net/url"

	"github.com/oslokommune/okctl-metrics-service/pkg/endpoints/metrics"
	"github.com/oslokommune/okctl-metrics-service/pkg/endpoints/metrics/types"
)

// Event shadows okctl metrics service Event type
type Event = metrics.Event

type (
	// Category defines the category of an event
	Category = types.Category
	// Action defines the action of an event
	Action = types.Action
)

type context struct {
	WarningOut io.Writer
	UserAgent  string
	APIURL     url.URL
}
