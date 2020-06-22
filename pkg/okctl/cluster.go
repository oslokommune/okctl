// Package okctl provides a high level interface for all okctl commands
package okctl

import (
	"fmt"
	"regexp"

	val "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/cfn/builder/vpc"
	"github.com/oslokommune/okctl/pkg/cfn/manager"
	"github.com/oslokommune/okctl/pkg/cfn/process"
	"github.com/oslokommune/okctl/pkg/config/repository"
	eksctlPkg "github.com/oslokommune/okctl/pkg/run/eksctl"
	"github.com/sirupsen/logrus"
)

const (
	defaultTimeOut = 5
	envMinLength   = 3
	envMaxLength   = 10
)

// CreateClusterOpts defines the inputs required for creating a
// new EKS cluster
type CreateClusterOpts struct {
	AWSAccountID string
	Environment  string
	Cidr         string
}

// LoggerContext creates a log context based on the provided inputs
func (o *CreateClusterOpts) LoggerContext(logger *logrus.Logger) *logrus.Entry {
	return logger.WithFields(logrus.Fields{
		"action":         "create",
		"component":      "cluster",
		"aws_account_id": o.AWSAccountID,
		"environment":    o.Environment,
		"cidr":           o.Cidr,
	})
}

// Valid determines if the provided cluster inputs are valid
func (o *CreateClusterOpts) Valid() error {
	return val.ValidateStruct(o,
		val.Field(&o.AWSAccountID,
			val.Required,
			val.Match(regexp.MustCompile("^[0-9]{12}$")).
				Error("must consist of 12 digits"),
		),
		val.Field(&o.Environment,
			val.Required,
			val.Length(envMinLength, envMaxLength),
		),
		val.Field(&o.Cidr,
			val.Required,
		),
	)
}

// DeleteClusterOpts defines the inputs required for removing an
// existing EKS cluster
type DeleteClusterOpts struct {
	Environment string
}

// LoggerContext creates a log context based on the provided inputs
func (o *DeleteClusterOpts) LoggerContext(logger *logrus.Logger) *logrus.Entry {
	return logger.WithFields(logrus.Fields{
		"environment": o.Environment,
	})
}

// Valid determines if the delete cluster opts are correct
func (o *DeleteClusterOpts) Valid() error {
	return val.ValidateStruct(o,
		val.Field(&o.Environment,
			val.Required,
			val.Length(envMinLength, envMaxLength),
		),
	)
}

// ClusterConfig knows how to populate an eksctl cluster config prior to creation
func ClusterConfig(name, region, cidr, awsAccountID string, m *manager.Manager, provider v1alpha1.CloudProvider) (*v1alpha1.ClusterConfig, error) {
	clusterConfig := v1alpha1.NewClusterConfig()

	clusterConfig.Metadata.Name = name
	clusterConfig.Metadata.Region = region
	clusterConfig.VPC.CIDR = cidr
	clusterConfig.IAM.FargatePodExecutionRolePermissionsBoundary = v1alpha1.PermissionsBoundaryARN(awsAccountID)
	clusterConfig.IAM.ServiceRolePermissionsBoundary = v1alpha1.PermissionsBoundaryARN(awsAccountID)

	return clusterConfig, m.Outputs(map[string]manager.ProcessOutputFn{
		"PrivateSubnetIds": process.Subnets(provider, clusterConfig.VPC.Subnets.Private),
		"PublicSubnetIds":  process.Subnets(provider, clusterConfig.VPC.Subnets.Public),
		"Vpc":              process.String(&clusterConfig.VPC.ID),
	})
}

// DeleteCluster removes an EKS cluster and all related resources
func (o *Okctl) DeleteCluster(opts DeleteClusterOpts) error {
	ctxLogger := opts.LoggerContext(o.Logger)

	ctxLogger.Info("Started EKS cluster deletion process")

	exists := o.HasCluster(opts.Environment)
	if !exists {
		return fmt.Errorf("couldn't find cluster for env: %s", opts.Environment)
	}

	eksctl, err := eksctlPkg.New(o.Logger, o.CredentialsProvider, o.BinariesProvider)
	if err != nil {
		return err
	}

	clusterConfig, err := o.ClusterConfig(opts.Environment)
	if err != nil {
		return err
	}

	err = eksctl.DeleteCluster(o.Err, clusterConfig)
	if err != nil {
		return err
	}

	err = o.DeleteClusterConfig(opts.Environment)
	if err != nil {
		return err
	}

	m := manager.New(o.Logger, o.CloudProvider)

	err = m.Delete(vpc.StackName(o.RepoData.Name, opts.Environment))
	if err != nil {
		return err
	}

	o.RepoData.Clusters = remove(opts.Environment, o.RepoData.Clusters)

	return nil
}

func remove(env string, clusters []repository.Cluster) []repository.Cluster {
	for i, c := range clusters {
		if c.Environment == env {
			clusters[i] = clusters[len(clusters)-1]
			return clusters[:len(clusters)-1]
		}
	}

	return clusters
}

// CreateCluster starts the creation of all resources related to an EKS cluster
// such as a VPC, etc.
func (o *Okctl) CreateCluster(opts CreateClusterOpts) error {
	var err error

	ctxLogger := opts.LoggerContext(o.Logger)

	ctxLogger.Info("Starting EKS cluster creation process")

	ctxLogger.Debugf("known repository data: %s", o.RepoData)

	m := manager.New(o.Logger, o.CloudProvider).
		WithBuilder(vpc.New(o.RepoData.Name, opts.Environment, opts.Cidr, o.Region()))

	err = m.CreateIfNotExists(defaultTimeOut)
	if err != nil {
		return err
	}

	o.RepoData.Clusters = remove(opts.Environment, o.RepoData.Clusters)

	o.RepoData.Clusters = append(o.RepoData.Clusters, repository.Cluster{
		Environment: opts.Environment,
		AWS: repository.AWS{
			AccountID: opts.AWSAccountID,
			Cidr:      opts.Cidr,
		},
	})

	err = o.WriteCurrentRepoData()
	if err != nil {
		return err
	}

	eksctl, err := eksctlPkg.New(o.Logger, o.CredentialsProvider, o.BinariesProvider)
	if err != nil {
		return err
	}

	clusterConfig, err := ClusterConfig(o.ClusterName(opts.Environment), o.Region(), opts.Cidr, opts.AWSAccountID, m, o.CloudProvider)
	if err != nil {
		return err
	}

	ctxLogger.Debugf("cluster config: %s", clusterConfig)

	err = eksctl.CreateCluster(o.Err, clusterConfig)
	if err != nil {
		return err
	}

	// Move these storage operations into cobra and PostRunE
	err = o.WriteClusterConfig(opts.Environment, clusterConfig)
	if err != nil {
		return err
	}

	return err
}
