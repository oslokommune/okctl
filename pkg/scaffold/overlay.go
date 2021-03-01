package scaffold

import (
	"encoding/json"
)

type OperationType int

const (
	OperationTypeReplace OperationType = iota
)

func operationTypeToString(t OperationType) string {
	switch t {
	case OperationTypeReplace:
		return "replace"
	default:
		return "n/a"
	}
}

type Operation struct {
	Type  OperationType `json:"op"`
	Path  string        `json:"path"`
	Value string        `json:"value"`
}

type Patch struct {
	Operations []Operation `json:",inline"`
}

func (p Patch) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.Operations)
}
