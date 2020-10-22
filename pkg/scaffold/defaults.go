package scaffold

import (
	"fmt"

	kaex "github.com/oslokommune/kaex/pkg/api"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

func createOkctlVolume(app kaex.Application, volume map[string]string) (corev1.PersistentVolumeClaim, error) {
	var (
		mountPath string
		size      string
	)

	for mountPath, size = range volume {
		break
	}

	pvc, err := kaex.CreatePersistentVolume(app, mountPath, size)
	if err != nil {
		return corev1.PersistentVolumeClaim{}, fmt.Errorf("error creating pvc: %w", err)
	}

	return pvc, nil
}

func createOkctlService(app kaex.Application) (corev1.Service, error) {
	service, err := kaex.CreateService(app)
	if err != nil {
		return corev1.Service{}, fmt.Errorf("error creating kaex service: %w", err)
	}

	service.Spec.Type = "NodePort"

	return service, nil
}

func createOkctlDeployment(app kaex.Application) (appsv1.Deployment, error) {
	deployment, err := kaex.CreateDeployment(app)
	if err != nil {
		return appsv1.Deployment{}, err
	}

	return deployment, nil
}
