// Package jsonpatch provides some convenience structs for building
// a type safe patch
package jsonpatch

import (
	"encoding/json"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// OperationType represents one of the json patch operation types defines in https://tools.ietf.org/html/rfc6902#section-4
type OperationType string

func (receiver OperationType) String() string {
	return string(receiver)
}

// Validate ensures the OperationType is a valid JSON patch type
func (receiver OperationType) Validate() error {
	return validation.Validate(&receiver,
		validation.In(
			OperationTypeAdd,
			OperationTypeRemove,
			OperationTypeReplace,
			OperationTypeMove,
			OperationTypeCopy,
			OperationTypeTest,
		),
	)
}

const (
	// OperationTypeAdd represents an add operation
	OperationTypeAdd OperationType = "add"
	// OperationTypeRemove represents a remove operation
	OperationTypeRemove OperationType = "remove"
	// OperationTypeReplace represents a replace operation
	OperationTypeReplace OperationType = "replace"
	// OperationTypeMove represents a move operation
	OperationTypeMove OperationType = "move"
	// OperationTypeCopy represents a copy operation
	OperationTypeCopy OperationType = "copy"
	// OperationTypeTest represents a test operation
	OperationTypeTest OperationType = "test"
)

// Operation represents a single patch operation, meaning an action on a kubernetes resource attribute
type Operation struct {
	Type  OperationType `json:"op"`
	Path  string        `json:"path"`
	Value interface{}   `json:"value"`
}

// Equals knows how to determine if two operations are equal
func (o Operation) Equals(target Operation) bool {
	if o.Type != target.Type {
		return false
	}

	if o.Path != target.Path {
		return false
	}

	if o.Value != target.Value {
		return false
	}

	return true
}

// Patch represents a kustomize patch.json file containing a list of patch operations
type Patch struct {
	Operations []Operation `json:",inline"`
}

// Inline provides a convenient way of inlining
// pre-serialised json
type Inline struct {
	Data []byte `json:",inline"`
}

// MarshalJSON implements the json encode marshaller interface
func (b *Inline) MarshalJSON() ([]byte, error) {
	return b.Data, nil
}

// New initializes a Patch struct
func New() *Patch {
	return &Patch{
		Operations: []Operation{},
	}
}

// Add adds a patch operation to the patch
func (p *Patch) Add(o ...Operation) *Patch {
	p.Operations = append(p.Operations, o...)

	return p
}

// HasOperation determines if a patch has a specific operation
func (p *Patch) HasOperation(operation Operation) bool {
	for _, op := range p.Operations {
		if op.Equals(operation) {
			return true
		}
	}

	return false
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
			Type:  string(operation.Type),
			Path:  operation.Path,
			Value: operation.Value,
		}
	}

	return json.Marshal(patch)
}
