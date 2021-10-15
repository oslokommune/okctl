// Package metrics exposes the main metrics api
package metrics

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
)

var ctx = context{ // nolint:gochecknoglobals
	WarningOut: os.Stdout,
	UserAgent:  "okctl",
	APIURL:     url.URL{Scheme: "https", Host: "metrics.kjoremiljo.oslo.systems"},
}

// Publish sends a metric event to the okctl-metrics-service
func Publish(event Event) {
	err := publishE(event)
	if err != nil {
		fmt.Fprintf(ctx.WarningOut, "Warning: %s", err.Error())
	}
}

// publishE sends a metric event to the okctl-metrics-service
func publishE(event Event) error {
	err := event.Validate()
	if err != nil {
		return fmt.Errorf("validating metrics event: %w", err)
	}

	rawEvent, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshalling event as JSON: %w", err)
	}

	err = postEvent(ctx.APIURL, ctx.UserAgent, rawEvent)
	if err != nil {
		return fmt.Errorf("POSTing event: %w", err)
	}

	return nil
}

// SetUserAgent configures the user agent of the Publish POST request
func SetUserAgent(newUserAgent string) {
	ctx.UserAgent = newUserAgent
}

// SetMetricsOut configures where to write error messages related to metrics publishing
func SetMetricsOut(writer io.Writer) {
	ctx.WarningOut = writer
}
