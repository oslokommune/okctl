package cfn

import (
	"fmt"
)

// Not sure about this approach, will need to read up.

// nolint: golint
const (
	DefaultStackNamePrefix                            = "okctl"
	DefaultStackNameVpcID                             = "vpc"
	DefaultStackNameExternalSecretsPolicyID           = "externalsecretspolicy"
	DefaultStackNameAutoscalerPolicyID                = "autoscalerpolicy"
	DefaultStackNameBlockstoragePolicyID              = "blockstoragepolicy"
	DefaultStackNameAlbIngressControllerPolicyID      = "albingresscontrollerpolicy"
	DefaultStackNameAWSLoadBalancerControllerPolicyID = "awsloadbalancercontrollerpolicy"
	DefaultStackNameCloudwatchDatasourceID            = "cloudwatchdatasource"
	DefaultStackNameFargateCloudwatchID               = "fargatecloudwatch"
	DefaultStackNameExternalDNSPolicyID               = "externaldns"
	DefaultStackNameDomainID                          = "domain"
	DefaultStackNameCertificateID                     = "certificate"
	DefaultStackNameIdentityPool                      = "identitypool"
	DefaultStackNameIdentityPoolClient                = "identitypoolclient"
	DefaultStackNameAliasRecordSet                    = "aliasrecordset"
	DefaultStackNameIdentityPoolUser                  = "identitypooluser"
	DefaultStackNameRDSPostgres                       = "rdspostgres"
	DefaultStackNameS3Bucket                          = "s3bucket"
)

// StackNamer knows how to name cloud formation stacks
type StackNamer struct{}

// NewStackNamer returns an initialised stack namer
func NewStackNamer() *StackNamer {
	return &StackNamer{}
}

// S3Bucket returns the stack name of an AWS S3 Bucket
func (n *StackNamer) S3Bucket(app, repository, env string) string {
	return fmt.Sprintf("%s-%s-%s-%s-%s",
		DefaultStackNamePrefix,
		DefaultStackNameS3Bucket,
		repository,
		env,
		app,
	)
}

// RDSPostgres returns the stack name of a RDS Postgres instance
func (n *StackNamer) RDSPostgres(app, repository, env string) string {
	return fmt.Sprintf("%s-%s-%s-%s-%s",
		DefaultStackNamePrefix,
		DefaultStackNameRDSPostgres,
		repository,
		env,
		app,
	)
}

// Vpc returns the stack name of a vpc
func (n *StackNamer) Vpc(repository, env string) string {
	return fmt.Sprintf("%s-%s-%s-%s",
		DefaultStackNamePrefix,
		DefaultStackNameVpcID,
		repository,
		env,
	)
}

// ExternalSecretsPolicy returns the stack name of an external secrets policy
func (n *StackNamer) ExternalSecretsPolicy(repository, env string) string {
	return fmt.Sprintf("%s-%s-%s-%s",
		DefaultStackNamePrefix,
		DefaultStackNameExternalSecretsPolicyID,
		repository,
		env,
	)
}

// AutoscalerPolicy returns the stack name of an autoscaler policy
func (n *StackNamer) AutoscalerPolicy(repository, env string) string {
	return fmt.Sprintf("%s-%s-%s-%s",
		DefaultStackNamePrefix,
		DefaultStackNameAutoscalerPolicyID,
		repository,
		env,
	)
}

// BlockstoragePolicy returns the stack name of an Blockstorage policy
func (n *StackNamer) BlockstoragePolicy(repository, env string) string {
	return fmt.Sprintf("%s-%s-%s-%s",
		DefaultStackNamePrefix,
		DefaultStackNameBlockstoragePolicyID,
		repository,
		env,
	)
}

// CloudwatchDatasource returns the stack name of an Blockstorage policy
func (n *StackNamer) CloudwatchDatasource(repository, env string) string {
	return fmt.Sprintf("%s-%s-%s-%s",
		DefaultStackNamePrefix,
		DefaultStackNameCloudwatchDatasourceID,
		repository,
		env,
	)
}

// FargateCloudwatch returns the stack name of an Blockstorage policy
func (n *StackNamer) FargateCloudwatch(repository, env string) string {
	return fmt.Sprintf("%s-%s-%s-%s",
		DefaultStackNamePrefix,
		DefaultStackNameFargateCloudwatchID,
		repository,
		env,
	)
}

// AlbIngressControllerPolicy returns the stack name of the policy
func (n *StackNamer) AlbIngressControllerPolicy(repository, env string) string {
	return fmt.Sprintf("%s-%s-%s-%s",
		DefaultStackNamePrefix,
		DefaultStackNameAlbIngressControllerPolicyID,
		repository,
		env,
	)
}

// AWSLoadBalancerControllerPolicy returns the stack name of the policy
func (n *StackNamer) AWSLoadBalancerControllerPolicy(repository, env string) string {
	return fmt.Sprintf("%s-%s-%s-%s",
		DefaultStackNamePrefix,
		DefaultStackNameAWSLoadBalancerControllerPolicyID,
		repository,
		env,
	)
}

// ExternalDNSPolicy returns the stack name of the policy
func (n *StackNamer) ExternalDNSPolicy(repository, env string) string {
	return fmt.Sprintf("%s-%s-%s-%s",
		DefaultStackNamePrefix,
		DefaultStackNameExternalDNSPolicyID,
		repository,
		env,
	)
}

// Domain returns the stack name of the domain
func (n *StackNamer) Domain(repository, env, subdomain string) string {
	return fmt.Sprintf("%s-%s-%s-%s-%s",
		DefaultStackNamePrefix,
		DefaultStackNameDomainID,
		repository,
		env,
		subdomain,
	)
}

// Certificate returns the stack name of the certificate
func (n *StackNamer) Certificate(repository, env, subdomain string) string {
	return fmt.Sprintf("%s-%s-%s-%s-%s",
		DefaultStackNamePrefix,
		DefaultStackNameCertificateID,
		repository,
		env,
		subdomain,
	)
}

// IdentityPool returns the stack name of the identity pool
func (n *StackNamer) IdentityPool(repository, env string) string {
	return fmt.Sprintf("%s-%s-%s-%s",
		DefaultStackNamePrefix,
		DefaultStackNameIdentityPool,
		repository,
		env,
	)
}

// IdentityPoolClient returns the stack name of the identity pool client
func (n *StackNamer) IdentityPoolClient(repository, env, purpose string) string {
	return fmt.Sprintf("%s-%s-%s-%s-%s",
		DefaultStackNamePrefix,
		DefaultStackNameIdentityPoolClient,
		repository,
		env,
		purpose,
	)
}

// IdentityPoolUser returns the stack name of the identity pool user
func (n *StackNamer) IdentityPoolUser(repository, env, email string) string {
	return fmt.Sprintf("%s-%s-%s-%s-%s",
		DefaultStackNamePrefix,
		DefaultStackNameIdentityPoolUser,
		repository,
		env,
		email,
	)
}

// AliasRecordSet returns the stack name of the alias record set
func (n *StackNamer) AliasRecordSet(repository, env, subdomain string) string {
	return fmt.Sprintf("%s-%s-%s-%s-%s",
		DefaultStackNamePrefix,
		DefaultStackNameAliasRecordSet,
		repository,
		env,
		subdomain,
	)
}
