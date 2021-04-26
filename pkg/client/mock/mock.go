// Package mock contains some convenience functions for getting data
package mock

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/helm"

	"helm.sh/helm/v3/pkg/release"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/clusterconfig"
)

// nolint: golint gosec
const (
	DefaultRegion                     = "eu-west-1"
	DefaultAWSAccountID               = "123456789012"
	DefaultClusterName                = "okctl-staging"
	DefaultDomain                     = "okctl-staging.oslo.systems"
	DefaultAuthDomain                 = "auth.okctl-staging.oslo.systems"
	DefaultDomainFilter               = DefaultDomain
	DefaultFQDN                       = "okctl-staging.oslo.systems."
	DefaultCallbackURL                = "https://argocd.okctl-staging.oslo.systems/callback"
	DefaultDatabaseEndpoint           = "https://db.okctl-staging.oslo.systems"
	DefaultHostedZoneID               = "Z0FAKE41FAKE6I841FAKE"
	DefaultAliasHostedZoneID          = "Z0FAKE67FAKE6I231FAKE"
	DefaultCertificateARN             = "arn:aws:acm:eu-west-1:123456789012:certificate/123456789012-1234-1234-1234-12345678"
	DefaultServiceAccountName         = "important-sa"
	DefaultPolicyARN                  = "arn:aws:iam::123456789012:policy/policy-name-with-path"
	DefaultNamespace                  = "kube-system"
	DefaultManifestName               = "okctl-cm"
	DefaultManifestType               = client.ManifestTypeConfigMap
	DefaultHelmReleaseName            = "okctl-helm-release"
	DefaultSecretParameterName        = "release-secret"
	DefaultSecretParameterVersion     = 1
	DefaultSecretParameterPath        = "/okctl/staging/release-secret"
	DefaultEmail                      = "bob@thebuilder.com"
	DefaultUserPoolID                 = "TYUJBFW3893FAKE"
	DefaultPurpose                    = "argocd"
	DefaultClientID                   = "gehu-fgerg432-ewge"
	DefaultClientSecret               = "0ef90weug09jfqh3rf"
	DefaultVpcID                      = "vpc-0e9801d129EXAMPLE"
	DefaultCidr                       = "192.168.0.0/20"
	DefaultPublicSubnetCidr           = "192.168.1.0/24"
	DefaultPublicSubnetID             = "gguhef789FAKE"
	DefaultPrivateSubnetCidr          = "192.168.2.0/24"
	DefaultPrivateSubnetID            = "e9e093ufFAKE"
	DefaultDatabaseSubnetCidr         = "192.168.3.0/24"
	DefaultDatabaseSubnetID           = "djfh093FAKE"
	DefaultDatabaseSubnetGroupName    = "okctl-staging-DatabaseGroup"
	DefaultAvailabilityZone           = "eu-west-1a"
	DefaultVersion                    = "1.18"
	DefaultPostgresName               = "backend"
	DefaultPostgresUserName           = "administrator"
	DefaultSecurityGroupID            = "sg-2979fue9FAKE"
	DefaultPostgresSecretARN          = "arn:aws:secretsmanager:eu-west-1:123456789012:secret:secret/administrator"
	DefaultPostgresSecretFriendlyName = "secrets/administrator"
	DefaultPostgresAdminName          = "okctl-staging-backend-admin"
	DefaultPostgresConfigMapName      = "okctl-staging-backend-cm"
	DefaultRoleARN                    = "arn:aws:iam::123456789012:role/okctl-staging-Role"
	DefaultPostgresLambdaFunctionARN  = "arn:aws:lambda:eu-west-1:123456789012:function:rotater"
	DefaultS3BucketName               = "rotater"
	DefaultPostgresDatabasePort       = 5432
	DefaultArgoDomain                 = "argocd.okctl-staging.oslo.systems"
	DefaultArgoURL                    = "https://argocd.okctl-staging.oslo.systems"
	DefaultMonitoringURL              = "https://grafana.okctl-staging.oslo.systems"
	DefaultGithubName                 = "okctl-iac"
	DefaultGithubOrg                  = "oslokommune"
	DefaultGithubFullName             = "oslokommune/okctl-iac"
	DefaultGithubURL                  = "git@github.com:oslokommune/okctl-iac"
	DefaultImage                      = "my-image"

	StackNameHostedZone         = "okctl-staging-oslo-systems-HostedZone"
	StackNameCertificate        = "okctl-staging-oslo-systems-Certificate"
	StackNameManagedPolicy      = "okctl-staging-ManagedPolicy"
	StackNameIdentityPool       = "okctl-staging-IdentityPool"
	StackNameIdentityPoolClient = "okctl-staging-IdentityPoolClient"
	StackNameIdentityPoolUser   = "okctl-staging-bobthebuilder-IdentityPoolUser"
	StackNameRecordSetAlias     = "okctl-staging-RecordSetAlias"
	StackNameVpc                = "okctl-staging-Vpc"
	StackNamePostgresDatabase   = "okctl-staging-backend-PostgresDatabase"
	StackNameRotaterBucket      = "okctl-staging-Rotater"
)

