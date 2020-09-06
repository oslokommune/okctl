// Package store provides interfaces for more easily interacting with
// a storage layer
package store

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/spf13/afero"
)

// Processed contains information about how the data
// has been processed.
type Processed struct {
	Format string
	Data   []byte
}

// PreProcessor defines a type for preprocessing of the
// data before it is stored
type PreProcessor interface {
	PreProcess(data interface{}) (*Processed, error)
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

// AddStoreBytes contains the fields require for adding
// a StoreBytes operation, useful when operations are
// created dynamically from a slice or map
type AddStoreBytes struct {
	Name    string
	Data    []byte
	Options []OperationOption
}

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

func (t *toJSON) PreProcess(data interface{}) (*Processed, error) {
	d, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal as json: %w", err)
	}

	return &Processed{
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

type toYAML struct{}

func (t *toYAML) PreProcess(data interface{}) (*Processed, error) {
	d, err := yaml.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal as yaml: %w", err)
	}

	return &Processed{
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

// fileSystemWork is invoked during task processing
type fileSystemWork func() (string, error)

// fileSystemTask represents a task to be performed
type fileSystemTask struct {
	Name string
	Type string
	Work fileSystemWork
}

// fileSystem contains state for the file system
// storage implementation
type fileSystem struct {
	BaseDir           string
	CreateDirectories bool
	OverwriteExisting bool
	fs                *afero.Afero
	report            *Report
	tasks             []*fileSystemTask
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
		fs:                fs,
		CreateDirectories: true,
		OverwriteExisting: true,
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

type fileSystemSetBaseDir struct {
	baseDir string
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

func (f *fileSystem) AlterStore(alterers ...Alterer) Operations {
	types := make([]string, len(alterers))
	for i, a := range alterers {
		types[i] = a.Type()
	}

	f.tasks = append(f.tasks, &fileSystemTask{
		Name: "n/a",
		Type: fmt.Sprintf("Alter[%s]", strings.Join(types, ", ")),
		Work: f.alterStore(alterers...),
	})

	return f
}

func (f *fileSystem) alterStore(alterers ...Alterer) fileSystemWork {
	return func() (string, error) {
		for _, a := range alterers {
			err := a.Alter(f)
			if err != nil {
				return "", fmt.Errorf("failed to apply alteration %s: %w", a.Type(), err)
			}
		}

		return "n/a", nil
	}
}

func (f *fileSystem) StoreStruct(name string, data interface{}, preProcessor PreProcessor, options ...OperationOption) Operations {
	f.tasks = append(f.tasks, &fileSystemTask{
		Name: name,
		Type: fmt.Sprintf("StoreStruct[preprocessing=%s]", preProcessor.Type()),
		Work: f.storeStruct(name, data, preProcessor, options...),
	})

	return f
}

func (f *fileSystem) storeStruct(name string, data interface{}, preProcessor PreProcessor, options ...OperationOption) fileSystemWork {
	return func() (string, error) {
		d, err := preProcessor.PreProcess(data)
		if err != nil {
			return "", fmt.Errorf("failed to preprocess data: %w", err)
		}

		return f.storeBytes(name, d.Data, options...)()
	}
}

func (f *fileSystem) StoreBytes(name string, data []byte, options ...OperationOption) Operations {
	f.tasks = append(f.tasks, &fileSystemTask{
		Name: name,
		Type: "StoreBytes",
		Work: f.storeBytes(name, data, options...),
	})

	return f
}

func (f *fileSystem) storeBytes(name string, data []byte, options ...OperationOption) fileSystemWork {
	return func() (string, error) {
		if !f.OverwriteExisting {
			exists, err := f.fs.Exists(path.Join(f.BaseDir, name))
			if err != nil {
				return "", fmt.Errorf("failed to determine if file exists: %w", err)
			}

			if exists {
				return "", fmt.Errorf("file '%s' exists and overwrite is disabled", path.Join(f.BaseDir, name))
			}
		}

		if f.CreateDirectories {
			err := f.fs.MkdirAll(f.BaseDir, 0o744)
			if err != nil {
				return "", fmt.Errorf("failed to create directories: %w", err)
			}
		}

		exists, err := f.fs.DirExists(f.BaseDir)
		if err != nil {
			return "", fmt.Errorf("failed to determine if directory exists: %w", err)
		}

		if !exists {
			return "", fmt.Errorf("directory does not exist '%s' and create directories disabled", f.BaseDir)
		}

		var fileMode os.FileMode = 0o644

		for _, option := range options {
			switch o := option().(type) {
			case os.FileMode:
				fileMode = o
			default:
				return "", fmt.Errorf("cannot process unknown operation option: %v", o)
			}
		}

		err = f.fs.WriteFile(path.Join(f.BaseDir, name), data, fileMode)
		if err != nil {
			return "", fmt.Errorf("failed to write file: %w", err)
		}

		return path.Join(f.BaseDir, name), nil
	}
}

func (f *fileSystem) Remove(name string, options ...OperationOption) Operations {
	f.tasks = append(f.tasks, &fileSystemTask{
		Name: name,
		Type: "Remove",
		Work: f.remove(name, options...),
	})

	return f
}

func (f *fileSystem) remove(name string, _ ...OperationOption) fileSystemWork {
	return func() (string, error) {
		exists, err := f.fs.Exists(path.Join(f.BaseDir, name))
		if err != nil {
			return "", fmt.Errorf("failed to determine if file exists: %w", err)
		}

		if !exists {
			return path.Join(f.BaseDir, name), nil
		}

		err = f.fs.Remove(path.Join(f.BaseDir, name))
		if err != nil {
			return "", fmt.Errorf("failed to remove file: %w", err)
		}

		return path.Join(f.BaseDir, name), nil
	}
}

func (f *fileSystem) Do() (*Report, error) {
	f.report.Actions = make([]Action, len(f.tasks))

	for i, task := range f.tasks {
		p, err := task.Work()
		if err != nil {
			return nil, fmt.Errorf("failed to process task %s(%s): %w", task.Type, task.Name, err)
		}

		f.report.Actions[i] = Action{
			Name:        task.Name,
			Path:        p,
			Type:        task.Type,
			Description: fmt.Sprintf("task.%d %s to file '%s' (path: %s)", i+1, task.Type, task.Name, p),
		}
	}

	return f.report, nil
}
