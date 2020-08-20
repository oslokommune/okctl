package cfn

import (
	"fmt"
	"regexp"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	cfPkg "github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/oslokommune/okctl/pkg/api/okctl.io/v1alpha1"
)

const (
	// CapabilityIam is required when a stack affects permissions
	CapabilityIam = cfPkg.CapabilityCapabilityIam
	// CapabilityNamedIam is required when a stack affects permissions
	CapabilityNamedIam = cfPkg.CapabilityCapabilityNamedIam
	// CapabilityAutoExpand is required when a stack contains macros
	CapabilityAutoExpand = cfPkg.CapabilityCapabilityAutoExpand

	stackDoesNotExitRegex = "^Stack with id (.*)%s(.*) does not exist$"
	defaultSleepTime      = 30
	awsErrValidationError = "ValidationError"
)

// Stack defines a single cloud formation stack
type Stack = cfPkg.Stack

// Runner stores state required for interacting with the AWS
// cloud formation API
type Runner struct {
	Provider v1alpha1.CloudProvider
}

// NewRunner returns a new runner
func NewRunner(provider v1alpha1.CloudProvider) *Runner {
	return &Runner{
		Provider: provider,
	}
}

// Exists returns true if a cloud formation stack already exists
func (r *Runner) Exists(stackName string) (bool, error) {
	req := &cfPkg.DescribeStacksInput{
		StackName: aws.String(stackName),
	}

	re, err := regexp.Compile(fmt.Sprintf(stackDoesNotExitRegex, stackName))
	if err != nil {
		return false, fmt.Errorf("failed to compile regex for stack existence: %w", err)
	}

	stack, err := r.Provider.CloudFormation().DescribeStacks(req)
	if err != nil {
		switch e := err.(type) {
		case awserr.Error:
			if e.Code() == awsErrValidationError && re.MatchString(e.Message()) {
				return false, nil
			}

			return false, err
		default:
			return false, err
		}
	}

	return r.StackStatusIsNotDeleted(stack.Stacks[0]), nil
}

// ProcessOutputFn defines a callback for handling output data
type ProcessOutputFn func(string) error