// nolint: golint gochecknoglobals
var (
	DefaultServiceAccountLabels = map[string]string{
		"aws-usage": "cluster-ops",
	}
)

// ID returns a fake id
func ID() api.ID {
	return api.ID{
		Region:       DefaultRegion,
		AWSAccountID: DefaultAWSAccountID,
		ClusterName:  DefaultClusterName,
	}
}

// CloudFormationTemplate just returns some bytes
func CloudFormationTemplate() []byte {
	return []byte(`AWSTemplateFormatVersion: 2010-09-09
`)
}

// NameServers returns a set of fake name servers
func NameServers() []string {
	return []string{
		"ns-158-fake.awsdns-19.com.",
		"ns-1389-fake.awsdns-45.org.",
		"ns-622-fake.awsdns-13.net.",
		"ns-1614-fake.awsdns-09.co.uk.",
	}
}

// HostedZone returns a fake hosted zone
func HostedZone() *client.HostedZone {
	return &client.HostedZone{
		ID:                     ID(),
		IsDelegated:            true,
		Primary:                true,
		Managed:                true,
		FQDN:                   DefaultFQDN,
		Domain:                 DefaultDomain,
		HostedZoneID:           DefaultHostedZoneID,
		NameServers:            NameServers(),
		StackName:              StackNameHostedZone,
		CloudFormationTemplate: CloudFormationTemplate(),
	}
}

// Certificate returns a fake certificate
func Certificate() *client.Certificate {
	return &client.Certificate{
		ID:                     ID(),
		FQDN:                   DefaultFQDN,
		Domain:                 DefaultDomain,
		HostedZoneID:           DefaultHostedZoneID,
		ARN:                    DefaultCertificateARN,
		StackName:              StackNameCertificate,
		CloudFormationTemplate: CloudFormationTemplate(),
	}
}

// ServiceAccountClusterConfig returns a fake cluster config
// for a service account
func ServiceAccountClusterConfig() *v1alpha5.ClusterConfig {
	c, _ := clusterconfig.NewServiceAccount(&clusterconfig.ServiceAccountArgs{
		ClusterName:            DefaultClusterName,
		Labels:                 DefaultServiceAccountLabels,
		Name:                   DefaultServiceAccountName,
		Namespace:              DefaultNamespace,
		PermissionsBoundaryArn: v1alpha1.PermissionsBoundaryARN(DefaultAWSAccountID),
		PolicyArn:              DefaultPolicyARN,
		Region:                 DefaultRegion,
	})

	return c
}

// ServiceAccount returns a fake service account
func ServiceAccount() *client.ServiceAccount {
	return &client.ServiceAccount{
		ID:        ID(),
		Name:      DefaultServiceAccountName,
		PolicyArn: DefaultPolicyARN,
		Config:    ServiceAccountClusterConfig(),
	}
}

