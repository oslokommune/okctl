// Package mock provides mocks for various components
package mock

import (
	"net"

	"github.com/oslokommune/okctl/pkg/cfn"
	"github.com/oslokommune/okctl/pkg/cfn/builder/output"
	"github.com/oslokommune/okctl/pkg/cfn/components/cluster"
	"github.com/oslokommune/okctl/pkg/cfn/components/vpc"
)

// DefaultStackName is the default cloud formation stack name
const DefaultStackName = "myStack"

// Vpc returns a valid cloud formation specification of a VPC
func Vpc() *vpc.VPC {
	_, network, _ := net.ParseCIDR("192.168.0.0/20")
	return vpc.New(cluster.New("test", "test"), network)
}

// Resources returns valid resources
func Resources() []cfn.ResourceNamer {
	return []cfn.ResourceNamer{
		Vpc(),
	}
}

// Outputs returns valid outputs
func Outputs() []cfn.Outputer {
	v := Vpc()

	return []cfn.Outputer{
		output.NewValue(v.Name(), v.Ref()),
	}
}

type builder struct {
	buildFn     func() error
	outputsFn   func() []cfn.Outputer
	resourcesFn func() []cfn.ResourceNamer
	stackNameFn func() string
}

// Build returns the mocked build response
func (b *builder) Build() error {
	return b.buildFn()
}

// Outputs returns the mocked outputs response
func (b *builder) Outputs() []cfn.Outputer {
	return b.outputsFn()
}

// Resources returns the mocked resources response
func (b *builder) Resources() []cfn.ResourceNamer {
	return b.resourcesFn()
}

// StackName returns the mocked stack name
func (b *builder) StackName() string {
	return b.stackNameFn()
}

// NewGoodBuilder returns a mocked builder with success set on all
func NewGoodBuilder() cfn.Builder {
	return &builder{
		buildFn: func() error {
			return nil
		},
		resourcesFn: Resources,
		outputsFn:   Outputs,
		stackNameFn: func() string {
			return DefaultStackName
		},
	}
}
