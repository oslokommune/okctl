package handlers

import "github.com/spf13/cobra"

// RunEHandler defines the interface for cobra command RunE handlers
type RunEHandler func(*cobra.Command, []string) error
