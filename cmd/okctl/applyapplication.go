package main

import (
	"bytes"
	"fmt"
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
				err = fmt.Errorf("failed validating options: %w", err)

				fmt.Fprint(o.Err, err)

				return nil
			}

			rawApplication, err := acquireRawApplication(opts.File)
			app, err := api.ParseApplication(rawApplication)
			if err != nil {
				err = fmt.Errorf("unable to parse application.yaml: %w", err)

				fmt.Fprint(o.Err, err)

				return nil
			}

			var kubernetesResourcesBuffer bytes.Buffer

			err = api.Expand(&kubernetesResourcesBuffer, app, false)

			switch opts.Output {
			case "stdout":
				fmt.Fprintf(o.Out, "%s", &kubernetesResourcesBuffer)
			case "files":
				err := writeResourcesToFile(o, app, kubernetesResourcesBuffer)
				if err != nil {
					fmt.Fprintf(o.Err, "%w", err)
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
			fmt.Errorf("failed to read stdin: %w", err)

			return "", err
		}

		return string(result), nil
	} else {
		result, err := ioutil.ReadFile(path)
		if err != nil {
			fmt.Errorf("failed to read file: %w", err)

			return "", err
		}

		return string(result), err
	}
}

func writeResourcesToFile(o *okctl.Okctl, app api.Application, kubernetesResources bytes.Buffer) error {
	rootPath := "./deployment"
	kubernetesResourcesPath := path.Join(rootPath, fmt.Sprintf("%s.yaml", app.Name))

	err := ioutil.WriteFile(kubernetesResourcesPath, kubernetesResources.Bytes(), 0644)
	if err != nil {
		fmt.Errorf("unable to write kubernetes resources to file: %w", err)

		return err
	}

	fmt.Fprintf(o.Out, "Kubernetes resources successfully saved to %s", kubernetesResourcesPath)

	return nil
}
