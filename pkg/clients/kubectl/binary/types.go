package binary

import (
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/binaries"
	"github.com/oslokommune/okctl/pkg/credentials"
	"github.com/spf13/afero"
)

type teardownFn func() error

type client struct {
	fs                  *afero.Afero
	binaryProvider      binaries.Provider
	credentialsProvider credentials.Provider
	cluster             v1alpha1.Cluster
}
