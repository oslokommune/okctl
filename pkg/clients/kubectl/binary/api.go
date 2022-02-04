package binary

import (
	"fmt"
	"io"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/binaries"
	"github.com/oslokommune/okctl/pkg/credentials"

	"github.com/oslokommune/okctl/pkg/clients/kubectl"
	"github.com/spf13/afero"
)

// Apply runs kubectl apply on a manifest
func (c client) Apply(manifest io.Reader) error {
	targetPath, teardowner, err := c.cacheReaderOnFs(manifest)
	if err != nil {
		return fmt.Errorf("caching manifest on file system: %w", err)
	}

	defer func() {
		_ = teardowner()
	}()

	err = c.runKubectlCommand("apply", []string{"apply", "-f", targetPath})
	if err != nil {
		return fmt.Errorf("applying manifest: %w", err)
	}

	return nil
}

// Delete runs kubectl delete on a manifest
func (c client) Delete(manifest io.Reader) error {
	targetPath, teardowner, err := c.cacheReaderOnFs(manifest)
	if err != nil {
		return fmt.Errorf("caching manifest on file system: %w", err)
	}

	defer func() {
		_ = teardowner()
	}()

	err = c.runKubectlCommand("delete", []string{"delete", "-f", targetPath})
	if err != nil {
		return fmt.Errorf("deleting manifest: %w", err)
	}

	return nil
}

// New returns an initialized kubectl binary client
func New(fs *afero.Afero, binaryProvider binaries.Provider, credentialsProvider credentials.Provider, cluster v1alpha1.Cluster) kubectl.Client {
	return &client{
		fs:                  fs,
		binaryProvider:      binaryProvider,
		credentialsProvider: credentialsProvider,
		cluster:             cluster,
	}
}
