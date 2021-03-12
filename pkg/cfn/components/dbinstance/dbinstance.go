// Package dbinstance knows how to build a cloud formation
// resource for an RDS DBInstance
package dbinstance

import (
	"fmt"

	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/awslabs/goformation/v4/cloudformation/rds"
	"github.com/oslokommune/okctl/pkg/cfn"
)

// DBInstance stores the state for a cloud formation dbinstance
type DBInstance struct {
	StoredName        string
	DBName            string
	DBSubnetGroupName string
	DBParameterGroup  cfn.NameReferencer
	MonitoringRole    cfn.Namer
	Admin             cfn.NameReferencer
	SecurityGroup     cfn.Namer
}

// NamedOutputs returns the outputs commonly used by other stacks or components
func (i *DBInstance) NamedOutputs() map[string]cloudformation.Output {
	return cfn.NewValueMap().
		Add(cfn.NewValue(i.Name(), i.Ref())).
		Add(cfn.NewValue("EndpointAddress", cloudformation.GetAtt(i.Name(), "Endpoint.Address"))).
		Add(cfn.NewValue("EndpointPort", cloudformation.GetAtt(i.Name(), "Endpoint.Port"))).
		NamedOutputs()
}

// Name returns the name of the resource
func (i *DBInstance) Name() string {
	return i.StoredName
}

// Ref returns a cloud formation intrinsic ref to the resource
func (i *DBInstance) Ref() string {
	return cloudformation.Ref(i.Name())
}

const (
	backupRetentionDays       = 7
	maxAllocatedStorageGB     = 100
	monitoringIntervalSeconds = 10
	insightsRetentionDays     = 7
)

// Resource returns the cloud formation resource for a dbinstance
//
// Supported postgres versions:
// - https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/CHAP_PostgreSQL.html#PostgreSQL.Concepts.General.version13
// Supported db instance types:
// - https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/Concepts.DBInstanceClass.html#Concepts.DBInstanceClass.Support
func (i *DBInstance) Resource() cloudformation.Resource {
	return &rds.DBInstance{
		AllocatedStorage:            "20", // Minimum 20 when gp2 is storage class
		AllowMajorVersionUpgrade:    false,
		AutoMinorVersionUpgrade:     true,
		BackupRetentionPeriod:       backupRetentionDays, // Enables automatic backup
		CopyTagsToSnapshot:          true,
		DBInstanceClass:             "db.t3.small", // Should this be configurable?
		DBName:                      i.DBName,
		DBParameterGroupName:        i.DBParameterGroup.Ref(),
		DBSubnetGroupName:           i.DBSubnetGroupName,
		DeleteAutomatedBackups:      true,
		DeletionProtection:          false,
		EnableCloudwatchLogsExports: []string{"postgresql", "upgrade"},
		EnablePerformanceInsights:   true,
		Engine:                      "postgres",
		EngineVersion:               "13.1",
		MasterUserPassword: cloudformation.Sub(
			fmt.Sprintf(`{{resolve:secretsmanager:${%s}::password}}`, i.Admin.Name()),
		),
		MasterUsername: cloudformation.Sub(
			fmt.Sprintf(`{{resolve:secretsmanager:${%s}::username}}`, i.Admin.Name()),
		),
		MaxAllocatedStorage:                maxAllocatedStorageGB,
		MonitoringInterval:                 monitoringIntervalSeconds,
		MonitoringRoleArn:                  cloudformation.GetAtt(i.MonitoringRole.Name(), "Arn"),
		MultiAZ:                            true,
		PreferredMaintenanceWindow:         "Mon:00:00-Mon:03:00",
		PreferredBackupWindow:              "03:00-06:00",
		PerformanceInsightsRetentionPeriod: insightsRetentionDays,
		Port:                               "5432",
		PubliclyAccessible:                 false,
		StorageEncrypted:                   true,
		StorageType:                        "gp2",
		UseDefaultProcessorFeatures:        true,
		VPCSecurityGroups: []string{
			cloudformation.GetAtt(i.SecurityGroup.Name(), "GroupId"),
		},
		AWSCloudFormationDependsOn: []string{
			i.DBParameterGroup.Name(),
			i.MonitoringRole.Name(),
		},
	}
}

// New creates a database instance
// - https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-rds-database-instance.html
//
// Things we should consider doing:
// - Configuring TLS for encrypted connection to the database
// - Enabling IAM database authentication
func New(
	resourceName, dbName, dbSubnetGroupName string,
	dbParameterGroup cfn.NameReferencer,
	monitoringRole cfn.Namer,
	admin cfn.NameReferencer,
	securityGroup cfn.Namer,
) *DBInstance {
	return &DBInstance{
		StoredName:        resourceName,
		DBName:            dbName,
		DBSubnetGroupName: dbSubnetGroupName,
		DBParameterGroup:  dbParameterGroup,
		MonitoringRole:    monitoringRole,
		Admin:             admin,
		SecurityGroup:     securityGroup,
	}
}
