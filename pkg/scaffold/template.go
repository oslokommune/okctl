package scaffold

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/spf13/afero"
)

// InterpolationOpts defines possible data to inject into the templates
type InterpolationOpts struct {
	Domain string
}

// GenerateOkctlAppTemplate generates an okctl appliction template
func GenerateOkctlAppTemplate(opts *InterpolationOpts) ([]byte, error) {
	interpolated, err := goTemplateToBytes(okctlAppTemplate, opts)
	if err != nil {
		return nil, err
	}

	return interpolated, nil
}

// SaveOkctlAppTemplate saves a byte array as an application.yaml file in the current directory
func SaveOkctlAppTemplate(fs *afero.Afero, path string, template []byte) error {
	applicationFile, err := fs.Create(path)
	if err != nil {
		return fmt.Errorf("error creating application.yaml: %w", err)
	}

	_, err = applicationFile.Write(template)
	if err != nil {
		return fmt.Errorf("error writing to application.yaml: %w", err)
	}

	err = applicationFile.Close()
	if err != nil {
		return fmt.Errorf("unable to close application.yaml after writing: %w", err)
	}

	return err
}

// goTemplateToBytes converts a Go template plus provided data to a string
func goTemplateToBytes(templateString string, data interface{}) ([]byte, error) {
	tmpl, err := template.New("t").Parse(templateString)
	if err != nil {
		return nil, err
	}

	tmplBuffer := new(bytes.Buffer)
	err = tmpl.Execute(tmplBuffer, data)

	if err != nil {
		return nil, err
	}

	return tmplBuffer.Bytes(), nil
}

const okctlAppTemplate = `# A name that identifies your app
name: my-app
# An URI for your app Docker image
image: docker.pkg.github.com/my-org/my-repo/my-package
# The version of your app which is available as an image
version: 0.0.1

# The URL your app should be available on
# Change to something other than https to disable configuring TLS
# Comment this out to avoid setting up an ingress
subDomain: my-app

# The port your app listens on
# Comment this out to avoid setting up a service (required if url is specified)
port: 3000

# How many replicas of your application should we scaffold
#replicas: 3 # 1 by default

# A namespace where your app will live
#namespace: my-namespace

# A Docker repository secret for pulling your image
#imagePullSecret: my-pull-secret-name

# The environment your app requires
#environment:
#  MY_VARIABLE: my-value

# Volumes to mount
#volumes:
#  - /path/to/mount/volume: # Requests 1Gi by default
#  - /path/to/mount/volume: 24Gi

# Annotations for your ingress
#ingress:
#  annotations:
#    nginx.ingress.kubernetes.io/cors-allow-origin: http://localhost:8080
#    cert-manager.io/cluster-issuer: letsencrypt-production
`
