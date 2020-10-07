package main

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/oslokommune/okctl/pkg/storage"

	"github.com/oslokommune/okctl/pkg/scaffold"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/oslokommune/kaex/pkg/api"
	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/spf13/cobra"
)

// ApplyApplicationOpts contains all the possible options for "apply application"
type ApplyApplicationOpts struct {
	File        string
	Output      string
	Environment string
}

// Validate the options for "apply application"
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
		RunE: func(cmd *cobra.Command, _ []string) error {
			err := opts.Validate()
			if err != nil {
				return fmt.Errorf("failed validating options: %w", err)
			}

			app, err := scaffold.ReadApplication(o, opts.File)
			if err != nil {
				return fmt.Errorf("unable to parse application.yaml: %w", err)
			}

			deployment, err := scaffold.NewApplicationDeployment(app, o, cmd, opts.Environment)
			if err != nil {
				return fmt.Errorf("error creating a new application deployment: %w", err)
			}

			err = handleOutput(o.Out, app, opts.Output, deployment)
			if err != nil {
				return fmt.Errorf("unable to generate output: %w", err)
			}

			return nil
		},
		PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
			return nil
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&opts.Environment, "environment", "e", "", "Specify what environment to use")
	flags.StringVarP(&opts.File, "file", "f", "", "Specify the file path. Use \"-\" for stdin")
	flags.StringVarP(&opts.Output, "output", "o", "files", "Specify the format of the result. Choices: files, stdout")

	return cmd
}

func handleOutput(writer io.Writer, app api.Application, outputFormat string, deployment *scaffold.ApplicationDeployment) error {
	switch outputFormat {
	case "stdout":
		err := deployment.Write(writer)
		if err != nil {
			return fmt.Errorf("error writing deployment resources: %w", err)
		}
	case "files":
		basePath := filepath.Join("./deployment", app.Name)
		resourceStorage := storage.NewFileSystemStorage(basePath)

		kubernetesFile, _ := resourceStorage.Create(".", fmt.Sprintf("%s.yaml", app.Name), 0o600)

		err := deployment.WriteKubernetesResources(kubernetesFile)
		if err != nil {
			return fmt.Errorf("unable to write kubernetes resources to file: %w", err)
		}

		argoFile, _ := resourceStorage.Create(".", fmt.Sprintf("%s-application.yaml", app.Name), 0o600)

		err = deployment.WriteArgoResources(argoFile)
		if err != nil {
			return fmt.Errorf("unable to write Argo Application to file: %w", err)
		}

		fmt.Fprintf(writer, "Kubernetes resources and Argo Application successfully saved to %s\n", basePath)
		fmt.Fprint(writer, "To deploy, first:\n\t1. Commit and push the changes to master\n")
		fmt.Fprintf(writer, "\t2. Run kubectl apply -f %s/%s-application.yaml\n", basePath, app.Name)
	}

	return nil
}
