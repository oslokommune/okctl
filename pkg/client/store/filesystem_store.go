package store

import (
	"container/list"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/oslokommune/okctl/pkg/config/constant"

	"github.com/spf13/afero"
)

type fileSystemReadType string

const (
	readTypeStruct fileSystemReadType = "struct"
	readTypeBytes  fileSystemReadType = "bytes"
)

// fileSystemWork represents the result of the work
type fileSystemWork struct {
	ReadType fileSystemReadType
	Path     string
	Bytes    []byte
	Struct   interface{}
	Fn       func() error
}

type fileSystemTaskType string

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
	Type      fileSystemTaskType
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

// WithWriteIfNotExists will write the data if the file doesn't exist
func WithWriteIfNotExists(data []byte) OperationOption {
	return func() interface{} {
		return data
	}
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
		return fmt.Errorf(constant.FileSystemCastError)
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
			return fmt.Errorf(constant.StructNotFoundError, name)
		}

		err := process(data, f)
		if err != nil {
			return fmt.Errorf(constant.ProcessStructError, err)
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
			return fmt.Errorf(constant.ReadFileError, err)
		}

		_, err = postProcessor.PostProcess(into, w.Bytes)
		if err != nil {
			return fmt.Errorf(constant.PostProcessDataError, err)
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

func (f *fileSystem) getBytes(name string, callback GetBytesCallback, options ...OperationOption) *fileSystemWork {
	work := &fileSystemWork{
		ReadType: readTypeBytes,
	}

	work.Fn = func() error {
		for _, option := range options {
			switch o := option().(type) {
			case []byte:
				exists, err := f.fs.Exists(path.Join(f.BaseDir, name))
				if err != nil {
					return fmt.Errorf(constant.FileExistsError, err)
				}

				if !exists {
					err = f.storeBytes(name, o).Fn()
					if err != nil {
						return err
					}
				}
			default:
				return fmt.Errorf(constant.UnknownOptionError, o)
			}
		}

		data, err := f.fs.ReadFile(path.Join(f.BaseDir, name))
		if err != nil {
			return fmt.Errorf(constant.ReadFileError, err)
		}

		if callback != nil {
			callback(name, data)
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
				return fmt.Errorf(constant.ApplyAlterationError, a.Type(), err)
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
			return fmt.Errorf(constant.PreProcessDataError, err)
		}

		w := f.storeBytes(name, d.Data, options...)

		err = w.Fn()
		if err != nil {
			return fmt.Errorf(constant.StoreBytesError, err)
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
				return fmt.Errorf(constant.CheckIfFileExistsError, err)
			}

			if exists {
				return fmt.Errorf(constant.CannotOverwriteFileError, path.Join(f.BaseDir, name))
			}
		}

		if f.CreateDirectories {
			err := f.fs.MkdirAll(f.BaseDir, 0o744)
			if err != nil {
				return fmt.Errorf(constant.CreateDirectoriesError, err)
			}
		}

		exists, err := f.fs.DirExists(f.BaseDir)
		if err != nil {
			return fmt.Errorf(constant.CheckIfDirectoryExistsError, err)
		}

		if !exists {
			return fmt.Errorf(constant.CannotOverwriteDirectoryError, f.BaseDir)
		}

		var fileMode os.FileMode = 0o644

		for _, option := range options {
			switch o := option().(type) {
			case os.FileMode:
				fileMode = o
			default:
				return fmt.Errorf(constant.ProcessUnkownOpreationError, o)
			}
		}

		err = f.fs.WriteFile(path.Join(f.BaseDir, name), data, fileMode)
		if err != nil {
			return fmt.Errorf(constant.WriteToFileError, err)
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

func (f *fileSystem) RemoveDir(name string, options ...OperationOption) Operations {
	f.tasks.PushBack(&fileSystemTask{
		Name: name,
		Type: taskTypeRemove,
		Work: f.removeDir(name, options...),
	})

	return f
}

func (f *fileSystem) remove(name string, _ ...OperationOption) *fileSystemWork {
	work := &fileSystemWork{}

	work.Fn = func() error {
		exists, err := f.fs.Exists(path.Join(f.BaseDir, name))
		if err != nil {
			return fmt.Errorf(constant.CheckIfFileExistsError, err)
		}

		work.Path = path.Join(f.BaseDir, name)

		if !exists {
			return nil
		}

		err = f.fs.Remove(path.Join(f.BaseDir, name))
		if err != nil {
			return fmt.Errorf(constant.RemoveFileError, err)
		}

		return nil
	}

	return work
}

func (f *fileSystem) removeDir(name string, _ ...OperationOption) *fileSystemWork {
	work := &fileSystemWork{}

	work.Fn = func() error {
		exists, err := f.fs.Exists(path.Join(f.BaseDir, name))
		if err != nil {
			return fmt.Errorf(constant.CheckIfDirectoryExistsError, err)
		}

		work.Path = path.Join(f.BaseDir, name)

		if !exists {
			return nil
		}

		isDir, err := f.fs.IsDir(work.Path)
		if err != nil {
			return err
		}

		isEmpty, err := f.fs.IsEmpty(work.Path)
		if err != nil {
			return err
		}

		if isDir && isEmpty {
			err = f.fs.Remove(path.Join(f.BaseDir, name))
			if err != nil {
				return fmt.Errorf(constant.RemoveDirectoryError, err)
			}
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
			return nil, fmt.Errorf(constant.FileSystemTaskCastError)
		}

		err := task.Work.Fn()
		if err != nil {
			return nil, fmt.Errorf(constant.ProcessTaskError, task.Type, task.Name, err)
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
