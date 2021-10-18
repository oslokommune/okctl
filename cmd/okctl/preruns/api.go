package preruns

import (
	"fmt"

	"github.com/oslokommune/okctl/pkg/config/load"
	"github.com/oslokommune/okctl/pkg/metrics"
	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/spf13/cobra"
)

// PreRunECombinator takes an arbitrary amount of preRunE functions and runs them all
func PreRunECombinator(preRunEers ...PreRunEer) PreRunEer {
	return func(cmd *cobra.Command, args []string) error {
		for _, runEer := range preRunEers {
			err := runEer(cmd, args)
			if err != nil {
				return err
			}
		}

		return nil
	}
}

// LoadUserData ensures the o.UserData is loaded or created
func LoadUserData(o *okctl.Okctl) PreRunEer {
	return func(cmd *cobra.Command, args []string) error {
		userDataNotFound := load.CreateOnUserDataNotFoundWithNoInput()

		o.UserDataLoader = load.UserDataFromFlagsEnvConfigDefaults(cmd, userDataNotFound)

		return o.LoadUserData()
	}
}

// InitializeMetrics initializes required metrics data
func InitializeMetrics(o *okctl.Okctl) PreRunEer {
	return func(cmd *cobra.Command, args []string) error {
		metrics.SetUserAgent(o.UserState.Metrics.UserAgent)
		metrics.SetMetricsOut(o.Out)

		err := metrics.SetAPIURL(o.UserState.Metrics.APIURL)
		if err != nil {
			return fmt.Errorf("setting metrics API URL: %w", err)
		}

		return nil
	}
}
