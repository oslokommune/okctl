package cfn

import (
	"fmt"

	"github.com/awslabs/goformation/v4/cloudformation"
)

// Composer is invoked by the builder to retrieve
// the cloud formation components
type Composer interface {
	Compose() (*Composition, error)
}

// Builder stores state required for rendering a cloud
// formation template
type Builder struct {
	composer Composer
	template *cloudformation.Template
}

// Composition contains the cloud formation components as
// constructed by a composer
type Composition struct {
	Outputs   []StackOutputer
	Resources []ResourceNamer
	Transform *cloudformation.Transform
}

// New returns a VPC builder
func New(composer Composer) *Builder {
	return &Builder{
		composer: composer,
		template: cloudformation.NewTemplate(),
	}
}

// Build collects all resources and outputs from a composer
// adds them to a cloud formation template and renders the
// template body
func (b *Builder) Build() ([]byte, error) {
	c, err := b.composer.Compose()
	if err != nil {
		return nil, err
	}

	for _, output := range c.Outputs {
		for key, value := range output.NamedOutputs() {
			if _, hasKey := b.template.Outputs[key]; hasKey {
				return nil, fmt.Errorf("already have output with name: %s", key)
			}

			b.template.Outputs[key] = value
		}
	}

	for _, resource := range c.Resources {
		if _, hasKey := b.template.Resources[resource.Name()]; hasKey {
			return nil, fmt.Errorf("already have resource with name: %s", resource.Name())
		}

		b.template.Resources[resource.Name()] = resource.Resource()
	}

	b.template.Transform = c.Transform

	return b.template.YAML()
}

// Ensure that the VPC builder implements the Builder interface
var _ StackBuilder = &Builder{}
