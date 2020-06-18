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
	OsDarwin = "darwin"
	OsLinux  = "linux"

	ArchAmd64 = "amd64"
)

type Data struct {
	User     User
	Host     Host
	Binaries []Binary
}

type User struct {
	ID       string
	Username string
}

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

type Binary struct {
	Name       string
	Version    string
	BufferSize string
	URLPattern string
	Archive    Archive
	Checksums  []Checksum
}

type Archive struct {
	Type   string
	Target string
}

type Checksum struct {
	Os     string
	Arch   string
	Type   string
	Digest string
}

type Host struct {
	Os   string
	Arch string
}

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
		Binaries: []Binary{
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
		},
	}
}

func (d *Data) Survey() error {
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
		return err
	}

	d.User.Username = answers.Username

	return d.User.Valid()
}

func (d *Data) YAML() ([]byte, error) {
	return yaml.Marshal(d)
}
