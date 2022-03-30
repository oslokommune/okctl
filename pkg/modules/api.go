package modules

import (
	"fmt"
	"path"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/spf13/afero"
)

// This should be moved to ~/.okctl/conf.yaml if prod START
const modulesRepositoryURL = "git@github.com:oslokommune/okctl-modules-poc.git"

// This should be moved to ~/.okctl/conf.yaml if prod END

// InstallModule downloads a module from the okctl modules repository and installs it into a directory
func InstallModule(fs *afero.Afero, moduleName string, destDirectory string) error {
	moduleFs := memfs.New()

	err := acquireModules(moduleFs)
	if err != nil {
		return fmt.Errorf("acquiring modules: %w", err)
	}

	err = copyModuleToFs(moduleFs, fs, moduleName, destDirectory)
	if err != nil {
		return fmt.Errorf("storing module: %w", err)
	}

	return nil
}

func acquireModules(fs billy.Filesystem) error {
	_, err := git.Clone(memory.NewStorage(), fs, &git.CloneOptions{
		URL:   modulesRepositoryURL,
		Depth: 1,
	})
	if err != nil {
		return fmt.Errorf("cloning modules repository: %w", err)
	}

	return nil
}

func copyModuleToFs(sourceFs billy.Filesystem, destFs *afero.Afero, moduleName string, destBaseDirectory string) error {
	files, err := sourceFs.ReadDir(moduleName)
	if err != nil {
		return fmt.Errorf("listing module files: %w", err)
	}

	moduleDir := path.Join(destBaseDirectory, moduleName)

	err = destFs.MkdirAll(moduleDir, 0o700)
	if err != nil {
		return fmt.Errorf("preparing directory: %w", err)
	}

	for _, file := range files {
		f, err := sourceFs.Open(path.Join(moduleName, file.Name()))
		if err != nil {
			return fmt.Errorf("opening module file: %w", err)
		}

		err = destFs.WriteReader(path.Join(moduleDir, file.Name()), f)
		if err != nil {
			return fmt.Errorf("storing module file: %w", err)
		}

		err = f.Close()
		if err != nil {
			return fmt.Errorf("closing module file: %w", err)
		}
	}

	return nil
}
