package validate

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/oslokommune/okctl/pkg/config"
)

const (
	OsDarwin = "darwin"
	OsLinux  = "linux"

	ArchAmd64 = "amd64"
)

// Host ensures that the operating system and architecture of the host machine
// are supported.
func Host(h config.Host) error {
	return validation.ValidateStruct(&h,
		validation.Field(&h.Arch,
			validation.Required,
			validation.In([]string{
				ArchAmd64,
			}),
		),
		validation.Field(&h.Os,
			validation.Required,
			validation.In([]string{
				OsDarwin,
				OsLinux,
			}),
		),
	)
}
