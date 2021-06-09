package common

import (
	"errors"

	"github.com/asdine/storm/v3"
	"github.com/oslokommune/okctl/pkg/controller/common/dependencytree"
)

// IsNotFound determines if an error is of the type storm ErrNotFound
func IsNotFound(_ interface{}, err error) bool {
	return errors.Is(err, storm.ErrNotFound)
}

// SetAllNodesAbsent sets the state as absent
func SetAllNodesAbsent(receiver *dependencytree.Node) {
	receiver.State = dependencytree.NodeStateAbsent
}

// BoolToState converts a boolean to a dependencytree.NodeState
func BoolToState(present bool) dependencytree.NodeState {
	if present {
		return dependencytree.NodeStatePresent
	}

	return dependencytree.NodeStateAbsent
}
