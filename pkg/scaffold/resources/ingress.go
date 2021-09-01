package resources

import (
	"fmt"
	"github.com/oslokommune/okctl/pkg/config/constant"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	networkingv1beta1 "k8s.io/api/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// CreateOkctlIngress creates an ingress customized for okctl
func CreateOkctlIngress(app v1alpha1.Application) (networkingv1beta1.Ingress, error) {
	ingress, err := createGenericIngress(app)
	if err != nil {
		return networkingv1beta1.Ingress{}, err
	}

	ingress.Spec.Rules[0].HTTP.Paths[0].Path = "/*"

	ingress.Annotations["kubernetes.io/ingress.class"] = "alb"
	ingress.Annotations["alb.ingress.kubernetes.io/scheme"] = "internet-facing"
	ingress.Annotations["alb.ingress.kubernetes.io/listen-ports"] = "[{\"HTTP\": 80}, {\"HTTPS\": 443}]"
	ingress.Annotations["alb.ingress.kubernetes.io/actions.ssl-redirect"] =
		`{"Type": "redirect", "RedirectConfig": { "Protocol": "HTTPS", "Port": "443", "StatusCode": "HTTP_301"}}`

	redirectPath := networkingv1beta1.HTTPIngressPath{
		Path: "/*",
		Backend: networkingv1beta1.IngressBackend{
			ServiceName: "ssl-redirect",
			ServicePort: intstr.IntOrString{
				Type:   intstr.String,
				StrVal: "use-annotation",
			},
		},
	}

	ingress.Spec.Rules[0].HTTP.Paths = append([]networkingv1beta1.HTTPIngressPath{redirectPath}, ingress.Spec.Rules[0].HTTP.Paths...)

	return ingress, nil
}

func generateDefaultIngress() networkingv1beta1.Ingress {
	return networkingv1beta1.Ingress{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Ingress",
			APIVersion: "networking.k8s.io/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Labels:      map[string]string{},
			Annotations: map[string]string{},
		},
		Spec: networkingv1beta1.IngressSpec{
			Rules: make([]networkingv1beta1.IngressRule, 1),
		},
	}
}

func createGenericIngress(app v1alpha1.Application) (networkingv1beta1.Ingress, error) {
	hostURL, err := app.URL()
	if err != nil {
		return networkingv1beta1.Ingress{}, fmt.Errorf(constant.GetApplicationURLError, err)
	}

	ingress := generateDefaultIngress()
	ingress.ObjectMeta.Namespace = app.Metadata.Namespace

	ingress.ObjectMeta.Name = app.Metadata.Name

	ingress.Spec.Rules = []networkingv1beta1.IngressRule{
		{
			Host: hostURL.Host,
			IngressRuleValue: networkingv1beta1.IngressRuleValue{
				HTTP: &networkingv1beta1.HTTPIngressRuleValue{
					Paths: []networkingv1beta1.HTTPIngressPath{{
						Path: "/",
						Backend: networkingv1beta1.IngressBackend{
							ServiceName: app.Metadata.Name,
							ServicePort: intstr.IntOrString{
								IntVal: defaultServiceListeningPort,
							},
						},
					}},
				},
			},
		},
	}

	return ingress, nil
}
