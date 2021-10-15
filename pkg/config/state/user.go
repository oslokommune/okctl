package state

import (
	"regexp"
	"runtime"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/google/uuid"
)

const (
	// OsDarwin is a constant for macos
	OsDarwin = "darwin"
	// OsLinux is a constant for a linux based os
	OsLinux = "linux"

	// ArchAmd64 represents all 64-bit systems
	ArchAmd64 = "amd64"

	defaultMetricsUserAgent = "okctl"
)

// User stores the state for the configuration
// of okctl itself
type User struct {
	User     UserInfo
	Host     Host
	Binaries []Binary
	Metrics  Metrics
}

// UserInfo stores state related to the user themselves
type UserInfo struct {
	ID       string
	Username string
}

// Validate returns no error if it passes all tests
func (u UserInfo) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.ID,
			validation.Required,
			is.UUIDv4,
		),
		validation.Field(&u.Username,
			validation.Required,
			validation.Match(regexp.MustCompile("^[a-z]{3}[0-9]{6}$")).Error("username must be in the form: yyyXXXXXX (y = letter, x = digit)"),
		),
	)
}

// Binary stores information on how a dependent CLI
// can be staged
type Binary struct {
	Name       string
	Version    string
	BufferSize string
	URLPattern string
	Archive    Archive
	Checksums  []Checksum
}

// Archive represents the compression type
type Archive struct {
	Type   string
	Target string
}

// Checksum represents the hashing algorithm and result
type Checksum struct {
	Os     string
	Arch   string
	Type   string
	Digest string
}

// Host represents the user system
type Host struct {
	Os   string
	Arch string
}

// Metrics exposes configuration of metrics
type Metrics struct {
	UserAgent string `json:"userAgent"`
}

// Validate determines if the host operating
// system is valid
func (h Host) Validate() error {
	return validation.ValidateStruct(&h,
		validation.Field(&h.Arch,
			validation.Required,
			validation.In(
				ArchAmd64,
			),
		),
		validation.Field(&h.Os,
			validation.Required,
			validation.In(
				OsDarwin,
				OsLinux,
			),
		),
	)
}

// NewUser returns the default configuration for the user state
func NewUser() *User {
	return &User{
		User: UserInfo{
			ID: uuid.New().String(),
		},
		Host: Host{
			Os:   runtime.GOOS,
			Arch: runtime.GOARCH,
		},
		Binaries: KnownBinaries(),
		Metrics: Metrics{
			UserAgent: defaultMetricsUserAgent,
		},
	}
}
