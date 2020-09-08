// Package store provides interfaces for more easily interacting with
// a storage layer
package store

// PreProcessed contains information about how the data
// has been processed.
type PreProcessed struct {
	Format string
	Data   []byte
}

// PreProcessor defines a type for preprocessing of the
// data before it is stored
type PreProcessor interface {
	PreProcess(data interface{}) (*PreProcessed, error)
	Type() string
}

// PostProcessed contains information about how the data
// has been processed.
type PostProcessed struct {
	Format string
	Data   interface{}
}

// PostProcessor defines a type for postprocessing of
// the data before it is returned
type PostProcessor interface {
	PostProcess(into interface{}, data []byte) (*PostProcessed, error)
	Type() string
}

// OperationOption provides a means to modify the behavior
// of an operation, by returning a type or interface that
// can be evaluated by the implementations.
type OperationOption func() interface{}

// Alterer defines the interface for performing a run-time
// alteration of the underlying store
type Alterer interface {
	Alter(implementation interface{}) error
	Type() string
}

// AddStoreStruct contains the fields required for adding
// a StoreStruct operation, useful when operations are
// created dynamically from a slice or map
type AddStoreStruct struct {
	Name         string
	Data         interface{}
	PreProcessor PreProcessor
	Options      []OperationOption
}

// AddStoreBytes contains the fields required for adding
// a StoreBytes operation, useful when operations are
// created dynamically from a slice or map
type AddStoreBytes struct {
	Name    string
	Data    []byte
	Options []OperationOption
}

// AddGetStruct contains the fields required for adding
// a GetStruct operation, useful when operations are
// created dynamically from a slice or map
type AddGetStruct struct {
	Name          string
	Into          interface{}
	PostProcessor PostProcessor
	Options       []OperationOption
}

// AddGetBytes contains the fields required for adding
// a GetBytes operation, useful when operations are
// created dynamically from a slice or map
type AddGetBytes struct {
	Name     string
	Callback GetBytesCallback
	Options  []OperationOption
}

// GetBytesCallback allows the receiver to process the
// data immediately without having to go via the report
type GetBytesCallback func(name string, data []byte) error

// ProcessGetStruct defines the function interface that allows for inline
// processing of a recently retrieved struct.
type ProcessGetStruct func(data interface{}, operations Operations) error

// Operations defines the functions that can be chained
// before being executed
type Operations interface {
	// AlterStore makes it possible to perform a modification of the store
	// during processing.
	AlterStore(alterers ...Alterer) Operations

	// StoreStruct will preprocess the data with the given preprocessor
	// and then store the processed data under the given name, how this
	// data is stored is given by the underlying implementation.
	// This method is chainable, so multiple actions may be performed.
	StoreStruct(name string, data interface{}, preProcessor PreProcessor, options ...OperationOption) Operations

	// AddStoreStruct makes it possible to add operations when
	// iterating over a map or slice.
	AddStoreStruct(operations ...AddStoreStruct) Operations

	// StoreBytes stores the provided data under the given name, how this
	// data is stored is given by the underlying implementation.
	// This method is chainable, so multiple actions may be performed.
	StoreBytes(name string, data []byte, options ...OperationOption) Operations

	// AddStoreBytes makes it possible to add operations when
	// iterating over a map or slice.
	AddStoreBytes(operations ...AddStoreBytes) Operations

	// Remove will delete the data under the given name, how this data is
	// removed is given by the underlying implementation
	// This method is chainable, so multiple actions may be performed.
	Remove(name string, options ...OperationOption) Operations

	// GetStruct will read the stored data and process it using the
	// provided post processor, the result will be stored in the report
	// and can be accessed afterwards.
	GetStruct(name string, into interface{}, postProcessor PostProcessor, options ...OperationOption) Operations

	// AddGetStruct makes it possible to add operations when
	// iterating over a map or slice
	AddGetStruct(operations ...AddGetStruct) Operations

	// ProcessGetStruct provides inline processing of a recently
	// retrieved and processed struct.
	ProcessGetStruct(name string, storeStruct ProcessGetStruct) Operations

	// GetBytes will read the stored data and store the result
	// in the report under the given name, and can be accessed
	// afterwards. Otherwise, it is possible to provided a
	// callback, which will be invoked if it is not nil.
	GetBytes(name string, callback GetBytesCallback, options ...OperationOption) Operations

	// AddGetBytes makes it possible to add operations when
	// iterating over a map or slice
	AddGetBytes(operations ...AddGetBytes) Operations

	// Do initiates the chain of operations and will either return
	// an error or a report containing the actions that were performed.
	Do() (*Report, error)
}

// Report contains information about what actions
// were performed
type Report struct {
	Type          string
	Configuration string
	Actions       []Action
	Data          map[string]interface{}
}

// Action contains information about the performed
// task
type Action struct {
	Name        string
	Path        string
	Type        string
	Description string
}
