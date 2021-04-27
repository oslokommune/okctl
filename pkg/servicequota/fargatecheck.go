package servicequota

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/service/cloudwatch"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/servicequotas"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
)

const (
	evaluationPeriodInHours = 1
)

// FargateCheck is used to check if you have enough Elastic Ips
type FargateCheck struct {
	provider v1alpha1.CloudProvider
	required int
}

// NewFargateCheck makes a new instance of check for Elastic Ips
func NewFargateCheck(required int, provider v1alpha1.CloudProvider) *FargateCheck {
	return &FargateCheck{
		provider: provider,
		required: required,
	}
}

// CheckAvailability determines if you have sufficient fargate pods
// nolint: funlen
func (e *FargateCheck) CheckAvailability() (*Result, error) {
	q, err := e.provider.ServiceQuotas().GetServiceQuota(&servicequotas.GetServiceQuotaInput{
		QuotaCode:   aws.String("L-790AF391"),
		ServiceCode: aws.String("fargate"),
	})
	if err != nil {
		return nil, fmt.Errorf("getting fargate on-demand pods quotas: %w", err)
	}

	now := time.Now()

	rounded := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		now.Hour(),
		0,
		0,
		0,
		now.Location(),
	)

	end := rounded
	start := rounded.Add(-evaluationPeriodInHours * time.Hour)
	period := int64(evaluationPeriodInHours * time.Hour.Seconds())

	data, err := e.provider.CloudWatch().GetMetricStatistics(&cloudwatch.GetMetricStatisticsInput{
		Dimensions: []*cloudwatch.Dimension{
			{
				Name:  aws.String("Service"),
				Value: aws.String("Fargate"),
			},
			{
				Name:  aws.String("Type"),
				Value: aws.String("Resource"),
			},
			{
				Name:  aws.String("Resource"),
				Value: aws.String("OnDemand"),
			},
			{
				Name:  aws.String("Class"),
				Value: aws.String("None"),
			},
		},
		EndTime:    aws.Time(end),
		MetricName: aws.String("ResourceCount"),
		Namespace:  aws.String("AWS/Usage"),
		Period:     aws.Int64(period),
		StartTime:  aws.Time(start),
		Statistics: []*string{
			aws.String("Maximum"),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("getting fargate on-demand pods utilisation: %w", err)
	}

	quota := int(*q.Quota.Value)
	count := int(*data.Datapoints[0].Maximum)
	available := quota - count

	return &Result{
		Required:    e.required,
		Available:   available,
		HasCapacity: e.required <= available,
		Description: "Fargate On-Demand Pods",
	}, nil
}
