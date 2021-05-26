package resources

import (
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CreateOkctlDeployment creates a deployment customized for okctl
func CreateOkctlDeployment(app v1alpha1.Application) appsv1.Deployment {
	deployment := generateDefaultDeployment()

	deployment.ObjectMeta.Name = app.Metadata.Name
	deployment.ObjectMeta.Namespace = app.Metadata.Namespace

	deployment.Spec.Replicas = &app.Replicas

	deployment.Spec.Selector.MatchLabels = map[string]string{
		"app": app.Metadata.Name,
	}

	if app.ImagePullSecret != "" {
		deployment.Spec.Template.Spec.ImagePullSecrets = []corev1.LocalObjectReference{
			{Name: app.ImagePullSecret},
		}
	}

	deployment.Spec.Template.ObjectMeta.Labels = map[string]string{
		"app": app.Metadata.Name,
	}

	deployment.Spec.Template.Spec.Volumes = createVolumes(app)
	deployment.Spec.Template.Spec.Containers = createContainers(app)

	return deployment
}

func generateDefaultDeployment() appsv1.Deployment {
	replicas := int32(1)

	return appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "",
					Annotations: nil,
				},
				Spec: corev1.PodSpec{
					DNSPolicy: corev1.DNSDefault,
					Volumes:   nil,
				},
			},
		},
	}
}

func createContainers(app v1alpha1.Application) []corev1.Container {
	var (
		envVars  = make([]corev1.EnvVar, len(app.Environment))
		envCount = 0
	)

	for key, value := range app.Environment {
		envVars[envCount] = corev1.EnvVar{Name: key, Value: value}

		envCount++
	}

	volumeMounts := make([]corev1.VolumeMount, len(app.Volumes))

	for index, volume := range app.Volumes {
		for path := range volume {
			volumeMounts[index] = corev1.VolumeMount{
				Name:      createPVCName(app, path),
				MountPath: path,
			}
		}
	}

	containers := []corev1.Container{{
		Name:         app.Metadata.Name,
		Env:          envVars,
		VolumeMounts: volumeMounts,
	}}

	return containers
}

func createVolumes(app v1alpha1.Application) []corev1.Volume {
	volumes := make([]corev1.Volume, len(app.Volumes))

	for index, volume := range app.Volumes {
		for path := range volume {
			volumes[index] = corev1.Volume{
				Name: createPVCName(app, path),
				VolumeSource: corev1.VolumeSource{
					PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
						ClaimName: createPVCName(app, path),
					},
				},
			}

			break
		}
	}

	return volumes
}
