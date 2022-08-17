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

	// https://kubernetes-sigs.github.io/aws-load-balancer-controller/v2.4/guide/tasks/ssl_redirect/
	pathType := networkingv1.PathTypeImplementationSpecific

	redirectPath := networkingv1.HTTPIngressPath{
		Path:     "/*",
		PathType: &pathType,
		Backend: networkingv1.IngressBackend{
			Service: &networkingv1.IngressServiceBackend{
				Name: "ssl-redirect",
				Port: networkingv1.ServiceBackendPort{
					Name: "use-annotation",
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
			APIVersion: apiVersion(),
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

func apiVersion() string {
	return fmt.Sprintf("%s/%s", networkingv1.SchemeGroupVersion.Group, networkingv1.SchemeGroupVersion.Version)
}

func createGenericIngress(app v1alpha1.Application) (networkingv1.Ingress, error) {
	hostURL, err := app.URL()
	if err != nil {
		return networkingv1.Ingress{}, fmt.Errorf("getting application URL: %w", err)
	}

	ingress := generateDefaultIngress()
	ingress.ObjectMeta.Namespace = app.Metadata.Namespace

	ingress.ObjectMeta.Name = app.Metadata.Name

	// https://kubernetes-sigs.github.io/aws-load-balancer-controller/v2.4/guide/tasks/ssl_redirect/
	pathType := networkingv1.PathTypeImplementationSpecific

	ingress.Spec.Rules = []networkingv1.IngressRule{
		{
			Host: hostURL.Host,
			IngressRuleValue: networkingv1.IngressRuleValue{
				HTTP: &networkingv1.HTTPIngressRuleValue{
					Paths: []networkingv1.HTTPIngressPath{{
						Path:     "/",
						PathType: &pathType,
						Backend: networkingv1.IngressBackend{
							Service: &networkingv1.IngressServiceBackend{
								Name: app.Metadata.Name,
								Port: networkingv1.ServiceBackendPort{
									Number: defaultServiceListeningPort,
								},
							},
							Resource: nil,
						},
					}},
				},
			},
		},
	}

	return ingress, nil
}
