// Package cfn defines interfaces for creating cloud formation
package cfn

import (
	"github.com/awslabs/goformation/v4/cloudformation"
)

// ResourceNameReferencer knows how to name, referencer and return a cloud formation resource
type ResourceNameReferencer interface {
	Resourcer
	Namer
	Referencer
}

// NameReferencer knows how to name and reference a cloud formation resource
type NameReferencer interface {
	Namer
	Referencer
}

// ResourceNamer knows how to name and return a cloud formation resource
type ResourceNamer interface {
	Resourcer
	Namer
}

// Resourcer knows how to return a cloud formation resource
type Resourcer interface {
	Resource() cloudformation.Resource
}

// Namer knows how to name a cloud formation resource
type Namer interface {
	Name() string
}

// Referencer knows how to create an intrinsic ref to a resource
type Referencer interface {
	Ref() string
}

// Outputer knows how to create cloud formation outputs
type Outputer interface {
	NamedOutputs() map[string]map[string]interface{}
}

// Builder knows how to create a cloud formation stack
type Builder interface {
	Build() error
	StackName() string
	Outputs() []Outputer
	Resources() []ResourceNamer
}
