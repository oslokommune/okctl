package okctl

import (
	"fmt"
	"regexp"

	val "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/cfn/builder/vpc"
	"github.com/oslokommune/okctl/pkg/cfn/manager"
	"github.com/oslokommune/okctl/pkg/cfn/process"
	"github.com/oslokommune/okctl/pkg/cloud"
	eksctlPkg "github.com/oslokommune/okctl/pkg/run/eksctl"
	"github.com/sirupsen/logrus"
)

const defaultTimeOut = 5

type CreateClusterOpts struct {
	AWSAccountID string
	Environment  string
	Cidr         string
}

func (o *CreateClusterOpts) LoggerContext(logger *logrus.Logger) *logrus.Entry {
	return logger.WithFields(logrus.Fields{
		"action":         "create",
		"component":      "cluster",
		"aws_account_id": o.AWSAccountID,
		"environment":    o.Environment,
		"cidr":           o.Cidr,
	})
}

func (o *CreateClusterOpts) Valid() error {
	const minLength = 3

	const maxLength = 10

	return val.ValidateStruct(o,
		val.Field(&o.AWSAccountID,
			val.Required,
			val.Match(regexp.MustCompile("^[0-9]{12}$")).
				Error("must consist of 12 digits"),
		),
		val.Field(&o.Environment,
			val.Required,
			val.Length(minLength, maxLength),
		),
		val.Field(&o.Cidr,
			val.Required,
		),
	)
}

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

func (o *Okctl) CreateCluster(opts CreateClusterOpts) error {
	ctxLogger := opts.LoggerContext(o.Logger)

	ctxLogger.Info("Starting EKS cluster creation process")

	for _, c := range o.RepoData.Clusters {
		if c.Environment == opts.Environment && c.AWS.AccountID == opts.AWSAccountID {
			return fmt.Errorf("cluster: %s, already exists", opts.Environment)
		}
	}

	builder := vpc.New(o.RepoData.Name, opts.Environment, opts.Cidr, o.Region())

	prov, err := cloud.New(o.Region(), o.CredentialsProvider)
	if err != nil {
		return err
	}

	m := manager.New(o.Logger, builder, prov.Provider)

	err = m.CreateIfNotExists(defaultTimeOut)
	if err != nil {
		return err
	}

	err = o.WriteCurrentRepoData()
	if err != nil {
		return err
	}

	eksctl, err := eksctlPkg.New(o.Logger, o.CredentialsProvider, o.BinariesProvider)
	if err != nil {
		return err
	}

	clusterConfig, err := ClusterConfig(o.ClusterName(opts.Environment), o.Region(), opts.Cidr, opts.AWSAccountID, m, prov.Provider)
	if err != nil {
		return err
	}

	err = eksctl.CreateCluster(o.Err, clusterConfig)
	if err != nil {
		return err
	}

	yaml, err := clusterConfig.YAML()
	if err != nil {
		return err
	}

	err = o.WriteToOutputDir(opts.Environment, "cluster.yml", yaml)
	if err != nil {
		return err
	}

	return err
}
