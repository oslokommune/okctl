package aws

import (
	"fmt"
	"strconv"

	"github.com/oslokommune/okctl/pkg/cfn/components"

	"github.com/aws/aws-sdk-go/aws"
	dynamodbapi "github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/cfn"
)

type kvStoreProvider struct {
	provider v1alpha1.CloudProvider
}

// CreateStore knows how to create a DynamoDB table
func (k kvStoreProvider) CreateStore(opts api.CreateStoreOpts) error {
	b := cfn.New(components.NewDynamoDBTableComposer(
		storeNameGenerator(string(opts.Name)),
		opts.Keys[0],
	))

	stackName := cfn.NewStackNamer().DynamoDBTable(opts.ClusterID.ClusterName, string(opts.Name))

	template, err := b.Build()
	if err != nil {
		return errors.E(err, "building CloudFormation template")
	}

	r := cfn.NewRunner(k.provider)

	err = r.CreateIfNotExists(opts.ClusterID.ClusterName, stackName, template, nil, defaultTimeOut)
	if err != nil {
		return errors.E(err, "creating CloudFormation template")
	}

	return nil
}

// DeleteStore knows how to delete a DynamoDB table
func (k kvStoreProvider) DeleteStore(opts api.DeleteStoreOpts) error {
	stackName := cfn.NewStackNamer().DynamoDBTable(opts.ClusterID.ClusterName, string(opts.Name))

	r := cfn.NewRunner(k.provider)

	err := r.Delete(stackName)
	if err != nil {
		return errors.E(err, "deleting DynamoDB Table CloudFormation template")
	}

	return nil
}

// InsertItem knows how to insert a string into a DynamoDB table
func (k kvStoreProvider) InsertItem(opts api.InsertItemOpts) error {
	attrs, err := keyValueStoreFieldAsAttributeValue(opts.Item.Fields)
	if err != nil {
		return fmt.Errorf("parsing item fields: %w", err)
	}

	_, err = k.provider.DynamoDB().PutItem(&dynamodbapi.PutItemInput{
		TableName: aws.String(storeNameGenerator(string(opts.TableName))),
		Item:      attrs,
	})
	if err != nil {
		return fmt.Errorf("putting item: %w", err)
	}

	return nil
}

// GetString knows how to retrieve a string from a DynamoDB table
func (k kvStoreProvider) GetString(opts api.GetStringOpts) (string, error) {
	result, err := k.provider.DynamoDB().GetItem(&dynamodbapi.GetItemInput{
		TableName: aws.String(storeNameGenerator(string(opts.TableName))),
		Key: map[string]*dynamodbapi.AttributeValue{
			opts.Selector.Key: {S: aws.String(opts.Selector.Value)},
		},
	})
	if err != nil {
		return "", errors.E(err, "getting item")
	}

	if result.Item == nil {
		return "", errors.E(errors.New("not found"), errors.NotExist)
	}

	return *result.Item[opts.Field].S, nil
}

// RemoveItem knows how to delete an item from a DynamoDB table
func (k kvStoreProvider) RemoveItem(opts api.DeleteItemOpts) error {
	attrs, err := keyValueStoreFieldAsAttributeValue(map[string]interface{}{
		opts.Field: opts.Key,
	})
	if err != nil {
		return errors.E(err, "converting to attribute_value")
	}

	_, err = k.provider.DynamoDB().DeleteItem(&dynamodbapi.DeleteItemInput{
		TableName: aws.String(storeNameGenerator(string(opts.TableName))),
		Key:       attrs,
	})
	if err != nil {
		return errors.E(err, "deleting item")
	}

	return nil
}

func storeNameGenerator(original string) string {
	return fmt.Sprintf("okctl-%s", original)
}

func keyValueStoreFieldAsAttributeValue(fields map[string]interface{}) (map[string]*dynamodbapi.AttributeValue, error) {
	attrs := make(map[string]*dynamodbapi.AttributeValue, len(fields))

	for key, value := range fields {
		attrs[key] = &dynamodbapi.AttributeValue{}

		switch v := value.(type) {
		case string:
			attrs[key].S = aws.String(v)
		case int:
			attrs[key].N = aws.String(strconv.Itoa(v))
		default:
			return nil, fmt.Errorf("value type not supported: %T", value)
		}
	}

	return attrs, nil
}

// NewDynamoDBKeyValueStoreCloudProvider initializes a new DynamoDB key/value store cloud provider
func NewDynamoDBKeyValueStoreCloudProvider(provider v1alpha1.CloudProvider) api.KeyValueStoreCloudProvider {
	return &kvStoreProvider{
		provider: provider,
	}
}
