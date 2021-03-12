// Package testing contains some helpers for testing the
// cloud formation templates of components
package testing

import (
	"testing"

	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/oslokommune/okctl/pkg/cfn"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
)

type nameReferencer struct {
	name string
}

// Name returns the resource name
func (r *nameReferencer) Name() string {
	return r.name
}

// Ref returns a cloud formation ref to the resource
func (r *nameReferencer) Ref() string {
	return cloudformation.Ref(r.Name())
}

// NewNameReferencer returns an initialised name referencer
func NewNameReferencer(name string) cfn.NameReferencer {
	return &nameReferencer{
		name: name,
	}
}

// CloudFormationTemplateTestCase contains the required fields
// for setting up a test
type CloudFormationTemplateTestCase struct {
	Name    string
	Golden  string
	Content cfn.ResourceNameOutputReferencer
}

// RunTests runs the cloud formation template tests
func RunTests(t *testing.T, testCases []CloudFormationTemplateTestCase) {
	for _, tc := range testCases {
		tc := tc

		t.Run(tc.Name, func(t *testing.T) {
			template := cloudformation.NewTemplate()
			template.Resources = map[string]cloudformation.Resource{
				tc.Content.Name(): tc.Content.Resource(),
			}

			for name, output := range tc.Content.NamedOutputs() {
				template.Outputs[name] = output
			}

			data, err := template.YAML()
			assert.NoError(t, err)

			g := goldie.New(t)
			g.Assert(t, tc.Golden, data)
		})
	}
}
