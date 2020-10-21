package state

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/client/store"
	"github.com/spf13/afero"
)

// HostedZoner defines the allowed actions on a hosted zone
type HostedZoner interface {
	SaveHostedZone(domain string, zone HostedZone) (*store.Report, error)
	GetHostedZone(domain string) HostedZone
	DeleteHostedZone(domain string) (*store.Report, error)
	GetHostedZones() map[string]HostedZone
	GetPrimaryHostedZone() *HostedZone
}

// Clusterer defines the allowed actions on a cluster
type Clusterer interface {
	SaveCluster(cluster Cluster) (*store.Report, error)
	DeleteCluster() (*store.Report, error)
	GetCluster() Cluster
}

// Vpcer defines the allowed actions on a vpc
type Vpcer interface {
	SaveVPC(vpc VPC) (*store.Report, error)
	DeleteVPC() (*store.Report, error)
	GetVPC() VPC
}

// Githuber defines the allowed actions on the github state
type Githuber interface {
	SaveGithub(github Github) (*store.Report, error)
	DeleteGithub() (*store.Report, error)
	GetGithub() Github
}

// Argocder defines the allowed actions on the argocd state
type Argocder interface {
	SaveArgoCD(cd ArgoCD) (*store.Report, error)
	DeleteArgoCD() (*store.Report, error)
	GetArgoCD() ArgoCD
}

// Certificater defines the allowed actions on the certificate state
type Certificater interface {
	SaveCertificate(certificate Certificate) (*store.Report, error)
	GetCertificate(domain string) Certificate
	GetCertificates() map[string]Certificate
}

// Metadataer defines the allowed action on the metadata state
type Metadataer interface {
	SaveMetadata(metadata Metadata) (*store.Report, error)
	GetMetadata() Metadata
}

// IdentityPooler defines the allowed actions on the identitypool state
type IdentityPooler interface {
	SaveIdentityPool(pool IdentityPool) (*store.Report, error)
	GetIdentityPool() IdentityPool
	SaveIdentityPoolClient(client IdentityPoolClient) (*store.Report, error)
	GetIdentityPoolClient(purpose string) IdentityPoolClient
	GetIdentityPoolUser(email string) IdentityPoolUser
	SaveIdentityPoolUser(client IdentityPoolUser) (*store.Report, error)
}

// RepositoryStateWithEnv provides actions for interacting with
// the state of a repository
type RepositoryStateWithEnv interface {
	HostedZoner
	Clusterer
	Vpcer
	Githuber
	Argocder
	Certificater
	Metadataer
	IdentityPooler
	GetClusterName() string
	Save() (*store.Report, error)
}

// SaverFn implements the storage operation of the state
type SaverFn func(r *Repository) (*store.Report, error)

// DefaultFileSystemSaver returns an initialised saver
func DefaultFileSystemSaver(repoFile, outputDir string, fs *afero.Afero) SaverFn {
	return func(r *Repository) (*store.Report, error) {
		return store.NewFileSystem(outputDir, fs).
			StoreStruct(repoFile, r, store.ToYAML()).
			Do()
	}
}

type repository struct {
	state   *Repository
	env     string
	saverFn SaverFn
}

func (r *repository) GetIdentityPoolUser(email string) IdentityPoolUser {
	pool := r.GetIdentityPool()

	if pool.Users == nil {
		return IdentityPoolUser{}
	}

	return pool.Users[email]
}

func (r *repository) SaveIdentityPoolUser(user IdentityPoolUser) (*store.Report, error) {
	cluster := r.GetCluster()

	if cluster.IdentityPool.Users == nil {
		cluster.IdentityPool.Users = map[string]IdentityPoolUser{}
	}

	cluster.IdentityPool.Users[user.Email] = user

	r.state.Clusters[r.env] = cluster

	return r.save()
}

// NewRepositoryStateWithEnv returns an initialised setter for a given env
func NewRepositoryStateWithEnv(env string, r *Repository, fn SaverFn) RepositoryStateWithEnv {
	return &repository{
		state:   r,
		env:     env,
		saverFn: fn,
	}
}

// SaveIdentityPoolClient saves the identity pool client
func (r *repository) SaveIdentityPoolClient(client IdentityPoolClient) (*store.Report, error) {
	cluster := r.GetCluster()

	if cluster.IdentityPool.Clients == nil {
		cluster.IdentityPool.Clients = map[string]IdentityPoolClient{}
	}

	cluster.IdentityPool.Clients[client.Purpose] = client

	r.state.Clusters[r.env] = cluster

	return r.save()
}

// GetIdentityPoolClient returns the identity pool client
func (r *repository) GetIdentityPoolClient(purpose string) IdentityPoolClient {
	pool := r.GetIdentityPool()

	if pool.Clients == nil {
		return IdentityPoolClient{}
	}

	return pool.Clients[purpose]
}

// SaveIdentityPool saves the identity pool and stores the state
func (r *repository) SaveIdentityPool(pool IdentityPool) (*store.Report, error) {
	cluster := r.GetCluster()
	cluster.IdentityPool = pool

	r.state.Clusters[r.env] = cluster

	return r.save()
}

func (r *repository) GetIdentityPool() IdentityPool {
	_ = r.GetCluster()
	return r.state.Clusters[r.env].IdentityPool
}

// SaveCertificate updates the state with the provided certificate
func (r *repository) SaveCertificate(certificate Certificate) (*store.Report, error) {
	_ = r.GetCertificates()
	r.state.Clusters[r.env].Certificates[certificate.Domain] = certificate

	return r.save()
}

