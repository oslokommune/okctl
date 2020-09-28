package argoapp

import (
	"github.com/oslokommune/kaex/pkg/api"
	argo "github.com/oslokommune/okctl/internal/third_party/argoproj/argo-cd/pkg/apis/application/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"path"
)

func generateDefaultArgoApp() *argo.Application {
	return &argo.Application{
		TypeMeta:   metav1.TypeMeta{
			Kind:       "Application",
			APIVersion: "argoproj.io/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:                       "",
		},
		Spec:       argo.ApplicationSpec{
			Source:               argo.ApplicationSource{
				RepoURL:        "",
				TargetRevision: "HEAD",
			},
			Destination:          argo.ApplicationDestination{
				Server:    "https://kubernetes.default.svc",
			},
			Project:              "default",
			SyncPolicy:           &argo.SyncPolicy{
				Automated:   &argo.SyncPolicyAutomated{
					Prune:    false,
					SelfHeal: false,
				},
			},
		},
	}
}

func CreateArgoApp(app api.Application, repositoryURL string) (*argo.Application, error) {
	argoApp := generateDefaultArgoApp()

	argoApp.ObjectMeta.Name = app.Name

	argoApp.Spec.Destination.Namespace = app.Name
	argoApp.Spec.Source.Path = path.Join("deployment", app.Name)
	argoApp.Spec.Source.RepoURL = repositoryURL

	return argoApp, nil
}
