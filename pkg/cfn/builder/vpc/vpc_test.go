package vpc_test

import (
	"testing"

	"github.com/oslokommune/okctl/pkg/cfn/builder/vpc"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
)

func TestVPC(t *testing.T) {
	builder := vpc.New("test", "test", "192.168.0.0/20", "eu-west-1")

	got, err := builder.Build()
	assert.NoError(t, err)

	g := goldie.New(t)
	g.Assert(t, "vpc-cloudformation.yaml", got)
}
