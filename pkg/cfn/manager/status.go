package manager

import (
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

// This is verbatim stolen from: https://github.com/weaveworks/eksctl/blob/master/pkg/cfn/manager/api.go

// StackStatusIsNotTransitional returns true if stack is in a steady state
func (*Manager) StackStatusIsNotTransitional(s *Stack) bool {
	for _, state := range nonTransitionalReadyStackStatuses() {
		if *s.StackStatus == state {
			return true
		}
	}

	return false
}

func nonTransitionalReadyStackStatuses() []string {
	return []string{
		cloudformation.StackStatusCreateComplete,
		cloudformation.StackStatusUpdateComplete,
		cloudformation.StackStatusRollbackComplete,
		cloudformation.StackStatusUpdateRollbackComplete,
	}
}

// StackStatusIsNotReady returns true if the stack is in a transitional state
func (*Manager) StackStatusIsNotReady(s *Stack) bool {
	for _, state := range nonReadyStackStatuses() {
		if *s.StackStatus == state {
			return true
		}
	}

	return false
}

func nonReadyStackStatuses() []string {
	return []string{
		cloudformation.StackStatusCreateInProgress,
		cloudformation.StackStatusCreateFailed,
		cloudformation.StackStatusRollbackInProgress,
		cloudformation.StackStatusRollbackFailed,
		cloudformation.StackStatusDeleteInProgress,
		cloudformation.StackStatusDeleteFailed,
		cloudformation.StackStatusUpdateInProgress,
		cloudformation.StackStatusUpdateCompleteCleanupInProgress,
		cloudformation.StackStatusUpdateRollbackInProgress,
		cloudformation.StackStatusUpdateRollbackFailed,
		cloudformation.StackStatusUpdateRollbackCompleteCleanupInProgress,
		cloudformation.StackStatusReviewInProgress,
	}
}

// StackStatusIsNotDeleted returns true if the stack exists in some form
func (*Manager) StackStatusIsNotDeleted(s *Stack) bool {
	for _, state := range allNonDeletedStackStatuses() {
		if *s.StackStatus == state {
			return true
		}
	}

	return false
}

func allNonDeletedStackStatuses() []string {
	return []string{
		cloudformation.StackStatusCreateInProgress,
		cloudformation.StackStatusCreateFailed,
		cloudformation.StackStatusCreateComplete,
		cloudformation.StackStatusRollbackInProgress,
		cloudformation.StackStatusRollbackFailed,
		cloudformation.StackStatusRollbackComplete,
		cloudformation.StackStatusDeleteInProgress,
		cloudformation.StackStatusDeleteFailed,
		cloudformation.StackStatusUpdateInProgress,
		cloudformation.StackStatusUpdateCompleteCleanupInProgress,
		cloudformation.StackStatusUpdateComplete,
		cloudformation.StackStatusUpdateRollbackInProgress,
		cloudformation.StackStatusUpdateRollbackFailed,
		cloudformation.StackStatusUpdateRollbackCompleteCleanupInProgress,
		cloudformation.StackStatusUpdateRollbackComplete,
		cloudformation.StackStatusReviewInProgress,
	}
}
