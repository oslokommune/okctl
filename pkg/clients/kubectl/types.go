package kubectl

import "io"

// Resource defines data required to identify a Kubernetes resource
type Resource struct {
	// Namespace defines the namespace where the resource resides
	Namespace string
	// Kind defines the resource kind to patch
	Kind string
	// Name defines the name of the resource to patch
	Name string
}

// PatchOpts defines required data for patch operations on Kubernetes resources
type PatchOpts struct {
	Resource
	// Patch defines the actual patch to apply
	Patch io.Reader
}

// Client defines functionality expected of a kubectl client
type Client interface {
	// Apply knows how to apply a manifest
	Apply(manifest io.Reader) error
	// Delete knows how to delete a manifest
	Delete(manifest io.Reader) error
	// Patch knows how to patch resources
	Patch(PatchOpts) error
	// Exists knows how to check the existence of a resource
	Exists(Resource) (bool, error)
}