// ManifestContent returns a fake ConfigMap
func ManifestContent() []byte {
	return []byte(fmt.Sprintf(`apiVersion: v1
kind: ConfigMap
metadata:
  name: %s
  namespace: %s
data:
  app.properties: |
    env=dev
`, DefaultManifestName, DefaultNamespace))
}

// KubernetesManifest returns a fake kubernetes manifest
func KubernetesManifest() *client.KubernetesManifest {
	return &client.KubernetesManifest{
		ID:        ID(),
		Name:      DefaultManifestName,
		Namespace: DefaultNamespace,
		Type:      DefaultManifestType,
		Content:   ManifestContent(),
	}
}

// ManagedPolicy returns a fake managed policy
func ManagedPolicy() *client.ManagedPolicy {
	return &client.ManagedPolicy{
		ID:                     ID(),
		StackName:              StackNameManagedPolicy,
		PolicyARN:              DefaultPolicyARN,
		CloudFormationTemplate: CloudFormationTemplate(),
	}
}

// HelmChart returns a fake helm chart
func HelmChart() *helm.Chart {
	return &helm.Chart{
		RepositoryName: "my-repo",
		RepositoryURL:  "https://something/repo",
		ReleaseName:    DefaultHelmReleaseName,
		Version:        "v1.0.0",
		Chart:          "my-chart",
		Namespace:      DefaultNamespace,
	}
}

// Helm returns a fake helm
func Helm() *client.Helm {
	return &client.Helm{
		ID: ID(),
		Release: &release.Release{
			Name: DefaultHelmReleaseName,
		},
		Chart: HelmChart(),
	}
}

// ExternalDNSKube returns a fake external dns kube
func ExternalDNSKube() *client.ExternalDNSKube {
	return &client.ExternalDNSKube{
		ID:           ID(),
		HostedZoneID: DefaultHostedZoneID,
		DomainFilter: DefaultDomainFilter,
		Manifests: map[string][]byte{
			"config-map.yaml": ManifestContent(),
		},
	}
}

// ExternalDNS returns a fake external dns
func ExternalDNS() *client.ExternalDNS {
	return &client.ExternalDNS{
		Kube: ExternalDNSKube(),
	}
}

// SecretParameter returns a fake secret parameter
func SecretParameter(content string) *client.SecretParameter {
	return &client.SecretParameter{
		ID:      ID(),
		Name:    DefaultSecretParameterName,
		Path:    DefaultSecretParameterPath,
		Version: DefaultSecretParameterVersion,
		Content: content,
	}
}

// RecordSetAlias returns a fake record set alias
func RecordSetAlias() *client.RecordSetAlias {
	return &client.RecordSetAlias{
		AliasDomain:            DefaultAuthDomain,
		AliasHostedZones:       DefaultAliasHostedZoneID,
		StackName:              StackNameRecordSetAlias,
		CloudFormationTemplate: CloudFormationTemplate(),
	}
}

// IdentityPool returns a fake identity pool
func IdentityPool() *client.IdentityPool {
	return &client.IdentityPool{
		ID:                      ID(),
		UserPoolID:              DefaultUserPoolID,
		AuthDomain:              DefaultAuthDomain,
		HostedZoneID:            DefaultHostedZoneID,
		StackName:               StackNameIdentityPool,
		CloudFormationTemplates: CloudFormationTemplate(),
		Certificate:             Certificate(),
		RecordSetAlias:          RecordSetAlias(),
	}
}

// IdentityPoolClient returns a fake identity pool client
func IdentityPoolClient() *client.IdentityPoolClient {
	return &client.IdentityPoolClient{
		ID:                      ID(),
		UserPoolID:              DefaultUserPoolID,
		Purpose:                 DefaultPurpose,
		CallbackURL:             DefaultCallbackURL,
		ClientID:                DefaultClientID,
		ClientSecret:            DefaultClientSecret,
		StackName:               StackNameIdentityPoolClient,
		CloudFormationTemplates: CloudFormationTemplate(),
	}
}

