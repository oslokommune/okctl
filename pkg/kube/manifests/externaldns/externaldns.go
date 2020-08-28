// Package externaldns provides kubernetes manifests for deploy external dns
package externaldns

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/api/rbac/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const requiredFsGroup = 65534

// ExternalDNS contains the state for apply the external-dns
// manifests to kubernetes
type ExternalDNS struct {
	Namespace    string
	DomainFilter string
	Version      string
	OwnerID      string
	FsGroup      int64
	RunAsNonRoot bool
	Replicas     int32
	Ctx          context.Context
}

// New returns an initialised external-dns state
func New(hostedZoneID, domainFilter string) *ExternalDNS {
	return &ExternalDNS{
		Namespace:    "kube-system",
		DomainFilter: domainFilter,
		Version:      "v0.7.3",
		OwnerID:      hostedZoneID,
		FsGroup:      requiredFsGroup,
		RunAsNonRoot: true,
		Replicas:     1,
		Ctx:          context.Background(),
	}
}

// CreateDeployment creates the external-dns Deployment manifest
func (e *ExternalDNS) CreateDeployment(clientSet kubernetes.Interface) (interface{}, error) {
	deployClient := clientSet.AppsV1().Deployments(e.Namespace)
	return deployClient.Create(e.Ctx, e.DeploymentManifest(), metav1.CreateOptions{})
}

// DeploymentManifest returns the deployment manifest
func (e *ExternalDNS) DeploymentManifest() *appsv1.Deployment {
	return &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "external-dns",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &e.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "external-dns",
				},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "external-dns",
					},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  "external-dns",
							Image: fmt.Sprintf("registry.opensource.zalan.do/teapot/external-dns:%s", e.Version),
							Args: []string{
								"--source=service",
								"--source=ingress",
								fmt.Sprintf("--domain-filter=%s", e.DomainFilter),
								"--provider=aws",
								"--aws-zone-type=public",
								"--registry=txt",
								fmt.Sprintf("--txt-owner-id=%s", e.OwnerID),
							},
						},
					},
					ServiceAccountName: "external-dns",
					SecurityContext: &v1.PodSecurityContext{
						RunAsNonRoot: &e.RunAsNonRoot,
						FSGroup:      &e.FsGroup,
					},
				},
			},
			Strategy:                appsv1.DeploymentStrategy{},
			MinReadySeconds:         0,
			RevisionHistoryLimit:    nil,
			Paused:                  false,
			ProgressDeadlineSeconds: nil,
		},
		Status: appsv1.DeploymentStatus{},
	}
}

// CreateClusterRole creates the cluster role manifest
func (e *ExternalDNS) CreateClusterRole(clientSet kubernetes.Interface) (interface{}, error) {
	clusterRoleClient := clientSet.RbacV1beta1().ClusterRoles()
	return clusterRoleClient.Create(e.Ctx, e.ClusterRoleManifest(), metav1.CreateOptions{})
}

// ClusterRoleManifest returns the cluster role manifest
func (e *ExternalDNS) ClusterRoleManifest() *v1beta1.ClusterRole {
	return &v1beta1.ClusterRole{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ClusterRole",
			APIVersion: "rbac.authorization.k8s.io/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "external-dns",
		},
		Rules: []v1beta1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"get", "watch", "list"},
				Verbs:     []string{"services", "endpoints", "pods"},
			},
			{
				APIGroups: []string{"extensions", "networking.k8s.io"},
				Resources: []string{"ingresses"},
				Verbs:     []string{"get", "watch", "list"},
			},
			{
				APIGroups: []string{""},
				Resources: []string{"nodes"},
				Verbs:     []string{"list", "watch"},
			},
		},
	}
}

// CreateClusterRoleBinding creates the cluster role binding manifest
func (e *ExternalDNS) CreateClusterRoleBinding(clientSet kubernetes.Interface) (interface{}, error) {
	clusterRoleBindingClient := clientSet.RbacV1beta1().ClusterRoleBindings()
	return clusterRoleBindingClient.Create(e.Ctx, e.ClusterRoleBindingManifest(), metav1.CreateOptions{})
}

// ClusterRoleBindingManifest returns the cluster role binding manifest
func (e *ExternalDNS) ClusterRoleBindingManifest() *v1beta1.ClusterRoleBinding {
	return &v1beta1.ClusterRoleBinding{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ClusterRoleBinding",
			APIVersion: "rbac.authorization.k8s.io/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "external-dns-viewer",
		},
		Subjects: []v1beta1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      "external-dns",
				Namespace: "kube-system",
			},
		},
		RoleRef: v1beta1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     "external-dns",
		},
	}
}
