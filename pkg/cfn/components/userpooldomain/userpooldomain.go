// Package userpooldomain provides functionality for setting
// up a domain with a user pool
// - https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-cognito-userpooldomain.html
package userpooldomain

import (
	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/awslabs/goformation/v4/cloudformation/cognito"
	"github.com/oslokommune/okctl/pkg/cfn"
)

// UserPoolDomain stores the state for a cloud formation
// cognito user pool domain
type UserPoolDomain struct {
	StoredName        string
	Domain            string
	UserPool          cfn.NameReferencer
	PlaceHolderRecord cfn.Namer
	CertificateARN    string
}

// Resource returns the cloud formation resource for a
// cognito user pool domain
func (d *UserPoolDomain) Resource() cloudformation.Resource {
	return &cognito.UserPoolDomain{
		CustomDomainConfig: &cognito.UserPoolDomain_CustomDomainConfigType{
			CertificateArn: d.CertificateARN,
		},
		Domain:     d.Domain,
		UserPoolId: d.UserPool.Ref(),
		AWSCloudFormationDependsOn: []string{
			d.UserPool.Name(),
			d.PlaceHolderRecord.Name(),
		},
	}
}

// Name returns the logical id of the resource
func (d *UserPoolDomain) Name() string {
	return d.StoredName
}

// New returns an initialised cognito user pool domain
func New(domain, certificateARN string, userPool cfn.NameReferencer, placeHolderRecord cfn.Namer) *UserPoolDomain {
	return &UserPoolDomain{
		StoredName:        "UserPoolDomain",
		Domain:            domain,
		UserPool:          userPool,
		CertificateARN:    certificateARN,
		PlaceHolderRecord: placeHolderRecord,
	}
}
