package kubeconfig

import (
	"encoding/base64"
	"fmt"
	"github.com/oslokommune/okctl/pkg/config/constant"
	"strings"

	"github.com/oslokommune/okctl/pkg/apis/eksctl.io/v1alpha5"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"

	"github.com/aws/aws-sdk-go/service/eks"

	"github.com/aws/aws-sdk-go/aws"
	"k8s.io/client-go/tools/clientcmd"
	clientCmdApi "k8s.io/client-go/tools/clientcmd/api"
)

// Config contains the
type Config struct {
	content clientCmdApi.Config
}

// Bytes returns the serialised kubeconfig
func (k *Config) Bytes() ([]byte, error) {
	return clientcmd.Write(k.content)
}

// Config returns the plain config
func (k *Config) Config() clientCmdApi.Config {
	return k.content
}

// Validate a configuration
func (k *Config) Validate() error {
	return clientcmd.Validate(k.content)
}

// Getter defines the kubeconf get interactions
type Getter interface {
	Get() (*Config, error)
}

type kubeConfig struct {
	provider v1alpha1.CloudProvider
	cfg      *v1alpha5.ClusterConfig
}

// New returns an initialised kubeconfig creator
func New(clusterConfig *v1alpha5.ClusterConfig, provider v1alpha1.CloudProvider) Getter {
	return &kubeConfig{
		provider: provider,
		cfg:      clusterConfig,
	}
}

// Get the kubeconfig by describing the EKS cluster
func (a *kubeConfig) Get() (*Config, error) {
	c, err := a.provider.EKS().DescribeCluster(&eks.DescribeClusterInput{
		Name: aws.String(a.cfg.Metadata.Name),
	})
	if err != nil {
		return nil, err
	}

	cluster := c.Cluster
	if cluster == nil {
		return nil, fmt.Errorf(constant.ClusterNilError)
	}

	var data []byte

	switch *cluster.Status {
	case eks.ClusterStatusCreating, eks.ClusterStatusDeleting, eks.ClusterStatusFailed:
		return nil, fmt.Errorf(constant.CreateKubeconfigWithInvalidStatusError, *cluster.Status)
	default:
		data, err = base64.StdEncoding.DecodeString(*cluster.CertificateAuthority.Data)
		if err != nil {
			return nil, fmt.Errorf(constant.DecodeCertificateAuthorityDataError, err)
		}
	}

	a.cfg.Status = &v1alpha5.ClusterStatus{
		Endpoint:                 aws.StringValue(cluster.Endpoint),
		CertificateAuthorityData: data,
		ARN:                      aws.StringValue(cluster.Arn),
		StackName:                aws.StringValue(cluster.Name),
	}

	kubeCfg := Create(GetUser(a.provider.PrincipalARN()), a.cfg)

	return &Config{
		content: kubeCfg,
	}, nil
}

// Create returns an initialised kubeconfig
func Create(username string, cfg *v1alpha5.ClusterConfig) clientCmdApi.Config {
	contextName := fmt.Sprintf("%s@%s", username, cfg.Metadata.String())

	return clientCmdApi.Config{
		Kind:        "Config",
		APIVersion:  "v1",
		Preferences: clientCmdApi.Preferences{},
		Clusters: map[string]*clientCmdApi.Cluster{
			cfg.Metadata.String(): {
				Server:                   cfg.Status.Endpoint,
				CertificateAuthorityData: cfg.Status.CertificateAuthorityData,
			},
		},
		AuthInfos: map[string]*clientCmdApi.AuthInfo{
			contextName: {
				Exec: &clientCmdApi.ExecConfig{
					APIVersion: "client.authentication.k8s.io/v1alpha1",
					Command:    "aws-iam-authenticator",
					Env: []clientCmdApi.ExecEnvVar{
						{
							Name:  "AWS_STS_REGIONAL_ENDPOINTS",
							Value: "regional",
						},
						{
							Name:  "AWS_DEFAULT_REGION",
							Value: cfg.Metadata.Region,
						},
						{
							Name:  "AWS_PROFILE",
							Value: "default",
						},
					},
					Args: []string{"token", "-i", cfg.Metadata.Name},
				},
			},
		},
		Contexts: map[string]*clientCmdApi.Context{
			contextName: {
				Cluster:  cfg.Metadata.String(),
				AuthInfo: contextName,
			},
		},
		CurrentContext: contextName,
	}
}

// GetUser returns the username part
// - This is stolen from eksctl
func GetUser(iamRoleARN string) string {
	usernameParts := strings.Split(iamRoleARN, "/")
	if len(usernameParts) > 1 {
		return usernameParts[len(usernameParts)-1]
	}

	return "iam-root-account"
}
