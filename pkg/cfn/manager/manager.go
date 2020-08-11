// Package manager knows how to interact with AWS cloud formation stacks
package manager

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	cfPkg "github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/oslokommune/okctl/pkg/api/okctl.io/v1alpha1"
)

const (
	//stackStatus              = "Stacks[].StackStatus"
	stackDoesNotExistPattern = "Stack with id %s does not exist"
	defaultSleepTime         = 30
	awsErrValidationError    = "ValidationError"
)

// Stack defines a single cloud formation stack
type Stack = cfPkg.Stack

// Manager stores state required for interacting with the AWS
// cloud formation API
type Manager struct {
	StackName    string
	TemplateBody []byte
	Provider     v1alpha1.CloudProvider
}

// New returns a new manager
func New(stackName string, templateBody []byte, provider v1alpha1.CloudProvider) *Manager {
	return &Manager{
		StackName:    stackName,
		TemplateBody: templateBody,
		Provider:     provider,
	}
}

// Exists returns true if a cloud formation stack already exists
func (m *Manager) Exists() (bool, error) {
	req := &cfPkg.DescribeStacksInput{
		StackName: aws.String(m.StackName),
	}

	stack, err := m.Provider.CloudFormation().DescribeStacks(req)
	if err != nil {
		switch e := err.(type) {
		case awserr.Error:
			if e.Code() == awsErrValidationError && fmt.Sprintf(stackDoesNotExistPattern, m.StackName) == e.Message() {
				return false, nil
			}

			return false, err
		default:
			return false, err
		}
	}

	return m.StackStatusIsNotDeleted(stack.Stacks[0]), nil
}

// ProcessOutputFn defines a callback for handling output data
type ProcessOutputFn func(string) error

// Outputs processes the cloud formation stacks given the provided processors
func (m *Manager) Outputs(processors map[string]ProcessOutputFn) error {
	stack, err := m.Provider.CloudFormation().DescribeStacks(&cfPkg.DescribeStacksInput{
		StackName: aws.String(m.StackName),
	})
	if err != nil {
		return err
	}

	for key, fn := range processors {
		for _, o := range stack.Stacks[0].Outputs {
			if *o.OutputKey == key {
				err = fn(*o.OutputValue)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// Ready returns true if the stack is in a valid steady state
func (m *Manager) Ready() (bool, error) {
	stack, err := m.Provider.CloudFormation().DescribeStacks(&cfPkg.DescribeStacksInput{
		StackName: aws.String(m.StackName),
	})
	if err != nil {
		return false, err
	}

	return m.StackStatusIsNotTransitional(stack.Stacks[0]), nil
}

func (m *Manager) existsAndReady() (bool, error) {
	exists, err := m.Exists()
	if err != nil {
		return false, err
	}

	if exists {
		ready, err := m.Ready()
		if err != nil {
			return false, err
		}

		if ready {
			return true, nil
		}

		return false, fmt.Errorf("stack: %s exists and is in a transitional state", m.StackName)
	}

	return false, nil
}

// Delete a cloud formation stack
func (m *Manager) Delete() error {
	_, err := m.Provider.CloudFormation().DeleteStack(&cfPkg.DeleteStackInput{
		StackName: aws.String(m.StackName),
	})
	if err != nil {
		return err
	}

	return m.watchDelete(m.StackName)
}

// CreateIfNotExists creates a cloud formation stack if none exists from before
func (m *Manager) CreateIfNotExists(timeout int64) error {
	yes, err := m.existsAndReady()
	if err != nil {
		return err
	}

	if yes {
		return nil
	}

	r, err := m.Provider.CloudFormation().CreateStack(&cfPkg.CreateStackInput{
		OnFailure:        aws.String(cfPkg.OnFailureDelete),
		StackName:        aws.String(m.StackName),
		TemplateBody:     aws.String(string(m.TemplateBody)),
		TimeoutInMinutes: aws.Int64(timeout),
	})
	if err != nil {
		return err
	}

	return m.watchCreate(r)
}

func (m *Manager) watchDelete(stackName string) error {
	for {
		stack, err := m.Provider.CloudFormation().DescribeStacks(&cfPkg.DescribeStacksInput{
			StackName: aws.String(stackName),
		})
		if err != nil {
			return err
		}

		if len(stack.Stacks) != 1 {
			return fmt.Errorf("expected 1 cloudformation stack to be deleted")
		}

		sleepTime := defaultSleepTime * time.Second

		switch *stack.Stacks[0].StackStatus {
		case cfPkg.StackStatusDeleteComplete:
			return nil
		case cfPkg.StackStatusDeleteFailed:
			return fmt.Errorf("failed to delete stack: %s", *stack.Stacks[0].StackStatusReason)
		case cfPkg.StackStatusDeleteInProgress:
			time.Sleep(sleepTime)
		default:
			return fmt.Errorf("wtf")
		}
	}
}

// Reimplement this as wait
func (m *Manager) watchCreate(r *cfPkg.CreateStackOutput) error {
	for {
		stack, err := m.Provider.CloudFormation().DescribeStacks(&cfPkg.DescribeStacksInput{
			StackName: r.StackId,
		})
		if err != nil {
			return err
		}

		if len(stack.Stacks) != 1 {
			return fmt.Errorf("expected 1 cloudformation stack to be created")
		}

		sleepTime := defaultSleepTime * time.Second

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
