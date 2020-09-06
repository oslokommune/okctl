// Package store provides interfaces for more easily interacting with
// a storage layer
package store

import (
	"container/list"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/spf13/afero"
)

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

type toJSON struct{}

func (t *toJSON) PreProcess(data interface{}) (*PreProcessed, error) {
	d, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal as json: %w", err)
	}

	return &PreProcessed{
		Format: t.Type(),
		Data:   d,
	}, nil
}

func (t *toJSON) Type() string {
	return "json"
}

// ToJSON provides a preprocessor that marshals a struct to json format
func ToJSON() PreProcessor {
	return &toJSON{}
}

type fromJSON struct{}

func (f *fromJSON) PostProcess(into interface{}, data []byte) (*PostProcessed, error) {
	err := json.Unmarshal(data, into)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal as json: %w", err)
	}

	return &PostProcessed{
		Format: f.Type(),
		Data:   into,
	}, nil
}

func (f *fromJSON) Type() string {
	return "json"
}

// FromJSON provides a postprocessor that unmarshals from yaml into a struct
func FromJSON() PostProcessor {
	return &fromJSON{}
}

type toYAML struct{}

func (t *toYAML) PreProcess(data interface{}) (*PreProcessed, error) {
	d, err := yaml.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal as yaml: %w", err)
	}

	return &PreProcessed{
		Format: t.Type(),
		Data:   d,
	}, nil
}

func (t *toYAML) Type() string {
	return "yaml"
}

// ToYAML provides a preprocessor that marshals a struct to yaml format
func ToYAML() PreProcessor {
	return &toYAML{}
}

type fromYAML struct{}

func (f *fromYAML) PostProcess(into interface{}, data []byte) (*PostProcessed, error) {
	err := yaml.Unmarshal(data, into)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal as yaml: %w", err)
	}

	return &PostProcessed{
		Format: f.Type(),
		Data:   into,
	}, nil
}

func (f *fromYAML) Type() string {
	return "yaml"
}

// FromYAML provides a postprocessor that unmarshals yaml data ino a struct
func FromYAML() PostProcessor {
	return &fromYAML{}
}

type readType string

const (
	readTypeStruct readType = "struct"
	readTypeBytes  readType = "bytes"
)

// fileSystemWork represents the result of the work
type fileSystemWork struct {
	ReadType readType
	Path     string
	Bytes    []byte
	Struct   interface{}
	Fn       func() error
}

type taskType string

const (
	taskTypeGetBytes         = "GetBytes"
	taskTypeGetStruct        = "GetStruct"
	taskTypeStoreStruct      = "StoreStruct"
	taskTypeStoreBytes       = "StoreBytes"
	taskTypeAlter            = "Alter"
	taskTypeRemove           = "Remove"
	taskTypeProcessGetStruct = "ProcessGetStruct"
)

// fileSystemTask represents a task to be performed
type fileSystemTask struct {
	Name      string
	Type      taskType
	Processor string
	Work      *fileSystemWork
}

// fileSystem contains state for the file system
// storage implementation
type fileSystem struct {
	BaseDir           string
	CreateDirectories bool
	OverwriteExisting bool
	fs                *afero.Afero
	report            *Report
	tasks             *list.List
}

// FileSystemOption defines how options to the file system
// can be provided
type FileSystemOption func(fs *fileSystem)

// FileSystemCreateDirectories ensures that directories are always
// created before storage operations are performed if set to
// true
func FileSystemCreateDirectories(v bool) FileSystemOption {
	return func(fs *fileSystem) {
		fs.CreateDirectories = v
	}
}

// FileSystemOverwriteExisting will ensure that existing files
// are overwritten if set to true
func FileSystemOverwriteExisting(v bool) FileSystemOption {
	return func(fs *fileSystem) {
		fs.OverwriteExisting = v
	}
}

// Ensure that fileSystem implements the Operations interface
var _ Operations = &fileSystem{}

