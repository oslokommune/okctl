package resources

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/kube/securitygrouppolicy/api/types/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CreateSecurityGroupPolicy creates an initialized security group policy based on an okctl application
func CreateSecurityGroupPolicy(app v1alpha1.Application) v1beta1.SecurityGroupPolicy {
	sgp := generateDefaultSecurityGroupPolicy()

	sgp.ObjectMeta.Name = app.Metadata.Name
	sgp.Spec.PodSelector.MatchLabels = map[string]string{
		"app": app.Metadata.Name,
	}

	return sgp
}

func generateDefaultSecurityGroupPolicy() v1beta1.SecurityGroupPolicy {
	return v1beta1.SecurityGroupPolicy{
		TypeMeta: metav1.TypeMeta{
			APIVersion: fmt.Sprintf("%s/%s", v1beta1.GroupName, v1beta1.GroupVersion),
			Kind:       v1beta1.SecurityGroupPolicyKind,
		},
		ObjectMeta: metav1.ObjectMeta{},
		Spec: v1beta1.SecurityGroupPolicySpec{
			PodSelector: v1beta1.SecurityGroupPolicyPodSelector{},
			SecurityGroups: v1beta1.SecurityGroupPolicySecurityGroups{
				GroupIDs: []string{},
			},
		},
	}
}
