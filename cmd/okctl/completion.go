package main

import (
	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/spf13/cobra"
)

// nolint: funlen
func buildCompletionCommand(o *okctl.Okctl) *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "completion [bash|zsh|fish|powershell]",
		Short:                 CompletionShortDescription,
		Long:                  CompletionLongDescription,
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.ExactValidArgs(1),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			var err error

			switch args[0] {
			case "bash":
				err = cmd.Root().GenBashCompletion(o.Out)
			case "zsh":
				err = cmd.Root().GenZshCompletion(o.Out)
			case "fish":
				err = cmd.Root().GenFishCompletion(o.Out, true)
			case "powershell":
				err = cmd.Root().GenPowerShellCompletion(o.Out)
			}

			return err
		},
	}

	return cmd
}
