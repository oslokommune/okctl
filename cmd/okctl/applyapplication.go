package main

import (
	"fmt"
	"io/ioutil"
	"path"

	"github.com/oslokommune/okctl/pkg/okctlapplication"

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

// nolint funlen
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

			app, err := okctlapplication.AcquireApplication(opts.File)
			if err != nil {
				return fmt.Errorf("unable to parse application.yaml: %w", err)
			}

			var iacRepoURL string

			cluster := okctlapplication.GetCluster(o, cmd, opts.Environment)
			if cluster != nil {
				for item := range cluster.Github.Repositories {
					iacRepoURL = cluster.Github.Repositories[item].GitURL

					break
				}
			}

			resources, err := okctlapplication.ConvertApplicationToResources(app, iacRepoURL)
			if err != nil {
				return fmt.Errorf("error expanding resources: %w", err)
			}

			err = handleOutput(o, app, opts.Output, resources)
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
	flags.StringVarP(&opts.Output, "output", "o", "files", "Specify how the format of the result. Choices: files, stdout")

	return cmd
}

func handleOutput(o *okctl.Okctl, app api.Application, outputFormat string, expandedApp *okctlapplication.ArgoCDDeploymentResources) error {
	switch outputFormat {
	case "stdout":
		fmt.Fprint(o.Out, &expandedApp.KubernetesResourcesBuffer)
		fmt.Fprint(o.Out, &expandedApp.ArgoAppBuffer)
	case "files":
		rootPath := "./deployment"
		kubernetesResourcesPath := path.Join(rootPath, fmt.Sprintf("%s.yaml", app.Name))
		argoAppResourcePath := path.Join(rootPath, fmt.Sprintf("%s-application.yaml", app.Name))

		err := ioutil.WriteFile(kubernetesResourcesPath, expandedApp.KubernetesResourcesBuffer.Bytes(), 0o600)
		if err != nil {
			return fmt.Errorf("unable to write kubernetes resources to file: %w", err)
		}

		err = ioutil.WriteFile(argoAppResourcePath, expandedApp.ArgoAppBuffer.Bytes(), 0o600)
		if err != nil {
			return fmt.Errorf("unable to write Argo Application to file: %w", err)
		}

		fmt.Fprintf(o.Out, "Kubernetes resources and Argo Application successfully saved to %s", rootPath)
	}

	return nil
}
