package okctl

import "github.com/oslokommune/okctl/pkg/config"

type Okctl struct {
	*config.Config
}

func New() *Okctl {
	return &Okctl{
		Config: config.New(),
	}
}
