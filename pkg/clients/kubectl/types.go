package kubectl

import "io"

// Client defines functionality expected of a kubectl client
type Client interface {
	Apply(manifest io.Reader) error
}
