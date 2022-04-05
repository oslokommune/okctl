package dryrun

import (
	"context"
	"fmt"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/client"
	"io"
)

type argocdService struct {
	out io.Writer
}

func (a argocdService) CreateArgoCD(_ context.Context, _ client.CreateArgoCDOpts) (*client.ArgoCD, error) {
	fmt.Fprintf(a.out, formatCreate("ArgoCD controller"))

	return &client.ArgoCD{}, nil
}

func (a argocdService) DeleteArgoCD(_ context.Context, _ client.DeleteArgoCDOpts) error {
	fmt.Fprintf(a.out, formatDelete("ArgoCD controller"))

	return nil
}

func (a argocdService) SetupApplicationsSync(_ context.Context, _ v1alpha1.Cluster) error {
	return nil
}
