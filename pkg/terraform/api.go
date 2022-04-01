package terraform

import (
	"errors"
	"fmt"
	"io/fs"
	"path"

	"github.com/oslokommune/okctl/pkg/clients/terraform"

	"github.com/oslokommune/okctl/pkg/logging"

	"github.com/oslokommune/okctl/pkg/clients/terraform/binary"
	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/spf13/afero"
)

// SyncDirectory knows how to terraform apply a certain directory
func SyncDirectory(o *okctl.Okctl, targetDirectory string) error {
	log := logging.GetLogger("terraform", "SyncDirectory")
	client := binary.New(o.BinariesProvider, o.CredentialsProvider, o.Declaration.Terraform.Version)

	synchronizable, err := isSynchronizable(o.FileSystem, targetDirectory)
	if err != nil {
		return fmt.Errorf("checking synchronizability: %w", err)
	}

	if !synchronizable {
		log.Debug(fmt.Sprintf("%s not synchronizable", targetDirectory))

		return nil
	}

	err = ensureInitialized(o.FileSystem, client, targetDirectory)
	if err != nil {
		return fmt.Errorf("ensuring initialization: %w", err)
	}

	err = ensureConfig(o.FileSystem, targetDirectory, o.Declaration.Terraform.Version, o.Declaration.Metadata.Region)
	if err != nil {
		return fmt.Errorf("ensuring config: %w", err)
	}

	err = client.Apply(targetDirectory)
	if err != nil {
		return fmt.Errorf("applying terraform: %w", err)
	}

	return nil
}

func ensureInitialized(filesystem *afero.Afero, tf terraform.Client, dir string) error {
	exists, err := filesystem.DirExists(path.Join(dir, ".terraform"))
	if err != nil {
		return fmt.Errorf("checking .terraform directory existence: %w", err)
	}

	if exists {
		return nil
	}

	err = tf.Initialize(dir)
	if err != nil {
		return fmt.Errorf("initializing: %w", err)
	}

	return nil
}

func isSynchronizable(fileSystem *afero.Afero, dir string) (bool, error) {
	exists, err := fileSystem.DirExists(dir)
	if err != nil {
		return false, fmt.Errorf("checking directory existence: %w", err)
	}

	if !exists {
		return false, nil
	}

	err = fileSystem.Walk(dir, func(currentPath string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if path.Ext(currentPath) == ".tf" {
			return errfoundTerraformFile
		}

		return nil
	})
	if err != nil {
		if errors.Is(err, errfoundTerraformFile) {
			return true, nil
		}

		return false, fmt.Errorf("walking directory: %w", err)
	}

	return true, nil
}

var errfoundTerraformFile = errors.New("found tf file")

func ensureConfig(fs *afero.Afero, dir string, version string, region string) error {
	configPath := path.Join(dir, "terraform.tf")

	exists, err := fs.Exists(configPath)
	if err != nil {
		return fmt.Errorf("checking terraform config existence: %w", err)
	}

	if exists {
		return nil
	}

	cfg, err := provisionConfig(version, region)
	if err != nil {
		return fmt.Errorf("provisioning config: %w", err)
	}

	err = fs.WriteReader(configPath, cfg)
	if err != nil {
		return fmt.Errorf("writing config: %w", err)
	}

	return nil
}
