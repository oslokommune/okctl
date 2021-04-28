package resources

import (
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const defaultServiceListeningPort = 80

// CreateOkctlService creates a service customized for okctl
func CreateOkctlService(app v1alpha1.Application) corev1.Service {
	service := generateDefaultService()

	service.ObjectMeta.Name = app.Metadata.Name
	service.ObjectMeta.Namespace = app.Metadata.Namespace

	service.Spec.Selector = map[string]string{
		"app": app.Metadata.Name,
	}

	service.Spec.Ports[0].TargetPort = intstr.IntOrString{
		IntVal: app.Port,
	}

	service.Spec.Type = "NodePort"

	return service
}

func generateDefaultService() corev1.Service {
	return corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{{Port: defaultServiceListeningPort}},
			Type:  "ClusterIP",
		},
	}
}
