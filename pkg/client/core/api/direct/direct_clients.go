package direct

import "github.com/oslokommune/okctl/pkg/client"

// ToolChain contains the direct clients
type ToolChain struct {
	AppPostgresIntegration client.ApplicationPostgresAPI
	Certificate            client.CertificateAPI
	Cluster                client.ClusterAPI
	Component              client.ComponentAPI
	ContainerRepo          client.ContainerRepositoryAPI
	ExternalDNS            client.ExternalDNSAPI
	Helm                   client.HelmAPI
	IdentityManager        client.IdentityManagerAPI
	ManagedPolicy          client.ManagedPolicyAPI
	Manifest               client.ManifestAPI
	Parameter              client.ParameterAPI
	SecuityGroup           client.SecurityGroupAPI
	ServiceAccount         client.ServiceAccountAPI
	Vpc                    client.VPCAPI
	Domain                 client.DomainAPI
}
