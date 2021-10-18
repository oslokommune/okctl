// Package restapi implements a way to publish metrics through a REST interface
package restapi

import (
	"net/url"
)

type client struct {
	apiURL url.URL
}
