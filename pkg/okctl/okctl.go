package okctl

import (
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/binaries"
	"github.com/oslokommune/okctl/pkg/config"
	"github.com/oslokommune/okctl/pkg/credentials"
)

// Okctl stores all state required for invoking commands
type Okctl struct {
	*config.Config

	CloudProvider       v1alpha1.CloudProvider
	BinariesProvider    binaries.Provider
	CredentialsProvider credentials.Provider
}

// New returns a new okctl instance
func New() *Okctl {
	return &Okctl{
		Config:              config.New(),
		BinariesProvider:    binaries.NewErrorProvider(),
		CredentialsProvider: credentials.NewErrorProvider(),
	}
}
