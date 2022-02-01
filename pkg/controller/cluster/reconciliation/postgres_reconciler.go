package reconciliation

import (
	"context"
	"fmt"
	"regexp"

	"github.com/mishudark/errors"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"

	"github.com/oslokommune/okctl/pkg/controller/common/reconciliation"

	"github.com/oslokommune/okctl/pkg/cfn"

	clientCore "github.com/oslokommune/okctl/pkg/client/core"

	"github.com/oslokommune/okctl/pkg/client"
)

const postgresReconcilerIdentifier = "Postgres"

type postgresReconciler struct {
	client client.ComponentService
}

type database struct {
	Name      string
	User      string
	Namespace string
}

// Reconcile knows how to do what is necessary to ensure the desired state is achieved
//nolint:funlen
func (z *postgresReconciler) Reconcile(ctx context.Context, meta reconciliation.Metadata, state *clientCore.StateHandlers) (reconciliation.Result, error) {
	vpcExists, err := state.Vpc.HasVPC(meta.ClusterDeclaration.Metadata.Name)
	if err != nil {
		return reconciliation.Result{}, fmt.Errorf("checking vpc existence: %w", err)
	}

	if !vpcExists {
		if meta.Purge {
			return reconciliation.Result{Requeue: false}, nil
		}

		return reconciliation.Result{Requeue: true}, nil
	}

	vpc, err := state.Vpc.GetVpc(cfn.NewStackNamer().Vpc(meta.ClusterDeclaration.Metadata.Name))
	if err != nil {
		return reconciliation.Result{}, fmt.Errorf("getting vpc: %w", err)
	}

	clusterID := reconciliation.ClusterMetaAsID(meta.ClusterDeclaration.Metadata)
	dbSubnetIDs := subnetsAsIDList(vpc.DatabaseSubnets)
	dbSubnetCIDRs := subnetsAsCIDRList(vpc.DatabaseSubnets)

	actionMap, err := z.determineActions(meta, state)
	if err != nil {
		return reconciliation.Result{}, fmt.Errorf("determining course of action: %w", err)
	}

	for db, action := range actionMap {
		switch action {
		case reconciliation.ActionCreate:
			_, err = z.client.CreatePostgresDatabase(ctx, client.CreatePostgresDatabaseOpts{
				ID:                clusterID,
				ApplicationName:   db.Name,
				UserName:          db.User,
				VpcID:             vpc.VpcID,
				DBSubnetGroupName: vpc.DatabaseSubnetsGroupName,
				DBSubnetIDs:       dbSubnetIDs,
				DBSubnetCIDRs:     dbSubnetCIDRs,
				Namespace:         db.Namespace,
			})
			if err != nil {
				return reconciliation.Result{}, fmt.Errorf("creating postgres database: %w", err)
			}
		case reconciliation.ActionDelete:
			err = z.client.DeletePostgresDatabase(ctx, client.DeletePostgresDatabaseOpts{
				ID:              clusterID,
				ApplicationName: db.Name,
				Namespace:       db.Namespace,
				VpcID:           vpc.VpcID,
			})
			if err != nil {
				return reconciliation.Result{}, fmt.Errorf("deleting database: %s, got: %w", db.Name, err)
			}
		case reconciliation.ActionNoop:
			continue
		}
	}

	return reconciliation.Result{}, nil
}

func (z *postgresReconciler) determineActions(meta reconciliation.Metadata, state *clientCore.StateHandlers) (map[database]reconciliation.Action, error) {
	actionMap := make(map[database]reconciliation.Action)

	indicatedDatabases := meta.ClusterDeclaration.Databases.Postgres

	if meta.Purge {
		indicatedDatabases = []v1alpha1.ClusterDatabasesPostgres{}
	}

	existingDatabases, err := state.Component.GetPostgresDatabases()
	if err != nil {
		return nil, fmt.Errorf("checking existing postgres databases: %w", err)
	}

	for _, stateDatabase := range existingDatabases {
		db := database{
			Name:      stateDatabase.ApplicationName,
			User:      stateDatabase.UserName,
			Namespace: stateDatabase.Namespace,
		}

		actionMap[db] = reconciliation.ActionDelete
	}

	for _, indicatedDatabase := range indicatedDatabases {
		db := database{
			Name:      indicatedDatabase.Name,
			User:      indicatedDatabase.User,
			Namespace: indicatedDatabase.Namespace,
		}

		_, ok := actionMap[db]
		if ok {
			actionMap[db] = reconciliation.ActionNoop
		} else {
			err = validatePostgresDatabaseName(db.Name)
			if err != nil {
				return nil, fmt.Errorf("invalid database name: %w", err)
			}
			actionMap[db] = reconciliation.ActionCreate
		}
	}

	return actionMap, nil
}

// String returns the identifier type
func (z *postgresReconciler) String() string {
	return postgresReconcilerIdentifier
}

// NewPostgresReconciler creates a new reconciler for the Postgres resource
func NewPostgresReconciler(client client.ComponentService) reconciliation.Reconciler {
	return &postgresReconciler{
		client: client,
	}
}

func subnetsAsCIDRList(subnets []client.VpcSubnet) []string {
	cidrs := make([]string, len(subnets))

	for index, subnet := range subnets {
		cidrs[index] = subnet.Cidr
	}

	return cidrs
}

func subnetsAsIDList(subnets []client.VpcSubnet) []string {
	ids := make([]string, len(subnets))

	for index, subnet := range subnets {
		ids[index] = subnet.ID
	}

	return ids
}

func validatePostgresDatabaseName(dbName string) error {
	// Two reserved names that cannot be used as dbName:
	if dbName == "db" || dbName == "database" {
		return errors.New("'db' and 'database' are reserved")
	}

	// #####
	// From the AWS console when creating a database:

	// 1 to 60 alphanumeric characters or hyphens
	const maxLengthDatabaseName int = 60

	length := len(dbName)
	if length > maxLengthDatabaseName {
		return errors.New("cannot be longer than 60 characters")
	}

	// First character must be a letter
	startsWithNumber, _ := regexp.MatchString(`^[0-9]`, dbName)
	if startsWithNumber {
		return errors.New("cannot start with a number")
	}

	// Can't end with a hyphen.
	endsWithHyphen, _ := regexp.MatchString(`-$`, dbName)
	if endsWithHyphen {
		return errors.New("cannot end with a hyphen")
	}

	// Can't contain two consecutive hyphens
	containsConsecutiveHyphens, _ := regexp.MatchString(`--`, dbName)
	if containsConsecutiveHyphens {
		return errors.New("cannot have two consecutive hyphens")
	}

	// #####
	// When creating a postgres database in aws console you can create it with capital letters,
	// but it will be stored in lowercase letters.
	// In addition: the process of creating S3 bucket for the database will result in a error
	// if the database name contains a uppercase letter
	// So: don't allow the use of uppercase letters in databases.postgres.name
	containsUpperCase, _ := regexp.MatchString(`[A-Z]`, dbName)
	if containsUpperCase {
		return errors.New("cannot contain uppercase letters")
	}

	return nil
}
