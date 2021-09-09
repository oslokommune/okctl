package securitygroup

import (
	"fmt"

	jsonpatch "github.com/evanphx/json-patch/v5"
	"github.com/oslokommune/okctl/pkg/api"
	patchpkg "github.com/oslokommune/okctl/pkg/jsonpatch"
	"sigs.k8s.io/yaml"
)

type patchAppendNewRuleOpts struct {
	OriginalTemplate []byte
	ResourceName     string
	RuleType         string
	Rule             api.Rule
}

type patchRemoveRuleOpts struct {
	OriginalTemplate []byte
	ResourceName     string

	Rule     api.Rule
	RuleType string
}

// patchAppendNewRule appends a new rule to an existing Cloudformation template with a SecurityGroup resource
func patchAppendNewRule(opts patchAppendNewRuleOpts) ([]byte, error) {
	index, err := acquireRuleIndex(opts.OriginalTemplate, opts.RuleType, opts.ResourceName, opts.Rule)
	if err != nil {
		return nil, fmt.Errorf("acquiring rule index: %w", err)
	}

	if index != -1 {
		return opts.OriginalTemplate, nil
	}

	patchJSON, err := patchpkg.New().Add(patchpkg.Operation{
		Type:  patchpkg.OperationTypeAdd,
		Path:  fmt.Sprintf("/Resources/%s/Properties/SecurityGroup%s/0", opts.ResourceName, opts.RuleType),
		Value: opts.Rule,
	}).MarshalJSON()
	if err != nil {
		return nil, err
	}

	jsonData, err := yaml.YAMLToJSON(opts.OriginalTemplate)
	if err != nil {
		return nil, fmt.Errorf("converting json to yaml: %w", err)
	}

	patch, err := jsonpatch.DecodePatch(patchJSON)
	if err != nil {
		return nil, fmt.Errorf("decoding patch: %w", err)
	}

	modified, err := patch.Apply(jsonData)
	if err != nil {
		return nil, fmt.Errorf("applying patch: %w", err)
	}

	return yaml.JSONToYAML(modified)
}

// patchRemoveExistingRule removes an existing rule from a CloudFormation template with a SecurityGroup resource
func patchRemoveExistingRule(opts patchRemoveRuleOpts) ([]byte, error) {
	index, err := acquireRuleIndex(opts.OriginalTemplate, opts.RuleType, opts.ResourceName, opts.Rule)
	if err != nil {
		return nil, fmt.Errorf("acquiring rule index: %w", err)
	}

	jsonData, err := yaml.YAMLToJSON(opts.OriginalTemplate)
	if err != nil {
		return nil, fmt.Errorf("converting json to yaml: %w", err)
	}

	patchJSON, err := patchpkg.New().Add(patchpkg.Operation{
		Type: patchpkg.OperationTypeRemove,
		Path: fmt.Sprintf("/Resources/%s/Properties/SecurityGroup%s/%d", opts.ResourceName, opts.RuleType, index),
	}).MarshalJSON()
	if err != nil {
		return nil, err
	}

	patch, err := jsonpatch.DecodePatch(patchJSON)
	if err != nil {
		return nil, fmt.Errorf("decoding patch: %w", err)
	}

	modified, err := patch.Apply(jsonData)
	if err != nil {
		return nil, fmt.Errorf("applying patch: %w", err)
	}

	return yaml.JSONToYAML(modified)
}
