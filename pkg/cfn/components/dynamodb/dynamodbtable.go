// Package dynamodb knows how to create a dynamodb table
// cloud formation resource
package dynamodb

import (
	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/awslabs/goformation/v4/cloudformation/dynamodb"
	"github.com/oslokommune/okctl/pkg/cfn"
)

//nolint:deadcode,varcheck
const (
	billingModeProvisioned   = "PROVISIONED"
	billingModePayPerRequest = "PAY_PER_REQUEST"
)

//nolint:deadcode,varcheck
const (
	attributeTypeString = "S"
	attributeTypeNumber = "N"
	attributeTypeBinary = "B"
)

//nolint:deadcode,varcheck
const (
	keyTypeHash  = "HASH"
	keyTypeRange = "RANGE"
)

// Table contains the state required for building
// the CloudFormation template
type Table struct {
	StoredName   string
	TableName    string
	PartitionKey string
}

// Resource returns the cloud formation template
func (receiver *Table) Resource() cloudformation.Resource {
	return &dynamodb.Table{
		AttributeDefinitions: []dynamodb.Table_AttributeDefinition{
			{
				AttributeName: receiver.PartitionKey,
				AttributeType: attributeTypeString,
			},
		},
		BillingMode: billingModeProvisioned,
		KeySchema: []dynamodb.Table_KeySchema{
			{
				AttributeName: receiver.PartitionKey,
				KeyType:       keyTypeHash,
			},
		},
		PointInTimeRecoverySpecification: &dynamodb.Table_PointInTimeRecoverySpecification{
			PointInTimeRecoveryEnabled: true,
		},
		ProvisionedThroughput: &dynamodb.Table_ProvisionedThroughput{
			ReadCapacityUnits:  1,
			WriteCapacityUnits: 1,
		},
		SSESpecification: &dynamodb.Table_SSESpecification{
			SSEEnabled: true,
		},
		TableName: receiver.TableName,
	}
}

// Name returns the name of the resource
func (receiver *Table) Name() string {
	return receiver.StoredName
}

// Ref returns an AWS intrinsic ref to the resource
func (receiver *Table) Ref() string {
	return cloudformation.Ref(receiver.Name())
}

// NamedOutputs returns the named outputs
func (receiver *Table) NamedOutputs() map[string]cloudformation.Output {
	return cfn.NewValue(receiver.Name(), receiver.Ref()).NamedOutputs()
}

// New returns an initialised DynamoDB cloud formation template
// - https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-dynamodb-table.html
func New(resourceName, tableName string, partitionKey string) *Table {
	return &Table{
		StoredName:   resourceName,
		TableName:    tableName,
		PartitionKey: partitionKey,
	}
}
