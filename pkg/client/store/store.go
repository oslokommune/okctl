// Package store provides interfaces for more easily interacting with
// a storage layer
package store

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

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

// Operations defines the functions that can be chained
// before being executed
type Operations interface {
	// StoreStruct will preprocess the data with the given preprocessor
	// and then store the processed data under the given name, how this
	// data is stored is given by the underlying implementation.
	// This method is chainable, so multiple actions may be performed.
	StoreStruct(name string, data interface{}, preProcessor PreProcessor, options ...OperationOption) Operations

	// StoreBytes stores the provided data under the given name, how this
	// data is stored is given by the underlying implementation.
	// This method is chainable, so multiple actions may be performed.
	StoreBytes(name string, data []byte, options ...OperationOption) Operations

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
	d, err := json.Marshal(data)
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
type fileSystemWork func() error

// fileSystemTask represents a task to be performed
type fileSystemTask struct {
	Name string
	Path string
	Type string
	Work fileSystemWork
}

// fileSystem contains state for the file system
// storage implementation
type fileSystem struct {
	baseDir           string
	fs                *afero.Afero
	createDirectories bool
	overwriteExisting bool
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
		fs.createDirectories = v
	}
}

// FileSystemOverwriteExisting will ensure that existing files
// are overwritten if set to true
func FileSystemOverwriteExisting(v bool) FileSystemOption {
	return func(fs *fileSystem) {
		fs.overwriteExisting = v
	}
}

// Ensure that fileSystem implements the Operations interface
var _ Operations = &fileSystem{}

// NewFileSystem returns an initialised filesystem store
func NewFileSystem(baseDir string, fs *afero.Afero, options ...FileSystemOption) Operations {
	f := &fileSystem{
		baseDir:           baseDir,
		fs:                fs,
		createDirectories: true,
		overwriteExisting: true,
	}

	for _, option := range options {
		option(f)
	}

	f.report = &Report{
		Type: "FileSystem",
		Configuration: fmt.Sprintf(
			"BaseDir: %s\nCreateDirectories: %t\nOverWriteExisting: %t\n",
			f.baseDir,
			f.createDirectories,
			f.overwriteExisting,
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

func (f *fileSystem) StoreStruct(name string, data interface{}, preProcessor PreProcessor, options ...OperationOption) Operations {
	f.tasks = append(f.tasks, &fileSystemTask{
		Name: name,
		Path: path.Join(f.baseDir, name),
		Type: fmt.Sprintf("StoreStruct[preprocessing=%s]", preProcessor.Type()),
		Work: f.storeStruct(name, data, preProcessor, options...),
	})

	return f
}

func (f *fileSystem) storeStruct(name string, data interface{}, preProcessor PreProcessor, options ...OperationOption) fileSystemWork {
	return func() error {
		d, err := preProcessor.PreProcess(data)
		if err != nil {
			return fmt.Errorf("failed to preprocess data: %w", err)
		}

		return f.storeBytes(name, d.Data, options...)()
	}
}

func (f *fileSystem) StoreBytes(name string, data []byte, options ...OperationOption) Operations {
	f.tasks = append(f.tasks, &fileSystemTask{
		Name: name,
		Path: path.Join(f.baseDir, name),
		Type: "StoreBytes",
		Work: f.storeBytes(name, data, options...),
	})

	return f
}

func (f *fileSystem) storeBytes(name string, data []byte, options ...OperationOption) fileSystemWork {
	return func() error {
		if !f.overwriteExisting {
			exists, err := f.fs.Exists(path.Join(f.baseDir, name))
			if err != nil {
				return fmt.Errorf("failed to determine if file exists: %w", err)
			}

			if exists {
				return fmt.Errorf("file: %s exists, and overwrite existing is: %t", path.Join(f.baseDir, name), f.overwriteExisting)
			}
		}

		if f.createDirectories {
			err := f.fs.MkdirAll(f.baseDir, 0o744)
			if err != nil {
				return fmt.Errorf("failed to create directories: %w", err)
			}
		}

		exists, err := f.fs.DirExists(f.baseDir)
		if err != nil {
			return fmt.Errorf("failed to determine if directory exists: %w", err)
		}

		if !exists {
			return fmt.Errorf("directory does not exist: %s", f.baseDir)
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

		err = f.fs.WriteFile(path.Join(f.baseDir, name), data, fileMode)
		if err != nil {
			return fmt.Errorf("failed to write file: %w", err)
		}

		return nil
	}
}

func (f *fileSystem) Remove(name string, options ...OperationOption) Operations {
	f.tasks = append(f.tasks, &fileSystemTask{
		Name: name,
		Path: path.Join(f.baseDir, name),
		Type: "Remove",
		Work: f.remove(name, options...),
	})

	return f
}

func (f *fileSystem) remove(name string, _ ...OperationOption) fileSystemWork {
	return func() error {
		exists, err := f.fs.Exists(path.Join(f.baseDir, name))
		if err != nil {
			return fmt.Errorf("failed to determine if file exists: %w", err)
		}

		if !exists {
			return nil
		}

		err = f.fs.Remove(path.Join(f.baseDir, name))
		if err != nil {
			return fmt.Errorf("failed to remove file: %w", err)
		}

		return nil
	}
}

func (f *fileSystem) Do() (*Report, error) {
	f.report.Actions = make([]Action, len(f.tasks))

	for i, task := range f.tasks {
		err := task.Work()
		if err != nil {
			return nil, fmt.Errorf("failed to process task: %s (%s), because: %w", task.Name, task.Path, err)
		}

		f.report.Actions[i] = Action{
			Name:        task.Name,
			Path:        task.Path,
			Type:        task.Type,
			Description: fmt.Sprintf("task #%d %s: %s (%s)", i, task.Type, task.Name, task.Path),
		}
	}

	return f.report, nil
}