// GetCertificate returns the certificate for the given domain
func (r *repository) GetCertificate(domain string) Certificate {
	certs := r.GetCertificates()
	return certs[domain]
}

// GetCertificates returns all certificates
func (r *repository) GetCertificates() map[string]Certificate {
	cluster := r.GetCluster()

	if cluster.Certificates == nil {
		cluster.Certificates = map[string]Certificate{}
	}

	r.state.Clusters[r.env] = cluster

	return cluster.Certificates
}

// SaveMetadata updates the metadata
func (r *repository) SaveMetadata(metadata Metadata) (*store.Report, error) {
	r.state.Metadata = metadata

	return r.save()
}

// GetMetadata returns the stored metadata
func (r *repository) GetMetadata() Metadata {
	return r.state.Metadata
}

// SaveArgoCD updates the argocd state
func (r *repository) SaveArgoCD(cd ArgoCD) (*store.Report, error) {
	cluster := r.GetCluster()
	cluster.ArgoCD = cd

	r.state.Clusters[r.env] = cluster

	return r.save()
}

// DeleteArgoCD removes the argocd state
func (r *repository) DeleteArgoCD() (*store.Report, error) {
	cluster := r.GetCluster()
	cluster.ArgoCD = ArgoCD{}
	r.state.Clusters[r.env] = cluster

	return r.save()
}

// GetArgoCD retrieves the argocd state
func (r *repository) GetArgoCD() ArgoCD {
	cluster := r.GetCluster()
	return cluster.ArgoCD
}

// SaveGithub stores the github state
func (r *repository) SaveGithub(github Github) (*store.Report, error) {
	cluster := r.GetCluster()
	cluster.Github = github
	r.state.Clusters[r.env] = cluster

	return r.save()
}

// DeleteGithub removes the github state
func (r *repository) DeleteGithub() (*store.Report, error) {
	cluster := r.GetCluster()
	cluster.Github = Github{}
	r.state.Clusters[r.env] = cluster

	return r.save()
}

// GetGithub retrieves the github state
func (r *repository) GetGithub() Github {
	cluster := r.GetCluster()
	return cluster.Github
}

// SaveVPC updates the state
func (r *repository) SaveVPC(vpc VPC) (*store.Report, error) {
	cluster := r.GetCluster()
	cluster.VPC = vpc
	r.state.Clusters[r.env] = cluster

	return r.save()
}

// DeleteVPC removes the vpc
func (r *repository) DeleteVPC() (*store.Report, error) {
	cluster := r.GetCluster()
	cluster.VPC = VPC{}
	r.state.Clusters[r.env] = cluster

	return r.save()
}

// GetVPC retrieves the VPC
func (r *repository) GetVPC() VPC {
	cluster := r.GetCluster()
	return cluster.VPC
}

// SaveCluster stores the cluster
func (r *repository) SaveCluster(cluster Cluster) (*store.Report, error) {
	_ = r.GetCluster()
	r.state.Clusters[r.env] = cluster

	return r.save()
}

// DeleteCluster removes the cluster state
func (r *repository) DeleteCluster() (*store.Report, error) {
	delete(r.state.Clusters, r.env)
	return r.save()
}

// SaveHostedZone stores the hosted zone
func (r *repository) SaveHostedZone(domain string, zone HostedZone) (*store.Report, error) {
	cluster := r.GetCluster()

	if cluster.HostedZone == nil {
		cluster.HostedZone = map[string]HostedZone{}
	}

	cluster.HostedZone[domain] = zone
	r.state.Clusters[r.env] = cluster

	return r.save()
}

// DeleteHostedZone removes the hosted zone from the state
func (r *repository) DeleteHostedZone(domain string) (*store.Report, error) {
	cluster := r.GetCluster()

	if cluster.HostedZone == nil {
		return &store.Report{}, nil
	}

	delete(cluster.HostedZone, domain)

	return r.save()
}

// GetHostedZone retrieves the hosted zone state
func (r *repository) GetHostedZone(domain string) HostedZone {
	cluster := r.GetCluster()

	if cluster.HostedZone == nil {
		cluster.HostedZone = map[string]HostedZone{}
	}

	r.state.Clusters[r.env] = cluster

	return r.state.Clusters[r.env].HostedZone[domain]
}

// GetHostedZones retrieves the hosted zones
func (r *repository) GetHostedZones() map[string]HostedZone {
	c := r.GetCluster()

	if c.HostedZone == nil {
		c.HostedZone = map[string]HostedZone{}
	}

	r.state.Clusters[r.env] = c

	return c.HostedZone
}

// GetPrimaryHostedZone returns the primary hosted zone if one exists
func (r *repository) GetPrimaryHostedZone() *HostedZone {
	for _, z := range r.GetHostedZones() {
		z := z

		if z.Primary && len(z.Domain) > 0 {
			return &z
		}
	}

	return nil
}

// GetCluster retrieves the cluster
func (r *repository) GetCluster() Cluster {
	if r.state.Clusters == nil {
		r.state.Clusters = map[string]Cluster{}
	}

	if _, ok := r.state.Clusters[r.env]; !ok {
		r.state.Clusters[r.env] = Cluster{
			Environment: r.env,
		}
	}

	return r.state.Clusters[r.env]
}

// GetClusterName returns the cluster name
func (r *repository) GetClusterName() string {
	cluster := r.GetCluster()
	return fmt.Sprintf("%s-%s", r.state.Metadata.Name, cluster.Environment)
}

func (r *repository) Save() (*store.Report, error) {
	return r.save()
}

func (r *repository) save() (*store.Report, error) {
	return r.saverFn(r.state)
}
