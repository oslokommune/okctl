package resources

import (
	"fmt"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

// CreateOkctlVolume creates persistent volume claims customized for okctl
func CreateOkctlVolume(app v1alpha1.Application, volume map[string]string) (corev1.PersistentVolumeClaim, error) {
	var (
		mountPath string
		size      string
	)

	for mountPath, size = range volume {
		break
	}

	pvc, err := createPersistentVolume(app, mountPath, size)
	if err != nil {
		return corev1.PersistentVolumeClaim{}, fmt.Errorf("error creating pvc: %w", err)
	}

	return pvc, nil
}


func createPersistentVolume(app v1alpha1.Application, path string, size string) (corev1.PersistentVolumeClaim, error) {
	volume := generateDefaultPVC()

	volume.ObjectMeta.Name = CreatePVCName(app, path)
	volume.ObjectMeta.Namespace = app.Metadata.Namespace

	capacity, err := createStorageRequest(size)
	if err != nil {
		return corev1.PersistentVolumeClaim{}, err
	}
	volume.Spec.Resources.Requests = capacity

	return volume, nil
}

func CreatePVCName(app v1alpha1.Application, path string) string {
	cleanPath := strings.Replace(path, "/", "", -1)

	return fmt.Sprintf("%s-%s", app.Metadata.Name, cleanPath)
}

func generateDefaultPVC() corev1.PersistentVolumeClaim {
	return corev1.PersistentVolumeClaim{
		TypeMeta: metav1.TypeMeta{
			Kind:       "PersistentVolumeClaim",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{},
		Spec: corev1.PersistentVolumeClaimSpec{
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceRequestsStorage: resource.Quantity{
						Format: "1Gi",
					},
				},
			},
			AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteMany},
		},
	}
}

func createStorageRequest(requestSize string) (corev1.ResourceList, error) {
	quantity, err := resource.ParseQuantity("1Gi")
	if requestSize != "" {
		quantity, err = resource.ParseQuantity(requestSize)

		if err != nil {
			return nil, err
		}
	}

	return corev1.ResourceList{
		corev1.ResourceStorage: quantity,
	}, nil
}

