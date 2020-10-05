package scaffold

import (
	"fmt"
	"path"

	"github.com/oslokommune/kaex/pkg/api"
	argo "github.com/oslokommune/okctl/internal/third_party/argoproj/argo-cd/pkg/apis/application/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	networkingv1 "k8s.io/api/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	corev1 "k8s.io/api/core/v1"
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

func createArgoApp(app api.Application, repositoryURL string) *argo.Application {
	argoApp := generateDefaultArgoApp()

	argoApp.ObjectMeta.Name = app.Name

	if app.Namespace != "" {
		argoApp.Spec.Destination.Namespace = app.Namespace
	}

	argoApp.Spec.Source.Path = fmt.Sprintf("%s/", path.Join("deployment", app.Name))

	if repositoryURL != "" {
		argoApp.Spec.Source.RepoURL = repositoryURL
	}

	return argoApp
}

func createOkctlVolume(app api.Application, volume map[string]string) (corev1.PersistentVolumeClaim, error) {
	var (
		mountPath string
		size      string
	)

	for mountPath, size = range volume {
		break
	}

	pvc, err := api.CreatePersistentVolume(app, mountPath, size)
	if err != nil {
		return corev1.PersistentVolumeClaim{}, fmt.Errorf("error creating pvc: %w", err)
	}

	return pvc, nil
}

func createOkctlService(app api.Application) (corev1.Service, error) {
	service, err := api.CreateService(app)
	if err != nil {
		return corev1.Service{}, fmt.Errorf("error creating kaex service: %w", err)
	}

	service.Spec.Type = "NodePort"

	return service, nil
}

func createOkctlIngress(app api.Application) (networkingv1.Ingress, error) {
	ingress, err := api.CreateIngress(app)
	if err != nil {
		return networkingv1.Ingress{}, err
	}

	ingress.Annotations["kubernetes.io/ingress.class"] = "alb"
	ingress.Annotations["alb.ingress.kubernetes.io/scheme"] = "internet-facing"

	return ingress, nil
}

func createOkctlDeployment(app api.Application) (appsv1.Deployment, error) {
	deployment, err := api.CreateDeployment(app)
	if err != nil {
		return appsv1.Deployment{}, err
	}

	return deployment, nil
}
