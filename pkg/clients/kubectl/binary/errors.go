package binary

import (
	"regexp"

	"github.com/oslokommune/okctl/pkg/clients/kubectl"
)

// errorHandler identifies known errors and returns a default error if error is not known
func errorHandler(err error, defaultError error) error {
	if isNotFoundErr(err) {
		return kubectl.ErrNotFound
	}

	return defaultError
}

var reNotFoundErr = regexp.MustCompile(`Error from server \(NotFound\)`)

func isNotFoundErr(err error) bool {
	return reNotFoundErr.MatchString(err.Error())
}
