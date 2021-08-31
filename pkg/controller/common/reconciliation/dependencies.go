package reconciliation

import (
	"fmt"
	"github.com/oslokommune/okctl/pkg/config/constant"

	clientCore "github.com/oslokommune/okctl/pkg/client/core"
)

// DependencyTestFn defines a function which tests if a dependency is met
type DependencyTestFn func() (bool, error)

// AssertDependencyExistence asserts that the existence of all the provided tests is as expected
func AssertDependencyExistence(expectExistence bool, tests ...DependencyTestFn) (bool, error) {
	for _, test := range tests {
		actualExistence, err := test()
		if err != nil {
			return true, fmt.Errorf(constant.CheckDepedencyError, err)
		}

		if expectExistence != actualExistence {
			return false, nil
		}
	}

	return true, nil
}

// GenerateClusterExistenceTest is a convenience function for creating a cluster existence test function
func GenerateClusterExistenceTest(state *clientCore.StateHandlers, clusterName string) DependencyTestFn {
	return func() (bool, error) {
		return state.Cluster.HasCluster(clusterName)
	}
}

// GeneratePrimaryDomainDelegationTest is a convenience function for creating a dependency test checking for
// primary hosted zone delegation status
func GeneratePrimaryDomainDelegationTest(state *clientCore.StateHandlers) DependencyTestFn {
	return func() (bool, error) {
		hasPrimaryHostedZone, err := state.Domain.HasPrimaryHostedZone()
		if err != nil {
			return false, fmt.Errorf(constant.CheckIfPrimaryHostedZoneExistsError, err)
		}

		if !hasPrimaryHostedZone {
			return false, nil
		}

		domain, err := state.Domain.GetPrimaryHostedZone()
		if err != nil {
			return false, fmt.Errorf(constant.GetPrimaryHostedZoneError, err)
		}

		return domain.IsDelegated, nil
	}
}
