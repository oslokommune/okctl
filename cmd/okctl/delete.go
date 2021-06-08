package main

import (
	"fmt"
	"io"
	"io/ioutil"

	"github.com/oslokommune/okctl/pkg/controller/cluster"

	"github.com/oslokommune/okctl/pkg/controller/cluster/reconciliation"
	common "github.com/oslokommune/okctl/pkg/controller/common/reconciliation"

	"github.com/oslokommune/okctl/pkg/spinner"

	"github.com/oslokommune/okctl/pkg/context"

	"github.com/AlecAivazis/survey/v2"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/spf13/cobra"
)

const (
	deleteClusterArgs = 0
)

func buildDeleteCommand(o *okctl.Okctl) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete commands",
	}

	deleteClusterCommand := buildDeleteClusterCommand(o)
	cmd.AddCommand(deleteClusterCommand)
	cmd.AddCommand(buildDeletePostgresCommand(o))

	return cmd
}

// DeleteClusterOpts contains the required inputs
type DeleteClusterOpts struct {
	AWSCredentialsType    string
	GithubCredentialsType string

	DisableSpinner bool
	Confirm        bool

	Region       string
	AWSAccountID string
	ClusterName  string
}

// Validate the inputs
func (o *DeleteClusterOpts) Validate() error {
	return validation.ValidateStruct(o,
		validation.Field(&o.AWSAccountID, validation.Required),
		validation.Field(&o.Region, validation.Required),
		validation.Field(&o.ClusterName, validation.Required),
	)
}

// nolint: gocyclo, funlen, gocognit
func buildDeleteClusterCommand(o *okctl.Okctl) *cobra.Command {
	opts := &DeleteClusterOpts{}

	cmd := &cobra.Command{
		Use:   "cluster",
		Short: "Delete a cluster",
		Long: `Delete all resources related to an EKS cluster,
including VPC, this is a highly destructive operation.`,
		Args: cobra.ExactArgs(deleteClusterArgs),
		PreRunE: func(_ *cobra.Command, args []string) error {
			o.AWSCredentialsType = opts.AWSCredentialsType
			o.GithubCredentialsType = opts.GithubCredentialsType

			err := o.Initialise()
			if err != nil {
				return fmt.Errorf("initialising: %w", err)
			}

			opts.Region = o.Declaration.Metadata.Region
			opts.AWSAccountID = o.Declaration.Metadata.AccountID
			opts.ClusterName = o.Declaration.Metadata.Name

			err = opts.Validate()
			if err != nil {
				return err
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			var spinnerWriter io.Writer
			if opts.DisableSpinner {
				spinnerWriter = ioutil.Discard
			} else {
				spinnerWriter = o.Err
			}

			spin, err := spinner.New("synchronizing", spinnerWriter)
			if err != nil {
				return fmt.Errorf("error creating spinner: %w", err)
			}

			id := api.ID{
				Region:       opts.Region,
				AWSAccountID: opts.AWSAccountID,
				ClusterName:  opts.ClusterName,
			}

			handlers := o.StateHandlers(o.StateNodes())

			services, err := o.ClientServices(handlers)
			if err != nil {
				return fmt.Errorf("error getting services: %w", err)
			}

			reconciliationManager := common.NewCompositeReconciler(spin,
				reconciliation.NewArgocdReconciler(services.ArgoCD, services.Github),
				reconciliation.NewAWSLoadBalancerControllerReconciler(services.AWSLoadBalancerControllerService),
				reconciliation.NewAutoscalerReconciler(services.Autoscaler),
				reconciliation.NewKubePrometheusStackReconciler(services.Monitoring),
				reconciliation.NewLokiReconciler(services.Monitoring),
				reconciliation.NewPromtailReconciler(services.Monitoring),
				reconciliation.NewTempoReconciler(services.Monitoring),
				reconciliation.NewBlockstorageReconciler(services.Blockstorage),
				reconciliation.NewClusterReconciler(services.Cluster),
				reconciliation.NewExternalDNSReconciler(services.ExternalDNS),
				reconciliation.NewExternalSecretsReconciler(services.ExternalSecrets),
				reconciliation.NewIdentityManagerReconciler(services.IdentityManager),
				reconciliation.NewVPCReconciler(services.Vpc),
				reconciliation.NewZoneReconciler(services.Domain),
				reconciliation.NewNameserverDelegationReconciler(services.NameserverHandler),
				reconciliation.NewNameserverDelegatedTestReconciler(services.Domain),
				reconciliation.NewUsersReconciler(services.IdentityManager),
				reconciliation.NewPostgresReconciler(services.Component),
				reconciliation.NewCleanupALBReconciler(o.CloudProvider),
				reconciliation.NewCleanupSGReconciler(o.CloudProvider),
				&reconciliation.PostgresGroupReconciler{},
				reconciliation.NewServiceQuotaReconciler(o.CloudProvider),
			)

			reconciliationManager.SetCommonMetadata(&common.CommonMetadata{
				Ctx:         o.Ctx,
				Out:         o.Out,
				ClusterID:   id,
				Declaration: o.Declaration,
			})

			synchronizeOpts := &cluster.SynchronizeOpts{
				Debug:                 o.Debug,
				Out:                   o.Out,
				DeleteAll:             true,
				ID:                    id,
				ClusterDeclaration:    o.Declaration,
				ReconciliationManager: reconciliationManager,
				State:                 handlers,
			}

			ready, err := checkIfReady(id.ClusterName, o, opts.Confirm)
			if err != nil || !ready {
				return err
			}

			err = cluster.DeleteCluster(synchronizeOpts)
			if err != nil {
				return fmt.Errorf("synchronizing declaration with state: %w", err)
			}

			return nil
		},
	}

	flags := cmd.Flags()

	flags.StringVarP(&opts.AWSCredentialsType,
		"aws-credentials-type",
		"a",
		context.AWSCredentialsTypeSAML,
		fmt.Sprintf(
			"The form of authentication to use for AWS. Possible values: [%s,%s]",
			context.AWSCredentialsTypeSAML,
			context.AWSCredentialsTypeAccessKey,
		),
	)
	flags.StringVarP(
		&opts.GithubCredentialsType,
		"github-credentials-type",
		"g",
		context.GithubCredentialsTypeDeviceAuthentication,
		fmt.Sprintf(
			"The form of authentication to use for Github. Possible values: [%s,%s]",
			context.GithubCredentialsTypeDeviceAuthentication,
			context.GithubCredentialsTypeToken,
		),
	)
	flags.BoolVar(
		&opts.DisableSpinner,
		"no-spinner",
		false,
		"disables progress spinner",
	)
	flags.BoolVarP(
		&opts.Confirm,
		"confirm",
		"y",
		false,
		"confirm all choices",
	)

	return cmd
}

func checkIfReady(clusterName string, o *okctl.Okctl, preConfirmed bool) (bool, error) {
	if preConfirmed {
		return true, nil
	}

	ready := false
	prompt := &survey.Confirm{
		Message: fmt.Sprintf("This will delete %s and all assosicated resources, are you sure you want to continue?", clusterName),
	}

	err := survey.AskOne(prompt, &ready)
	if err != nil {
		return false, err
	}

	if !ready {
		_, err = fmt.Fprintf(o.Err, "user wasn't ready to continue, aborting.")
		if err != nil {
			return false, err
		}

		return false, err
	}

	return true, nil
}
