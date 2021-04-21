package resources

import (
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	intstrutil "k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/kubernetes/staging/src/k8s.io/apimachinery/pkg/util/intstr"
	"net/url"
	"strings"
)

// CreateOkctlIngress creates an ingress customized for okctl
func CreateOkctlIngress(app v1alpha1.Application) (networkingv1.Ingress, error) {
	ingress, err := createIngress(app)
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
			Service:  &networkingv1.IngressServiceBackend{
				Name: "ssl-redirect",
				Port: networkingv1.ServiceBackendPort{
					Name:   "",
					Number: intstrutil.Parse("use-annotation").StrVal,
				},
			},
			Resource: nil,
		},
	}
	//ServiceName: "ssl-redirect",
	//ServicePort: intstrutil.Parse("use-annotation"),

	ingress.Spec.Rules[0].HTTP.Paths = append([]networkingv1.HTTPIngressPath{redirectPath}, ingress.Spec.Rules[0].HTTP.Paths...)

	return ingress, nil
}

func generateDefaultIngress() networkingv1.Ingress {
	return networkingv1.Ingress{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Ingress",
			APIVersion: "networking.k8s.io/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        "name",
			Labels:      nil,
			Annotations: nil,
		},
		Spec: networkingv1.IngressSpec{
			Rules: make([]networkingv1.IngressRule, 1),
		},
	}
}


func createIngress(app v1alpha1.Application) (networkingv1.Ingress, error) {
	hostUrl, err := url.Parse(app.Url)
	if err != nil {
		return networkingv1.Ingress{}, err
	}

	ingress := generateDefaultIngress()
	ingress.ObjectMeta.Namespace = app.Metadata.Namespace

	ingress.ObjectMeta.Name = app.Metadata.Name
	ingress.ObjectMeta.Annotations = app.Ingress.Annotations

	ingress.Spec.Rules[0] = networkingv1.IngressRule{
		Host: hostUrl.Host,
		IngressRuleValue: networkingv1.IngressRuleValue{
			HTTP: &networkingv1.HTTPIngressRuleValue{
				Paths: []networkingv1.HTTPIngressPath{{
					Path: "/",
					Backend: networkingv1.IngressBackend{
						ServiceName: app.Metadata.Name,
						ServicePort: intstr.IntOrString{
							IntVal: 80,
						},
					},
				}},
			},
		},
	}

	if hostUrl.Scheme == "https" {
		ingress.Spec.TLS = []networkingv1.IngressTLS{
			{
				Hosts: []string{
					hostUrl.Host,
				},
				SecretName: strings.Join([]string{
					app.Metadata.Name,
					"tls",
				}, "-"),
			},
		}
	}

	return ingress, nil
}
