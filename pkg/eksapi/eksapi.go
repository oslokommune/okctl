// Package eksapi provides some convenience functionality
// for retrieving information about an AWS EKS cluster
package eksapi

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/config/constant"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
)

// EKSAPI contains the required state for interacting
// with the AWS EKS API
type EKSAPI struct {
	clusterName string
	provider    v1alpha1.CloudProvider
}

// New returns an initialised EKS API client
func New(clusterName string, provider v1alpha1.CloudProvider) *EKSAPI {
	return &EKSAPI{
		clusterName: clusterName,
		provider:    provider,
	}
}

// FargateProfilePodExecutionRoleARN retrieves the Fargate profile pod execution role ARN
func (f *EKSAPI) FargateProfilePodExecutionRoleARN(profile string) (string, error) {
	p, err := f.provider.EKS().DescribeFargateProfile(&eks.DescribeFargateProfileInput{
		ClusterName:        aws.String(f.clusterName),
		FargateProfileName: aws.String(profile),
	})
	if err != nil {
		return "", fmt.Errorf(constant.GetFargateProfileError, err)
	}

	return *p.FargateProfile.PodExecutionRoleArn, nil
}