// Outputs processes the cloud formation stacks given the provided processors
func (r *Runner) Outputs(stackName string, processors map[string]ProcessOutputFn) error {
	stack, err := r.Provider.CloudFormation().DescribeStacks(&cfPkg.DescribeStacksInput{
		StackName: aws.String(stackName),
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
func (r *Runner) Ready(stackName string) (bool, error) {
	stack, err := r.Provider.CloudFormation().DescribeStacks(&cfPkg.DescribeStacksInput{
		StackName: aws.String(stackName),
	})
	if err != nil {
		return false, err
	}

	return r.StackStatusIsNotTransitional(stack.Stacks[0]), nil
}

func (r *Runner) existsAndReady(stackName string) (bool, error) {
	exists, err := r.Exists(stackName)
	if err != nil {
		return false, err
	}

	if exists {
		ready, err := r.Ready(stackName)
		if err != nil {
			return false, err
		}

		if ready {
			return true, nil
		}

		return false, fmt.Errorf("stack: %s exists and is in a transitional state", stackName)
	}

	return false, nil
}

// Delete a cloud formation stack
func (r *Runner) Delete(stackName string) error {
	_, err := r.Provider.CloudFormation().DeleteStack(&cfPkg.DeleteStackInput{
		StackName: aws.String(stackName),
	})
	if err != nil {
		return err
	}

	return r.watchDelete(stackName)
}

// CreateIfNotExists creates a cloud formation stack if none exists from before
func (r *Runner) CreateIfNotExists(stackName string, template []byte, capabilities []string, timeout int64) error {
	yes, err := r.existsAndReady(stackName)
	if err != nil {
		return err
	}

	if yes {
		return nil
	}

	stackInput := &cfPkg.CreateStackInput{
		OnFailure:        aws.String(cfPkg.OnFailureDelete),
		StackName:        aws.String(stackName),
		TemplateBody:     aws.String(string(template)),
		TimeoutInMinutes: aws.Int64(timeout),
	}

	for _, c := range capabilities {
		c := c
		stackInput.Capabilities = append(stackInput.Capabilities, &c)
	}

	_, err = r.Provider.CloudFormation().CreateStack(stackInput)
	if err != nil {
		return err
	}

	return r.watchCreate(stackName)
}

func (r *Runner) watchDelete(stackName string) error {
	re, err := regexp.Compile(fmt.Sprintf(stackDoesNotExitRegex, stackName))
	if err != nil {
		return fmt.Errorf("failed to compile regex for stack existence: %w", err)
	}

	for {
		// https://docs.aws.amazon.com/sdk-for-go/api/service/cloudformation/#CloudFormation.DeleteStack
		// Deleted stacks do not show up in the DescribeStacks API if the deletion has been completed successfully.
		stack, err := r.Provider.CloudFormation().DescribeStacks(&cfPkg.DescribeStacksInput{
			StackName: aws.String(stackName),
		})
		if err != nil {
			switch e := err.(type) {
			case awserr.Error:
				if e.Code() == awsErrValidationError && re.MatchString(e.Message()) {
					return nil
				}

				return err
			default:
				return err
			}
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
// - Should provide a context, where we can wait for some timeout
// - How should we handle describe that fails?
func (r *Runner) watchCreate(stackName string) error {
	for {
		stack, err := r.Provider.CloudFormation().DescribeStacks(&cfPkg.DescribeStacksInput{
			StackName: aws.String(stackName),
		})

		if err != nil {
			return fmt.Errorf("failed to describe stack after create: %w", err)
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

// This is verbatim stolen from: https://github.com/weaveworks/eksctl/blob/master/pkg/cfn/manager/api.go

// StackStatusIsNotTransitional returns true if stack is in a steady state
func (*Runner) StackStatusIsNotTransitional(s *Stack) bool {
	for _, state := range nonTransitionalReadyStackStatuses() {
		if *s.StackStatus == state {
			return true
		}
	}

	return false
}

func nonTransitionalReadyStackStatuses() []string {
	return []string{
		cfPkg.StackStatusCreateComplete,
		cfPkg.StackStatusUpdateComplete,
		cfPkg.StackStatusRollbackComplete,
		cfPkg.StackStatusUpdateRollbackComplete,
	}
}

// StackStatusIsNotReady returns true if the stack is in a transitional state
func (*Runner) StackStatusIsNotReady(s *Stack) bool {
	for _, state := range nonReadyStackStatuses() {
		if *s.StackStatus == state {
			return true
		}
	}

	return false
}

func nonReadyStackStatuses() []string {
	return []string{
		cfPkg.StackStatusCreateInProgress,
		cfPkg.StackStatusCreateFailed,
		cfPkg.StackStatusRollbackInProgress,
		cfPkg.StackStatusRollbackFailed,
		cfPkg.StackStatusDeleteInProgress,
		cfPkg.StackStatusDeleteFailed,
		cfPkg.StackStatusUpdateInProgress,
		cfPkg.StackStatusUpdateCompleteCleanupInProgress,
		cfPkg.StackStatusUpdateRollbackInProgress,
		cfPkg.StackStatusUpdateRollbackFailed,
		cfPkg.StackStatusUpdateRollbackCompleteCleanupInProgress,
		cfPkg.StackStatusReviewInProgress,
	}
}

// StackStatusIsNotDeleted returns true if the stack exists in some form
func (*Runner) StackStatusIsNotDeleted(s *Stack) bool {
	for _, state := range allNonDeletedStackStatuses() {
		if *s.StackStatus == state {
			return true
		}
	}

	return false
}

func allNonDeletedStackStatuses() []string {
	return []string{
		cfPkg.StackStatusCreateInProgress,
		cfPkg.StackStatusCreateFailed,
		cfPkg.StackStatusCreateComplete,
		cfPkg.StackStatusRollbackInProgress,
		cfPkg.StackStatusRollbackFailed,
		cfPkg.StackStatusRollbackComplete,
		cfPkg.StackStatusDeleteInProgress,
		cfPkg.StackStatusDeleteFailed,
		cfPkg.StackStatusUpdateInProgress,
		cfPkg.StackStatusUpdateCompleteCleanupInProgress,
		cfPkg.StackStatusUpdateComplete,
		cfPkg.StackStatusUpdateRollbackInProgress,
		cfPkg.StackStatusUpdateRollbackFailed,
		cfPkg.StackStatusUpdateRollbackCompleteCleanupInProgress,
		cfPkg.StackStatusUpdateRollbackComplete,
		cfPkg.StackStatusReviewInProgress,
	}
}
