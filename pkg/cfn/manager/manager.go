// Package manager knows how to interact with AWS cloud formation stacks
package manager

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	cfPkg "github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/cfn"
	"github.com/sirupsen/logrus"
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
	Logger   *logrus.Entry
	Builder  cfn.Builder
	Provider v1alpha1.CloudProvider
	Template *cloudformation.Template
}

// New returns a new manager
func New(logger *logrus.Logger, provider v1alpha1.CloudProvider) *Manager {
	return &Manager{
		Logger:   logger.WithFields(logrus.Fields{}),
		Provider: provider,
		Template: cloudformation.NewTemplate(),
	}
}

// WithBuilder adds a builder to cloud formation stack manager
func (m *Manager) WithBuilder(builder cfn.Builder) *Manager {
	m.Builder = builder

	m.Logger = m.Logger.WithFields(logrus.Fields{
		"cloud_formation_stack_name": builder.StackName(),
	})

	return m
}

// Exists returns true if a cloud formation stack already exists
func (m *Manager) Exists() (bool, error) {
	stack, err := m.Provider.CloudFormation().DescribeStacks(&cfPkg.DescribeStacksInput{
		StackName: aws.String(m.Builder.StackName()),
	})
	if err != nil {
		switch e := err.(type) {
		case awserr.Error:
			if e.Code() == awsErrValidationError && fmt.Sprintf(stackDoesNotExistPattern, m.Builder.StackName()) == e.Message() {
				return false, nil
			}
		default:
			return false, err
		}
	}

	// Is this check really necessary?
	return m.StackStatusIsNotDeleted(stack.Stacks[0]), nil
}

// ProcessOutputFn defines a callback for handling output data
type ProcessOutputFn func(string) error

// Outputs processes the cloud formation stacks given the provided processors
func (m *Manager) Outputs(processors map[string]ProcessOutputFn) error {
	stack, err := m.Provider.CloudFormation().DescribeStacks(&cfPkg.DescribeStacksInput{
		StackName: aws.String(m.Builder.StackName()),
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
		StackName: aws.String(m.Builder.StackName()),
	})
	if err != nil {
		return false, err
	}

	return m.StackStatusIsNotTransitional(stack.Stacks[0]), nil
}

func (m *Manager) existsAndReady() error {
	exists, err := m.Exists()
	if err != nil {
		return err
	}

	ready, err := m.Ready()
	if err != nil {
		return err
	}

	if exists {
		if ready {
			return nil
		}

		return fmt.Errorf("stack: %s exists and is in a transitional state", m.Builder.StackName())
	}

	return nil
}

func (m *Manager) collectResources() error {
	for _, resource := range m.Builder.Resources() {
		if _, hasKey := m.Template.Resources[resource.Name()]; hasKey {
			return fmt.Errorf("already have resource with name: %s", resource.Name())
		}

		m.Template.Resources[resource.Name()] = resource.Resource()
	}

	return nil
}

func (m *Manager) collectOutputs() error {
	for _, output := range m.Builder.Outputs() {
		for key, value := range output.NamedOutputs() {
			if _, hasKey := m.Template.Outputs[key]; hasKey {
				return fmt.Errorf("already have output with name: %s", key)
			}

			m.Template.Outputs[key] = value
		}
	}

	return nil
}

// Delete a cloud formation stack
func (m *Manager) Delete(stackName string) error {
	_, err := m.Provider.CloudFormation().DeleteStack(&cfPkg.DeleteStackInput{
		StackName: aws.String(stackName),
	})
	if err != nil {
		return err
	}

	return m.watchDelete(stackName)
}

// CreateIfNotExists creates a cloud formation stack if none exists from before
func (m *Manager) CreateIfNotExists(timeout int64) error {
	err := m.existsAndReady()
	if err != nil {
		return err
	}

	err = m.Builder.Build()
	if err != nil {
		return err
	}

	err = m.collectResources()
	if err != nil {
		return err
	}

	err = m.collectOutputs()
	if err != nil {
		return err
	}

	body, err := m.Template.YAML()
	if err != nil {
		return err
	}

	r, err := m.Provider.CloudFormation().CreateStack(&cfPkg.CreateStackInput{
		OnFailure:        aws.String(cfPkg.OnFailureDelete),
		StackName:        aws.String(m.Builder.StackName()),
		TemplateBody:     aws.String(string(body)),
		TimeoutInMinutes: aws.Int64(timeout),
	})
	if err != nil {
		return err
	}

	return m.watchCreate(r)
}

func (m *Manager) watchDelete(stackName string) error {
	m.Logger.Info("Stack deletion request sent to AWS")

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
			m.Logger.Info("Waiting for stack deletion to complete.. sleeping for 30 seconds")
			time.Sleep(sleepTime)
		default:
			return fmt.Errorf("wtf")
		}
	}
}

// Reimplement this as wait
func (m *Manager) watchCreate(r *cfPkg.CreateStackOutput) error {
	m.Logger.Info("Stack creation request sent to AWS")

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
			m.Logger.Info("Waiting for stack creation to complete.. sleeping for 30 seconds")
			time.Sleep(sleepTime)
		default:
			return fmt.Errorf("wtf")
		}
	}
}

// YAML returns the cloud formation template in a yaml serialisation
func (m *Manager) YAML() ([]byte, error) {
	return m.Template.YAML()
}

// JSON returns the cloud formation template in a json serialisation
func (m *Manager) JSON() ([]byte, error) {
	return m.Template.JSON()
}
