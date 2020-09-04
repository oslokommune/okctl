// Package client provides convenience functions for invoking API operations
package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/sanity-io/litter"
)

const (
	targetVpcs                               = "vpcs/"
	targetClusters                           = "clusters/"
	targetExternalSecretsPolicy              = "managedpolicies/externalsecrets/"
	targetExternalSecretsServiceAccount      = "serviceaccounts/externalsecrets/"
	targetExternalSecretsHelm                = "helm/externalsecrets/"
	targetAlbIngressControllerPolicy         = "managedpolicies/albingresscontroller/"
	targetAlbIngressControllerServiceAccount = "serviceaccounts/albingresscontroller/"
	targetAlbIngressControllerHelm           = "helm/albingresscontroller/"
	targetExternalDNSPolicy                  = "managedpolicies/externaldns/"
	targetExternalDNSServiceAccount          = "serviceaccounts/externaldns/"
	targetDomain                             = "domains/"
	targetKubeExternalDNS                    = "kube/externaldns/"
	targetCertificate                        = "certificates/"
	targetParameterSecret                    = "parameters/secret/"
	targetHelmArgoCD                         = "helm/argocd/"
	targetKubeExternalSecret                 = "kube/externalsecrets/"
)

// Creator contains all the creation operations towards the API
type Creator interface {
	CreateCluster(opts *api.ClusterCreateOpts) (*api.Cluster, error)
	CreateVpc(opts *api.CreateVpcOpts) (*api.Vpc, error)
	CreateExternalSecretsPolicy(opts *api.CreateExternalSecretsPolicyOpts) (*api.ManagedPolicy, error)
	CreateAlbIngressControllerPolicy(opts *api.CreateAlbIngressControllerPolicyOpts) (*api.ManagedPolicy, error)
	CreateExternalDNSPolicy(opts *api.CreateExternalDNSPolicyOpts) (*api.ManagedPolicy, error)
	CreateExternalSecretsServiceAccount(opts *api.CreateExternalSecretsServiceAccountOpts) (*api.ServiceAccount, error)
	CreateAlbIngressControllerServiceAccount(opts *api.CreateAlbIngressControllerServiceAccountOpts) (*api.ServiceAccount, error)
	CreateExternalDNSServiceAccount(opts *api.CreateExternalDNSServiceAccountOpts) (*api.ServiceAccount, error)
	CreateExternalSecretsHelmChart(opts *api.CreateExternalSecretsHelmChartOpts) (*api.Helm, error)
	CreateAlbIngressControllerHelmChart(opts *api.CreateAlbIngressControllerHelmChartOpts) (*api.Helm, error)
	CreateArgoCD(opts *api.CreateArgoCDOpts) (*api.Helm, error)
	CreateExternalDNSKubeDeployment(opts *api.CreateExternalDNSKubeDeploymentOpts) (*api.Kube, error)
	CreateExternalSecrets(opts *api.CreateExternalSecretsOpts) (*api.Kube, error)
	CreateDomain(opts *api.CreateDomainOpts) (*api.Domain, error)
	CreateCertificate(opts *api.CreateCertificateOpts) (*api.Certificate, error)
	CreateSecret(opts *api.CreateSecretOpts) (*api.SecretParameter, error)
}

// Destroyer contains all destructive operations towards the API
type Destroyer interface {
	DeleteCluster(opts *api.ClusterDeleteOpts) error
	DeleteVpc(opts *api.DeleteVpcOpts) error
}

// Ensure that the Client implements the Creator interface
var _ Creator = &Client{}

// Ensure that the Client implements the Destroyer interface
var _ Destroyer = &Client{}

// Client stores state for invoking API operations
type Client struct {
	BaseURL  string
	Client   *http.Client
	Progress io.Writer
	Debug    bool
}

