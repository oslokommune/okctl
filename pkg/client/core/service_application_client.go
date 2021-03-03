package core

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"

	"github.com/oslokommune/okctl/pkg/spinner"

	kaex "github.com/oslokommune/kaex/pkg/api"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
	clientFilesystem "github.com/oslokommune/okctl/pkg/client/core/store/filesystem"
	"github.com/oslokommune/okctl/pkg/client/store"
	"github.com/oslokommune/okctl/pkg/scaffold"
)

type applicationService struct {
	fs     *afero.Afero
	spin   spinner.Spinner
	paths  clientFilesystem.Paths
	cert   client.CertificateService
	store  client.ApplicationStore
	report client.ApplicationReport
}

func createCertificateFn(ctx context.Context, certService client.CertificateService, id *api.ID, hostedZoneID string) func(domain string) (string, error) {
	return func(fqdn string) (string, error) {
		cert, certFnErr := certService.CreateCertificate(ctx, api.CreateCertificateOpts{
			ID:           *id,
			FQDN:         fqdn,
			Domain:       fqdn,
			HostedZoneID: hostedZoneID,
		})
		if certFnErr != nil {
			return "", certFnErr
		}

		return cert.CertificateARN, nil
	}
}

func writeSuccessMessage(writer io.Writer, applicationName string, argoCDResourcePath string) {
	fmt.Fprintf(writer, "Successfully scaffolded %s\n", applicationName)
	fmt.Fprintln(writer, "To deploy your application:")
	fmt.Fprintln(writer, "\t1. Commit and push the changes done by okctl")
	fmt.Fprintf(writer, "\t2. Run kubectl apply -f %s\n", argoCDResourcePath)
	fmt.Fprintf(writer, "If using an ingress, it can take up to five minutes for the routing to configure")
}

// ScaffoldApplication turns a file path into Kubernetes resources
func (s *applicationService) ScaffoldApplication(ctx context.Context, opts *client.ScaffoldApplicationOpts) error {
	err := opts.Validate()
	if err != nil {
		return err
	}

	err = s.spin.Start("Scaffolding application")

	defer func() {
		err = s.spin.Stop()
	}()

	app, err := inferApplicationFromStdinOrFile(opts.In, opts.ApplicationFilePath)
	if err != nil {
		return err
	}

	applicationDir := path.Join(s.paths.BaseDir, app.Name)
	applicationDir = strings.Replace(applicationDir, opts.RepoDir+"/", "", 1)
	certFn := createCertificateFn(ctx, s.cert, opts.ID, opts.HostedZoneID)

	deployment, err := scaffold.NewApplicationDeployment(*app, certFn, opts.IACRepoURL, applicationDir)
	if err != nil {
		return fmt.Errorf("error creating a new application deployment: %w", err)
	}

	var kubernetesResources, argoCDResource bytes.Buffer

	err = deployment.WriteKubernetesResources(&kubernetesResources)
	if err != nil {
		return err
	}

	err = deployment.WriteArgoResources(&argoCDResource)
	if err != nil {
		return err
	}

	applicationScaffold := &client.ScaffoldedApplication{
		ApplicationName:     app.Name,
		KubernetesResources: kubernetesResources.Bytes(),
		ArgoCDResource:      argoCDResource.Bytes(),
	}

	report, err := s.store.SaveApplication(applicationScaffold)
	if err != nil {
		return err
	}

	err = s.report.ReportCreateApplication(applicationScaffold, []*store.Report{report})
	if err != nil {
		return err
	}

	writeSuccessMessage(opts.Out, app.Name, path.Join(opts.RepoDir, applicationDir, fmt.Sprintf("%s-application.yaml", app.Name)))

	return nil
}

// NewApplicationService initializes a new Scaffold application service
func NewApplicationService(
	fs *afero.Afero,
	spin spinner.Spinner,
	paths clientFilesystem.Paths,
	cert client.CertificateService,
	store client.ApplicationStore,
	state client.ApplicationReport,
) client.ApplicationService {
	return &applicationService{
		fs:     fs,
		spin:   spin,
		paths:  paths,
		cert:   cert,
		store:  store,
		report: state,
	}
}

func inferApplicationFromStdinOrFile(stdin io.Reader, fs *afero.Afero, path string) (app client.OkctlApplication, err error) {
	var inputReader io.Reader

	switch path {
	case "-":
		inputReader = stdin
	default:
		inputReader, err = fs.Open(filepath.Clean(path))
		if err != nil {
			return app, fmt.Errorf("opening application file: %w", err)
		}
	}

	app, err := kaex.ParseApplication(inputReader)
	if err != nil {
		return nil, fmt.Errorf("unable to parse application.yaml: %w", err)
	}

	return &app, nil
}
