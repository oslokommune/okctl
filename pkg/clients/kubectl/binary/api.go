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

// Patch applies patches to Kubernetes resources
func (c client) Patch(opts kubectl.PatchOpts) error {
	rawPatch, err := io.ReadAll(opts.Patch)
	if err != nil {
		return fmt.Errorf("reading patch: %w", err)
	}

	err = c.runKubectlCommand("patch", []string{
		"--namespace", opts.Namespace,
		"patch",
		opts.Kind, opts.Name,
		"--patch", string(rawPatch),
		"--type", "json",
	})
	if err != nil {
		return fmt.Errorf("calling kubectl: %w", err)
	}

	return nil
}

// Exists returns true if resource is found, false if not
func (c client) Exists(resource kubectl.Resource) (bool, error) {
	err := c.runKubectlCommand("exists", []string{
		"--namespace", resource.Namespace,
		"get",
		resource.Kind,
		resource.Name,
	})
	if err != nil {
		if isNotFoundErr(err) {
			return false, nil
		}

		return false, fmt.Errorf("checking resource existence: %w", err)
	}

	return true, nil
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
