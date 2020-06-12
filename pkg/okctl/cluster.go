package okctl

import (
	"fmt"
	"os"
	"regexp"

	val "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/oslokommune/okctl/pkg/cfn/builder/vpc"
	"github.com/oslokommune/okctl/pkg/cfn/manager"
	"github.com/oslokommune/okctl/pkg/cloud"
	"github.com/oslokommune/okctl/pkg/login"
	"github.com/oslokommune/okctl/pkg/stage"
	"github.com/oslokommune/okctl/pkg/storage"
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

	appDataDir, err := o.GetAppDataDir()
	if err != nil {
		return err
	}

	store := storage.NewFileSystemStorage(appDataDir)

	stagers, err := stage.FromConfig(o.AppData.Binaries, o.AppData.Host, store)
	if err != nil {
		return err
	}

	for _, s := range stagers {
		err = s.Run()
		if err != nil {
			return err
		}
	}

	for _, c := range o.RepoData.Clusters {
		if c.Environment == opts.Environment && c.AWS.AccountID == opts.AWSAccountID {
			return fmt.Errorf("cluster: %s, already exists", opts.Environment)
		}
	}

	if o.NoInput {
		return fmt.Errorf("create cluster requires user input for now")
	}

	interactiveLogin, err := login.Interactive(opts.AWSAccountID, o.Region(), o.Username())
	if err != nil {
		return err
	}

	builder := vpc.New(o.RepoData.Name, opts.Environment, opts.Cidr, o.Region())

	provider, err := cloud.New(o.Region(), interactiveLogin)
	if err != nil {
		return err
	}

	m := manager.New(o.Logger, builder, provider.Provider)

	ctxLogger.Info("Creating EKS compatible VPC")

	err = m.Create(builder.StackName(), 5)
	if err != nil {
		return err
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

	ctxLogger.Warn("We do not support creating the EKS cluster itself yet, exiting")

	return nil
}
