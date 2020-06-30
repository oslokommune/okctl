// +build integration

package manager

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	gfn "github.com/awslabs/goformation/v4/cloudformation"
	"github.com/davecgh/go-spew/spew"
	"github.com/oslokommune/okctl/pkg/api/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/cfn/builder/vpc"
	"github.com/stretchr/testify/assert"
)

func NewCloudFormationSession(t *testing.T) *cloudformation.CloudFormation {
	assert.NotEmpty(t, os.Getenv("AWS_ACCESS_KEY_ID"))
	assert.NotEmpty(t, os.Getenv("AWS_SECRET_ACCESS_KEY"))

	sess, err := session.NewSession(
		&aws.Config{
			Region: aws.String(v1alpha1.RegionEuWest1),
		},
	)
	assert.NoError(t, err)

	return cloudformation.New(sess)
}

func NewVPC(t *testing.T) string {
	builder := vpc.New("test", "dev", "172.16.10.0/20", "eu-west-1")

	err := builder.Build()
	assert.NoError(t, err)

	resources := builder.Resources()

	template := gfn.NewTemplate()
	for _, resource := range resources {
		template.Resources[resource.Name()] = resource.Resource()
	}

	got, err := template.YAML()
	assert.NoError(t, err)

	return string(got)
}

func TestValidate(t *testing.T) {
	body := NewVPC(t)
	cf := NewCloudFormationSession(t)

	res, err := cf.ValidateTemplate(&cloudformation.ValidateTemplateInput{
		TemplateBody: &body,
		TemplateURL:  nil,
	})
	assert.NoError(t, err)
	log.Println(spew.Sdump(res))
}

func TestApply(t *testing.T) {
	body := NewVPC(t)
	cf := NewCloudFormationSession(t)

	result, err := cf.CreateStack(&cloudformation.CreateStackInput{
		OnFailure:        aws.String("DO_NOTHING"),
		StackName:        aws.String("test-eks-vpc"),
		TemplateBody:     &body,
		TimeoutInMinutes: aws.Int64(5),
	})

	assert.NoError(t, err)

	doCleanup := false
	defer func() {
		if doCleanup {
			_, err := cf.DeleteStack(&cloudformation.DeleteStackInput{
				StackName: result.StackId,
			})
			assert.NoError(t, err)
		}
	}()

Loop:
	for {
		stack, err := cf.DescribeStacks(&cloudformation.DescribeStacksInput{
			NextToken: nil,
			StackName: result.StackId,
		})
		assert.NoError(t, err)
		assert.Len(t, stack.Stacks, 1)

		assert.NotNil(t, stack.Stacks[0].StackStatus)

		switch *stack.Stacks[0].StackStatus {
		case cloudformation.StackStatusCreateComplete:
			log.Println("success")
			break Loop
		case cloudformation.StackStatusCreateFailed:
			log.Println(spew.Sdump(stack))
			assert.Fail(t, "failed to create stack")
			break Loop
		case cloudformation.StackStatusCreateInProgress:
			log.Println("still creating, sleeping..")
			time.Sleep(5 * time.Second)
		}
	}

}
