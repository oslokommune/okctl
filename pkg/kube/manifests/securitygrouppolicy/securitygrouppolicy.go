// Package securitygrouppolicy provides a convenient way of
// interacting with the SecurityGroupPolicy CRD
package securitygrouppolicy

import (
	"context"

	v1beta1client "github.com/oslokommune/okctl/pkg/kube/securitygrouppolicy/clientset/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/oslokommune/okctl/pkg/kube/securitygrouppolicy/api/types/v1beta1"
	restclient "k8s.io/client-go/rest"
)

// SecurityGroupPolicy contains the required state
type SecurityGroupPolicy struct {
	Name      string
	Namespace string
	Manifest  *v1beta1.SecurityGroupPolicy

	Config *restclient.Config
	Ctx    context.Context
}

// New returns an initialised client for interacting with security group policies
func New(name, namespace string, manifest *v1beta1.SecurityGroupPolicy, config *restclient.Config) *SecurityGroupPolicy {
	return &SecurityGroupPolicy{
		Name:      name,
		Namespace: namespace,
		Manifest:  manifest,
		Config:    config,
		Ctx:       context.Background(),
	}
}

// Create the security group policy
func (s *SecurityGroupPolicy) Create() (*v1beta1.SecurityGroupPolicy, error) {
	client, err := v1beta1client.NewForConfig(s.Config)
	if err != nil {
		return nil, err
	}

	p, err := client.SecurityGroupPolicy(s.Namespace).Create(s.Ctx, s.Manifest)
	if err != nil {
		return nil, err
	}

	return p, nil
}

// Delete the security group policy
func (s *SecurityGroupPolicy) Delete() error {
	client, err := v1beta1client.NewForConfig(s.Config)
	if err != nil {
		return err
	}

	policy := metav1.DeletePropagationForeground

	return client.SecurityGroupPolicy(s.Namespace).Delete(s.Ctx, s.Name, metav1.DeleteOptions{
		PropagationPolicy: &policy,
	})
}

// Manifest returns a SecurityGroupPolicy manifest
func Manifest(name, namespace string, matchLabels map[string]string, securityGroups []string) *v1beta1.SecurityGroupPolicy {
	return &v1beta1.SecurityGroupPolicy{
		TypeMeta: metav1.TypeMeta{
			Kind:       "SecurityGroupPolicy",
			APIVersion: "vpcresources.k8s.aws/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: v1beta1.SecurityGroupPolicySpec{
			PodSelector: v1beta1.SecurityGroupPolicyPodSelector{
				MatchLabels: matchLabels,
			},
			SecurityGroups: v1beta1.SecurityGroupPolicySecurityGroups{
				GroupIDs: securityGroups,
			},
		},
	}
}
