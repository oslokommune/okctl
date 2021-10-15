package preruns

import "github.com/spf13/cobra"

// PreRunEer defines the function signature the Cobra library expects
type PreRunEer func(cmd *cobra.Command, args []string) error