// New returns a client that wraps the common API operations
func New(debug bool, progress io.Writer, serverURL string) *Client {
	return &Client{
		Progress: progress,
		BaseURL:  serverURL,
		Client:   &http.Client{},
		Debug:    debug,
	}
}

// CreateExternalSecrets invokes the external secrets creation operation
func (c *Client) CreateExternalSecrets(opts *api.CreateExternalSecretsOpts) (*api.Kube, error) {
	into := &api.Kube{}
	return into, c.DoPost(targetKubeExternalSecret, opts, into)
}

// CreateArgoCD invokes the argocd creation operation
func (c *Client) CreateArgoCD(opts *api.CreateArgoCDOpts) (*api.Helm, error) {
	into := &api.Helm{}
	return into, c.DoPost(targetHelmArgoCD, opts, into)
}

// CreateSecret invokes the secret creation operation
func (c *Client) CreateSecret(opts *api.CreateSecretOpts) (*api.SecretParameter, error) {
	into := &api.SecretParameter{}
	return into, c.DoPost(targetParameterSecret, opts, into)
}

// CreateCertificate invokes the certificate creation operation
func (c *Client) CreateCertificate(opts *api.CreateCertificateOpts) (*api.Certificate, error) {
	into := &api.Certificate{}
	return into, c.DoPost(targetCertificate, opts, into)
}

// CreateExternalDNSKubeDeployment invokes the external dns kube deployment
func (c *Client) CreateExternalDNSKubeDeployment(opts *api.CreateExternalDNSKubeDeploymentOpts) (*api.Kube, error) {
	into := &api.Kube{}
	return into, c.DoPost(targetKubeExternalDNS, opts, into)
}

// CreateDomain invokes the domain creation
func (c *Client) CreateDomain(opts *api.CreateDomainOpts) (*api.Domain, error) {
	into := &api.Domain{}
	return into, c.DoPost(targetDomain, opts, into)
}

// CreateExternalDNSPolicy invokes the external dns policy creation
func (c *Client) CreateExternalDNSPolicy(opts *api.CreateExternalDNSPolicyOpts) (*api.ManagedPolicy, error) {
	into := &api.ManagedPolicy{}
	return into, c.DoPost(targetExternalDNSPolicy, opts, into)
}

// CreateExternalDNSServiceAccount invokes the external dns service account creation
func (c *Client) CreateExternalDNSServiceAccount(opts *api.CreateExternalDNSServiceAccountOpts) (*api.ServiceAccount, error) {
	into := &api.ServiceAccount{}
	return into, c.DoPost(targetExternalDNSServiceAccount, opts, into)
}

// CreateAlbIngressControllerHelmChart invokes the alb ingress controller helm chart creator
func (c *Client) CreateAlbIngressControllerHelmChart(opts *api.CreateAlbIngressControllerHelmChartOpts) (*api.Helm, error) {
	into := &api.Helm{}
	return into, c.DoPost(targetAlbIngressControllerHelm, opts, into)
}

// CreateAlbIngressControllerServiceAccount invokes the alb ingress controller service account creator
func (c *Client) CreateAlbIngressControllerServiceAccount(opts *api.CreateAlbIngressControllerServiceAccountOpts) (*api.ServiceAccount, error) {
	into := &api.ServiceAccount{}
	return into, c.DoPost(targetAlbIngressControllerServiceAccount, opts, into)
}

// CreateAlbIngressControllerPolicy invokes the alb policy creator
func (c *Client) CreateAlbIngressControllerPolicy(opts *api.CreateAlbIngressControllerPolicyOpts) (*api.ManagedPolicy, error) {
	into := &api.ManagedPolicy{}
	return into, c.DoPost(targetAlbIngressControllerPolicy, opts, into)
}

// CreateExternalSecretsHelmChart invokes the external secrets helm chart operation
func (c *Client) CreateExternalSecretsHelmChart(opts *api.CreateExternalSecretsHelmChartOpts) (*api.Helm, error) {
	into := &api.Helm{}
	return into, c.DoPost(targetExternalSecretsHelm, opts, into)
}

