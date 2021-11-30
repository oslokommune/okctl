package externalsecret

import typesv1 "github.com/oslokommune/okctl/pkg/kube/externalsecret/api/types/v1"

// SecretManifestOpts contains data used to create an external secret
type SecretManifestOpts struct {
	Name               string
	Namespace          string
	BackendType        string
	Annotations        map[string]string
	Labels             map[string]string
	Data               []typesv1.ExternalSecretData
	StringDataTemplate map[string]interface{}
}