// NewFileSystem returns an initialised filesystem store
func NewFileSystem(baseDir string, fs *afero.Afero, options ...FileSystemOption) Operations {
	f := &fileSystem{
		BaseDir:           baseDir,
		CreateDirectories: true,
		OverwriteExisting: true,
		fs:                fs,
		tasks:             list.New(),
	}

	for _, option := range options {
		option(f)
	}

	f.report = &Report{
		Type: "FileSystem",
		Configuration: fmt.Sprintf(
			"CreateDirectories: %t\nOverWriteExisting: %t\n",
			f.CreateDirectories,
			f.OverwriteExisting,
		),
		Data: map[string]interface{}{},
	}

	return f
}

// WithFilePermissionsMode will set the permissions of the file
func WithFilePermissionsMode(mode os.FileMode) OperationOption {
	return func() interface{} {
		return mode
	}
}

// SetBaseDir to the provided value
func SetBaseDir(baseDir string) Alterer {
	return &fileSystemSetBaseDir{
		baseDir: baseDir,
	}
}

func (f *fileSystemSetBaseDir) Alter(implementation interface{}) error {
	to, ok := implementation.(*fileSystem)
	if !ok {
		return fmt.Errorf("could not cast implemenation to *fileSystem")
	}

	to.BaseDir = f.baseDir

	return nil
}

func (f *fileSystemSetBaseDir) Type() string {
	return "SetBaseDir"
}

type fileSystemSetBaseDir struct {
	baseDir string
}

func (f *fileSystem) ProcessGetStruct(name string, storeStruct ProcessGetStruct) Operations {
	f.tasks.PushBack(&fileSystemTask{
		Name: name,
		Type: taskTypeProcessGetStruct,
		Work: f.processGetStruct(name, storeStruct),
	})

	return f
}

func (f *fileSystem) processGetStruct(name string, process ProcessGetStruct) *fileSystemWork {
	work := &fileSystemWork{}
	work.Fn = func() error {
		data, hasKey := f.report.Data[name]
		if !hasKey {
			return fmt.Errorf("failed to retrieve struct: no such name '%s'", name)
		}

		err := process(data, f)
		if err != nil {
			return fmt.Errorf("failed to process struct: %w", err)
		}

		return nil
	}

	return work
}

func (f *fileSystem) AddGetStruct(operations ...AddGetStruct) Operations {
	for _, o := range operations {
		f.GetStruct(o.Name, o.Into, o.PostProcessor, o.Options...)
	}

	return f
}

func (f *fileSystem) AddGetBytes(operations ...AddGetBytes) Operations {
	for _, o := range operations {
		f.GetBytes(o.Name, nil, o.Options...)
	}

	return f
}

func (f *fileSystem) GetStruct(name string, into interface{}, postProcessor PostProcessor, options ...OperationOption) Operations {
	f.tasks.PushBack(&fileSystemTask{
		Name:      name,
		Type:      taskTypeGetStruct,
		Processor: fmt.Sprintf("postprocessor=%s", postProcessor.Type()),
		Work:      f.getStruct(name, into, postProcessor, options...),
	})

	return f
}

func (f *fileSystem) getStruct(name string, into interface{}, postProcessor PostProcessor, options ...OperationOption) *fileSystemWork {
	work := &fileSystemWork{
		ReadType: readTypeStruct,
	}

	work.Fn = func() error {
		w := f.getBytes(name, nil, options...)

		err := w.Fn()
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}

		_, err = postProcessor.PostProcess(into, w.Bytes)
		if err != nil {
			return fmt.Errorf("failed to postprocess data: %w", err)
		}

		work.Path = w.Path
		work.Struct = into

		return nil
	}

	return work
}

func (f *fileSystem) GetBytes(name string, callback GetBytesCallback, options ...OperationOption) Operations {
	f.tasks.PushBack(&fileSystemTask{
		Name: name,
		Type: taskTypeGetBytes,
		Work: f.getBytes(name, callback, options...),
	})

	return f
}

