package commands

import (
	"fmt"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/oslokommune/okctl/pkg/upgrade/clusterversion"
	"github.com/oslokommune/okctl/pkg/version"
)

// ValidateBinaryVsClusterVersionErr contains the error message for invalid binary vs cluster version
const ValidateBinaryVsClusterVersionErr = "validating binary against cluster version: %w"

// SaveClusterVersionErr contains the error message for error with saving cluster version
const SaveClusterVersionErr = "saving cluster version: %w"

// ValidateBinaryEqualsClusterVersion is a wrapper for clusterVersioner.ValidateBinaryEqualsClusterVersion
func ValidateBinaryEqualsClusterVersion(o *okctl.Okctl) error {
	clusterID := api.ID{
		Region:       o.Declaration.Metadata.Region,
		AWSAccountID: o.Declaration.Metadata.AccountID,
		ClusterName:  o.Declaration.Metadata.Name,
	}

	clusterVersioner := clusterversion.New(
		o.Out,
		clusterID,
		o.StateHandlers(o.StateNodes()).Upgrade,
	)

	err := clusterVersioner.ValidateBinaryEqualsClusterVersion(version.GetVersionInfo().Version)
	if err != nil {
		return fmt.Errorf(ValidateBinaryVsClusterVersionErr, err)
	}

	return nil
}

// ValidateBinaryVersionNotLessThanClusterVersion is a wrapper for clusterVersioner.ValidateBinaryVersionNotLessThanClusterVersion
func ValidateBinaryVersionNotLessThanClusterVersion(o *okctl.Okctl) error {
	clusterID := api.ID{
		Region:       o.Declaration.Metadata.Region,
		AWSAccountID: o.Declaration.Metadata.AccountID,
		ClusterName:  o.Declaration.Metadata.Name,
	}

	clusterVersioner := clusterversion.New(
		o.Out,
		clusterID,
		o.StateHandlers(o.StateNodes()).Upgrade,
	)

	err := clusterVersioner.ValidateBinaryVersionNotLessThanClusterVersion(version.GetVersionInfo().Version)
	if err != nil {
		return fmt.Errorf(ValidateBinaryVsClusterVersionErr, err)
	}

	return nil
}
