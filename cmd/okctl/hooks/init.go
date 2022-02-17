package hooks

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/oslokommune/okctl/cmd/okctl/auth"
	"github.com/oslokommune/okctl/pkg/config/load"
	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/spf13/cobra"
)

// InitializeEnvironment - initialize everything required for ALL operations
func InitializeEnvironment(o *okctl.Okctl) RunEer {
	return func(cmd *cobra.Command, _ []string) error {
		if cmd.Name() == cobra.ShellCompRequestCmd {
			return nil
		}

		var err error

		err = initializeAuthentication(o)
		if err != nil {
			return fmt.Errorf("initializing authentication: %w", err)
		}

		err = loadUserData(o, cmd)
		if err != nil {
			return fmt.Errorf("loading application data: %w", err)
		}

		initializeIo(o, cmd)

		return nil
	}
}

func loadUserData(o *okctl.Okctl, cmd *cobra.Command) error {
	userDataNotFound := load.CreateOnUserDataNotFound()

	o.UserDataLoader = load.UserDataFromFlagsEnvConfigDefaults(cmd, userDataNotFound)

	return o.LoadUserData()
}

func initializeAuthentication(o *okctl.Okctl) error {
	err := auth.ValidateCredentialTypes()
	if err != nil {
		return err
	}

	auth.EnableServiceUserAuthentication(o)

	return nil
}

func initializeIo(o *okctl.Okctl, cmd *cobra.Command) {
	o.Out = cmd.OutOrStdout()
	o.Err = cmd.OutOrStderr()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM) //nolint

	go func() {
		<-c

		err := GracefullyTearDownState(o)
		if err != nil {
			fmt.Fprintf(cmd.ErrOrStderr(), "gracefully up state: %s", err.Error())
		}

		os.Exit(1)
	}()
}
