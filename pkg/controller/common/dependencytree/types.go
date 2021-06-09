package dependencytree

// NodeType defines what type of resource a Node represents
type NodeType string

const (
	// NodeTypeGroup represents a node that has no actions associated with it. For now, only the root node
	NodeTypeGroup NodeType = "group"
	// NodeTypeZone represents a HostedZone resource
	NodeTypeZone NodeType = "hosted-zone"
	// NodeTypeVPC represents a VPC resource
	NodeTypeVPC NodeType = "vpc"
	// NodeTypeCluster represents a EKS cluster resource
	NodeTypeCluster NodeType = "cluster"
	// NodeTypeExternalSecrets represents an External Secrets resource
	NodeTypeExternalSecrets NodeType = "external-secrets"
	// NodeTypeAutoscaler represents an autoscaler resource
	NodeTypeAutoscaler NodeType = "autoscaler"
	// NodeTypeBlockstorage represents a blockstorage resource
	NodeTypeBlockstorage NodeType = "blockstorage"
	// NodeTypeKubePromStack represents a kubernetes-prometheus-stack resource
	NodeTypeKubePromStack NodeType = "kubernetes-prometheus-stack"
	// NodeTypeLoki represents a loki resource
	NodeTypeLoki NodeType = "loki"
	// NodeTypePromtail represents a promtail deployment
	NodeTypePromtail NodeType = "promtail"
	// NodeTypeTempo represents a Tempo deployment
	NodeTypeTempo NodeType = "tempo"
	// NodeTypeAWSLoadBalancerController represents an AWS load balancer controller resource
	NodeTypeAWSLoadBalancerController NodeType = "aws-load-balancer-controller"
	// NodeTypeExternalDNS represents an External DNS resource
	NodeTypeExternalDNS NodeType = "external-dns"
	// NodeTypeIdentityManager represents a Identity Manager resource
	NodeTypeIdentityManager NodeType = "identity-manager"
	// NodeTypeArgoCD represents an ArgoCD resource
	NodeTypeArgoCD NodeType = "argocd"
	// NodeTypeNameserverDelegator represents delegation of nameservers for a HostedZone
	NodeTypeNameserverDelegator NodeType = "nameserver-delegator"
	// NodeTypeNameserversDelegatedTest represents testing if nameservers has been successfully delegated
	NodeTypeNameserversDelegatedTest NodeType = "nameserver-delegator-test"
	// NodeTypeUsers represents the users we want to add to the cognito user pool
	NodeTypeUsers NodeType = "users"
	// NodeTypePostgres represents the postgres databases we want to add to the cluster
	NodeTypePostgres NodeType = "postgres"
	// NodeTypePostgresInstance represents a postgres instance
	NodeTypePostgresInstance NodeType = "postgres-instance"
	// NodeTypeApplication represents an okctl application resource
	NodeTypeApplication NodeType = "application"
	// NodeTypeCleanupALB represents a cleanup of ALBs
	NodeTypeCleanupALB NodeType = "cleanup-alb"
	// NodeTypeCleanupSG represents a cleanup of SecurityGroups
	NodeTypeCleanupSG NodeType = "cleanup-sg"
	// NodeTypeServiceQuota represents a service quota check
	NodeTypeServiceQuota NodeType = "service-quota"
	// NodeTypeContainerRepository represents a container repository
	NodeTypeContainerRepository NodeType = "container-repository"
)

// NodeState defines what state the resource is in, used to infer what action to take
type NodeState int

const (
	// NodeStateNoop represents a state where no action is needed. E.g.: if the desired state of the
	// resource conforms with the actual state
	NodeStateNoop NodeState = iota
	// NodeStatePresent represents the state where the resource exists
	NodeStatePresent
	// NodeStateAbsent represents the state where the resource does not exist
	NodeStateAbsent
)

// Node represents a component of the cluster and its dependencies
type Node struct {
	Type     NodeType
	State    NodeState
	Data     interface{}
	Children []*Node
}

// ApplyFn is a kind of function we can run on all the nodes in a Node tree with the ApplyFunction() function
type ApplyFn func(receiver *Node)

// ApplyFnWithTarget is a kind of function we can run on all the nodes in a Node tree with the ApplyFunction() function
type ApplyFnWithTarget func(receiver *Node, target *Node)