func (f *fileSystem) getBytes(name string, callback GetBytesCallback, _ ...OperationOption) *fileSystemWork {
	work := &fileSystemWork{
		ReadType: readTypeBytes,
	}

	work.Fn = func() error {
		data, err := f.fs.ReadFile(path.Join(f.BaseDir, name))
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}

		if callback != nil {
			err = callback(name, data)
			if err != nil {
				return fmt.Errorf("callback failed: %w", err)
			}
		}

		work.Path = path.Join(f.BaseDir, name)
		work.Bytes = data

		return nil
	}

	return work
}

func (f *fileSystem) AddStoreStruct(operations ...AddStoreStruct) Operations {
	for _, o := range operations {
		f.StoreStruct(o.Name, o.Data, o.PreProcessor, o.Options...)
	}

	return f
}

func (f *fileSystem) AddStoreBytes(operations ...AddStoreBytes) Operations {
	for _, o := range operations {
		f.StoreBytes(o.Name, o.Data, o.Options...)
	}

	return f
}

func (f *fileSystem) AlterStore(alterers ...Alterer) Operations {
	types := make([]string, len(alterers))
	for i, a := range alterers {
		types[i] = fmt.Sprintf("alterer=%s", a.Type())
	}

	f.tasks.PushBack(&fileSystemTask{
		Type:      taskTypeAlter,
		Processor: strings.Join(types, ", "),
		Work:      f.alterStore(alterers...),
	})

	return f
}

func (f *fileSystem) alterStore(alterers ...Alterer) *fileSystemWork {
	work := &fileSystemWork{}

	work.Fn = func() error {
		for _, a := range alterers {
			err := a.Alter(f)
			if err != nil {
				return fmt.Errorf("failed to apply alteration %s: %w", a.Type(), err)
			}
		}

		return nil
	}

	return work
}

func (f *fileSystem) StoreStruct(name string, data interface{}, preProcessor PreProcessor, options ...OperationOption) Operations {
	f.tasks.PushBack(&fileSystemTask{
		Name:      name,
		Type:      taskTypeStoreStruct,
		Processor: fmt.Sprintf("preprocessing=%s", preProcessor.Type()),
		Work:      f.storeStruct(name, data, preProcessor, options...),
	})

	return f
}

func (f *fileSystem) storeStruct(name string, data interface{}, preProcessor PreProcessor, options ...OperationOption) *fileSystemWork {
	work := &fileSystemWork{}

	work.Fn = func() error {
		d, err := preProcessor.PreProcess(data)
		if err != nil {
			return fmt.Errorf("failed to preprocess data: %w", err)
		}

		w := f.storeBytes(name, d.Data, options...)

		err = w.Fn()
		if err != nil {
			return fmt.Errorf("failed to store bytes: %w", err)
		}

		work.Path = w.Path

		return nil
	}

	return work
}

func (f *fileSystem) StoreBytes(name string, data []byte, options ...OperationOption) Operations {
	f.tasks.PushBack(&fileSystemTask{
		Name: name,
		Type: taskTypeStoreBytes,
		Work: f.storeBytes(name, data, options...),
	})

	return f
}

