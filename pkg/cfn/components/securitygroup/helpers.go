package securitygroup

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/oslokommune/okctl/pkg/api"

	"sigs.k8s.io/yaml"
)

type serializableRule struct {
	Description           string     `json:"Description"`
	FromPort              int        `json:"FromPort"`
	ToPort                int        `json:"ToPort"`
	CidrIP                string     `json:"CidrIp,omitempty"`
	Protocol              string     `json:"IpProtocol"`
	SourceSecurityGroupID stringOrFn `json:"SourceSecurityGroupId,omitempty"`
}

func acquireRuleIndex(template []byte, ruleType string, resourceName string, rule api.Rule) (int, error) {
	var temp struct {
		Resources map[string]struct {
			Properties struct {
				Ingresses []serializableRule `json:"SecurityGroupIngress"`
				Egresses  []serializableRule `json:"SecurityGroupEgress"`
			} `json:"Properties"`
		} `json:"Resources"`
	}

	err := yaml.Unmarshal(template, &temp)
	if err != nil {
		return -1, fmt.Errorf("unmarshalling template: %w", err)
	}

	var elements []api.Rule

	switch ruleType {
	case "Ingress":
		elements = make([]api.Rule, len(temp.Resources[resourceName].Properties.Ingresses))

		for index, r := range temp.Resources[resourceName].Properties.Ingresses {
			elements[index] = serializableToAPIRule(r)
		}
	case "Egress":
		elements = make([]api.Rule, len(temp.Resources[resourceName].Properties.Egresses))

		for index, r := range temp.Resources[resourceName].Properties.Egresses {
			elements[index] = serializableToAPIRule(r)
		}
	}

	for index, item := range elements {
		if rule.Equals(item) {
			return index, nil
		}
	}

	return -1, nil
}

func serializableToAPIRule(sr serializableRule) api.Rule {
	return api.Rule{
		Description:           sr.Description,
		FromPort:              sr.FromPort,
		ToPort:                sr.ToPort,
		CidrIP:                sr.CidrIP,
		Protocol:              sr.Protocol,
		SourceSecurityGroupID: sr.SourceSecurityGroupID.String(),
	}
}

/*
   Due to Cloudformation allowing values to be either a string or a reference (map[string][]string) i.e. an Fn::GetAtt.
   we needed a way to unmarshall both. The following code implements a way to handle both cases.
*/
const getAttKey = "Fn::GetAtt"

// stringOrFn enables parsing values that can be either a string value or a CFN function
type stringOrFn struct {
	StringValue string
	FnValue     map[string][]string
}

// UnmarshalJSON knows how to handle a value that can be either a string or a CFN function reference
func (s *stringOrFn) UnmarshalJSON(bytes []byte) error {
	var value interface{}

	if err := yaml.Unmarshal(bytes, &value); err != nil {
		return err
	}

	switch v := value.(type) {
	case string:
		s.StringValue = v
	case map[string]interface{}:
		items := make([]string, 0)

		for _, item := range v[getAttKey].([]interface{}) {
			items = append(items, item.(string))
		}

		s.FnValue = map[string][]string{getAttKey: items}
	default:
		return fmt.Errorf("parsing value: %v of type %v", v, reflect.TypeOf(v))
	}

	return nil
}

// String generates a string generalized for ease of comparison
func (s stringOrFn) String() string {
	if s.StringValue != "" {
		return s.StringValue
	}

	if _, ok := s.FnValue[getAttKey]; !ok {
		return ""
	}

	return fmt.Sprintf("Fn::GetAtt: %s", strings.Join(s.FnValue[getAttKey], ","))
}
