package okctl

import (
	"github.com/oslokommune/okctl/pkg/binaries"
	"github.com/oslokommune/okctl/pkg/config"
	"github.com/oslokommune/okctl/pkg/credentials"
)

type Okctl struct {
	*config.Config

	BinariesProvider    binaries.Provider
	CredentialsProvider credentials.Provider
}

func New() *Okctl {
	return &Okctl{
		Config:              config.New(),
		BinariesProvider:    binaries.NewErrorProvider(),
		CredentialsProvider: credentials.NewErrorProvider(),
	}
}
