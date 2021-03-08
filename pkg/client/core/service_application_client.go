package core

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"path"
	"path/filepath"

	"github.com/oslokommune/okctl/pkg/config"
	"github.com/oslokommune/okctl/pkg/spinner"
	"github.com/spf13/afero"

	"sigs.k8s.io/yaml"

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

func (s *applicationService) createCertificate(ctx context.Context, id *api.ID, hostedZoneID, fqdn string) (string, error) {
	cert, certFnErr := s.cert.CreateCertificate(ctx, api.CreateCertificateOpts{
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

func writeSuccessMessage(writer io.Writer, applicationName string, argoCDResourcePath string) {
	fmt.Fprintf(writer, "\nSuccessfully scaffolded %s\n", applicationName)
	fmt.Fprintln(writer, "To deploy your application:")
	fmt.Fprintln(writer, "\t1. Commit and push the changes done by okctl")
	fmt.Fprintf(writer, "\t2. Run kubectl apply -f %s\n", argoCDResourcePath)
	fmt.Fprintf(writer, "If using an ingress, it can take up to five minutes for the routing to configure\n")
}

// ScaffoldApplication turns a file path into Kubernetes resources
// nolint: funlen
func (s *applicationService) ScaffoldApplication(ctx context.Context, opts *client.ScaffoldApplicationOpts) error {
	err := opts.Validate()
	if err != nil {
		return err
	}

	err = s.spin.Start("Scaffolding application")

	defer func() {
		err = s.spin.Stop()
	}()

	okctlApp, err := inferApplicationFromStdinOrFile(opts.In, s.fs, opts.ApplicationFilePath)
	if err != nil {
		return err
	}

	// See function comment
	app := okctlApplicationToKaexApplication(okctlApp, opts.HostedZoneDomain)

	relativeApplicationDir := path.Join(opts.OutputDir, config.DefaultApplicationsOutputDir, okctlApp.Name)
	relativeArgoCDSourcePath := path.Join(relativeApplicationDir, config.DefaultApplicationOverlayDir, opts.ID.Environment)

	base, err := scaffold.GenerateApplicationBase(*app, opts.IACRepoURL, relativeArgoCDSourcePath)
	if err != nil {
		return fmt.Errorf("error creating a new application deployment: %w", err)
	}

	certArn, err := s.createCertificate(
		ctx,
		opts.ID,
		opts.HostedZoneID,
		fmt.Sprintf("%s.%s", okctlApp.SubDomain, opts.HostedZoneDomain),
	)
	if err != nil {
		return fmt.Errorf("create certificate: %w", err)
	}

	overlay, err := scaffold.GenerateApplicationOverlay(okctlApp, opts.HostedZoneDomain, certArn)
	if err != nil {
		return fmt.Errorf("generating application overlay: %w", err)
	}

	applicationScaffold := &client.ScaffoldedApplication{
		ApplicationName:      app.Name,
		Environment:          opts.ID.Environment,
		BaseKustomization:    base.Kustomization,
		Deployment:           base.Deployment,
		Service:              base.Service,
		Ingress:              base.Ingress,
		Volume:               base.Volumes,
		OverlayKustomization: overlay.Kustomization,
		ArgoCDResource:       base.ArgoApplication,
		IngressPatch:         overlay.IngressPatch,
		ServicePatch:         overlay.ServicePatch,
		DeploymentPatch:      overlay.DeploymentPatch,
	}

	report, err := s.store.SaveApplication(applicationScaffold)
	if err != nil {
		return err
	}

	err = s.report.ReportCreateApplication(applicationScaffold, []*store.Report{report})
	if err != nil {
		return err
	}

	argoCDResourcePath := path.Join(relativeApplicationDir, "argocd-application.yaml")
	writeSuccessMessage(opts.Out, app.Name, argoCDResourcePath)

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

func inferApplicationFromStdinOrFile(stdin io.Reader, fs *afero.Afero, path string) (client.OkctlApplication, error) {
	var (
		err         error
		app         client.OkctlApplication
		inputReader io.Reader
	)

	switch path {
	case "-":
		inputReader = stdin
	default:
		inputReader, err = fs.Open(filepath.Clean(path))
		if err != nil {
			return app, fmt.Errorf("opening application file: %w", err)
		}
	}

	var buf []byte

	buf, err = ioutil.ReadAll(inputReader)
	if err != nil {
		return app, fmt.Errorf("reading application file: %w", err)
	}

	err = yaml.Unmarshal(buf, &app)
	if err != nil {
		return app, fmt.Errorf("parsing application yaml: %w", err)
	}

	return app, nil
}

// I'm assuming we'll be making enough customizations down the line to have our own okctlApplication, but for now
// mapping it to a Kaex application works fine
func okctlApplicationToKaexApplication(okctlApp client.OkctlApplication, primaryHostedZoneDomain string) (kaexApp *kaex.Application) {
	kaexApp = &kaex.Application{
		Name:            okctlApp.Name,
		Namespace:       okctlApp.Namespace,
		Image:           okctlApp.Image,
		Version:         okctlApp.Version,
		ImagePullSecret: okctlApp.ImagePullSecret,
		Url:             fmt.Sprintf("%s.%s", okctlApp.SubDomain, primaryHostedZoneDomain),
		Port:            okctlApp.Port,
		Replicas:        okctlApp.Replicas,
		Environment:     okctlApp.Environment,
		Volumes:         okctlApp.Volumes,
	}

	return kaexApp
}
