// Package application provides functionality for interacting with okctl application data
package application

import (
	"regexp"
	"runtime"

	"github.com/AlecAivazis/survey/v2"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/google/uuid"
	"sigs.k8s.io/yaml"
)

const (
	// OsDarwin is a constant for macos
	OsDarwin = "darwin"
	// OsLinux is a constant for a linux based os
	OsLinux = "linux"

	// ArchAmd64 represents all 64-bit systems
	ArchAmd64 = "amd64"
)

// Data stores the state for the configuration
// of okctl itself
type Data struct {
	User     User
	Host     Host
	Binaries []Binary
}

// User stores state related to the user themselves
type User struct {
	ID       string
	Username string
}

// Valid returns no error if it passes all tests
func (u User) Valid() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.ID,
			validation.Required,
			is.UUIDv4,
		),
		validation.Field(&u.Username,
			validation.Required,
			validation.Match(regexp.MustCompile("^byr[0-9]{6}$")).Error("username must be in the form: byrXXXXXX"),
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

// Valid determines if the host operating
// system is valid
func (h Host) Valid() error {
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

// New Data returns the default configuration for the application
func New() *Data {
	return &Data{
		User: User{
			ID: uuid.New().String(),
		},
		Host: Host{
			Os:   runtime.GOOS,
			Arch: runtime.GOARCH,
		},
		Binaries: KnownBinaries(),
	}
}

// KnownBinaries returns a list of known binaries
func KnownBinaries() []Binary {
	return []Binary{
		{
			Name:       "kubectl",
			Version:    "1.16.8",
			BufferSize: "100mb",
			URLPattern: "https://amazon-eks.s3.us-west-2.amazonaws.com/#{ver}/2020-04-16/bin/${os}/${arch}/kubectl",
			Checksums: []Checksum{
				{
					Os:     "darwin",
					Arch:   "amd64",
					Type:   "sha256",
					Digest: "6e8439099c5a7d8d2f8f550f2f04301f9b0bb229a5f7c56477743a2cd11de2aa",
				},
				{
					Os:     "linux",
					Arch:   "amd64",
					Type:   "sha256",
					Digest: "e29544e1334f68e81546b8c8774c2484cbf82e8e5723d2a7e654f8a8fd79a7b2",
				},
			},
		},
		{
			Name:       "eksctl",
			Version:    "0.21.0",
			BufferSize: "100mb",
			URLPattern: "https://github.com/weaveworks/eksctl/releases/download/#{ver}/eksctl_#{os}_#{arch}.tar.gz",
			Archive: Archive{
				Type:   ".tar.gz",
				Target: "eksctl",
			},
			Checksums: []Checksum{
				{
					Os:     "darwin",
					Arch:   "amd64",
					Type:   "sha256",
					Digest: "3cdcbb1792bb131cc0ed944cbfc51dd6f1b2261a480436efc6f8124dea7c8c14",
				},
				{
					Os:     "linux",
					Arch:   "amd64",
					Type:   "sha256",
					Digest: "4573bca35af67fa002fb722b4d41fae2224a27576619ed2f1e269dd7bd15c214",
				},
			},
		},
	}
}

// Survey starts an interactive survey for fetching configuration
// information from the end user
func (d *Data) Survey() (*Data, error) {
	qs := []*survey.Question{
		{
			Name: "username",
			Prompt: &survey.Input{
				Message: "Your username:",
				Help:    "This is your byr user, e.g., byrXXXXXX. We store it in the application configuration so you don't have to enter it each time.",
			},
		},
	}

	answers := struct {
		Username string
	}{}

	err := survey.Ask(qs, &answers)
	if err != nil {
		return nil, err
	}

	d.User.Username = answers.Username

	return d, d.User.Valid()
}

// YAML returns the data serialised in a yaml reperesentation
func (d *Data) YAML() ([]byte, error) {
	return yaml.Marshal(d)
}
