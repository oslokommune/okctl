// Package testutils provides test utilities for upgrade
package testutils

import "github.com/oslokommune/okctl/pkg/client"

type upgradeStateMock struct {
	upgrades        map[string]*client.Upgrade
	originalVersion string
}

func (m *upgradeStateMock) GetUpgrades() ([]*client.Upgrade, error) {
	upgrades := make([]*client.Upgrade, len(m.upgrades))

	i := 0

	for _, u := range m.upgrades {
		upgrades[i] = u
		i++
	}

	return upgrades, nil
}

func (m *upgradeStateMock) SaveUpgrade(upgrade *client.Upgrade) error {
	m.upgrades[upgrade.Version] = upgrade
	return nil
}

func (m *upgradeStateMock) GetUpgrade(version string) (*client.Upgrade, error) {
	u, ok := m.upgrades[version]
	if !ok {
		return nil, client.ErrUpgradeNotFound
	}

	return u, nil
}

func (m *upgradeStateMock) SaveOriginalClusterVersionIfNotExists(originalOkctlVersion *client.OriginalClusterVersion) error {
	if len(m.originalVersion) == 0 {
		m.originalVersion = originalOkctlVersion.Value
	}

	return nil
}

func (m *upgradeStateMock) GetOriginalClusterVersion() (*client.OriginalClusterVersion, error) {
	if len(m.originalVersion) == 0 {
		return nil, client.ErrOriginalClusterVersionNotFound
	}

	return &client.OriginalClusterVersion{
		Value: m.originalVersion,
	}, nil
}

//goland:noinspection GoUnusedParameter
func (m *upgradeStateMock) SaveClusterVersionInfo(version *client.ClusterVersion) error {
	panic("implement me")
}

func (m *upgradeStateMock) GetClusterVersionInfo() (*client.ClusterVersion, error) {
	panic("implement me")
}

// MockUpgradeState returns a mocked upgrade state
func MockUpgradeState() client.UpgradeState {
	return &upgradeStateMock{
		upgrades: make(map[string]*client.Upgrade),
	}
}
