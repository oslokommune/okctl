package cfn

import (
	"fmt"

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

type manager struct {
	t *cloudformation.Template
}

func New() *manager {
	return &manager{
		t: cloudformation.NewTemplate(),
	}
}

func (d *manager) Add(resources ...ResourceNamer) error {
	for _, resource := range resources {
		if _, hasKey := d.t.Resources[resource.Name()]; hasKey {
			return fmt.Errorf("already have resource with name: %s", resource.Name())
		}

		d.t.Resources[resource.Name()] = resource.Resource()
	}

	return nil
}

func (d *manager) YAML() ([]byte, error) {
	return d.t.YAML()
}
