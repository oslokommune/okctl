package scaffold

import (
	"fmt"

	kaex "github.com/oslokommune/kaex/pkg/api"
	networkingv1 "k8s.io/api/networking/v1beta1"
	intstrutil "k8s.io/apimachinery/pkg/util/intstr"
)

/*
CertificateCreatorFn handles generating a certificate for a certain host

Parameters:
fdqn string: The fully qualified domain name for the application e.g. my-app.my-cluster.oslo.systems

Returns:
certificateARN string: The certificate identifier e.g. arn:aws:acm:eu-west-1:42141252131:certificate/183509ca-0664-43b0-bb36-10b70c837597
*/
type CertificateCreatorFn func(fqdn string) (certificateARN string, err error)

func createOkctlIngress(app kaex.Application) (*networkingv1.Ingress, error) {
	ingress, err := kaex.CreateIngress(app)
	if err != nil {
		return nil, err
	}

	ingress.Spec.Rules[0].HTTP.Paths[0].Path = "/*"

	if ingress.Annotations == nil {
		ingress.Annotations = map[string]string{}
	}

	ingress.Annotations["kubernetes.io/ingress.class"] = "alb"
	ingress.Annotations["alb.ingress.kubernetes.io/scheme"] = "internet-facing"

	if ingress.Spec.TLS != nil {
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
	}

	return &ingress, nil
}

func createOkctlIngressOverlay(fn CertificateCreatorFn, host string) (patch Patch, err error) {
	certArn, err := fn(host)
	if err != nil {
		return Patch{}, fmt.Errorf("creating certificate: %w", err)
	}

	patch = Patch{Operations: []Operation{
		{
			Type: OperationTypeAdd,
			Path: "/metadata/annotations",
			Value: map[string]string{
				"alb.ingress.kubernetes.io/certificate-arn": certArn,
			},
		},
		{
			Type:  OperationTypeReplace,
			Path:  "/spec/rules/0/host",
			Value: "", // TOOD: acquire env specific host
		},
	}}

	return patch, nil
}
