package kubectl

import "io"

// Client defines functionality expected of a kubectl client
type Client interface {
	// Apply knows how to apply manifests to a cluster
	Apply(manifest io.Reader) error
}
