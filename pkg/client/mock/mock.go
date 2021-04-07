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
	DefaultRegion                 = "eu-west-1"
	DefaultAWSAccountID           = "123456789012"
	DefaultEnvironment            = "staging"
	DefaultRepository             = "okctl"
	DefaultClusterName            = "okctl-staging"
	DefaultDomain                 = "okctl-staging.oslo.systems"
	DefaultAuthDomain             = "auth.okctl-staging.oslo.systems"
	DefaultDomainFilter           = DefaultDomain
	DefaultFQDN                   = "okctl-staging.oslo.systems."
	DefaultCallbackURL            = "https://argocd.okctl-staging.oslo.systems/callback"
	DefaultHostedZoneID           = "Z0FAKE41FAKE6I841FAKE"
	DefaultAliasHostedZoneID      = "Z0FAKE67FAKE6I231FAKE"
	DefaultCertificateARN         = "arn:aws:acm:eu-west-1:123456789012:certificate/123456789012-1234-1234-1234-12345678"
	DefaultServiceAccountName     = "important-sa"
	DefaultPolicyARN              = "arn:aws:iam::123456789012:policy/policy-name-with-path"
	DefaultNamespace              = "kube-system"
	DefaultManifestName           = "okctl-cm"
	DefaultManifestType           = client.ManifestTypeConfigMap
	DefaultHelmReleaseName        = "okctl-helm-release"
	DefaultExternalDNSName        = "external-dns"
	DefaultSecretParameterName    = "release-secret"
	DefaultSecretParameterVersion = 1
	DefaultSecretParameterPath    = "/okctl/staging/release-secret"
	DefaultEmail                  = "bob@thebuilder.com"
	DefaultUserPoolID             = "TYUJBFW3893FAKE"
	DefaultPurpose                = "argocd"
	DefaultClientID               = "gehu-fgerg432-ewge"
	DefaultClientSecret           = "0ef90weug09jfqh3rf"

	StackNameHostedZone         = "okctl-staging-oslo-systems-HostedZone"
	StackNameCertificate        = "okctl-staging-oslo-systems-Certificate"
	StackNameManagedPolicy      = "okctl-staging-ManagedPolicy"
	StackNameIdentityPool       = "okctl-staging-IdentityPool"
	StackNameIdentityPoolClient = "okctl-staging-IdentityPoolClient"
	StackNameIdentityPoolUser   = "okctl-staging-bobthebuilder-IdentityPoolUser"
	StackNameRecordSetAlias     = "okctl-staging-RecordSetAlias"
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
		Environment:  DefaultEnvironment,
		Repository:   DefaultRepository,
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
		Name: DefaultExternalDNSName,
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
