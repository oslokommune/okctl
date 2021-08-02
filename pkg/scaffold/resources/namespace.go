package resources

import (
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CreateNamespace creates a namespace declaration based on an okctl application
func CreateNamespace(app v1alpha1.Application) corev1.Namespace {
	namespace := generateDefaultNamespace()

	namespace.ObjectMeta.Name = app.Metadata.Namespace

	return namespace
}

func generateDefaultNamespace() corev1.Namespace {
	return corev1.Namespace{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Namespace",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{},
	}
}
