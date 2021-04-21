package resources

import (
	argo "github.com/oslokommune/okctl/internal/third_party/argoproj/argo-cd/pkg/apis/application/v1alpha1"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func generateDefaultArgoApp() argo.Application {
	return argo.Application{
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
				RepoURL:        "git@github.com:<organization>/<infrastructure as code repository URL>", // TODO: Remove placeholder
				TargetRevision: "HEAD",
			},
			Destination: argo.ApplicationDestination{
				Server:    "https://kubernetes.default.svc",
				Namespace: "<namespace your app should run in>", // TODO: Remove placeholder
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

// CreateArgoApp creates an ArgoCD definition customized for okctl
func CreateArgoApp(app v1alpha1.Application, sourceRepositoryURL string, sourceRepositoryPath string) argo.Application {
	argoApp := generateDefaultArgoApp()

	argoApp.ObjectMeta.Name = app.Metadata.Name

	if app.Metadata.Namespace != "" {
		argoApp.Spec.Destination.Namespace = app.Metadata.Namespace
	}

	argoApp.Spec.Source.Path = sourceRepositoryPath
	argoApp.Spec.Source.RepoURL = sourceRepositoryURL

	return argoApp
}
