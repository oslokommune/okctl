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
	DefaultStackNameContainerRepository               = "containerrepository"
	DefaultStackNameSecurityGroup                     = "sg"
)

// StackNamer knows how to name cloud formation stacks
type StackNamer struct{}

// NewStackNamer returns an initialised stack namer
func NewStackNamer() *StackNamer {
	return &StackNamer{}
}

// S3Bucket returns the stack name of an AWS S3 Bucket
func (n *StackNamer) S3Bucket(app, clusterName string) string {
	return fmt.Sprintf("%s-%s-%s-%s",
		DefaultStackNamePrefix,
		DefaultStackNameS3Bucket,
		clusterName,
		app,
	)
}

// RDSPostgres returns the stack name of a RDS Postgres instance
func (n *StackNamer) RDSPostgres(app, clusterName string) string {
	return fmt.Sprintf("%s-%s-%s-%s",
		DefaultStackNamePrefix,
		DefaultStackNameRDSPostgres,
		clusterName,
		app,
	)
}

// Vpc returns the stack name of a vpc
func (n *StackNamer) Vpc(clusterName string) string {
	return fmt.Sprintf("%s-%s-%s",
		DefaultStackNamePrefix,
		DefaultStackNameVpcID,
		clusterName,
	)
}

// ExternalSecretsPolicy returns the stack name of an external secrets policy
func (n *StackNamer) ExternalSecretsPolicy(clusterName string) string {
	return fmt.Sprintf("%s-%s-%s",
		DefaultStackNamePrefix,
		DefaultStackNameExternalSecretsPolicyID,
		clusterName,
	)
}

// AutoscalerPolicy returns the stack name of an autoscaler policy
func (n *StackNamer) AutoscalerPolicy(clusterName string) string {
	return fmt.Sprintf("%s-%s-%s",
		DefaultStackNamePrefix,
		DefaultStackNameAutoscalerPolicyID,
		clusterName,
	)
}

// BlockstoragePolicy returns the stack name of an Blockstorage policy
func (n *StackNamer) BlockstoragePolicy(clusterName string) string {
	return fmt.Sprintf("%s-%s-%s",
		DefaultStackNamePrefix,
		DefaultStackNameBlockstoragePolicyID,
		clusterName,
	)
}

// CloudwatchDatasource returns the stack name of an Blockstorage policy
func (n *StackNamer) CloudwatchDatasource(clusterName string) string {
	return fmt.Sprintf("%s-%s-%s",
		DefaultStackNamePrefix,
		DefaultStackNameCloudwatchDatasourceID,
		clusterName,
	)
}

// FargateCloudwatch returns the stack name of an Blockstorage policy
func (n *StackNamer) FargateCloudwatch(clusterName string) string {
	return fmt.Sprintf("%s-%s-%s",
		DefaultStackNamePrefix,
		DefaultStackNameFargateCloudwatchID,
		clusterName,
	)
}

// AlbIngressControllerPolicy returns the stack name of the policy
func (n *StackNamer) AlbIngressControllerPolicy(clusterName string) string {
	return fmt.Sprintf("%s-%s-%s",
		DefaultStackNamePrefix,
		DefaultStackNameAlbIngressControllerPolicyID,
		clusterName,
	)
}

// AWSLoadBalancerControllerPolicy returns the stack name of the policy
func (n *StackNamer) AWSLoadBalancerControllerPolicy(clusterName string) string {
	return fmt.Sprintf("%s-%s-%s",
		DefaultStackNamePrefix,
		DefaultStackNameAWSLoadBalancerControllerPolicyID,
		clusterName,
	)
}

// ExternalDNSPolicy returns the stack name of the policy
func (n *StackNamer) ExternalDNSPolicy(clusterName string) string {
	return fmt.Sprintf("%s-%s-%s",
		DefaultStackNamePrefix,
		DefaultStackNameExternalDNSPolicyID,
		clusterName,
	)
}

// Domain returns the stack name of the domain
func (n *StackNamer) Domain(clusterName, subdomain string) string {
	return fmt.Sprintf("%s-%s-%s-%s",
		DefaultStackNamePrefix,
		DefaultStackNameDomainID,
		clusterName,
		subdomain,
	)
}

// Certificate returns the stack name of the certificate
func (n *StackNamer) Certificate(clusterName, subdomain string) string {
	return fmt.Sprintf("%s-%s-%s-%s",
		DefaultStackNamePrefix,
		DefaultStackNameCertificateID,
		clusterName,
		subdomain,
	)
}

// IdentityPool returns the stack name of the identity pool
func (n *StackNamer) IdentityPool(clusterName string) string {
	return fmt.Sprintf("%s-%s-%s",
		DefaultStackNamePrefix,
		DefaultStackNameIdentityPool,
		clusterName,
	)
}

// IdentityPoolClient returns the stack name of the identity pool client
func (n *StackNamer) IdentityPoolClient(clusterName, purpose string) string {
	return fmt.Sprintf("%s-%s-%s-%s",
		DefaultStackNamePrefix,
		DefaultStackNameIdentityPoolClient,
		clusterName,
		purpose,
	)
}

// IdentityPoolUser returns the stack name of the identity pool user
func (n *StackNamer) IdentityPoolUser(clusterName, email string) string {
	return fmt.Sprintf("%s-%s-%s-%s",
		DefaultStackNamePrefix,
		DefaultStackNameIdentityPoolUser,
		clusterName,
		email,
	)
}

// AliasRecordSet returns the stack name of the alias record set
func (n *StackNamer) AliasRecordSet(clusterName, subdomain string) string {
	return fmt.Sprintf("%s-%s-%s-%s",
		DefaultStackNamePrefix,
		DefaultStackNameAliasRecordSet,
		clusterName,
		subdomain,
	)
}

// ContainerRepository returns the stack name of the container repository
func (n *StackNamer) ContainerRepository(clusterName, imageName string) string {
	return fmt.Sprintf("%s-%s-%s-%s",
		DefaultStackNamePrefix,
		DefaultStackNameContainerRepository,
		clusterName,
		imageName,
	)
}

// SecurityGroup returns the stack name of the security group cfn stack
func (n *StackNamer) SecurityGroup(clusterName, securityGroupName string) string {
	return fmt.Sprintf("%s-%s-%s-%s",
		DefaultStackNamePrefix,
		DefaultStackNameSecurityGroup,
		clusterName,
		securityGroupName,
	)
}
