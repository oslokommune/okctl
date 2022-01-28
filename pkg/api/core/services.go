package core

import "github.com/oslokommune/okctl/pkg/api"

// Services defines all available services
type Services struct {
	Cluster                    api.ClusterService
	Vpc                        api.VpcService
	ManagedPolicy              api.ManagedPolicyService
	ServiceAccount             api.ServiceAccountService
	Helm                       api.HelmService
	Kube                       api.KubeService
	Domain                     api.DomainService
	Certificate                api.CertificateService
	Parameter                  api.ParameterService
	IdentityManager            api.IdentityManagerService
	ComponentService           api.ComponentService
	ContainerRepositoryService api.ContainerRepositoryService
	SecurityGroupService       api.SecurityGroupService
}
