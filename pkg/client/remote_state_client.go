package client

import (
	"io"

	"github.com/oslokommune/okctl/pkg/api"
)

// Stater defines functionality for transferring remote state back and forth
type Stater interface {
	// Upload knows how to upload state to a remote location
	Upload(clusterID api.ID, reader io.Reader) error
	// Download knows how to download state from a remote location
	Download(clusterID api.ID) (io.Reader, error)
}

// Locker defines functionality for ensuring only one client mutates state at a time
type Locker interface {
	// AcquireStateLock knows how to activate a locking mechanism preventing others from mutating state
	AcquireStateLock(clusterID api.ID) error
	// ReleaseStateLock knows how to deactivate a locking mechanism allowing others to mutate state
	ReleaseStateLock(clusterID api.ID) error
}

// RemoteStateService defines expected functionality in a remote state service implementation
type RemoteStateService interface {
	// Purge knows how to tear down all traces of remote state
	Purge(clusterID api.ID) error
	Stater
	Locker
}
