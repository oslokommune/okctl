package load

import "github.com/oslokommune/okctl/pkg/config"

type DataNotFoundFn func(*config.Config) error

type DataNotFoundErr struct {
	err error
}

func (e *DataNotFoundErr) Error() string {
	return e.err.Error()
}
