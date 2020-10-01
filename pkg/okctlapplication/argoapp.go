// Package okctlapplication knows how to create an Argo Application populated with sensible data
package okctlapplication

import (
	"fmt"
	"path"

	"github.com/oslokommune/kaex/pkg/api"
	argo "github.com/oslokommune/okctl/internal/third_party/argoproj/argo-cd/pkg/apis/application/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func generateDefaultArgoApp() *argo.Application {
	return &argo.Application{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Application",
			APIVersion: "argoproj.io/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "",
			Namespace: "argocd",
		},
		Spec: argo.ApplicationSpec{
			Source: argo.ApplicationSource{
				RepoURL:        "<iac repo url>",
				TargetRevision: "HEAD",
			},
			Destination: argo.ApplicationDestination{
				Server:    "https://kubernetes.default.svc",
				Namespace: "<namespace your app should run in>",
			},
			Project: "default",
			SyncPolicy: &argo.SyncPolicy{
				Automated: &argo.SyncPolicyAutomated{
					Prune:    false,
					SelfHeal: false,
				},
			},
		},
	}
}

// CreateArgoApp creates an Argo Application struct based on a set of default values combined with data input
func CreateArgoApp(app api.Application, repositoryURL string) (*argo.Application, error) {
	argoApp := generateDefaultArgoApp()

	argoApp.ObjectMeta.Name = app.Name

	if app.Namespace != "" {
		argoApp.Spec.Destination.Namespace = app.Namespace
	}

	argoApp.Spec.Source.Path = fmt.Sprintf("%s/", path.Join("deployment", app.Name))
	argoApp.Spec.Source.RepoURL = repositoryURL

	return argoApp, nil
}
