package direct

import "github.com/oslokommune/okctl/pkg/client"

// ToolChain contains the direct clients
type ToolChain struct {
	Helm   client.HelmAPI
	Domain client.DomainAPI
}