// IdentityPoolUser returns a fake identity pool user
func IdentityPoolUser() *client.IdentityPoolUser {
	return &client.IdentityPoolUser{
		ID:                     ID(),
		Email:                  DefaultEmail,
		UserPoolID:             DefaultUserPoolID,
		StackName:              StackNameIdentityPoolUser,
		CloudFormationTemplate: CloudFormationTemplate(),
	}
}

// PublicSubnets returns fake public subnets
func PublicSubnets() []client.VpcSubnet {
	return []client.VpcSubnet{
		{
			ID:               DefaultPublicSubnetID,
			Cidr:             DefaultPublicSubnetCidr,
			AvailabilityZone: DefaultAvailabilityZone,
		},
	}
}

// PublicSubnetsAPI returns fake public subnets
func PublicSubnetsAPI() []api.VpcSubnet {
	return []api.VpcSubnet{
		{
			ID:               DefaultPublicSubnetID,
			Cidr:             DefaultPublicSubnetCidr,
			AvailabilityZone: DefaultAvailabilityZone,
		},
	}
}

// PrivateSubnets returns fake public subnets
func PrivateSubnets() []client.VpcSubnet {
	return []client.VpcSubnet{
		{
			ID:               DefaultPrivateSubnetID,
			Cidr:             DefaultPrivateSubnetCidr,
			AvailabilityZone: DefaultAvailabilityZone,
		},
	}
}

// PrivateSubnetsAPI returns fake private subnets
func PrivateSubnetsAPI() []api.VpcSubnet {
	return []api.VpcSubnet{
		{
			ID:               DefaultPrivateSubnetID,
			Cidr:             DefaultPrivateSubnetCidr,
			AvailabilityZone: DefaultAvailabilityZone,
		},
	}
}

// DatabaseSubnets returns fake public subnets
func DatabaseSubnets() []client.VpcSubnet {
	return []client.VpcSubnet{
		{
			ID:               DefaultDatabaseSubnetID,
			Cidr:             DefaultDatabaseSubnetCidr,
			AvailabilityZone: DefaultAvailabilityZone,
		},
	}
}

// Vpc returns a fake vpc
func Vpc() *client.Vpc {
	return &client.Vpc{
		ID:                       ID(),
		StackName:                StackNameVpc,
		CloudFormationTemplate:   CloudFormationTemplate(),
		VpcID:                    DefaultVpcID,
		Cidr:                     DefaultCidr,
		PublicSubnets:            PublicSubnets(),
		PrivateSubnets:           PrivateSubnets(),
		DatabaseSubnets:          DatabaseSubnets(),
		DatabaseSubnetsGroupName: DefaultDatabaseSubnetGroupName,
	}
}

// ClusterConfig returns a fake cluster config
func ClusterConfig() *v1alpha5.ClusterConfig {
	c, _ := clusterconfig.New(&clusterconfig.Args{
		ClusterName:            DefaultClusterName,
		PermissionsBoundaryARN: v1alpha1.PermissionsBoundaryARN(DefaultAWSAccountID),
		PrivateSubnets:         PrivateSubnetsAPI(),
		PublicSubnets:          PublicSubnetsAPI(),
		Region:                 DefaultRegion,
		Version:                DefaultVersion,
		VpcCidr:                DefaultCidr,
		VpcID:                  DefaultVpcID,
	})

	return c
}

// Cluster returns a fake cluster
func Cluster() *client.Cluster {
	return &client.Cluster{
		ID:     ID(),
		Name:   DefaultClusterName,
		Config: ClusterConfig(),
	}
}

// RotaterBucket returns a fake rotater bucket
func RotaterBucket() *client.S3Bucket {
	return &client.S3Bucket{
		Name:                   DefaultS3BucketName,
		StackName:              StackNameRotaterBucket,
		CloudFormationTemplate: string(CloudFormationTemplate()),
	}
}

