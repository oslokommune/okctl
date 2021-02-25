package filesystem

import (
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/client/store"
	"github.com/spf13/afero"
)

// KubePromStack contains the state written to the outputs
type KubePromStack struct {
	ID                     api.ID
	CertificateARN         string
	Hostname               string
	AuthHostname           string
	ClientID               string
	SecretsConfigName      string
	SecretsCookieSecretKey string
	SecretsClientSecretKey string
	SecretsAdminUserKey    string
	SecretsAdminPassKey    string
}

type monitoringStore struct {
	lokiPaths         Paths
	kubePromPaths     Paths
	kubePromHelmPaths Paths
	fs                *afero.Afero
}

func (s *monitoringStore) SaveLoki(loki *client.Loki) (*store.Report, error) {
	chart := &Helm{
		ID: loki.ID,
	}

	report, err := store.NewFileSystem(s.lokiPaths.BaseDir, s.fs).
		StoreStruct(s.lokiPaths.OutputFile, chart, store.ToJSON()).
		StoreStruct(s.lokiPaths.ReleaseFile, loki.Chart.Release, store.ToJSON()).
		StoreStruct(s.lokiPaths.ChartFile, loki.Chart.Chart, store.ToJSON()).
		Do()
	if err != nil {
		return nil, err
	}

	return report, nil
}

func (s *monitoringStore) RemoveLoki(_ api.ID) (*store.Report, error) {
	return store.NewFileSystem(s.lokiPaths.BaseDir, s.fs).
		RemoveDir("").
		Do()
}

func (s *monitoringStore) RemoveKubePromStack(_ api.ID) (*store.Report, error) {
	return store.NewFileSystem(s.kubePromHelmPaths.BaseDir, s.fs).
		RemoveDir("").
		Do()
}

func (s *monitoringStore) SaveKubePromStack(stack *client.KubePromStack) (*store.Report, error) {
	kubeProm := &KubePromStack{
		ID:                     stack.ID,
		CertificateARN:         stack.CertificateARN,
		Hostname:               stack.Hostname,
		AuthHostname:           stack.AuthHostname,
		ClientID:               stack.ClientID,
		SecretsConfigName:      stack.SecretsConfigName,
		SecretsCookieSecretKey: stack.SecretsCookieSecretKey,
		SecretsClientSecretKey: stack.SecretsClientSecretKey,
		SecretsAdminUserKey:    stack.SecretsAdminUserKey,
		SecretsAdminPassKey:    stack.SecretsAdminPassKey,
	}

	chart := &Helm{
		ID: stack.ID,
	}

	report, err := store.NewFileSystem(s.kubePromPaths.BaseDir, s.fs).
		// Outputs
		StoreStruct(s.kubePromPaths.OutputFile, kubeProm, store.ToJSON()).
		// Chart
		AlterStore(store.SetBaseDir(s.kubePromHelmPaths.BaseDir)).
		StoreStruct(s.kubePromHelmPaths.OutputFile, chart, store.ToJSON()).
		StoreStruct(s.kubePromHelmPaths.ReleaseFile, stack.Chart.Release, store.ToJSON()).
		StoreStruct(s.kubePromHelmPaths.ChartFile, stack.Chart.Chart, store.ToJSON()).
		Do()
	if err != nil {
		return nil, err
	}

	return report, nil
}

// NewMonitoringStore returns an initialised store
func NewMonitoringStore(
	lokiPaths, kubePromHelmPaths, kubePromPaths Paths,
	fs *afero.Afero,
) client.MonitoringStore {
	return &monitoringStore{
		lokiPaths:         lokiPaths,
		kubePromPaths:     kubePromPaths,
		kubePromHelmPaths: kubePromHelmPaths,
		fs:                fs,
	}
}
