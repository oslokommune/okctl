package cfn_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/oslokommune/okctl/pkg/api/mock"

	"github.com/oslokommune/okctl/pkg/config"

	"github.com/oslokommune/okctl/pkg/cfn"
	"github.com/oslokommune/okctl/pkg/cfn/components"
	"github.com/oslokommune/okctl/pkg/cloud"
	"github.com/oslokommune/okctl/pkg/integration"
	"github.com/stretchr/testify/assert"
)

// nolint: gochecknoglobals
var localstackTemplateBody = `AWSTemplateFormatVersion: '2010-09-09'
Description: Simple CloudFormation Test Template
Resources:
  S3Bucket:
    Type: AWS::S3::Bucket
    Properties:
      AccessControl: PublicRead
      BucketName: cf-test-bucket-1
      NotificationConfiguration:
        LambdaConfigurations:
        - Event: "s3:ObjectCreated:*"
          Function: aws:arn:lambda:test:testfunc
        QueueConfigurations:
        - Event: "s3:ObjectDeleted:*"
          Queue: aws:arn:sqs:test:testqueue
          Filter:
            S3Key:
              S3KeyFilter:
                Rules:
                  - { Name: name1, Value: value1 }
                  - { Name: name2, Value: value2 }
      Tags:
        - Key: foobar
          Value:
            Ref: SQSQueue
  SQSQueue:
    Type: AWS::SQS::Queue
    Properties:
      QueueName: cf-test-queue-1
      Tags:
        - Key: key1
          Value: value1
        - Key: key2
          Value: value2
  SNSTopic:
    Type: AWS::SNS::Topic
    Properties:
      TopicName: { "Fn::Join": [ "", [ { "Ref": "AWS::StackName" }, "-test-topic-1-1" ] ] }
      Tags:
        - Key: foo
          Value:
            Ref: S3Bucket
        - Key: bar
          Value: { "Fn::GetAtt": ["S3Bucket", "Arn"] }
  TopicSubscription:
    Type: AWS::SNS::Subscription
    Properties:
      Protocol: sqs
      TopicArn: !Ref SNSTopic
      Endpoint: !GetAtt SQSQueue.QueueArn
      FilterPolicy:
        eventType:
          - created
  KinesisStream:
    Type: AWS::Kinesis::Stream
    Properties:
      Name: cf-test-stream-1
  SQSQueueNoNameProperty:
    Type: AWS::SQS::Queue
  TestParam:
    Type: AWS::SSM::Parameter
    Properties:
      Name: cf-test-param-1
      Description: test param 1
      Type: String
      Value: value123
      Tags:
        tag1: value1
  ApiGatewayRestApi:
    Type: AWS::ApiGateway::RestApi
    Properties:
      Name: test-api
  GatewayResponseUnauthorized:
    Type: AWS::ApiGateway::GatewayResponse
    Properties:
      RestApiId:
        Ref: ApiGatewayRestApi
      ResponseType: UNAUTHORIZED
      ResponseTemplates:
        application/json: '{"errors":[{"message":"Custom text!", "extra":"Some extra info"}]}'
  GatewayResponseDefault500:
    Type: AWS::ApiGateway::GatewayResponse
    Properties:
      RestApiId:
        Ref: ApiGatewayRestApi
      ResponseType: DEFAULT_5XX
      ResponseTemplates:
        application/json: '{"errors":[{"message":$context.error.messageString}]}'`

//nolint:funlen
func TestTemplates(t *testing.T) {
	config.SkipUnlessIntegration(t)

	s := integration.NewLocalstack()
	err := s.Create(3 * time.Minute)
	assert.NoError(t, err)

	defer func() {
		err = s.Cleanup()
		assert.NoError(t, err)
	}()

	provider, err := cloud.NewFromSession("eu-west-1", "", s.AWSSession())
	assert.NoError(t, err)

	testCases := []struct {
		name         string
		stackName    string
		templateBody []byte
	}{
		{
			name:         "Empty template should work",
			stackName:    "emptyStack",
			templateBody: []byte("{\"Resources\":{}}"),
		},
		{
			name:         "Localstack cf template should work",
			stackName:    "localstack",
			templateBody: []byte(localstackTemplateBody),
		},
		{
			name:      "ExternalSecretsPolicy should be valid",
			stackName: "externalSecrets",
			templateBody: func() []byte {
				b, err := cfn.New(components.NewExternalSecretsPolicyComposer("cluster")).Build()
				assert.NoError(t, err)

				return b
			}(),
		},
	}

	for _, tc := range testCases {
		tc := tc

		// If test fails, grab the docker logs I think
		t.Run(tc.name, func(t *testing.T) {
			err := cfn.NewRunner(provider.Provider).CreateIfNotExists(
				mock.DefaultVersionInfo(),
				"myCluster",
				tc.stackName,
				tc.templateBody,
				nil,
				30,
			)
			assert.NoError(t, err)

			if err != nil {
				logs, err := s.Logs()
				assert.NoError(t, err)
				_, err = fmt.Fprintln(os.Stdout, logs)
				assert.NoError(t, err)
			}
		})
	}
}
