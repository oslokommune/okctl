// Package metrics exposes the main metrics api
package metrics

import (
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
// This function wraps the main publish function and prints out a warning if there is an error. It seems drastic to
// crash the program due to a metrics fault
func Publish(event Event) {
	err := publishE(event)
	if err != nil {
		fmt.Fprintf(ctx.WarningOut, "Warning: %s", err.Error())
	}
}

// SetAPIURL configures the API URL for the Publish POST request
func SetAPIURL(newAPIURL string) error {
	parsedURL, err := url.Parse(newAPIURL)
	if err != nil {
		return fmt.Errorf("parsing URL: %w", err)
	}

	ctx.APIURL = *parsedURL

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
