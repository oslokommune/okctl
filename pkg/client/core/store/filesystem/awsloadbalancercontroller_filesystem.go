package filesystem

// nolint: dupl

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/api"

	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/client/store"
	"github.com/spf13/afero"
)

type awsLoadBalancerControllerStore struct {
	policy         Paths
	serviceAccount Paths
	chart          Paths
	fs             *afero.Afero
}

func (s *awsLoadBalancerControllerStore) RemoveAWSLoadBalancerController(_ api.ID) (*store.Report, error) {
	report, err := store.NewFileSystem(s.policy.BaseDir, s.fs).
		Remove(s.policy.OutputFile).
		Remove(s.policy.CloudFormationFile).
		AlterStore(store.SetBaseDir(s.serviceAccount.BaseDir)).
		Remove(s.serviceAccount.OutputFile).
		Remove(s.serviceAccount.ConfigFile).
		AlterStore(store.SetBaseDir(s.chart.BaseDir)).
		Remove(s.chart.OutputFile).
		Remove(s.chart.ReleaseFile).
		Remove(s.chart.ChartFile).
		RemoveDir("").
		Do()
	if err != nil {
		return nil, err
	}

	return report, nil
}

// NewAWSLoadBalancerControllerStore returns an initialised store
func NewAWSLoadBalancerControllerStore(policy, serviceAccount, chart Paths, fs *afero.Afero) client.AWSLoadBalancerControllerStore {
	return &awsLoadBalancerControllerStore{
		policy:         policy,
		serviceAccount: serviceAccount,
		chart:          chart,
		fs:             fs,
	}
}

func (s *awsLoadBalancerControllerStore) SaveAWSLoadBalancerController(c *client.AWSLoadBalancerController) (*store.Report, error) {
	policy := &ManagedPolicy{
		ID:        c.Policy.ID,
		StackName: c.Policy.StackName,
		PolicyARN: c.Policy.PolicyARN,
	}

	account := &ServiceAccount{
		ID:        c.ServiceAccount.ID,
		PolicyArn: c.ServiceAccount.PolicyArn,
	}

	chart := &Helm{
		ID: c.Chart.ID,
	}

	report, err := store.NewFileSystem(s.policy.BaseDir, s.fs).
		// Policy
		StoreStruct(s.policy.OutputFile, policy, store.ToJSON()).
		StoreBytes(s.policy.CloudFormationFile, c.Policy.CloudFormationTemplate).
		// ServiceAccount
		AlterStore(store.SetBaseDir(s.serviceAccount.BaseDir)).
		StoreStruct(s.serviceAccount.OutputFile, account, store.ToJSON()).
		StoreStruct(s.serviceAccount.ConfigFile, c.ServiceAccount.Config, store.ToJSON()).
		// Chart
		AlterStore(store.SetBaseDir(s.chart.BaseDir)).
		StoreStruct(s.chart.OutputFile, chart, store.ToJSON()).
		StoreStruct(s.chart.ReleaseFile, c.Chart.Release, store.ToJSON()).
		StoreStruct(s.chart.ChartFile, c.Chart.Chart, store.ToJSON()).
		Do()
	if err != nil {
		return nil, fmt.Errorf("failed to store aws load balancer controller: %w", err)
	}

	return report, nil
}
