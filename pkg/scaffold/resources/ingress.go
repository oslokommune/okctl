package resources

import (
	kaex "github.com/oslokommune/kaex/pkg/api"
	networkingv1 "k8s.io/api/networking/v1beta1"
	intstrutil "k8s.io/apimachinery/pkg/util/intstr"
)

// CreateOkctlIngress creates an ingress customized for okctl
func CreateOkctlIngress(app kaex.Application) (networkingv1.Ingress, error) {
	ingress, err := kaex.CreateIngress(app)
	if err != nil {
		return networkingv1.Ingress{}, err
	}

	ingress.Spec.Rules[0].HTTP.Paths[0].Path = "/*"

	if ingress.Annotations == nil {
		ingress.Annotations = map[string]string{}
	}

	ingress.Annotations["kubernetes.io/ingress.class"] = "alb"
	ingress.Annotations["alb.ingress.kubernetes.io/scheme"] = "internet-facing"
	ingress.Annotations["alb.ingress.kubernetes.io/listen-ports"] = "[{\"HTTP\": 80}, {\"HTTPS\": 443}]"
	ingress.Annotations["alb.ingress.kubernetes.io/actions.ssl-redirect"] =
		`{"Type": "redirect", "RedirectConfig": { "Protocol": "HTTPS", "Port": "443", "StatusCode": "HTTP_301"}}`

	redirectPath := networkingv1.HTTPIngressPath{
		Path: "/*",
		Backend: networkingv1.IngressBackend{
			ServiceName: "ssl-redirect",
			ServicePort: intstrutil.Parse("use-annotation"),
		},
	}

	ingress.Spec.Rules[0].HTTP.Paths = append([]networkingv1.HTTPIngressPath{redirectPath}, ingress.Spec.Rules[0].HTTP.Paths...)

	return ingress, nil
}
