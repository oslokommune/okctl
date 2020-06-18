package okctl

import (
	"fmt"
	"os"
	"path"
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
	return val.ValidateStruct(o,
		val.Field(&o.AWSAccountID,
			val.Required,
			val.Match(regexp.MustCompile("^[0-9]{12}$")).
				Error("must consist of 12 digits"),
		),
		val.Field(&o.Environment,
			val.Required,
			// nolint
			val.Length(3, 10),
		),
		val.Field(&o.Cidr,
			val.Required,
		),
	)
}

// nolint
func (o *Okctl) CreateCluster(opts CreateClusterOpts) error {
	ctxLogger := opts.LoggerContext(o.Logger)

	ctxLogger.Info("Starting EKS cluster creation process")

	for _, c := range o.RepoData.Clusters {
		if c.Environment == opts.Environment && c.AWS.AccountID == opts.AWSAccountID {
			return fmt.Errorf("cluster: %s, already exists", opts.Environment)
		}
	}

	builder := vpc.New(o.RepoData.Name, opts.Environment, opts.Cidr, o.Region())

	provider, err := cloud.New(o.Region(), o.CredentialsProvider)
	if err != nil {
		return err
	}

	m := manager.New(o.Logger, builder, provider.Provider)

	ctxLogger.Info("Creating EKS compatible VPC")

	ready, err := m.Ready()
	if err != nil {
		return err
	}

	if !ready {
		err = m.Create(5)
		if err != nil {
			return err
		}
	}

	// TODO: Move this stuff out of here..
	repoDataPath, err := o.GetRepoDataPath()
	if err != nil {
		return err
	}

	f, err := o.FileSystem.OpenFile(repoDataPath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer func() {
		err = f.Close()
	}()

	d, err := o.RepoData.YAML()
	if err != nil {
		return err
	}
	_, err = f.Write(d)
	if err != nil {
		return err
	}

	ctxLogger.Info("Done creating EKS VPC")

	eksctl, err := eksctlPkg.New(o.Logger, o.CredentialsProvider, o.BinariesProvider)
	if err != nil {
		return err
	}

	clusterConfig := v1alpha1.NewClusterConfig()

	clusterConfig.Metadata.Name = o.ClusterName(opts.Environment)
	clusterConfig.Metadata.Region = o.Region()
	clusterConfig.VPC.CIDR = opts.Cidr
	clusterConfig.IAM.FargatePodExecutionRolePermissionsBoundary = v1alpha1.PermissionsBoundaryARN(opts.AWSAccountID)
	clusterConfig.IAM.ServiceRolePermissionsBoundary = v1alpha1.PermissionsBoundaryARN(opts.AWSAccountID)

	err = m.Outputs(map[string]manager.ProcessOutputFn{
		"PrivateSubnetIds": process.Subnets(provider.Provider, clusterConfig.VPC.Subnets.Private),
		"PublicSubnetIds":  process.Subnets(provider.Provider, clusterConfig.VPC.Subnets.Public),
		"Vpc":              process.String(&clusterConfig.VPC.ID),
	})
	if err != nil {
		return err
	}

	yaml, err := clusterConfig.YAML()
	if err != nil {
		return err
	}

	ctxLogger.Debug(string(yaml))

	err = eksctl.CreateCluster(o.Err, clusterConfig)
	if err != nil {
		return err
	}

	outDir, err := o.GetRepoOutputDir(opts.Environment)
	if err != nil {
		return err
	}

	err = os.MkdirAll(outDir, 0744)
	if err != nil {
		return err
	}

	file, err := o.FileSystem.OpenFile(path.Join(outDir, "cluster.yml"), os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	err = file.Close()

	defer func() {
		err = file.Close()
	}()

	_, err = file.Write(yaml)
	if err != nil {
		return err
	}

	return err
}
