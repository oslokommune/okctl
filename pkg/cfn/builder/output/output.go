package output

import (
	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/oslokommune/okctl/pkg/cfn"
)

type Outputer interface {
	cfn.Outputer
	Outputs() map[string]interface{}
	Name() string
}

type Joined struct {
	StoredName string
	Values     []string
}

func (j *Joined) NamedOutputs() map[string]map[string]interface{} {
	return map[string]map[string]interface{}{
		j.Name(): j.Outputs(),
	}
}

func (j *Joined) Outputs() map[string]interface{} {
	return map[string]interface{}{
		"NewValue": cloudformation.Join(",", j.Values),
	}
}

func (j *Joined) Name() string {
	return j.StoredName
}

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

type Value struct {
	StoredName string
	Value      string
}

func (v *Value) NamedOutputs() map[string]map[string]interface{} {
	return map[string]map[string]interface{}{
		v.Name(): v.Outputs(),
	}
}

func (v *Value) Outputs() map[string]interface{} {
	return map[string]interface{}{
		"NewValue": v.Value,
	}
}

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
