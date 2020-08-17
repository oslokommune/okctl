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

// StackOutputer ensures that the receiver will be
// returned a set of named output values that can
// be associated with a cloud formation stack and extracted
// after stack creation:
// - https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/outputs-section-structure.html
type StackOutputer interface {
	NamedOutputs() map[string]map[string]interface{}
}

// StackBuilder knows how to create a cloud formation stack
type StackBuilder interface {
	Build() ([]byte, error)
}
