package main

import (
	"bytes"
	"fmt"
	"github.com/oslokommune/okctl/pkg/argoapp"
	"io/ioutil"
	"os"
	"path"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/oslokommune/kaex/pkg/api"
	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/spf13/cobra"
)

type ApplyApplicationOpts struct {
	File   string
	Output string
}

func (o *ApplyApplicationOpts) Validate() error {
	return validation.ValidateStruct(o,
		validation.Field(&o.File, validation.Required),
		validation.Field(&o.Output, validation.In("files", "stdout")),
	)
}

func buildApplyApplicationCommand(o *okctl.Okctl) *cobra.Command {
	opts := &ApplyApplicationOpts{}

	cmd := &cobra.Command{
		Use:   "application",
		Short: "Applies an application.yaml to the IAC repo",
		RunE: func(_ *cobra.Command, _ []string) error {
			err := opts.Validate()
			if err != nil {
				return fmt.Errorf("failed validating options: %w", err)
			}

			rawApplication, err := acquireRawApplication(opts.File)
			app, err := api.ParseApplication(rawApplication)
			if err != nil {
				return fmt.Errorf("unable to parse application.yaml: %w", err)
			}

			var kubernetesResourcesBuffer bytes.Buffer
			var argoAppBuffer bytes.Buffer

			err = api.Expand(&kubernetesResourcesBuffer, app, false)
			if err != nil {
				return fmt.Errorf("error expanding application %w", err)
			}

			argoApp, err := argoapp.CreateArgoApp(app, "git@github.com:<organization>/<repository>")
			if err != nil {
				return fmt.Errorf("error creating ArgoApp from application.yaml: %w", err)
			}
			err = api.WriteResource(&argoAppBuffer, argoApp)
			if err != nil {
				return fmt.Errorf("error writing ArgoApp to buffer", err)
			}

			switch opts.Output {
				case "stdout":
					fmt.Fprint(o.Out, &kubernetesResourcesBuffer)
					fmt.Fprint(o.Out, &argoAppBuffer)
				case "files":
					err := writeResourcesToFile(o, app, kubernetesResourcesBuffer, argoAppBuffer)
					if err != nil {
									  return fmt.Errorf("error writing deployment resources: %w", err)
									  }
			}

			return nil
		},
		PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
			return nil
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&opts.File, "file", "f", "", "Specify the file path. Use \"-\" for stdin")
	flags.StringVarP(&opts.Output, "output", "o", "files", "Specify how the format of the result. Choices: files, stdout")

	return cmd
}

func acquireRawApplication(path string) (string, error) {
	if path == "-" {
		result, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return "", fmt.Errorf("failed to read stdin: %w", err)
		}

		return string(result), nil
	}

	result, err := ioutil.ReadFile(path)
	if err != nil {
			  return "", fmt.Errorf("failed to read file: %w", err)
			  }

	return string(result), err
}

func writeResourcesToFile(o *okctl.Okctl, app api.Application, kubernetesResources bytes.Buffer, argoAppBuffer bytes.Buffer) error {
	rootPath := "./deployment"
	kubernetesResourcesPath := path.Join(rootPath, fmt.Sprintf("%s.yaml", app.Name))
	argoAppResourcePath := path.Join(rootPath, fmt.Sprintf("%s-application.yaml", app.Name))

	err := ioutil.WriteFile(kubernetesResourcesPath, kubernetesResources.Bytes(), 0644)
	if err != nil {
		return fmt.Errorf("unable to write kubernetes resources to file: %w", err)
	}

	err = ioutil.WriteFile(argoAppResourcePath, argoAppBuffer.Bytes(), 0644)
	if err != nil {
		return fmt.Errorf("unable to write Argo Application to file: %w", err)
	}

	fmt.Fprintf(o.Out, "Kubernetes resources and Argo Application successfully saved to %s", rootPath)

	return nil
}
