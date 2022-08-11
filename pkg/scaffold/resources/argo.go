package resources

import (
	argo "github.com/oslokommune/okctl/internal/third_party/argoproj/argo-cd/pkg/apis/application/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GenerateDefaultArgoApp knows how to construct a default ArgoCD application
func GenerateDefaultArgoApp() argo.Application {
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
				RepoURL:        "",
				TargetRevision: "HEAD",
			},
			Destination: argo.ApplicationDestination{
				Server:    "https://kubernetes.default.svc",
				Namespace: "default",
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
func CreateArgoApp(name string, namespace string, sourceRepositoryURL string, sourceRepositoryPath string) argo.Application {
	argoApp := GenerateDefaultArgoApp()

	argoApp.ObjectMeta.Name = name

	if namespace != "" {
		argoApp.Spec.Destination.Namespace = namespace
	}

	argoApp.Spec.Source.Path = sourceRepositoryPath
	argoApp.Spec.Source.RepoURL = sourceRepositoryURL

	return argoApp
}
