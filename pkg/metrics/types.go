package metrics

import (
	"io"
	"net/url"

	metricsapi "github.com/oslokommune/okctl-metrics-service/pkg/endpoints/metrics"
)

// Event shadows okctl metrics service Event type
type Event = metricsapi.Event

type (
	// Category defines the category of an event
	Category = metricsapi.Category
	// Action defines the action of an event
	Action = metricsapi.Action
)

type context struct {
	WarningOut io.Writer
	UserAgent  string
	APIURL     url.URL
}

const (
	// LabelStart indicates the start of something
	LabelStart = "start"
	// LabelEnd indicates the end of something
	LabelEnd = "end"
)
