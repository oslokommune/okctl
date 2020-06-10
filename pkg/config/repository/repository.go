package repository

import "github.com/AlecAivazis/survey/v2"

type Data struct {
	Name     string
	Region   string
	BaseDir  string
	Clusters []Cluster
}

type Cluster struct {
	Name string
	AWS  AWS
}

type AWS struct {
	Account int
}

func New() *Data {
	return &Data{}
}

func (d *Data) Survey() error {
	qs := []*survey.Question{
		{
			Name: "name",
			Prompt: &survey.Input{
				Message: "Project name:",
				Help:    "This name will be used as a prefix for resources created in AWS, etc.",
			},
		},
		{
			Name: "region",
			Prompt: &survey.Select{
				Message: "Choose AWS region:",
				Options: []string{"eu-west-1"},
				Default: "eu-west-1",
				Help:    "This is the region that AWS resources will be created in",
			},
		},
		{
			Name: "basedir",
			Prompt: &survey.Input{
				Message: "Base output directory:",
				Default: "deployment",
				Help:    "This the location where all generated files will be written to",
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
	d.BaseDir = answers.BaseDir

	return nil
}
