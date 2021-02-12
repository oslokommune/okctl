package main

import (
	"github.com/oslokommune/okctl/pkg/okctl"
	"github.com/spf13/cobra"
)

// nolint: funlen
func buildCompletionCommand(o *okctl.Okctl) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate completion script",
		Long: `To load completions:

Bash:

	$ source <(okctl completion bash)

	# To load completions for each session, execute once:
	# Linux:
	$ okctl completion bash > /etc/bash_completion.d/okctl
	# macOS:
	$ okctl completion bash > /usr/local/etc/bash_completion.d/okctl

Zsh:

	# If shell completion is not already enabled in your environment,
	# you will need to enable it.  You can execute the following once:

	$ echo "autoload -U compinit; compinit" >> ~/.zshrc

	# To load completions for each session, execute once:
	$ okctl completion zsh > "${fpath[1]}/_okctl"

	# You will need to start a new shell for this setup to take effect.

fish:

	$ okctl completion fish | source

	# To load completions for each session, execute once:
	$ okctl completion fish > ~/.config/fish/completions/okctl.fish

PowerShell:

	PS> okctl completion powershell | Out-String | Invoke-Expression

	# To load completions for every new session, run:
	PS> okctl completion powershell > okctl.ps1
	# and source this file from your PowerShell profile.
	`,
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
