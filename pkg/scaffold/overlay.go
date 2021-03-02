package scaffold

import (
	"encoding/json"
)

type OperationType int

const (
	// OperationTypeAdd represents an add operation
	OperationTypeAdd OperationType = iota
	// OperationTypeRemove represents a remove operation
	OperationTypeRemove
	// OperationTypeReplace represents a replace operation
	OperationTypeReplace
	// OperationTypeMove represents a move operation
	OperationTypeMove
	// OperationTypeCopy represents a copy operation
	OperationTypeCopy
	// OperationTypeTest represents a test operation
	OperationTypeTest
)

func operationTypeToString(t OperationType) string {
	switch t {
	case OperationTypeAdd:
		return "add"
	case OperationTypeRemove:
		return "remove"
	case OperationTypeReplace:
		return "replace"
	case OperationTypeMove:
		return "move"
	case OperationTypeCopy:
		return "copy"
	case OperationTypeTest:
		return "test"
	default:
		return "n/a"
	}
}

type Operation struct {
	Type  OperationType `json:"op"`
	Path  string        `json:"path"`
	Value interface{}   `json:"value"`
}

type Patch struct {
	Operations []Operation `json:",inline"`
}

// MarshalJSON knows how to turn a Patch into a kustomize patch.json
func (p Patch) MarshalJSON() ([]byte, error) {
	type serializedOperation struct {
		Type  string      `json:"op"`
		Path  string      `json:"path"`
		Value interface{} `json:"value"`
	}

	patch := make([]serializedOperation, len(p.Operations))

	for index, operation := range p.Operations {
		patch[index] = serializedOperation{
			Type:  operationTypeToString(operation.Type),
			Path:  operation.Path,
			Value: operation.Value,
		}
	}

	return json.Marshal(patch)
}
