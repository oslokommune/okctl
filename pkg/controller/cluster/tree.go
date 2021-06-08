package cluster

import "github.com/oslokommune/okctl/pkg/controller/common/dependencytree"

// CreateResourceDependencyTree creates a tree
func CreateResourceDependencyTree() *dependencytree.Node {
	root := dependencytree.NewNode(dependencytree.NodeTypeServiceQuota)

	primaryHostedZoneNode := dependencytree.NewNode(dependencytree.NodeTypeZone)
	root.AppendChild(primaryHostedZoneNode)

	nameserverDelegationNode := dependencytree.NewNode(dependencytree.NodeTypeNameserverDelegator)
	primaryHostedZoneNode.AppendChild(nameserverDelegationNode)

	vpcNode := dependencytree.NewNode(dependencytree.NodeTypeVPC)
	primaryHostedZoneNode.AppendChild(vpcNode)

	cleanupSGNode := dependencytree.NewNode(dependencytree.NodeTypeCleanupSG)
	vpcNode.AppendChild(cleanupSGNode)

	clusterNode := dependencytree.NewNode(dependencytree.NodeTypeCluster)
	vpcNode.AppendChild(clusterNode)

	clusterNode.AppendChild(
		dependencytree.NewNode(dependencytree.NodeTypeCleanupALB),
		dependencytree.NewNode(dependencytree.NodeTypeExternalSecrets),
		dependencytree.NewNode(dependencytree.NodeTypeAutoscaler),
		dependencytree.NewNode(dependencytree.NodeTypeBlockstorage),
		dependencytree.NewNode(dependencytree.NodeTypeAWSLoadBalancerController),
		dependencytree.NewNode(dependencytree.NodeTypeExternalDNS),
		dependencytree.NewNode(dependencytree.NodeTypePostgres),
	)

	// All resources that requires SSL / a certificate needs the delegatedNameserversConfirmedNode as a dependency
	delegatedNameserversConfirmedNode := dependencytree.NewNode(dependencytree.NodeTypeNameserversDelegatedTest)
	clusterNode.AppendChild(delegatedNameserversConfirmedNode)

	identityProviderNode := dependencytree.NewNode(dependencytree.NodeTypeIdentityManager)
	delegatedNameserversConfirmedNode.AppendChild(identityProviderNode)

	identityProviderNode.AppendChild(
		dependencytree.NewNode(dependencytree.NodeTypeArgoCD),
		dependencytree.NewNode(dependencytree.NodeTypeUsers),
	)

	kubePromStack := dependencytree.NewNode(dependencytree.NodeTypeKubePromStack)
	identityProviderNode.AppendChild(kubePromStack)

	// This is not strictly required, but to a large extent it doesn't make much sense to setup Loki before
	// we have setup grafana.
	loki := dependencytree.NewNode(dependencytree.NodeTypeLoki)
	// Similarly, it doesn't make sense to install promtail without loki
	loki.AppendChild(dependencytree.NewNode(dependencytree.NodeTypePromtail))

	kubePromStack.AppendChild(
		loki,
		dependencytree.NewNode(dependencytree.NodeTypeTempo),
	)

	return root
}
