package resources

import (
	"encoding/json"
)

// OperationType represents one of the json patch operation types defines in https://tools.ietf.org/html/rfc6902#section-4
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

// Operation represents a single patch operation, meaning an action on a kubernetes resource attribute
type Operation struct {
	Type  OperationType `json:"op"`
	Path  string        `json:"path"`
	Value interface{}   `json:"value"`
}

// Patch represents a kustomize patch.json file containing a list of patch operations
type Patch struct {
	Operations []Operation `json:",inline"`
}

// NewPatch initializes a Patch struct
func NewPatch() *Patch {
	return &Patch{
		Operations: []Operation{},
	}
}

// AddOperations adds a patch operation to the patch
func (p *Patch) AddOperations(o ...Operation) {
	p.Operations = append(p.Operations, o...)
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
