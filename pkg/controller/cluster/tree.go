package cluster

import "github.com/oslokommune/okctl/pkg/controller/common/resourcetree"

// CreateResourceDependencyTree creates a tree
func CreateResourceDependencyTree() *resourcetree.ResourceNode {
	root := resourcetree.NewNode(resourcetree.ResourceNodeTypeServiceQuota)

	primaryHostedZoneNode := resourcetree.NewNode(resourcetree.ResourceNodeTypeZone)
	root.AppendChild(primaryHostedZoneNode)

	nameserverDelegationNode := resourcetree.NewNode(resourcetree.ResourceNodeTypeNameserverDelegator)
	primaryHostedZoneNode.AppendChild(nameserverDelegationNode)

	vpcNode := resourcetree.NewNode(resourcetree.ResourceNodeTypeVPC)
	primaryHostedZoneNode.AppendChild(vpcNode)

	cleanupSGNode := resourcetree.NewNode(resourcetree.ResourceNodeTypeCleanupSG)
	vpcNode.AppendChild(cleanupSGNode)

	clusterNode := resourcetree.NewNode(resourcetree.ResourceNodeTypeCluster)
	vpcNode.AppendChild(clusterNode)

	clusterNode.AppendChild(
		resourcetree.NewNode(resourcetree.ResourceNodeTypeCleanupALB),
		resourcetree.NewNode(resourcetree.ResourceNodeTypeExternalSecrets),
		resourcetree.NewNode(resourcetree.ResourceNodeTypeAutoscaler),
		resourcetree.NewNode(resourcetree.ResourceNodeTypeBlockstorage),
		resourcetree.NewNode(resourcetree.ResourceNodeTypeAWSLoadBalancerController),
		resourcetree.NewNode(resourcetree.ResourceNodeTypeExternalDNS),
		resourcetree.NewNode(resourcetree.ResourceNodeTypePostgres),
	)

	// All resources that requires SSL / a certificate needs the delegatedNameserversConfirmedNode as a dependency
	delegatedNameserversConfirmedNode := resourcetree.NewNode(resourcetree.ResourceNodeTypeNameserversDelegatedTest)
	clusterNode.AppendChild(delegatedNameserversConfirmedNode)

	identityProviderNode := resourcetree.NewNode(resourcetree.ResourceNodeTypeIdentityManager)
	delegatedNameserversConfirmedNode.AppendChild(identityProviderNode)

	identityProviderNode.AppendChild(
		resourcetree.NewNode(resourcetree.ResourceNodeTypeArgoCD),
		resourcetree.NewNode(resourcetree.ResourceNodeTypeUsers),
	)

	kubePromStack := resourcetree.NewNode(resourcetree.ResourceNodeTypeKubePromStack)
	identityProviderNode.AppendChild(kubePromStack)

	// This is not strictly required, but to a large extent it doesn't make much sense to setup Loki before
	// we have setup grafana.
	loki := resourcetree.NewNode(resourcetree.ResourceNodeTypeLoki)
	// Similarly, it doesn't make sense to install promtail without loki
	loki.AppendChild(resourcetree.NewNode(resourcetree.ResourceNodeTypePromtail))

	kubePromStack.AppendChild(
		loki,
		resourcetree.NewNode(resourcetree.ResourceNodeTypeTempo),
	)

	return root
}
