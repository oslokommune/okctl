package core

import (
	"github.com/oslokommune/okctl/pkg/breeze"
	"github.com/oslokommune/okctl/pkg/client"
)

// StateNodes contains all state storage nodes
type StateNodes struct {
	ArgoCD              breeze.Client
	Certificate         breeze.Client
	Cluster             breeze.Client
	Domain              breeze.Client
	ExternalDNS         breeze.Client
	Github              breeze.Client
	Manifest            breeze.Client
	Parameter           breeze.Client
	Vpc                 breeze.Client
	IdentityManager     breeze.Client
	Monitoring          breeze.Client
	Component           breeze.Client
	Helm                breeze.Client
	ManagedPolicy       breeze.Client
	ServiceAccount      breeze.Client
	ContainerRepository breeze.Client
	Upgrade             breeze.Client
}

// StateHandlers contains the state handlers
type StateHandlers struct {
	Helm                      client.HelmState
	ManagedPolicy             client.ManagedPolicyState
	ServiceAccount            client.ServiceAccountState
	Certificate               client.CertificateState
	IdentityManager           client.IdentityManagerState
	Github                    client.GithubState
	Manifest                  client.ManifestState
	Vpc                       client.VPCState
	Parameter                 client.ParameterState
	Domain                    client.DomainState
	ExternalDNS               client.ExternalDNSState
	Cluster                   client.ClusterState
	Component                 client.ComponentState
	Monitoring                client.MonitoringState
	ArgoCD                    client.ArgoCDState
	ContainerRepository       client.ContainerRepositoryState
	Loki                      client.LokiState
	Promtail                  client.PromtailState
	Tempo                     client.TempoState
	Autoscaler                client.AutoscalerState
	AWSLoadBalancerController client.AWSLoadBalancerControllerState
	Blockstorage              client.BlockstorageState
	ExternalSecrets           client.ExternalSecretsState
	Upgrade                   client.UpgradeState
}

// Services contains all client-side services
type Services struct {
	AWSLoadBalancerControllerService client.AWSLoadBalancerControllerService
	ArgoCD                           client.ArgoCDService
	ApplicationService               client.ApplicationService
	ApplicationManifestService       client.ApplicationManifestService
	ApplicationPostgresService       client.ApplicationPostgresService
	Certificate                      client.CertificateService
	Cluster                          client.ClusterService
	Domain                           client.DomainService
	ExternalDNS                      client.ExternalDNSService
	ExternalSecrets                  client.ExternalSecretsService
	Github                           client.GithubService
	Manifest                         client.ManifestService
	NameserverHandler                client.NSRecordDelegationService
	Parameter                        client.ParameterService
	Vpc                              client.VPCService
	IdentityManager                  client.IdentityManagerService
	Autoscaler                       client.AutoscalerService
	Blockstorage                     client.BlockstorageService
	Monitoring                       client.MonitoringService
	Component                        client.ComponentService
	Helm                             client.HelmService
	ManagedPolicy                    client.ManagedPolicyService
	ServiceAccount                   client.ServiceAccountService
	ContainerRepository              client.ContainerRepositoryService
	RemoteState                      client.RemoteStateService
}
