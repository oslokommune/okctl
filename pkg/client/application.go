package client

import (
	"context"
	"io"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client/store"
)

// ScaffoldApplicationOpts contains information necessary to scaffold application resources
type ScaffoldApplicationOpts struct {
	In io.Reader
	Out io.Writer

	ApplicationFilePath string
	RepoDir string

	Id           *api.ID
	HostedZoneID string
	IACRepoURL   string
}

// Validate ensures presented data is valid
func (o *ScaffoldApplicationOpts) Validate() error {
	return validation.ValidateStruct(o,
		validation.Field(&o.ApplicationFilePath, validation.Required),
		validation.Field(&o.Id, validation.Required),
		validation.Field(&o.HostedZoneID, validation.Required),
		validation.Field(&o.IACRepoURL, validation.Required),
	)
}

// ScaffoldedApplication contains information required by ApplicationStore and ApplicationReport
type ScaffoldedApplication struct {
	ApplicationName string

	KubernetesResources []byte
	ArgoCDResource      []byte
}

// ApplicationService applies the scaffolding API and produces the requested resources
type ApplicationService interface {
	ScaffoldApplication(context.Context, *ScaffoldApplicationOpts) error
}

// ApplicationStore handles writing deployment resources persistent storage
type ApplicationStore interface {
	SaveApplication(*ScaffoldedApplication) (*store.Report, error)
	RemoveApplication(string) (*store.Report, error)
}

// ApplicationReport handles writing output and progress
type ApplicationReport interface {
	ReportCreateApplication(*ScaffoldedApplication, []*store.Report) error
	ReportDeleteApplication([]*store.Report) error
}
