package filesystem

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/client/store"
	"github.com/spf13/afero"
)

type externalSecretsStore struct {
	policy         Paths
	serviceAccount Paths
	chart          Paths
	fs             *afero.Afero
}

func (s *externalSecretsStore) SaveExternalSecrets(e *client.ExternalSecrets) (*store.Report, error) {
	policy := &ManagedPolicy{
		ID:        e.Policy.ID,
		StackName: e.Policy.StackName,
		PolicyARN: e.Policy.PolicyARN,
	}

	account := &ServiceAccount{
		ID:        e.ServiceAccount.ID,
		PolicyArn: e.ServiceAccount.PolicyArn,
	}

	chart := &Helm{
		ID: e.Chart.ID,
	}

	report, err := store.NewFileSystem(s.policy.BaseDir, s.fs).
		// Policy
		StoreStruct(s.policy.OutputFile, policy, store.ToJSON()).
		StoreBytes(s.policy.CloudFormationFile, e.Policy.CloudFormationTemplate).
		// ServiceAccount
		AlterStore(store.SetBaseDir(s.serviceAccount.BaseDir)).
		StoreStruct(s.serviceAccount.OutputFile, account, store.ToJSON()).
		StoreStruct(s.serviceAccount.ConfigFile, e.ServiceAccount.Config, store.ToJSON()).
		// Chart
		AlterStore(store.SetBaseDir(s.chart.BaseDir)).
		StoreStruct(s.chart.OutputFile, chart, store.ToJSON()).
		StoreStruct(s.chart.ReleaseFile, e.Chart.Release, store.ToJSON()).
		StoreStruct(s.chart.ChartFile, e.Chart.Chart, store.ToJSON()).
		Do()
	if err != nil {
		return nil, fmt.Errorf("failed to store external secrets: %w", err)
	}

	return report, nil
}

// NewExternalSecretsStore returns an initialised store
func NewExternalSecretsStore(policy, serviceAccount, chart Paths, fs *afero.Afero) client.ExternalSecretsStore {
	return &externalSecretsStore{
		policy:         policy,
		serviceAccount: serviceAccount,
		chart:          chart,
		fs:             fs,
	}
}
