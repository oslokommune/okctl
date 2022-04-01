package terraform

import (
	"bytes"
	"fmt"
	"io"
	"text/template"
)

func provisionConfig(version string, region string) (io.Reader, error) {
	t, err := template.New("").Parse(configTemplate)
	if err != nil {
		return nil, fmt.Errorf("parsing template: %w", err)
	}

	buf := bytes.Buffer{}

	err = t.Execute(&buf, struct {
		Version string
		Region  string
	}{
		Version: version,
		Region:  region,
	})
	if err != nil {
		return nil, fmt.Errorf("building config: %w", err)
	}

	return &buf, nil
}

const configTemplate = `
terraform {
  required_version = "~> {{ .Version }}"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.6"
    }
  }

}

provider "aws" {
  region = "{{ .Region }}"
}
`
