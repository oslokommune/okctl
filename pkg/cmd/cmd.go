package cmd

import (
	"bytes"
	"text/template"
)

// GoTemplateToString converts a Go template plus provided data to a string
func GoTemplateToString(templateString string, data interface{}) (string, error) {
	tmpl, err := template.New("t").Parse(templateString)
	if err != nil {
		return "", err
	}

	tmplBuffer := new(bytes.Buffer)
	err = tmpl.Execute(tmplBuffer, data)

	if err != nil {
		return "", err
	}

	return tmplBuffer.String(), nil
}
