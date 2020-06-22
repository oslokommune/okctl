// Package output provides functionality for creating cloud formation output
package output

import (
	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/oslokommune/okctl/pkg/cfn"
)

// Outputer defines the required interface for
// fetching output information
type Outputer interface {
	cfn.Outputer
	Outputs() map[string]interface{}
	Name() string
}

// Joined stores state for creating an intrinsic join
type Joined struct {
	StoredName string
	Values     []string
}

// NamedOutputs returns the named outputs
func (j *Joined) NamedOutputs() map[string]map[string]interface{} {
	return map[string]map[string]interface{}{
		j.Name(): j.Outputs(),
	}
}

// Outputs returns the outputs only
func (j *Joined) Outputs() map[string]interface{} {
	return map[string]interface{}{
		"Value": cloudformation.Join(",", j.Values),
	}
}

// Name returns the name of the output
func (j *Joined) Name() string {
	return j.StoredName
}

// Add a value to the outputs that should be joined
func (j *Joined) Add(v ...string) *Joined {
	j.Values = append(j.Values, v...)

	return j
}

// NewJoined is a helper for creating cloud formation Joined data
// in the output of a cloud formation stack
//
// https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/intrinsic-function-reference-join.html
func NewJoined(name string) *Joined {
	return &Joined{
		StoredName: name,
	}
}

// Value stores the state for creating an output
type Value struct {
	StoredName string
	Value      string
}

// NamedOutputs returns the named cloud formation outputs
func (v *Value) NamedOutputs() map[string]map[string]interface{} {
	return map[string]map[string]interface{}{
		v.Name(): v.Outputs(),
	}
}

// Outputs returns only the cloud formation outputs
func (v *Value) Outputs() map[string]interface{} {
	return map[string]interface{}{
		"Value": v.Value,
	}
}

// Name returns the name given to the output
func (v *Value) Name() string {
	return v.StoredName
}

// NewValue is a helper for creating Value outputs in
// a cloud formation stack
func NewValue(name, v string) *Value {
	return &Value{
		StoredName: name,
		Value:      v,
	}
}