// CreateExternalSecretsServiceAccount invokes the external secrets service account operation
func (c *Client) CreateExternalSecretsServiceAccount(opts *api.CreateExternalSecretsServiceAccountOpts) (*api.ServiceAccount, error) {
	into := &api.ServiceAccount{}
	return into, c.DoPost(targetExternalSecretsServiceAccount, opts, into)
}

// CreateExternalSecretsPolicy invokes the external secrets policy create operation
func (c *Client) CreateExternalSecretsPolicy(opts *api.CreateExternalSecretsPolicyOpts) (*api.ManagedPolicy, error) {
	into := &api.ManagedPolicy{}
	return into, c.DoPost(targetExternalSecretsPolicy, opts, into)
}

// CreateVpc invokes the vpc create operation
func (c *Client) CreateVpc(opts *api.CreateVpcOpts) (*api.Vpc, error) {
	into := &api.Vpc{}
	return into, c.DoPost(targetVpcs, opts, into)
}

// DeleteVpc invokes the vpc delete operation
func (c *Client) DeleteVpc(opts *api.DeleteVpcOpts) error {
	return c.DoDelete(targetVpcs, opts)
}

// CreateCluster invokes the cluster create operation
func (c *Client) CreateCluster(opts *api.ClusterCreateOpts) (*api.Cluster, error) {
	into := &api.Cluster{}
	return into, c.DoPost(targetClusters, opts, into)
}

// DeleteCluster invokes the cluster delete operation
func (c *Client) DeleteCluster(opts *api.ClusterDeleteOpts) error {
	return c.DoDelete(targetClusters, opts)
}

// DoPost sends a POST request to the given endpoint
func (c *Client) DoPost(endpoint string, body interface{}, into interface{}) error {
	return c.Do(http.MethodPost, endpoint, body, into)
}

// DoDelete sends a DELETE request to the given endpoint
func (c *Client) DoDelete(endpoint string, body interface{}) error {
	return c.Do(http.MethodDelete, endpoint, body, nil)
}

// Do performs the request
func (c *Client) Do(method, endpoint string, body interface{}, into interface{}) error {
	if c.Debug {
		_, err := fmt.Fprintf(c.Progress, "client (method: %s, endpoint: %s) starting request: %s", method, endpoint, litter.Sdump(body))
		if err != nil {
			return fmt.Errorf("failed to write debug output: %w", err)
		}
	}

	data, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("%s: %w", pretty("failed to marshal data for", method, endpoint), err)
	}

	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", c.BaseURL, endpoint), bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("%s: %w", pretty("failed to create request for", method, endpoint), err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return fmt.Errorf("%s: %w", pretty("request failed for", method, endpoint), err)
	}

	out, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("%s: %w", pretty("failed to read response for", method, endpoint), err)
	}

	if resp.StatusCode >= 400 { // nolint: gomnd
		return fmt.Errorf("request failed with %s, because: %s", http.StatusText(resp.StatusCode), string(out))
	}

	defer func() {
		err = resp.Body.Close()
	}()

	if into != nil {
		if c.Debug {
			_, err = fmt.Fprintf(c.Progress, "client (method: %s, endpoint: %s) received data: %s", method, endpoint, out)
			if err != nil {
				return fmt.Errorf("failed to write debug output: %w", err)
			}
		}

		err = json.Unmarshal(out, into)
		if err != nil {
			return fmt.Errorf("failed to parse response: %w", err)
		}
	}

	_, err = io.Copy(c.Progress, strings.NewReader(string(out)))
	if err != nil {
		return fmt.Errorf("%s: %w", pretty("failed to write progress for", method, endpoint), err)
	}

	return nil
}

func pretty(msg, method, endpoint string) string {
	return fmt.Sprintf("%s: %s, %s", msg, method, endpoint)
}
