package core

import (
	stormpkg "github.com/asdine/storm/v3"
	"github.com/oslokommune/okctl/pkg/client"
)

// StateNodes contains all state storage nodes
type StateNodes struct {
	ArgoCD              stormpkg.Node
	Certificate         stormpkg.Node
	Cluster             stormpkg.Node
	Domain              stormpkg.Node
	ExternalDNS         stormpkg.Node
	Github              stormpkg.Node
	Manifest            stormpkg.Node
	Parameter           stormpkg.Node
	Vpc                 stormpkg.Node
	IdentityManager     stormpkg.Node
	Monitoring          stormpkg.Node
	Component           stormpkg.Node
	Helm                stormpkg.Node
	ManagedPolicy       stormpkg.Node
	ServiceAccount      stormpkg.Node
	ContainerRepository stormpkg.Node
}

// StateHandlers contains the state handlers
type StateHandlers struct {
	Helm                client.HelmState
	ManagedPolicy       client.ManagedPolicyState
	ServiceAccount      client.ServiceAccountState
	Certificate         client.CertificateState
	IdentityManager     client.IdentityManagerState
	Github              client.GithubState
	Manifest            client.ManifestState
	Vpc                 client.VPCState
	Parameter           client.ParameterState
	Domain              client.DomainState
	ExternalDNS         client.ExternalDNSState
	Cluster             client.ClusterState
	Component           client.ComponentState
	Monitoring          client.MonitoringState
	ArgoCD              client.ArgoCDState
	ContainerRepository client.ContainerRepositoryState
}

// Services contains all client-side services
type Services struct {
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
}
