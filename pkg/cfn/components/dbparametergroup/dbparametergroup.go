// Package dbparametergroup knows how to create a cloud formation
// for a db parameter group
package dbparametergroup

import (
	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/awslabs/goformation/v4/cloudformation/rds"
)

// DBParameterGroup stores the state for a cloud formation dbparametergroup
type DBParameterGroup struct {
	StoredName string
	Parameters map[string]string
}

// NamedOutputs returns the outputs for the resource
func (g *DBParameterGroup) NamedOutputs() map[string]cloudformation.Output {
	return nil
}

// Name returns the name of the resource
func (g *DBParameterGroup) Name() string {
	return g.StoredName
}

// Resource returns the cloud formation resource for a dbparametergroup
// For an exhaustive list of parameters:
// - https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/Appendix.PostgreSQL.CommonDBATasks.html#Appendix.PostgreSQL.CommonDBATasks.Parameters
func (g *DBParameterGroup) Resource() cloudformation.Resource {
	return &rds.DBParameterGroup{
		Family:     "postgres13",
		Parameters: g.Parameters,
	}
}

// Ref returns a cloud formation intrinsic ref to the resource
func (g *DBParameterGroup) Ref() string {
	return cloudformation.Ref(g.Name())
}

// New returns an initialised dbparametergroup
// - https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-rds-dbparametergroup.html
func New(resourceName string, parameters map[string]string) *DBParameterGroup {
	return &DBParameterGroup{
		StoredName: resourceName,
		Parameters: parameters,
	}
}
