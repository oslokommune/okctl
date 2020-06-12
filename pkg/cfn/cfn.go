package cfn

import (
	"github.com/awslabs/goformation/v4/cloudformation"
)

type ResourceNameReferencer interface {
	Resourcer
	Namer
	Referencer
}

type NameReferencer interface {
	Namer
	Referencer
}

type ResourceNamer interface {
	Resourcer
	Namer
}

type Resourcer interface {
	Resource() cloudformation.Resource
}

type Namer interface {
	Name() string
}

type Referencer interface {
	Ref() string
}

type Builder interface {
	Build() ([]ResourceNamer, error)
	StackName() string
}
