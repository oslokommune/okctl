package reconciler

import (
	"github.com/oslokommune/okctl/pkg/controller/resourcetree"
)

// noopMetadata contains data known at initialization. Usually information from the desired state
type noopMetadata struct{}

// noopResourceState contains data that potentially can only be known at runtime. E.g.: state only known after an
// external resource has been created
type noopResourceState struct{}

// NoopReconciler handles reconciliation for dummy nodes (e.g. the root node) and acts as a template for other
// reconcilers
type NoopReconciler struct{}

// SetCommonMetadata knows how to store common metadata on the reconciler. This should do nothing if common metadata is
// not needed
func (receiver *NoopReconciler) SetCommonMetadata(_ *resourcetree.CommonMetadata) {}

// Reconcile knows how to create, update and delete the relevant resource
func (receiver *NoopReconciler) Reconcile(node *resourcetree.ResourceNode) (*ReconcilationResult, error) {
	//metadata, ok := node.Metadata.(noopMetadata)
	//if !ok {
	//	return nil, errors.New("could not cast Noop metadata")
	//}
	//
	//state, ok := node.ResourceState.(noopResourceState)
	//if !ok {
	//	return nil, errors.New("could not cast Noop resource state")
	//}

	switch node.State {
	case resourcetree.ResourceNodeStatePresent:
		// Create a resource
	case resourcetree.ResourceNodeStateAbsent:
		// Delete a resource
	}

	return &ReconcilationResult{Requeue: false}, nil
}
