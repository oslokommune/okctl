package argoapp

import (
	"github.com/oslokommune/kaex/pkg/api"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"path"
)

func generateDefaultArgoApp() *Application {
	return &Application{
		TypeMeta:   metav1.TypeMeta{
			Kind:       "Application",
			APIVersion: "argoproj.io/v1alpha1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:                       "",
		},
		Spec:       ApplicationSpec{
			Source:               ApplicationSource{
				RepoURL:        "",
				TargetRevision: "HEAD",
			},
			Destination:          ApplicationDestination{
				Server:    "https://kubernetes.default.svc",
			},
			Project:              "default",
			SyncPolicy:           &SyncPolicy{
				Automated:   &SyncPolicyAutomated{
					Prune:    false,
					SelfHeal: false,
				},
			},
		},
	}
}

func CreateArgoApp(app api.Application, repositoryURL string) (*Application, error) {
	argoApp := generateDefaultArgoApp()

	argoApp.ObjectMeta.Name = app.Name

	argoApp.Spec.Destination.Namespace = app.Name
	argoApp.Spec.Source.Path = path.Join("deployment", app.Name)
	argoApp.Spec.Source.RepoURL = repositoryURL

	return argoApp, nil
}
