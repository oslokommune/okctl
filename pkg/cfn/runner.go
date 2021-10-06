package cfn

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	merrors "github.com/mishudark/errors"

	"github.com/oslokommune/okctl/pkg/version"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	cfPkg "github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
)

const (
	// CapabilityNamedIam is required when a stack affects permissions
	CapabilityNamedIam = cfPkg.CapabilityCapabilityNamedIam
	// CapabilityAutoExpand is required when a stack contains macros
	CapabilityAutoExpand = cfPkg.CapabilityCapabilityAutoExpand

	defaultSleepTime      = 30
	awsErrValidationError = "ValidationError"
)

var stackDoesNotExistRe = regexp.MustCompile("^Stack with id (.*) does not exist$")

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

	stack, err := r.Provider.CloudFormation().DescribeStacks(req)
	if err != nil {
		switch e := err.(type) {
		case awserr.Error:
			if e.Code() == awsErrValidationError && stackDoesNotExistRe.MatchString(e.Message()) {
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
func (r *Runner) CreateIfNotExists(
	versionInfo version.Info, clusterName, stackName string, template []byte, capabilities []string, timeout int64,
) error {
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
		Tags: []*cfPkg.Tag{
			{
				Key:   aws.String(v1alpha1.OkctlVersionTag),
				Value: aws.String(versionInfo.Version),
			},
			{
				Key:   aws.String(v1alpha1.OkctlCommitTag),
				Value: aws.String(versionInfo.ShortCommit),
			},
			{
				Key:   aws.String(v1alpha1.OkctlManagedTag),
				Value: aws.String("true"),
			},
			{
				Key:   aws.String(v1alpha1.OkctlClusterNameTag),
				Value: aws.String(clusterName),
			},
		},
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

// Get fetches an existing CloudFormation stack
func (r *Runner) Get(stackName string) (cfPkg.Stack, error) {
	response, err := r.Provider.CloudFormation().DescribeStacks(&cfPkg.DescribeStacksInput{
		StackName: aws.String(stackName),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() { //nolint:gocritic
			case awsErrValidationError:
				if stackDoesNotExistRe.MatchString(aerr.Message()) {
					return cfPkg.Stack{}, merrors.E(err, "stack not found", merrors.NotExist)
				}
			}
		}

		return cfPkg.Stack{}, fmt.Errorf("describing stack %s: %w", stackName, err)
	}

	return *response.Stacks[0], nil
}

// GetTemplate fetches an existing CloudFormation stack template
func (r *Runner) GetTemplate(stackName string) ([]byte, error) {
	response, err := r.Provider.CloudFormation().GetTemplate(&cfPkg.GetTemplateInput{
		StackName:     aws.String(stackName),
		TemplateStage: aws.String(cfPkg.TemplateStageOriginal),
	})
	if err != nil {
		return []byte{}, fmt.Errorf("describing stack %s: %w", stackName, err)
	}

	return []byte(*response.TemplateBody), nil
}

// Update updates an existing CloudFormation template
func (r *Runner) Update(stackName string, template []byte) error {
	response, err := r.Provider.CloudFormation().DescribeStacks(&cfPkg.DescribeStacksInput{
		StackName: aws.String(stackName),
	})
	if err != nil {
		return fmt.Errorf("describing stack %s: %w", stackName, err)
	}

	if len(response.Stacks) > 1 {
		return fmt.Errorf("too many stacks selected")
	}

	existingStack := response.Stacks[0]

	stackInput := &cfPkg.UpdateStackInput{
		StackName:    aws.String(stackName),
		TemplateBody: aws.String(string(template)),
		Capabilities: existingStack.Capabilities,
		Tags:         existingStack.Tags,
	}

	_, err = r.Provider.CloudFormation().UpdateStack(stackInput)
	if err != nil {
		return fmt.Errorf("updating CloudFormation stack: %w", err)
	}

	return r.watchUpdate(stackName)
}

func (r *Runner) watchDelete(stackName string) error {
	for {
		// https://docs.aws.amazon.com/sdk-for-go/api/service/cloudformation/#CloudFormation.DeleteStack
		// Deleted stacks do not show up in the DescribeStacks API if the deletion has been completed successfully.
		stack, err := r.Provider.CloudFormation().DescribeStacks(&cfPkg.DescribeStacksInput{
			StackName: aws.String(stackName),
		})
		if err != nil {
			switch e := err.(type) {
			case awserr.Error:
				if e.Code() == awsErrValidationError && stackDoesNotExistRe.MatchString(e.Message()) {
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
			return r.detailedErr(stack.Stacks[0])
		case cfPkg.StackStatusDeleteInProgress:
			time.Sleep(sleepTime)
		default:
			return r.detailedErr(stack.Stacks[0])
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
			return r.detailedErr(stack.Stacks[0])
		case cfPkg.StackStatusCreateInProgress:
			time.Sleep(sleepTime)
		default:
			return r.detailedErr(stack.Stacks[0])
		}
	}
}

func (r *Runner) watchUpdate(stackName string) error {
	for {
		stack, err := r.Provider.CloudFormation().DescribeStacks(&cfPkg.DescribeStacksInput{
			StackName: aws.String(stackName),
		})
		if err != nil {
			return fmt.Errorf("describing stack after update: %w", err)
		}

		if len(stack.Stacks) != 1 {
			return fmt.Errorf("expected 1 cloudformation stack to be updated")
		}

		sleepTime := defaultSleepTime * time.Second

		switch *stack.Stacks[0].StackStatus {
		case cfPkg.StackStatusUpdateInProgress:
			time.Sleep(sleepTime)
		case cfPkg.StackStatusUpdateComplete:
			return nil
		case cfPkg.StackStatusUpdateRollbackInProgress:
			time.Sleep(sleepTime)
		case cfPkg.StackStatusUpdateRollbackComplete:
			return r.detailedErr(stack.Stacks[0])
		default:
			return r.detailedErr(stack.Stacks[0])
		}
	}
}

func (r *Runner) detailedErr(stack *Stack) error {
	events, err := r.failedEvents(*stack.StackId)
	if err != nil {
		return fmt.Errorf("getting failed events: %w", err)
	}

	failures := make([]string, len(events))
	for i, e := range events {
		failures[i] = e.String()
	}

	reason := "unknown"
	if stack.StackStatusReason != nil {
		reason = *stack.StackStatusReason
	}

	return fmt.Errorf("stack: %s, failed events: %s",
		reason,
		strings.Join(failures, "\n"),
	)
}

// StackEvent state
type StackEvent struct {
	Status string
	Reason string
	Type   string
}

// String returns a string representation
func (e StackEvent) String() string {
	return fmt.Sprintf("%s: %s", e.Type, e.Reason)
}

func (r *Runner) failedEvents(stackName string) ([]StackEvent, error) {
	opts := &cfPkg.DescribeStackEventsInput{
		StackName: &stackName,
	}

	var failedEvents []StackEvent

	for {
		output, err := r.Provider.CloudFormation().DescribeStackEvents(opts)
		if err != nil {
			return nil, err
		}

		for _, event := range output.StackEvents {
			if *event.ResourceStatus == cfPkg.ResourceStatusCreateFailed {
				failedEvents = append(failedEvents, StackEvent{
					Status: *event.ResourceStatus,
					Reason: *event.ResourceStatusReason,
					Type:   *event.ResourceType,
				})
			}
		}

		if opts.NextToken == nil {
			break
		}

		opts.NextToken = output.NextToken
	}

	return failedEvents, nil
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
