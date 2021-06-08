package common

import (
	"errors"

	"github.com/asdine/storm/v3"
	"github.com/oslokommune/okctl/pkg/controller/common/resourcetree"
)

// IsNotFound determines if an error is of the type storm ErrNotFound
func IsNotFound(_ interface{}, err error) bool {
	return errors.Is(err, storm.ErrNotFound)
}

// SetAllNodesAbsent sets the state as absent
func SetAllNodesAbsent(receiver *resourcetree.ResourceNode) {
	receiver.State = resourcetree.ResourceNodeStateAbsent
}

// BoolToState converts a boolean to a resourcetree.ResourceNodeState
func BoolToState(present bool) resourcetree.ResourceNodeState {
	if present {
		return resourcetree.ResourceNodeStatePresent
	}

	return resourcetree.ResourceNodeStateAbsent
}
