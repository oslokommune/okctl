package hooks

import "github.com/spf13/cobra"

const logComponent = "hooks"

// RunEer defines the function signature the Cobra library expects
type RunEer func(cmd *cobra.Command, args []string) error
