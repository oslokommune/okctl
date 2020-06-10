package manager

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	cfPkg "github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/oslokommune/okctl/pkg/cfn"
)

type manager struct {
	b cfn.Builder
	t *cloudformation.Template
	c *sts.Credentials
}

func New(builder cfn.Builder, credentials *sts.Credentials) *manager {
	return &manager{
		b: builder,
		t: cloudformation.NewTemplate(),
		c: credentials,
	}
}

func (m *manager) Create(stackName string, timeout int64) error {
	resources, err := m.b.Build()
	if err != nil {
		return err
	}

	for _, resource := range resources {
		if _, hasKey := m.t.Resources[resource.Name()]; hasKey {
			return fmt.Errorf("already have resource with name: %s", resource.Name())
		}

		m.t.Resources[resource.Name()] = resource.Resource()
	}

	sess, err := session.NewSession(
		&aws.Config{
			Credentials: credentials.NewStaticCredentials(
				*m.c.AccessKeyId,
				*m.c.SecretAccessKey,
				*m.c.SessionToken,
			),
			Region: aws.String(endpoints.EuWest1RegionID),
		},
	)
	if err != nil {
		return err
	}

	cf := cfPkg.New(sess)

	body, err := m.t.YAML()
	if err != nil {
		return err
	}

	r, err := cf.CreateStack(&cfPkg.CreateStackInput{
		OnFailure:        aws.String(cfPkg.OnFailureDelete),
		StackName:        aws.String(stackName),
		TemplateBody:     aws.String(string(body)),
		TimeoutInMinutes: aws.Int64(timeout),
	})
	if err != nil {
		return err
	}

	return m.watchCreate(cf, r)
}

func (m *manager) watchCreate(cf *cfPkg.CloudFormation, r *cfPkg.CreateStackOutput) error {
	for {
		stack, err := cf.DescribeStacks(&cfPkg.DescribeStacksInput{
			NextToken: nil,
			StackName: r.StackId,
		})
		if err != nil {
			return err
		}

		if len(stack.Stacks) != 1 {
			return fmt.Errorf("expected 1 cloudformation stack to be created")
		}

		// nolint
		sleepTime := 5 * time.Second

		switch *stack.Stacks[0].StackStatus {
		case cfPkg.StackStatusCreateComplete:
			return nil
		case cfPkg.StackStatusCreateFailed:
			return fmt.Errorf("failed to create stack: %s", *stack.Stacks[0].StackStatusReason)
		case cfPkg.StackStatusCreateInProgress:
			time.Sleep(sleepTime)
		default:
			return fmt.Errorf("wtf")
		}
	}
}

func (m *manager) YAML() ([]byte, error) {
	return m.t.YAML()
}

func (m *manager) JSON() ([]byte, error) {
	return m.t.JSON()
}
