// Package policydocument implements the IAM policy document:
// - https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_policies.html
// Note: we have not attempted to implement the complete logic, this functionality
// only covers the subset we require
package policydocument

import (
	"encoding/json"

	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/google/uuid"
)

// EffectType enumerates valid effects a policy has on a resource:
// - https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_policies_elements_effect.html
type EffectType string

// ConditionOperatorType enumerates valid condition operators:
// - https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_policies_elements_condition_operators.html
type ConditionOperatorType string

const (
	// Version is the current version of the policy language,
	// and you should always include a Version element and set it to 2012-10-17:
	// - https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_policies_elements_version.html
	Version = "2012-10-17"

	// EffectTypeAllow allows access to a resource
	EffectTypeAllow EffectType = "Allow"
	// EffectTypeDeny denies access to a resource
	EffectTypeDeny EffectType = "Deny"

	// ConditionOperatorTypeStringEquals checks if the string matches exactly
	ConditionOperatorTypeStringEquals ConditionOperatorType = "StringEquals"
	// ConditionOperatorTypeNull checks if the key value exists
	ConditionOperatorTypeNull ConditionOperatorType = "Null"

	// Pseudo parameters enumerates some of the available pseudo parameters:
	// - https://docs.amazonaws.cn/en_us/AWSCloudFormation/latest/UserGuide/pseudo-parameter-reference.html

	// PseudoParamRegion will return a string representing the AWS region
	PseudoParamRegion string = "AWS::Region"
	// PseudoParamAccountID will return a string containing the AWS account id
	PseudoParamAccountID string = "AWS::AccountId"
)

// PolicyDocument provides some structure around IAM policy documents:
// - https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_policies.html
type PolicyDocument struct {
	Version   string           `json:"Version"`
	ID        string           `json:"Id,omitempty"`
	Statement []StatementEntry `json:"Statement"`
}

// StatementEntry ...
type StatementEntry struct {
	Sid       string                                      `json:"Sid,omitempty"`
	Effect    EffectType                                  `json:"Effect"`
	Action    []string                                    `json:"Action"`
	Resource  []string                                    `json:"Resource"`
	Condition map[ConditionOperatorType]map[string]string `json:"Condition,omitempty"`
}

// JSON returns the json marshalled version of the policy document
func (d *PolicyDocument) JSON() ([]byte, error) {
	return json.MarshalIndent(d, "", "    ")
}

// ID creates a UUID for use with the policy document id field:
// - https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_policies_elements_id.html
func ID() string {
	return uuid.New().String()
}

// AwsRegionRef returns a reference to the AWS region pseudo parameter
func AwsRegionRef() string {
	return cloudformation.Ref(PseudoParamRegion)
}

// AwsAccountIDRef returns a reference to the AWS account ID pseudo parameter
func AwsAccountIDRef() string {
	return cloudformation.Ref(PseudoParamAccountID)
}
