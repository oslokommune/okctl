package dryrun

import (
	"github.com/oslokommune/okctl/pkg/client/core"
	"io"
)

// ClientServices knows how to create a service collection containing services that only report their actions
func ClientServices(writer io.Writer) *core.Services {
	return &core.Services{
		Domain:                           domainClient{out: writer},
		Vpc:                              vpcService{out: writer},
		NameserverHandler:                nameserverService{out: writer},
		Cluster:                          clusterService{out: writer},
		Autoscaler:                       autoscalerService{out: writer},
		AWSLoadBalancerControllerService: loadbalancerService{out: writer},
		Blockstorage:                     blockstorageService{out: writer},
		ExternalDNS:                      externalDNSService{out: writer},
		ExternalSecrets:                  externalSecretsService{out: writer},
		IdentityManager:                  identityManagerService{out: writer},
		ArgoCD:                           argocdService{out: writer},
		Github:                           githubService{out: writer},
		Monitoring:                       monitoringService{out: writer},
		Component:                        postgresService{out: writer},
	}
}
