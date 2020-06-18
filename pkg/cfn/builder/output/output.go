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

type joined struct {
	name   string
	values []string
}

func (j *joined) NamedOutputs() map[string]map[string]interface{} {
	return map[string]map[string]interface{}{
		j.Name(): j.Outputs(),
	}
}

func (j *joined) Outputs() map[string]interface{} {
	return map[string]interface{}{
		"Value": cloudformation.Join(",", j.values),
	}
}

func (j *joined) Name() string {
	return j.name
}

func (j *joined) Add(v ...string) *joined {
	j.values = append(j.values, v...)

	return j
}

func Joined(name string) *joined {
	return &joined{
		name: name,
	}
}

type value struct {
	name  string
	value string
}

func (v *value) NamedOutputs() map[string]map[string]interface{} {
	return map[string]map[string]interface{}{
		v.Name(): v.Outputs(),
	}
}

func (v *value) Outputs() map[string]interface{} {
	return map[string]interface{}{
		"Value": v.value,
	}
}

func (v *value) Name() string {
	return v.name
}

func Value(name, v string) *value {
	return &value{
		name:  name,
		value: v,
	}
}
