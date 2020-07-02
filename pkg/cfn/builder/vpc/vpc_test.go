package vpc_test

import (
	"testing"

	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/oslokommune/okctl/pkg/cfn/builder/vpc"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
)

func TestVPC(t *testing.T) {
	builder := vpc.New("test", "test", "192.168.0.0/20", "eu-west-1")

	err := builder.Build()
	assert.NoError(t, err)

	resources := builder.Resources()

	template := cloudformation.NewTemplate()
	for _, resource := range resources {
		template.Resources[resource.Name()] = resource.Resource()
	}

	got, err := template.YAML()
	assert.NoError(t, err)

	g := goldie.New(t)
	g.Assert(t, "vpc-cloudformation.yaml", got)
}
