package v1alpha1

import (
	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
)

type CloudProvider interface {
	CloudFormation() cloudformationiface.CloudFormationAPI
	Region() string
}
