package core

import (
	"context"
	"fmt"
	"path"

	"github.com/mishudark/errors"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/config/constant"
	"github.com/oslokommune/okctl/pkg/jsonpatch"
	"github.com/oslokommune/okctl/pkg/scaffold"
	"github.com/spf13/afero"
)

const (
	defaultArgoCDApplicationManifestPermissions = 0o644 // u+rw g+r o+r
	defaultArgoCDApplicationManifestFilename    = "argocd-application.yaml"
)

type applicationService struct {
	certificateService    client.CertificateService
	appManifestService    client.ApplicationManifestService
	fs                    *afero.Afero
	absoluteRepositoryDir string
}

func (s *applicationService) createCertificate(ctx context.Context, id *api.ID, hostedZoneID, fqdn string) (string, error) {
	cert, certFnErr := s.certificateService.CreateCertificate(ctx, client.CreateCertificateOpts{
		ID:           *id,
		FQDN:         fqdn,
		Domain:       fqdn,
		HostedZoneID: hostedZoneID,
	})
	if certFnErr != nil {
		return "", certFnErr
	}

	return cert.ARN, nil
}

// ScaffoldApplication turns a file path into Kubernetes resources
// nolint: funlen
func (s *applicationService) ScaffoldApplication(ctx context.Context, opts *client.ScaffoldApplicationOpts) error {
	err := opts.Validate()
	if err != nil {
		return err
	}

	manifestSaver := generateManifestSaver(ctx, s.appManifestService, opts.Application.Metadata.Name)
	patchSaver := generatePatchSaver(ctx, s.appManifestService, opts.Cluster.Metadata.Name, opts.Application.Metadata.Name)

	err = scaffold.GenerateApplicationBase(scaffold.GenerateApplicationBaseOpts{
		SaveManifest: manifestSaver,
		Application:  opts.Application,
	})
	if err != nil {
		return fmt.Errorf("generating application base resources: %w", err)
	}

	certArn := ""

	if opts.Application.HasIngress() {
		certArn, err = s.createCertificate(
			ctx,
			&api.ID{
				Region:       opts.Cluster.Metadata.Region,
				AWSAccountID: opts.Cluster.Metadata.AccountID,
				ClusterName:  opts.Cluster.Metadata.Name,
			},
			opts.HostedZoneID,
			fmt.Sprintf("%s.%s", opts.Application.SubDomain, opts.Cluster.ClusterRootDomain),
		)
		if err != nil {
			return fmt.Errorf("creating certificate: %w", err)
		}
	}

	err = scaffold.GenerateApplicationOverlay(scaffold.GenerateApplicationOverlayOpts{
		SavePatch:      patchSaver,
		Application:    opts.Application,
		Domain:         opts.Cluster.ClusterRootDomain,
		CertificateARN: certArn,
	})
	if err != nil {
		return fmt.Errorf("generating application overlay: %w", err)
	}

	return nil
}

// CreateArgoCDApplicationManifest creates necessary files for the ArgoCD integration
func (s *applicationService) CreateArgoCDApplicationManifest(opts client.CreateArgoCDApplicationManifestOpts) error {
	relativeOverlayDir := path.Join(
		opts.Cluster.Github.OutputPath,
		constant.DefaultApplicationsOutputDir,
		opts.Application.Metadata.Name,
		constant.DefaultApplicationOverlayDir,
		opts.Cluster.Metadata.Name,
	)

	absoluteOverlayDir := path.Join(s.absoluteRepositoryDir, relativeOverlayDir)
	absoluteArgoCDApplicationManifestPath := path.Join(absoluteOverlayDir, defaultArgoCDApplicationManifestFilename)

	err := s.fs.MkdirAll(absoluteOverlayDir, 0o700)
	if err != nil {
		return errors.E(err, "ensuring overlay directory: %w", err)
	}

	err = scaffold.GenerateArgoCDApplicationManifest(scaffold.GenerateArgoCDApplicationManifestOpts{
		Saver: func(content []byte) error {
			return s.fs.WriteFile(absoluteArgoCDApplicationManifestPath, content, defaultArgoCDApplicationManifestPermissions)
		},
		Application:                   opts.Application,
		IACRepoURL:                    opts.Cluster.Github.URL(),
		RelativeApplicationOverlayDir: relativeOverlayDir,
	})
	if err != nil {
		return errors.E(err, "generating ArgoCD Application manifest")
	}

	return nil
}

// DeleteArgoCDApplicationManifest removes necessary files related to the ArgoCD integration
func (s *applicationService) DeleteArgoCDApplicationManifest(opts client.DeleteArgoCDApplicationManifestOpts) error {
	relativeOverlayDir := path.Join(
		opts.Cluster.Github.OutputPath,
		constant.DefaultApplicationsOutputDir,
		opts.Application.Metadata.Name,
		constant.DefaultApplicationOverlayDir,
		opts.Cluster.Metadata.Name,
	)

	absoluteOverlayDir := path.Join(s.absoluteRepositoryDir, relativeOverlayDir)
	absoluteArgoCDApplicationManifestPath := path.Join(absoluteOverlayDir, defaultArgoCDApplicationManifestFilename)

	err := s.fs.Remove(absoluteArgoCDApplicationManifestPath)
	if err != nil {
		return errors.E(err, "removing ArgoCD application manifest")
	}

	return nil
}

func generateManifestSaver(ctx context.Context, service client.ApplicationManifestService, applicationName string) scaffold.ManifestSaver {
	return func(filename string, content []byte) error {
		return service.SaveManifest(ctx, client.SaveManifestOpts{
			ApplicationName: applicationName,
			Filename:        filename,
			Content:         content,
		})
	}
}

func generatePatchSaver(ctx context.Context, service client.ApplicationManifestService, clusterName, applicationName string) scaffold.PatchSaver {
	return func(kind string, patch jsonpatch.Patch) error {
		return service.SavePatch(ctx, client.SavePatchOpts{
			ApplicationName: applicationName,
			ClusterName:     clusterName,
			Kind:            kind,
			Patch:           patch,
		})
	}
}

// NewApplicationService initializes a new Scaffold application service
func NewApplicationService(
	fs *afero.Afero,
	certificateService client.CertificateService,
	appManifestService client.ApplicationManifestService,
	absoluteRepositoryDir string,
) client.ApplicationService {
	return &applicationService{
		certificateService:    certificateService,
		appManifestService:    appManifestService,
		fs:                    fs,
		absoluteRepositoryDir: absoluteRepositoryDir,
	}
}
