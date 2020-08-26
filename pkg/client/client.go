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

	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/api"
)

const (
	targetVpcs                               = "vpcs/"
	targetClusters                           = "clusters/"
	targetClusterConfigs                     = "clusterconfigs/"
	targetExternalSecretsPolicy              = "managedpolicies/externalsecrets/"
	targetExternalSecretsServiceAccount      = "serviceaccounts/externalsecrets/"
	targetExternalSecretsHelm                = "helm/externalsecrets/"
	targetAlbIngressControllerPolicy         = "managedpolicies/albingresscontroller/"
	targetAlbIngressControllerServiceAccount = "serviceaccounts/albingresscontroller/"
	targetAlbIngressControllerHelm           = "helm/albingresscontroller/"
	targetExternalDNSPolicy                  = "managedpolicies/externaldns/"
	targetExternalDNSServiceAccount          = "serviceaccounts/externaldns/"
	targetDomain                             = "domains/"
	targetKubeExternalDNS                    = "/kube/externaldns/"
)

// Cluster client API calls
type Cluster interface {
	CreateCluster(opts *api.ClusterCreateOpts) error
	DeleteCluster(opts *api.ClusterDeleteOpts) error
}

// ClusterConfig client API calls
type ClusterConfig interface {
	CreateClusterConfig(opts *api.CreateClusterConfigOpts) error
}

// Vpc client API calls
type Vpc interface {
	CreateVpc(opts *api.CreateVpcOpts) error
	DeleteVpc(opts *api.DeleteVpcOpts) error
}

// ManagedPolicy API calls
type ManagedPolicy interface {
	CreateExternalSecretsPolicy(opts *api.CreateExternalSecretsPolicyOpts) (*api.ManagedPolicy, error)
	CreateAlbIngressControllerPolicy(opts *api.CreateAlbIngressControllerPolicyOpts) (*api.ManagedPolicy, error)
	CreateExternalDNSPolicy(opts *api.CreateExternalDNSPolicyOpts) (*api.ManagedPolicy, error)
}

// ServiceAccount API calls
type ServiceAccount interface {
	CreateExternalSecretsServiceAccount(opts *api.CreateExternalSecretsServiceAccountOpts) error
	CreateAlbIngressControllerServiceAccount(opts *api.CreateAlbIngressControllerServiceAccountOpts) error
	CreateExternalDNSServiceAccount(opts *api.CreateExternalDNSServiceAccountOpts) error
}

// Helm API calls
type Helm interface {
	CreateExternalSecretsHelmChart(opts *api.CreateExternalSecretsHelmChartOpts) (*api.Helm, error)
	CreateAlbIngressControllerHelmChart(opts *api.CreateAlbIngressControllerHelmChartOpts) (*api.Helm, error)
}

// Kube API calls
type Kube interface {
	CreateExternalDNSKubeDeployment(opts *api.CreateExternalDNSKubeDeploymentOpts) (*api.Kube, error)
}

// Domain API calls
type Domain interface {
	CreateDomain(opts *api.CreateDomainOpts) (*api.Domain, error)
}

// Client stores state for invoking API operations
type Client struct {
	BaseURL  string
	Client   *http.Client
	Progress io.Writer
}

// New returns a client that wraps the common API operations
func New(progress io.Writer, serverURL string) *Client {
	return &Client{
		Progress: progress,
		BaseURL:  serverURL,
		Client:   &http.Client{},
	}
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
func (c *Client) CreateExternalDNSServiceAccount(opts *api.CreateExternalDNSServiceAccountOpts) error {
	return c.DoPost(targetExternalDNSServiceAccount, opts, nil)
}

// CreateAlbIngressControllerHelmChart invokes the alb ingress controller helm chart creator
func (c *Client) CreateAlbIngressControllerHelmChart(opts *api.CreateAlbIngressControllerHelmChartOpts) (*api.Helm, error) {
	into := &api.Helm{}
	return into, c.DoPost(targetAlbIngressControllerHelm, opts, into)
}

// CreateAlbIngressControllerServiceAccount invokes the alb ingress controller service account creator
func (c *Client) CreateAlbIngressControllerServiceAccount(opts *api.CreateAlbIngressControllerServiceAccountOpts) error {
	return c.DoPost(targetAlbIngressControllerServiceAccount, opts, nil)
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
func (c *Client) CreateExternalSecretsServiceAccount(opts *api.CreateExternalSecretsServiceAccountOpts) error {
	return c.DoPost(targetExternalSecretsServiceAccount, opts, nil)
}

// CreateExternalSecretsPolicy invokes the external secrets policy create operation
func (c *Client) CreateExternalSecretsPolicy(opts *api.CreateExternalSecretsPolicyOpts) (*api.ManagedPolicy, error) {
	into := &api.ManagedPolicy{}
	return into, c.DoPost(targetExternalSecretsPolicy, opts, into)
}

// CreateClusterConfig invokes the cluster config create operation
func (c *Client) CreateClusterConfig(opts *api.CreateClusterConfigOpts) error {
	return c.DoPost(targetClusterConfigs, opts, nil)
}

// CreateVpc invokes the vpc create operation
func (c *Client) CreateVpc(opts *api.CreateVpcOpts) error {
	return c.DoPost(targetVpcs, opts, nil)
}

// DeleteVpc invokes the vpc delete operation
func (c *Client) DeleteVpc(opts *api.DeleteVpcOpts) error {
	return c.DoDelete(targetVpcs, opts)
}

// CreateCluster invokes the cluster create operation
func (c *Client) CreateCluster(opts *api.ClusterCreateOpts) error {
	return c.DoPost(targetClusters, opts, nil)
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
	data, err := json.Marshal(body)
	if err != nil {
		return errors.E(err, pretty("failed to marshal data for", method, endpoint))
	}

	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", c.BaseURL, endpoint), bytes.NewReader(data))
	if err != nil {
		return errors.E(err, pretty("failed to create request for", method, endpoint))
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return errors.E(err, pretty("request failed for", method, endpoint))
	}

	out, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.E(err, pretty("failed to read response for", method, endpoint))
	}

	defer func() {
		err = resp.Body.Close()
	}()

	if into != nil {
		err = json.Unmarshal(out, into)
		if err != nil {
			return fmt.Errorf("failed to parse response: %w", err)
		}
	}

	_, err = io.Copy(c.Progress, strings.NewReader(string(out)))
	if err != nil {
		return errors.E(err, pretty("failed to write progress for", method, endpoint))
	}

	return nil
}

func pretty(msg, method, endpoint string) string {
	return fmt.Sprintf("%s: %s, %s", msg, method, endpoint)
}