// PostgresDatabase returns a fake postgres database
func PostgresDatabase() *client.PostgresDatabase {
	return &client.PostgresDatabase{
		ID:                           ID(),
		ApplicationName:              DefaultPostgresName,
		UserName:                     DefaultPostgresUserName,
		StackName:                    StackNamePostgresDatabase,
		AdminSecretFriendlyName:      DefaultPostgresSecretFriendlyName,
		EndpointAddress:              DefaultDatabaseEndpoint,
		EndpointPort:                 DefaultPostgresDatabasePort,
		OutgoingSecurityGroupID:      DefaultSecurityGroupID,
		SecretsManagerAdminSecretARN: DefaultPostgresSecretARN,
		LambdaPolicyARN:              DefaultPolicyARN,
		LambdaRoleARN:                DefaultRoleARN,
		LambdaFunctionARN:            DefaultPostgresLambdaFunctionARN,
		CloudFormationTemplate:       string(CloudFormationTemplate()),
		Namespace:                    DefaultNamespace,
		AdminSecretName:              DefaultPostgresAdminName,
		AdminSecretARN:               DefaultPostgresSecretARN,
		DatabaseConfigMapName:        DefaultPostgresConfigMapName,
		RotaterBucket:                RotaterBucket(),
	}
}

// ArgoCD returns a fake argo cd
func ArgoCD() *client.ArgoCD {
	return &client.ArgoCD{
		ID:             ID(),
		ArgoDomain:     DefaultArgoDomain,
		ArgoURL:        DefaultArgoURL,
		AuthDomain:     DefaultAuthDomain,
		Certificate:    Certificate(),
		IdentityClient: IdentityPoolClient(),
		PrivateKey:     KubernetesManifest(),
		Secret:         KubernetesManifest(),
		ClientSecret:   SecretParameter("okay"),
		SecretKey:      SecretParameter("something"),
		Chart:          Helm(),
	}
}

// KubePromStack returns a fake kube prom stack
func KubePromStack() *client.KubePromStack {
	return &client.KubePromStack{
		ID:                                ID(),
		AuthHostname:                      DefaultAuthDomain,
		CertificateARN:                    DefaultCertificateARN,
		ClientID:                          DefaultClientID,
		FargateCloudWatchPolicyARN:        DefaultPolicyARN,
		FargateProfilePodExecutionRoleARN: DefaultRoleARN,
		Hostname:                          DefaultMonitoringURL,
		SecretsAdminPassKey:               "admin-pass",
		SecretsAdminUserKey:               "admin-user",
		SecretsClientSecretKey:            "client-secret",
		SecretsConfigName:                 "secrets-cm",
		SecretsCookieSecretKey:            "cookie-secret",
		Certificate:                       Certificate(),
		Chart:                             Helm(),
		ExternalSecret:                    KubernetesManifest(),
		IdentityPoolClient:                IdentityPoolClient(),
	}
}

// GithubSecret returns a fake github secret
func GithubSecret() *client.GithubSecret {
	return &client.GithubSecret{
		Name:    "okctl-iac-privatekey",
		Path:    "some/path/okctl-iac-privatekey",
		Version: 2, // nolint: gomnd
	}
}

// GithubDeployKey return a fake github deploy key
func GithubDeployKey() *client.GithubDeployKey {
	return &client.GithubDeployKey{
		Organisation:     DefaultGithubOrg,
		Repository:       DefaultGithubName,
		Identifier:       234567, // nolint: gomnd
		Title:            "okctl-iac-deploykey",
		PublicKey:        "ssh-rsa y390uf30uf03",
		PrivateKeySecret: GithubSecret(),
	}
}

// GithubRepository returns a fake github repository
func GithubRepository() *client.GithubRepository {
	return &client.GithubRepository{
		ID:           ID(),
		Organisation: DefaultGithubURL,
		Repository:   DefaultGithubName,
		FullName:     DefaultGithubFullName,
		GitURL:       DefaultGithubURL,
		DeployKey:    GithubDeployKey(),
	}
}

// ContainerRepository returns a fake container repository
func ContainerRepository() *client.ContainerRepository {
	return &client.ContainerRepository{
		ClusterID:              ID(),
		ImageName:              "my-image",
		StackName:              DefaultImage,
		CloudFormationTemplate: string(CloudFormationTemplate()),
	}
}
