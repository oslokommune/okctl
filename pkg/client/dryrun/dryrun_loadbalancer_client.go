package dryrun

import (
	"context"
	"fmt"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
	"io"
)

type loadbalancerService struct {
	out io.Writer
}

func (l loadbalancerService) CreateAWSLoadBalancerController(_ context.Context, _ client.CreateAWSLoadBalancerControllerOpts) (*client.AWSLoadBalancerController, error) {
	fmt.Fprintf(l.out, formatCreate("AWS LoadBalancer controller"))

	return &client.AWSLoadBalancerController{}, nil
}

func (l loadbalancerService) DeleteAWSLoadBalancerController(_ context.Context, _ api.ID) error {
	fmt.Fprintf(l.out, formatDelete("AWS LoadBalancer controller"))

	return nil
}
