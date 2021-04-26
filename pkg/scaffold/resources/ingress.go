package resources

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CreateOkctlIngress creates an ingress customized for okctl
func CreateOkctlIngress(app v1alpha1.Application) (networkingv1.Ingress, error) {
	ingress, err := createGenericIngress(app)
	if err != nil {
		return networkingv1.Ingress{}, err
	}

	ingress.Spec.Rules[0].HTTP.Paths[0].Path = "/*"

	ingress.Annotations["kubernetes.io/ingress.class"] = "alb"
	ingress.Annotations["alb.ingress.kubernetes.io/scheme"] = "internet-facing"
	ingress.Annotations["alb.ingress.kubernetes.io/listen-ports"] = "[{\"HTTP\": 80}, {\"HTTPS\": 443}]"
	ingress.Annotations["alb.ingress.kubernetes.io/actions.ssl-redirect"] =
		`{"Type": "redirect", "RedirectConfig": { "Protocol": "HTTPS", "Port": "443", "StatusCode": "HTTP_301"}}`

	redirectPath := networkingv1.HTTPIngressPath{
		Path: "/*",
		Backend: networkingv1.IngressBackend{
			Service: &networkingv1.IngressServiceBackend{
				Name: "ssl-redirect",
				Port: networkingv1.ServiceBackendPort{
					Name: "use-annotation", // TODO: Test
				},
			},
		},
	}

	ingress.Spec.Rules[0].HTTP.Paths = append([]networkingv1.HTTPIngressPath{redirectPath}, ingress.Spec.Rules[0].HTTP.Paths...)

	return ingress, nil
}

func generateDefaultIngress() networkingv1.Ingress {
	return networkingv1.Ingress{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Ingress",
			APIVersion: "networking.k8s.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Labels:      map[string]string{},
			Annotations: map[string]string{},
		},
		Spec: networkingv1.IngressSpec{
			Rules: make([]networkingv1.IngressRule, 1),
		},
	}
}

func createGenericIngress(app v1alpha1.Application) (networkingv1.Ingress, error) {
	hostURL, err := app.URL()
	if err != nil {
		return networkingv1.Ingress{}, fmt.Errorf("getting application URL: %w", err)
	}

	ingress := generateDefaultIngress()
	ingress.ObjectMeta.Namespace = app.Metadata.Namespace

	ingress.ObjectMeta.Name = app.Metadata.Name

	ingress.Spec.Rules = []networkingv1.IngressRule{
		{
			Host: hostURL.Host,
			IngressRuleValue: networkingv1.IngressRuleValue{
				HTTP: &networkingv1.HTTPIngressRuleValue{
					Paths: []networkingv1.HTTPIngressPath{{
						Path: "/",
						Backend: networkingv1.IngressBackend{
							Service: &networkingv1.IngressServiceBackend{
								Name: app.Metadata.Name,
								Port: networkingv1.ServiceBackendPort{
									Number: defaultServiceListeningPort,
								},
							},
						},
					}},
				},
			},
		},
	}

	return ingress, nil
}