func (f *fileSystem) storeBytes(name string, data []byte, options ...OperationOption) *fileSystemWork {
	work := &fileSystemWork{}

	work.Fn = func() error {
		if !f.OverwriteExisting {
			exists, err := f.fs.Exists(path.Join(f.BaseDir, name))
			if err != nil {
				return fmt.Errorf("failed to determine if file exists: %w", err)
			}

			if exists {
				return fmt.Errorf("file '%s' exists and overwrite is disabled", path.Join(f.BaseDir, name))
			}
		}

		if f.CreateDirectories {
			err := f.fs.MkdirAll(f.BaseDir, 0o744)
			if err != nil {
				return fmt.Errorf("failed to create directories: %w", err)
			}
		}

		exists, err := f.fs.DirExists(f.BaseDir)
		if err != nil {
			return fmt.Errorf("failed to determine if directory exists: %w", err)
		}

		if !exists {
			return fmt.Errorf("directory does not exist '%s' and create directories disabled", f.BaseDir)
		}

		var fileMode os.FileMode = 0o644

		for _, option := range options {
			switch o := option().(type) {
			case os.FileMode:
				fileMode = o
			default:
				return fmt.Errorf("cannot process unknown operation option: %v", o)
			}
		}

		err = f.fs.WriteFile(path.Join(f.BaseDir, name), data, fileMode)
		if err != nil {
			return fmt.Errorf("failed to write file: %w", err)
		}

		work.Path = path.Join(f.BaseDir, name)

		return nil
	}

	return work
}

func (f *fileSystem) Remove(name string, options ...OperationOption) Operations {
	f.tasks.PushBack(&fileSystemTask{
		Name: name,
		Type: taskTypeRemove,
		Work: f.remove(name, options...),
	})

	return f
}

func (f *fileSystem) remove(name string, _ ...OperationOption) *fileSystemWork {
	work := &fileSystemWork{}

	work.Fn = func() error {
		exists, err := f.fs.Exists(path.Join(f.BaseDir, name))
		if err != nil {
			return fmt.Errorf("failed to determine if file exists: %w", err)
		}

		work.Path = path.Join(f.BaseDir, name)

		if !exists {
			return nil
		}

		err = f.fs.Remove(path.Join(f.BaseDir, name))
		if err != nil {
			return fmt.Errorf("failed to remove file: %w", err)
		}

		return nil
	}

	return work
}

func (f *fileSystem) Do() (*Report, error) {
	i := 0

	for e := f.tasks.Front(); e != nil; e = e.Next() {
		task, ok := e.Value.(*fileSystemTask)
		if !ok {
			return nil, fmt.Errorf("failed to cast task to *fileSystemTask")
		}

		err := task.Work.Fn()
		if err != nil {
			return nil, fmt.Errorf("failed to process task %s(%s): %w", task.Type, task.Name, err)
		}

		switch task.Work.ReadType {
		case readTypeBytes:
			f.report.Data[task.Name] = task.Work.Bytes
		case readTypeStruct:
			f.report.Data[task.Name] = task.Work.Struct
		}

		description := fmt.Sprintf("task.%d %s", i+1, task.Type)
		actionType := string(task.Type)

		switch task.Type {
		case taskTypeGetBytes:
			description = fmt.Sprintf("%s from file '%s' (path: %s)", description, task.Name, task.Work.Path)
		case taskTypeGetStruct:
			description = fmt.Sprintf("%s[%s] from file '%s' (path: %s)", description, task.Processor, task.Name, task.Work.Path)
			actionType = fmt.Sprintf("%s[%s]", actionType, task.Processor)
		case taskTypeStoreStruct:
			description = fmt.Sprintf("%s[%s] to file '%s' (path: %s)", description, task.Processor, task.Name, task.Work.Path)
			actionType = fmt.Sprintf("%s[%s]", actionType, task.Processor)
		case taskTypeStoreBytes:
			description = fmt.Sprintf("%s to file '%s' (path: %s)", description, task.Name, task.Work.Path)
		case taskTypeAlter:
			description = fmt.Sprintf("%s[%s]", description, task.Processor)
			actionType = fmt.Sprintf("%s[%s]", actionType, task.Processor)
		case taskTypeRemove:
			description = fmt.Sprintf("%s file '%s' (path: %s)", description, task.Name, task.Work.Path)
		case taskTypeProcessGetStruct:
			description = fmt.Sprintf("%s on name '%s", description, task.Name)
		}

		f.report.Actions = append(f.report.Actions, Action{
			Name:        task.Name,
			Path:        task.Work.Path,
			Type:        actionType,
			Description: description,
		})

		i++
	}

	return f.report, nil
}
