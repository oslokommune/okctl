package kubectl

import "io"

// Client defines functionality expected of a kubectl client
type Client interface {
	// Apply knows how to apply a manifest
	Apply(manifest io.Reader) error
	// Delete knows how to delete a manifest
	Delete(manifest io.Reader) error
}
