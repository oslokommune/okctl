package resources

import (
	kaex "github.com/oslokommune/kaex/pkg/api"
	argo "github.com/oslokommune/okctl/internal/third_party/argoproj/argo-cd/pkg/apis/application/v1alpha1"
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
				RepoURL:        "git@github.com:<organization>/<infrastructure as code repository URL>",
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

// CreateArgoApp creates an ArgoCD definition customized for okctl
func CreateArgoApp(app kaex.Application, sourceRepositoryURL string, sourceRepositoryPath string) argo.Application {
	argoApp := generateDefaultArgoApp()

	argoApp.ObjectMeta.Name = app.Name

	if app.Namespace != "" {
		argoApp.Spec.Destination.Namespace = app.Namespace
	}

	argoApp.Spec.Source.Path = sourceRepositoryPath
	argoApp.Spec.Source.RepoURL = sourceRepositoryURL

	return argoApp
}
