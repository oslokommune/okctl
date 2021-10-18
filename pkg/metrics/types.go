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
	Category metricsapi.Category
	// Action defines the action of an event
	Action metricsapi.Action
)

// Categories
const (
	// CategoryCluster represents metrics associated with cluster manipulation
	CategoryCluster = metricsapi.CategoryCluster
	// CategoryApplication represents metrics associated with application manipulation
	CategoryApplication = metricsapi.CategoryApplication
)

// Actions
const (
	// ActionScaffold represents scaffolding a resource
	ActionScaffold = metricsapi.ActionScaffold
	// ActionApply represents applying a resource
	ActionApply = metricsapi.ActionApply
	// ActionDelete represents deleting a resource
	ActionDelete = metricsapi.ActionDelete
)

type context struct {
	WarningOut io.Writer
	UserAgent  string
	APIURL     url.URL
}
