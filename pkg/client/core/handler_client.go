package core

import "github.com/oslokommune/okctl/pkg/client"

// Services contains all client-side services
type Services struct {
	ALBIngressController             client.ALBIngressControllerService
	AWSLoadBalancerControllerService client.AWSLoadBalancerControllerService
	ArgoCD                           client.ArgoCDService
	ApplicationService               client.ApplicationService
	Certificate                      client.CertificateService
	Cluster                          client.ClusterService
	Domain                           client.DomainService
	ExternalDNS                      client.ExternalDNSService
	ExternalSecrets                  client.ExternalSecretsService
	Github                           client.GithubService
	Manifest                         client.ManifestService
	NameserverHandler                client.NameserverRecordDelegationService
	Parameter                        client.ParameterService
	Vpc                              client.VPCService
	IdentityManager                  client.IdentityManagerService
	Autoscaler                       client.AutoscalerService
	Blockstorage                     client.BlockstorageService
	KubePromStack                    client.KubePromStackService
}
