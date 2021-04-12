package storm

import (
	"errors"

	stormpkg "github.com/asdine/storm/v3"
	"github.com/oslokommune/okctl/pkg/client"
)

type monitoringState struct {
	node stormpkg.Node
}

// KubePromStack contains storm compatible state
type KubePromStack struct {
	Metadata `storm:"inline"`
	Name     string `storm:"unique"`

	ID                                ID
	AuthHostname                      string
	CertificateARN                    string
	ClientID                          string
	FargateCloudWatchPolicyARN        string
	FargateProfilePodExecutionRoleARN string
	Hostname                          string
	SecretsAdminPassKey               string
	SecretsAdminUserKey               string
	SecretsClientSecretKey            string
	SecretsConfigName                 string
	SecretsCookieSecretKey            string
}

// NewKubePromStack returns storm compatible state
func NewKubePromStack(s *client.KubePromStack, meta Metadata) *KubePromStack {
	return &KubePromStack{
		Metadata:                          meta,
		Name:                              "kubepromstack",
		ID:                                NewID(s.ID),
		AuthHostname:                      s.AuthHostname,
		CertificateARN:                    s.CertificateARN,
		ClientID:                          s.ClientID,
		FargateCloudWatchPolicyARN:        s.FargateCloudWatchPolicyARN,
		FargateProfilePodExecutionRoleARN: s.FargateProfilePodExecutionRoleARN,
		Hostname:                          s.Hostname,
		SecretsAdminPassKey:               s.SecretsAdminPassKey,
		SecretsAdminUserKey:               s.SecretsAdminUserKey,
		SecretsClientSecretKey:            s.SecretsClientSecretKey,
		SecretsConfigName:                 s.SecretsConfigName,
		SecretsCookieSecretKey:            s.SecretsCookieSecretKey,
	}
}

// Convert to client.KubePromStack
func (s *KubePromStack) Convert() *client.KubePromStack {
	return &client.KubePromStack{
		ID:                                s.ID.Convert(),
		AuthHostname:                      s.AuthHostname,
		CertificateARN:                    s.CertificateARN,
		ClientID:                          s.ClientID,
		FargateCloudWatchPolicyARN:        s.FargateCloudWatchPolicyARN,
		FargateProfilePodExecutionRoleARN: s.FargateProfilePodExecutionRoleARN,
		Hostname:                          s.Hostname,
		SecretsAdminPassKey:               s.SecretsAdminPassKey,
		SecretsAdminUserKey:               s.SecretsAdminUserKey,
		SecretsClientSecretKey:            s.SecretsClientSecretKey,
		SecretsConfigName:                 s.SecretsConfigName,
		SecretsCookieSecretKey:            s.SecretsCookieSecretKey,
	}
}

func (m *monitoringState) SaveKubePromStack(stack *client.KubePromStack) error {
	return m.node.Save(NewKubePromStack(stack, NewMetadata()))
}

func (m *monitoringState) RemoveKubePromStack() error {
	s := &KubePromStack{}

	err := m.node.One("Name", "kubepromstack", s)
	if err != nil {
		if errors.Is(err, stormpkg.ErrNotFound) {
			return nil
		}

		return err
	}

	return m.node.DeleteStruct(s)
}

func (m *monitoringState) GetKubePromStack() (*client.KubePromStack, error) {
	s := &KubePromStack{}

	err := m.node.One("Name", "kubepromstack", s)
	if err != nil {
		return nil, err
	}

	return s.Convert(), nil
}

// NewMonitoringState returns an initialised state client
func NewMonitoringState(node stormpkg.Node) client.MonitoringState {
	return &monitoringState{
		node: node,
	}
}
