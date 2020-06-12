package repository

import (
	"regexp"

	"github.com/AlecAivazis/survey/v2"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type Data struct {
	Name      string
	Region    string
	OutputDir string
	Clusters  []Cluster
}

func (d *Data) Validate() error {
	return validation.ValidateStruct(d,
		validation.Field(&d.Name,
			validation.Required,
		),
		validation.Field(&d.Region,
			validation.Required,
			validation.In(func() []interface{} {
				var o []interface{}
				for _, r := range v1alpha1.SupportedRegions() {
					o = append(o, r)
				}
				return o
			}()...),
		),
		validation.Field(&d.OutputDir,
			validation.Required,
		),
		validation.Field(&d.Clusters),
	)
}

type Cluster struct {
	Environment string
	AWS         AWS
}

const (
	envMinLength = 3
	envMaxLength = 10
)

func (c Cluster) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Environment,
			validation.Required,
			validation.Length(envMinLength, envMaxLength),
		),
		validation.Field(&c.AWS),
	)
}

type AWS struct {
	AccountID string
}

func (a AWS) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.AccountID,
			validation.Required,
			validation.Match(regexp.MustCompile("^[0-9]{12}$")),
		),
	)
}

func New() *Data {
	return &Data{}
}

func (d *Data) Survey() error {
	qs := []*survey.Question{
		{
			Name: "name",
			Prompt: &survey.Input{
				Message: "Name:",
				Help:    "A descriptive name, e.g., team or project, used among other things to prefix AWS resources",
			},
		},
		{
			Name: "region",
			Prompt: &survey.Select{
				Message: "Choose AWS region:",
				Options: v1alpha1.SupportedRegions(),
				Help:    "The AWS region resources will be created in",
			},
		},
		{
			Name: "basedir",
			Prompt: &survey.Input{
				Message: "Output directory:",
				Default: "infrastructure",
				Help:    "Path in the repository where generated files are stored",
			},
		},
	}

	answers := struct {
		Name    string
		Region  string
		BaseDir string
	}{}

	err := survey.Ask(qs, &answers)
	if err != nil {
		return err
	}

	d.Name = answers.Name
	d.Region = answers.Region
	d.OutputDir = answers.BaseDir

	return errors.Wrap(d.Validate(), "failed to validate repository data")
}

func (d *Data) YAML() ([]byte, error) {
	return yaml.Marshal(d)
}
