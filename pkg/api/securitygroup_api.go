package api

import (
	"context"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// RuleType defines wether the rule is an inbound or outbound type rule
type RuleType string

const (
	minimumSecurityGroupDescriptionLength = 0
	maximumSecurityGroupDescriptionLength = 255
	minimumPort                           = 0
	maximumPort                           = 65535
	// RuleTypeIngress represents an inbound rule
	RuleTypeIngress RuleType = "Ingress"
	// RuleTypeEgress represents an outbound rule
	RuleTypeEgress RuleType = "Egress"
	// RuleProtocolAll configures a rule to represent all types of traffic
	RuleProtocolAll = "-1"
	// RuleProtocolTCP configures a rule to represent TCP based traffic
	RuleProtocolTCP = "tcp"
	// RuleProtocolUDP configures a rule to represent UDP based traffic
	RuleProtocolUDP = "udp"
)

var (
	reIPv4CIDR        = regexp.MustCompile(`^([0-9]{1,3}\.){3}[0-9]{1,3}(/([0-9]|[1-2][0-9]|3[0-2]))$`)
	reSecurityGroupID = regexp.MustCompile(`sg-[0-9a-z]+`)
	// Regexp interpreted from https://docs.aws.amazon.com/vpc/latest/userguide/VPC_SecurityGroups.html
	// reSGName = regexp.MustCompile(`[a-zA-Z0-9\s._\-:/()#,@[]+=&;{}!$*]+`)
)

// Rule defines an opening in a Security Group
type Rule struct {
	Description           string `json:"Description"`
	FromPort              int    `json:"FromPort"`
	ToPort                int    `json:"ToPort"`
	CidrIP                string `json:"CidrIp,omitempty"`
	Protocol              string `json:"IpProtocol"`
	SourceSecurityGroupID string `json:"SourceSecurityGroupId,omitempty"`
}

// Validate ensures the required data is existent and correct
func (r Rule) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Description,
			validation.Length(minimumSecurityGroupDescriptionLength, maximumSecurityGroupDescriptionLength),
		),
		validation.Field(&r.FromPort, validation.Min(minimumPort), validation.Max(maximumPort)),
		validation.Field(&r.ToPort, validation.Min(minimumPort), validation.Max(maximumPort)),
		validation.Field(&r.Protocol, validation.In(RuleProtocolTCP, RuleProtocolUDP, RuleProtocolAll), validation.Required),
		validation.Field(&r.CidrIP,
			validation.Match(reIPv4CIDR),
			validation.When(r.SourceSecurityGroupID == "",
				validation.Required.Error("required when SourceSecurityGroupID is empty"),
			).Else(validation.Empty.Error("must be blank if SourceSecurityGroupID is specified")),
		),
		validation.Field(&r.SourceSecurityGroupID,
			validation.Match(reSecurityGroupID),
			validation.When(r.CidrIP == "",
				validation.Required.Error("required when CidrIP is empty"),
			).Else(validation.Empty.Error("must be blank if CidrIP is specified")),
		),
	)
}

// Equals knows if two rules can be considered equal
func (r Rule) Equals(target Rule) bool {
	if r.CidrIP != target.CidrIP {
		return false
	}

	if r.SourceSecurityGroupID != target.SourceSecurityGroupID {
		return false
	}

	if r.Protocol != target.Protocol {
		return false
	}

	if r.FromPort != target.FromPort {
		return false
	}

	if r.ToPort != target.ToPort {
		return false
	}

	return true
}

// SecurityGroup defines an AWS Security Group
type SecurityGroup struct {
	ID            string
	InboundRules  []Rule
	OutboundRules []Rule
}

// Validate ensures the required data is existent and correct
func (s SecurityGroup) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.ID, validation.Required),
	)
}

// CreateSecurityGroupOpts defines required data to create a Security Group
type CreateSecurityGroupOpts struct {
	ClusterID     ID
	VPCID         string
	Name          string
	Description   string
	InboundRules  []Rule
	OutboundRules []Rule
}

// Validate ensures the required data is existent and correct
func (s CreateSecurityGroupOpts) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.Name, validation.Required),
		validation.Field(&s.InboundRules, validation.NotNil),
		validation.Field(&s.OutboundRules, validation.NotNil),
	)
}

// DeleteSecurityGroupOpts defines required data to delete a Security Group
type DeleteSecurityGroupOpts struct {
	ClusterName string
	Name        string
}

// Validate ensures the required data is existent and correct
func (d DeleteSecurityGroupOpts) Validate() error {
	return validation.ValidateStruct(&d,
		validation.Field(&d.Name, validation.Required),
	)
}

// GetSecurityGroupOpts defines required data for fetching an existing security group
type GetSecurityGroupOpts struct {
	Name        string
	ClusterName string
}

// Validate ensures the required data is existent and correct
func (s GetSecurityGroupOpts) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.Name, validation.Required),
		validation.Field(&s.ClusterName, validation.Required),
	)
}

// AddRuleOpts defines required data for adding a rule to an existing security group
type AddRuleOpts struct {
	ClusterName               string
	SecurityGroupStackName    string
	SecurityGroupResourceName string
	RuleType                  RuleType
	Rule                      Rule
}

// Validate ensures the required data is existent and correct
func (a AddRuleOpts) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.ClusterName, validation.Required),
		validation.Field(&a.SecurityGroupStackName, validation.Required),
		validation.Field(&a.RuleType, validation.Required, validation.In(RuleTypeIngress, RuleTypeEgress)),
		validation.Field(&a.Rule, validation.Required),
	)
}

// RemoveRuleOpts defines required data for removing a rule from an existing security group
type RemoveRuleOpts struct {
	ClusterName               string
	SecurityGroupStackName    string
	SecurityGroupResourceName string
	RuleType                  RuleType
	Rule                      Rule
}

// Validate ensures the required data is existent and correct
func (a RemoveRuleOpts) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.ClusterName, validation.Required),
		validation.Field(&a.SecurityGroupStackName, validation.Required),
		validation.Field(&a.SecurityGroupResourceName, validation.Required),
		validation.Field(&a.RuleType, validation.Required, validation.In(RuleTypeIngress, RuleTypeEgress)),
		validation.Field(&a.Rule, validation.Required),
	)
}

// SecurityGroupCRUDer defines CRUD operations available for Security Groups
type SecurityGroupCRUDer interface {
	CreateSecurityGroup(ctx context.Context, opts CreateSecurityGroupOpts) (SecurityGroup, error)
	GetSecurityGroup(ctx context.Context, opts GetSecurityGroupOpts) (SecurityGroup, error)
	DeleteSecurityGroup(ctx context.Context, opts DeleteSecurityGroupOpts) error
}

// SecurityGroupRuleCRUDer defines CRUD operations available for rules on existing Security Groups
type SecurityGroupRuleCRUDer interface {
	AddRule(ctx context.Context, opts AddRuleOpts) (Rule, error)
	RemoveRule(ctx context.Context, opts RemoveRuleOpts) error
}

// SecurityGroupService provides the service layer
type SecurityGroupService interface {
	SecurityGroupCRUDer
	SecurityGroupRuleCRUDer
}

// SecurityGroupCloudProvider provides the cloud provider layer
type SecurityGroupCloudProvider interface {
	SecurityGroupCRUDer
	SecurityGroupRuleCRUDer
}
