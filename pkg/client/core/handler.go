package core

import "github.com/oslokommune/okctl/pkg/client"

// Services contains all client-side services
type Services struct {
	ALBIngressController client.ALBIngressControllerService
	ArgoCD               client.ArgoCDService
	Certificate          client.CertificateService
	Cluster              client.ClusterService
	Domain               client.DomainService
	ExternalDNS          client.ExternalDNSService
	ExternalSecrets      client.ExternalSecretsService
	Github               client.GithubService
	Manifest             client.ManifestService
	Parameter            client.ParameterService
	Vpc                  client.VPCService
}
